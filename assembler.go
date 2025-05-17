package assembler

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (a *Assembler) Assemble(filename string, outputFolder string) error {
	if a.Token == nil {
		a.Token = NewToken(global, "", nil)
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	//Preprocess File
	var lines []string = Preprocess(file)

	actualParent := a.Token
	for _, line := range lines {
		lineParts := strings.Split(line, " ")
		lineParts = removeEmptyStrings(lineParts)

		a.lineNumber++
		if len(lineParts) == 0 {
			continue
		}
		actualParent, err = a.Parse(lineParts, actualParent)
		if err != nil {
			panic("LINE " + strconv.Itoa(a.lineNumber) + " " + err.Error())
		}
	}

	str := printTokenTree(a.Token, 0)
	// Ensure output folder exists
	if outputFolder != "" {
		if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
			err = os.MkdirAll(outputFolder, 0755)
			if err != nil {
				// handle error (e.g. log and return)
				panic(fmt.Sprintf("Failed to create output folder: %v", err))
			}
		}
	}

	// Define output file path
	outputPath := filepath.Join(outputFolder, "output.parser")

	// Write to file
	ferr := os.WriteFile(outputPath, []byte(str), 0644)
	if ferr != nil {
		println(ferr.Error())
	} else {
		println("parsing saved to ./testfile/output.parser")
	}
	fmt.Println("blocks sizes")
	fmt.Print("instructions (4bytes/instr): ")
	fmt.Println(instructionCount)
	fmt.Print("instructions count: ")
	fmt.Println(instructionCount / 32)
	fmt.Print("variables: ")
	fmt.Println(variableCount)
	fmt.Print("constants: ")
	fmt.Println(constantCount)
	fmt.Print("strings: ")
	fmt.Println(stringCount)

	prog := compile(a.Token)
	bytes := BuildELFFile(prog)
	os.WriteFile("./testfile/output.exe", *bytes, 0644)
	return nil
}

func printTokenTree(t *Token, depth int) string {
	str := ""
	for i := 0; i < depth; i++ {
		str += "	"
	}
	if t.tokenType == instruction {
		str += getTType(t.tokenType) + " " + t.value + " " + getOpCode(t.opPair.opType) + "\n"
	} else {
		str += getTType(t.tokenType) + " " + t.value + "\n"
	}
	for _, child := range t.children {
		str += printTokenTree(child, depth+1)
	}
	return str
}

func removeEmptyStrings(arr []string) []string {
	result := make([]string, 0, len(arr))

	for _, str := range arr {
		if str != "" && str != "," {
			result = append(result, str)
		}
	}

	return result
}

func cleanupStr(str string) string {
	return strings.ReplaceAll(strings.TrimSpace(str), ",", "")
}

func (a *Assembler) Parse(lineParts []string, parent *Token) (*Token, error) {
	ln := cleanupStr(lineParts[0])

	if ln[0] == '.' {
		if parent.tokenType == global || parent.tokenType == section || parent.tokenType == globalLabel || parent.tokenType == constant {
			if ln == ".section" {
				if len(lineParts) == 2 {
					// section
					tk := NewToken(section, lineParts[1], a.Token)
					a.Token.children = append(a.Token.children, tk)
					return tk, nil
				} else {
					return parent, nil
				}
			}
			if parent.value == ".text" {
				if ln == ".globl" {
					//symbol main entry point
					tk := NewToken(entrypoint, ln, a.Token)
					tk.children = []*Token{NewToken(symbol, cleanupStr(lineParts[1]), tk)}
					parent.children = append(parent.children, tk)
					return parent, nil
				}
			} else if ln[len(ln)-1] == ':' {
				// constant
				tk := NewToken(constant, ln, a.Token)
				a.Token.children = append(a.Token.children, tk)
				if len(lineParts) > 1 {
					var finalLinePart2 string = strings.TrimSpace(strings.Join(lineParts[2:], " "))
					varS, isN := getVarSize(lineParts[1])
					if isN {
						constantCount += varS
					} else {
						finalLinePart2 = finalLinePart2[1 : len(finalLinePart2)-1]
						stringCount += len(finalLinePart2)*8 - 8
					}
					tk.children = []*Token{NewToken(varSize, cleanupStr(lineParts[1]), tk), NewToken(varValue, finalLinePart2, tk)}
				}
				return tk, nil
			} else if len(lineParts) == 1 {
				// section
				tk := NewToken(section, ln, a.Token)
				a.Token.children = append(a.Token.children, tk)
				return tk, nil
			} else {
				// variable defined with label on line before
				var finalLinePart2 string = strings.TrimSpace(strings.Join(lineParts[1:], " "))

				varS, isN := getVarSize(lineParts[0])
				if isN {
					if parent.tokenType == constant {
						constantCount += varS
					} else {
						variableCount += varS
					}
				} else {
					finalLinePart2 = finalLinePart2[1 : len(finalLinePart2)-1]
					stringCount += len(finalLinePart2)*8 - 8
				}
				if parent.tokenType == section {
					parent = parent.children[len(parent.children)-1]
					//change parent from globallabel (theoretically) to varLabel
					// if parent.parent is nil then we are at root.
				} else if parent.tokenType != constant && parent.parent != nil {
					parent.tokenType = varLabel
					parent.parent.children = append(parent.parent.children, a.Token.children[len(a.Token.children)-1])
					//remove child from a
					a.Token.children = a.Token.children[:len(a.Token.children)-1]
				} else if parent.tokenType != constant {
					parent = a.Token.children[len(a.Token.children)-1]
					parent.tokenType = varLabel
				}
				//add var size and value
				parent.children = []*Token{NewToken(varSize, cleanupStr(lineParts[0]), parent), NewToken(varValue, finalLinePart2, parent)}
				if parent.parent != nil {
					return parent.parent, nil
				} else {
					return parent, nil
				}
			}
		} else {
			// variable / constant size
			parent.children = append(parent.children, NewToken(varValue, ln, parent))
			return parent, nil
		}
		return parent, errors.New(". found in line " + ln + " but not matching anything")
	}

	//either label or var
	if ln[len(ln)-1] == ':' {
		if parent.tokenType == global || parent.tokenType == constant || parent.tokenType == section || parent.tokenType == globalLabel || parent.tokenType == localLabel {
			if ln[0] <= '0' && ln[0] >= '9' {
				//local label
				tk := NewToken(localLabel, ln, parent)
				parent.children = append(parent.children, tk)
				return tk, nil
			} else if len(lineParts) == 1 {
				//globalLabel
				tk := NewToken(globalLabel, ln, parent)
				a.Token.children = append(a.Token.children, tk)
				return tk, nil
			} else {
				//vars
				//increase instruction count to keep track of variable Size
				var finalLinePart2 string = strings.TrimSpace(strings.Join(lineParts[2:], " "))
				varS, isN := getVarSize(lineParts[1])
				if isN {
					variableCount += varS
				} else {
					finalLinePart2 = finalLinePart2[1 : len(finalLinePart2)-1]
					stringCount += len(finalLinePart2)*8 - 8
				}
				tk := NewToken(varLabel, ln, parent)
				//add var size and value
				tk.children = []*Token{NewToken(varSize, cleanupStr(lineParts[1]), tk), NewToken(varValue, finalLinePart2, tk)}
				parent.children = append(parent.children, tk)
				return parent, nil
			}
		}
	}
	// either variable value or code line
	instructionType, ok := InstructionToOpType[ln]
	if !ok {
		return parent, errors.New("Unknown instruction type: '" + ln + "'")
	}

	ptk := NewToken(instruction, ln, parent, &instructionType)
	parent.children = append(parent.children, ptk)
	lineParts = lineParts[1:]
	var err error
	//increase instruction count to keep track of machineCode Size
	instructionCount += 4

	if ln == "ebreak" || ln == "ecall" {
		return parent, nil
	}
	var newArr []string
	for _, li := range lineParts {
		splitArr := strings.Split(li, ",")
		if len(splitArr) != 1 {
			newArr = append(newArr, splitArr...)
		} else {
			newArr = append(newArr, li)
		}
	}
	lineParts = newArr
	switch instructionType.opType {
	case R:
		err = ParseRegisters(lineParts, ptk)
	case I:
		err = LexIType(lineParts, ptk)
	case S:
		err = LexSType(lineParts, ptk)
	case B:
		err = LexBType(lineParts, ptk)
	case U:
		err = LexUType(lineParts, ptk)
	case J:
		err = LexJType(lineParts, ptk)
	case CI:
		err = LexJType(lineParts, ptk)
	case CSS:
		err = LexJType(lineParts, ptk)
	case CL:
		err = LexJType(lineParts, ptk)
	case CJ:
		err = LexJType(lineParts, ptk)
	case CR:
		err = ParseRegisters(lineParts, ptk)
	case CB:
		err = LexJType(lineParts, ptk)
	case CIW:
		err = LexJType(lineParts, ptk)
	case CS:
		err = LexJType(lineParts, ptk)
	default:
		panic("unhandled default case")
	}

	if err != nil {
		return parent, err
	}
	return parent, nil
}

func ParseRegisters(strArr []string, parent *Token) error { //todo check for errors
	var numInst = 0
	for _, str := range strArr {
		if len(str) == 0 {
			continue
		}
		parent.children = append(parent.children, NewToken(register, strings.TrimSpace(str), parent))
		numInst++
	}
	return nil
}

func LexIType(strArr []string, parent *Token) error {
	err := ParseRegisters(strArr[:len(strArr)-1], parent)
	if err != nil {
		return err
	}

	adjustedVal, err := parseIntValue(strArr[len(strArr)-1])
	if err != nil {
		var vals = strings.Split(strArr[len(strArr)-1], "(")

		if len(vals) != 2 {
			cleanedStr := cleanupStr(vals[0])
			if _, ok := matchTokenValid(cleanedStr); ok == nil {
				parent.children = append(parent.children, NewToken(register, cleanedStr, parent))
			} else {
				parent.children = append(parent.children, NewToken(varValue, vals[0], parent))
			}
		} else {
			child := NewToken(complexValue, strArr[len(strArr)-1], parent)
			adjustedVal, err = parseIntValue(cleanupStr(vals[0]))
			if err != nil {

				if strings.Contains(vals[0], "%") {

					child.children = append(child.children, NewToken(modifier, cleanupStr(vals[0]), child))
				} else {
					return errors.New("I TYPE: Offset is not a number " + vals[0])
				}
			} else {
				child.children = append(child.children, NewToken(literal, strconv.Itoa(adjustedVal), child))
			}

			cleanedStr := cleanupStr(vals[1][:len(vals[1])-1])
			adjustedVal, err = parseIntValue(cleanedStr)
			if err == nil {
				child.value = strconv.Itoa(adjustedVal)
				child.children = append(child.children, NewToken(literal, cleanedStr, child))
			} else if strings.Contains(vals[1], ".") {
				child.children = append(child.children, NewToken(constantValue, cleanedStr, child))
			} else if _, ok := matchTokenValid(cleanedStr); ok == nil {
				child.children = append(child.children, NewToken(register, cleanedStr, child))
			} else {
				child.children = append(child.children, NewToken(varValue, cleanedStr, child))
			}
			parent.children = append(parent.children, child)
		}
	} else {
		parent.children = append(parent.children, NewToken(literal, strconv.Itoa(adjustedVal), parent))
	}

	return nil
}

func LexSType(strArr []string, parent *Token) error {
	err := ParseRegisters(strArr[:len(strArr)-1], parent)
	if err != nil {
		return err
	}

	// check if format offset(register)
	var vals = strings.Split(strArr[len(strArr)-1], "(")

	if len(vals) != 2 {
		cleanedStr := cleanupStr(vals[0])
		if _, ok := matchTokenValid(cleanedStr); ok == nil {
			parent.children = append(parent.children, NewToken(register, cleanedStr, parent))
		} else {
			parent.children = append(parent.children, NewToken(varValue, vals[0], parent))
		}
	} else {
		child := NewToken(complexValue, strArr[len(strArr)-1], parent)
		adjustedVal, err := parseIntValue(cleanupStr(vals[0]))
		if err != nil {

			if strings.Contains(vals[0], "%") {

				child.children = append(child.children, NewToken(modifier, cleanupStr(vals[0]), child))
			} else {
				return errors.New("S TYPE: Offset is not a number " + vals[0])
			}
		} else {
			child.children = append(child.children, NewToken(literal, strconv.Itoa(adjustedVal), child))
		}

		cleanupStr := cleanupStr(vals[1][:len(vals[1])-1])
		adjustedVal, err = parseIntValue(cleanupStr)
		if err == nil {
			child.value = strconv.Itoa(adjustedVal)
			child.children = append(child.children, NewToken(literal, cleanupStr, child))
		} else if strings.Contains(vals[1], ".") {
			child.children = append(child.children, NewToken(constantValue, cleanupStr, child))
		} else if _, ok := matchTokenValid(cleanupStr); ok == nil {
			child.children = append(child.children, NewToken(register, cleanupStr, child))
		} else {
			child.children = append(child.children, NewToken(varValue, cleanupStr, child))
		}

		parent.children = append(
			parent.children,
			child)
	}
	return nil
}

func LexBType(strArr []string, parent *Token) error {
	return LexIType(strArr, parent)
}

func LexUType(strArr []string, parent *Token) error {
	err := ParseRegisters(strArr[:len(strArr)-1], parent)
	if err != nil {
		return err
	}

	adjustedVal, err := parseIntValue(strArr[len(strArr)-1])
	if err != nil {
		var vals = strings.Split(strArr[len(strArr)-1], "(")

		if len(vals) != 2 {
			cleanedStr := cleanupStr(vals[0])
			if _, ok := matchTokenValid(cleanedStr); ok == nil {
				parent.children = append(parent.children, NewToken(register, cleanedStr, parent))
			} else {
				parent.children = append(parent.children, NewToken(varValue, vals[0], parent))
			}
		} else {
			child := NewToken(complexValue, strArr[len(strArr)-1], parent)
			adjustedVal, err = parseIntValue(cleanupStr(vals[0]))
			if err != nil {
				if strings.Contains(vals[0], "%") {
					child.children = append(child.children, NewToken(modifier, cleanupStr(vals[0]), child))
				} else {
					return errors.New("I TYPE: Offset is not a number " + vals[0])
				}
			} else {
				child.children = append(child.children, NewToken(literal, strconv.Itoa(adjustedVal), child))
			}

			cleanupStr := cleanupStr(vals[1][:len(vals[1])-1])
			adjustedVal, err = parseIntValue(cleanupStr)
			if err == nil {
				child.value = strconv.Itoa(adjustedVal)
				child.children = append(child.children, NewToken(literal, cleanupStr, child))
			} else if strings.Contains(vals[1], ".") {
				child.children = append(child.children, NewToken(constantValue, cleanupStr, child))
			} else if _, ok := matchTokenValid(cleanupStr); ok == nil {
				child.children = append(child.children, NewToken(register, cleanupStr, child))
			} else {
				child.children = append(child.children, NewToken(varValue, cleanupStr, child))
			}
			parent.children = append(parent.children, child)
		}
	} else {
		parent.children = append(parent.children, NewToken(literal, strconv.Itoa(adjustedVal), parent))
	}

	return nil
}

func LexJType(strArr []string, parent *Token) error {
	return LexUType(strArr, parent)
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (a *Assembler) Assemble(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	actualParent := a.Token
	for scanner.Scan() {
		line := scanner.Text()

		lineIndx := strings.Index(line, "#")
		if lineIndx != -1 {
			line = line[:lineIndx]
		}

		lineParts := strings.Split(line, " ") //     num1: .word 10        # First number
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
	println("parsing saved to ./testfile/output.parser")
	fmt.Println("blocks sizes")
	fmt.Print("instructions: ")
	fmt.Println(instructionCount)
	fmt.Print("variables: ")
	fmt.Println(variableCount)
	fmt.Print("constants: ")
	fmt.Println(constantCount)
	fmt.Print("strings: ")
	fmt.Println(stringCount)
	os.WriteFile("./testfile/output.parser", []byte(str), 0644)

	prog := compile(a.Token)
	bytes := BuildELFFile(prog)
	fmt.Printf("% x", *bytes)
	os.WriteFile("./testfile/output.go", *bytes, 0644)
	return scanner.Err()
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
		// todo handle either size of var (.word) or section (.text) or constant (.LC0:)
		if parent.tokenType == global || parent.tokenType == section || parent.tokenType == globalLabel || parent.tokenType == constant {
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
				return tk, nil
			} else if len(lineParts) == 1 {
				// section
				tk := NewToken(section, ln, a.Token)
				a.Token.children = append(a.Token.children, tk)
				return tk, nil
			} else {
				// constant variable defined without label

				varS, isN := getVarSize(lineParts[0])
				if isN {
					constantCount += varS
				} else {
					stringCount += len(strings.TrimSpace(strings.Join(lineParts[1:], " ")))*8 - 8
				}
				//add var size and value
				parent.children = []*Token{NewToken(varSize, cleanupStr(lineParts[0]), parent), NewToken(varValue, cleanupStr(strings.Join(lineParts[1:], " ")), parent)}
				return parent, nil
			}
		} else {
			// variable / constant size
			parent.children = append(parent.children, NewToken(varValue, ln, parent))
			return parent, nil
		}
		return parent, errors.New(". found but not matching")
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
				varS, isN := getVarSize(lineParts[1])
				if isN {
					variableCount += varS
				} else {
					stringCount += len(strings.TrimSpace(strings.Join(lineParts[2:], " ")))*8 - 8
				}
				tk := NewToken(varLabel, ln, parent)
				//add var size and value
				tk.children = []*Token{NewToken(varSize, cleanupStr(lineParts[1]), tk), NewToken(varValue, cleanupStr(strings.Join(lineParts[2:], " ")), tk)}
				parent.children = append(parent.children, tk)
				return parent, nil
			}
		} else if parent.tokenType == section || parent.tokenType == constant {
			//vars
			//increase instruction count to keep track of variable Size
			varS, isN := getVarSize(lineParts[1])
			if isN {
				variableCount += varS
			} else {
				stringCount += len(strings.TrimSpace(strings.Join(lineParts[2:], " ")))*8 - 8
			}

			tk := NewToken(varLabel, ln, parent)
			//add var size and value
			tk.children = []*Token{NewToken(varSize, cleanupStr(lineParts[1]), tk), NewToken(varValue, cleanupStr(strings.Join(lineParts[2:], " ")), tk)}
			parent.children = append(parent.children, tk)
			return parent, nil
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
	instructionCount += 32
	if ln == "ebreak" || ln == "ecall" || ln == "call" {
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

	_, err = strconv.Atoi(strArr[len(strArr)-1])
	if err != nil {
		var vals = strings.Split(strArr[len(strArr)-1], "(")

		if len(vals) != 2 {
			parent.children = append(parent.children, NewToken(varLabel, vals[0], parent))
		} else {
			child := NewToken(complexValue, strArr[len(strArr)-1], parent)
			_, err = strconv.Atoi(cleanupStr(vals[0]))
			if err != nil {

				if strings.Contains(vals[0], "%") {

					child.children = append(child.children, NewToken(modifier, cleanupStr(vals[0]), child))
				} else {
					return errors.New("I TYPE: Offset is not a number " + vals[0])
				}
			} else {
				child.children = append(child.children, NewToken(literal, cleanupStr(vals[0]), child))
			}

			if strings.Contains(vals[1], ".") {
				child.children = append(child.children, NewToken(constantValue, cleanupStr(vals[1][:len(vals[1])-1]), child))
			} else {
				child.children = append(child.children, NewToken(register, cleanupStr(vals[1][:len(vals[1])-1]), child))
			}
			parent.children = append(parent.children, child)
		}
	} else {
		parent.children = append(parent.children, NewToken(literal, strArr[len(strArr)-1], parent))
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
		parent.children = append(parent.children, NewToken(varLabel, vals[0], parent))
	} else {
		child := NewToken(complexValue, strArr[len(strArr)-1], parent)
		_, err = strconv.Atoi(cleanupStr(vals[0]))
		if err != nil {

			if strings.Contains(vals[0], "%") {

				child.children = append(child.children, NewToken(modifier, cleanupStr(vals[0]), child))
			} else {
				return errors.New("S TYPE: Offset is not a number " + vals[0])
			}
		} else {
			child.children = append(child.children, NewToken(literal, cleanupStr(vals[0]), child))
		}

		if strings.Contains(vals[1], ".") {
			child.children = append(child.children, NewToken(constantValue, cleanupStr(vals[1][:len(vals[1])-1]), child))
		} else {
			child.children = append(child.children, NewToken(register, cleanupStr(vals[1][:len(vals[1])-1]), child))
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

	_, err = strconv.Atoi(strArr[len(strArr)-1])
	if err != nil {
		var vals = strings.Split(strArr[len(strArr)-1], "(")

		if len(vals) != 2 {
			parent.children = append(parent.children, NewToken(varLabel, vals[0], parent))
		} else {
			child := NewToken(complexValue, strArr[len(strArr)-1], parent)
			_, err = strconv.Atoi(cleanupStr(vals[0]))
			if err != nil {

				if strings.Contains(vals[0], "%") {

					child.children = append(child.children, NewToken(modifier, cleanupStr(vals[0]), child))
				} else {
					return errors.New("I TYPE: Offset is not a number " + vals[0])
				}
			} else {
				child.children = append(child.children, NewToken(literal, cleanupStr(vals[0]), child))
			}

			if strings.Contains(vals[1], ".") {
				child.children = append(child.children, NewToken(constantValue, cleanupStr(vals[1][:len(vals[1])-1]), child))
			} else {
				child.children = append(child.children, NewToken(register, cleanupStr(vals[1][:len(vals[1])-1]), child))
			}
			parent.children = append(parent.children, child)
		}
	} else {
		parent.children = append(parent.children, NewToken(literal, strArr[len(strArr)-1], parent))
	}

	return nil
}

func LexJType(strArr []string, parent *Token) error {
	return LexUType(strArr, parent)
}

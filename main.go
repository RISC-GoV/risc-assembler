package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: assembler <input.s>")
		os.Exit(1)
	}
	asm := Assembler{}
	asm.Token = NewToken(global, "", nil)
	if err := asm.Assemble(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

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

	return scanner.Err()
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
			if ln[len(ln)-1] == ':' {
				println("constant " + ln)
				// constant
				tk := NewToken(constant, ln, parent)
				parent.children = append(parent.children, tk)
				return tk, nil
			} else {
				println("section " + ln)
				// section
				tk := NewToken(section, ln, parent)
				parent.children = append(parent.children, tk)
				return tk, nil
			}
		} else if parent.tokenType == section || parent.tokenType == globalLabel || parent.tokenType == constant {
			if ln == ".globl" {
				println("global " + ln)
				tk := NewToken(entrypoint, ln, parent)
				tk.children = []*Token{NewToken(globalLabel, cleanupStr(lineParts[1]), parent)}
				parent.children = append(parent.children, tk)
				return parent, nil
			}
		} else {
			// variable / constant size
			println("constant " + ln)
			parent.children = append(parent.children, NewToken(varValue, ln, parent))
			return parent, nil
		}
		return parent, errors.New(". found but not matching")
	}

	// either label or var
	if ln[len(ln)-1] == ':' {
		if parent.tokenType == global || parent.tokenType == section && parent.value == ".text" {
			// code
			println("globalLabel " + ln)
			tk := NewToken(globalLabel, ln, parent)
			parent.children = append(parent.children, tk)
			return tk, nil
		} else if parent.tokenType == section {
			// vars
			println("varLabel " + ln)
			tk := NewToken(varLabel, ln, parent)
			// add var size and value
			tk.children = []*Token{NewToken(varSize, cleanupStr(lineParts[1]), parent), NewToken(varValue, cleanupStr(lineParts[2]), parent)}
			parent.children = append(parent.children, tk)
			return parent, nil
		} else {
			// local label
			println("localLabel " + ln)
			tk := NewToken(localLabel, ln, parent)
			parent.children = append(parent.children, tk)
			return tk, nil
		}
	}
	// todo from here, either variable value or code line
	instructionType, ok := InstructionToOpType[ln]
	if !ok {
		return parent, errors.New("Unknown instruction type: '" + ln + "'")
	}

	parent.children = append(parent.children, NewToken(instruction, ln, parent))
	var err error
	switch instructionType.opType {
	case R:
		err = ParseRegisters(lineParts, parent)
		break
	case I:
		err = LexIType(lineParts, parent)
		break
	case S:
		err = LexSType(lineParts, parent)
		break
	case B:
		err = LexBType(lineParts, parent)
		break
	case U:
		err = LexUType(lineParts, parent)
		break
	case J:
		err = LexJType(lineParts, parent)
		break
	case CI:
		err = LexJType(lineParts, parent)
		break
	case CSS:
		err = LexJType(lineParts, parent)
		break
	case CL:
		err = LexJType(lineParts, parent)
		break
	case CJ:
		err = LexJType(lineParts, parent)
		break
	case CR:
		err = ParseRegisters(lineParts, parent)
		break
	case CB:
		err = LexJType(lineParts, parent)
		break
	case CIW:
		err = LexJType(lineParts, parent)
		break
	case CS:
		err = LexJType(lineParts, parent)
		break
	default:
		panic("unhandled default case")
	}

	if err != nil {
		return parent, err
	}
	return parent, nil
}

func ParseRegisters(strArr []string, parent *Token) error { // todo check for errors
	numInst := 0
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
		parent.children = append(parent.children, NewToken(varLabel, strArr[len(strArr)-1], parent))
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
	vals := strings.Split(strArr[len(strArr)-1], "(")
	if len(vals) != 2 { // todo
		return errors.New("S TYPE: Last arg is not of format 'offset(register)' " + strArr[len(strArr)-1])
	}

	_, err = strconv.Atoi(vals[0])
	if err != nil {
		return errors.New("S TYPE: Offset is not a number " + vals[0])
	}

	child := NewToken(complexValue, strArr[len(strArr)-1], parent)
	child.children = append(child.children, NewToken(literal, vals[0], child))
	child.children = append(child.children, NewToken(literal, strings.TrimSpace(vals[1][:len(vals[1])-1]), child))

	parent.children = append(
		parent.children,
		child)
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
		return errors.New("U TYPE: Last arg is not of type int " + strArr[len(strArr)-1])
	}

	parent.children = append(parent.children, NewToken(literal, strArr[len(strArr)-1], parent))

	return nil
}

func LexJType(strArr []string, parent *Token) error {
	return LexUType(strArr, parent)
}

func extractHex(immVal string) []byte {
	// Decode the hexadecimal string to bytes
	hexBytes, err := hex.DecodeString(immVal)
	if err != nil {
		return nil
	}
	return hexBytes
}

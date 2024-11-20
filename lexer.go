package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func LexFile(inputFile string, outputFile string) { //todo
	data, err := os.ReadFile(inputFile)
	check(err)
	_ = LexString(string(data))
	//err = os.WriteFile(outputFile, output, 0111)
	check(err)
}

func LexString(str string) []*Token {
	var tokenArr []*Token
	var buffer bytes.Buffer
	var parentToken *Token
	var functionToken *Token
	for _, char := range str {
		if char == '\n' {
			var line = buffer.String()
			if len(line) == 0 {
				continue
			}
			if line[0] == '.' {
				parentToken = NewToken(section, line[1:])
				if parentToken != nil {
					tokenArr = append(tokenArr, parentToken)
				}
			} else if parentToken != nil {
				var outToken *Token
				var err error
				if functionToken != nil {
					outToken, err = LexLine(line, functionToken)
				} else {
					outToken, err = LexLine(line, parentToken)
				}
				check(err)
				if outToken != nil {
					if outToken.tokenType == globalLabel {
						functionToken = outToken
					} else if outToken.tokenType == ret {
						functionToken = nil
					}
				}
			}
			buffer.Reset()
		}
		buffer.Grow(1)
		buffer.WriteRune(char)
	}
	return tokenArr
}

func LexLine(linePart string, parent *Token) (*Token, error) { //todo
	var token *Token
	if len(linePart) == 0 {
		return token, nil
	}
	if strings.TrimSpace(linePart) == "ret" {
		token = NewToken(ret, linePart)
		parent.children = append(parent.children, token)
		return token, nil
	}
	if linePart[0] == '.' {
		token = NewToken(labelSection, linePart[1:])
		parent.children = append(parent.children, token)
		return token, nil
	}
	if linePart[len(linePart)-1] == ':' {
		_, err := strconv.Atoi(linePart[1 : len(linePart)-1])
		if err != nil {
			token = NewToken(globalLabel, linePart[:len(linePart)-1])
			parent.children = append(parent.children, token)
			return token, nil
		}
		token = NewToken(localLabel, linePart[1:len(linePart)-1])
		parent.children = append(parent.children, token)
		return token, nil
	}

	var strArr []string = strings.Split(linePart, ",")
	if len(strings.Split(strArr[0], " ")) > 1 {
		var splitStr = strings.Split(strArr[0], " ")
		strArr = append([]string{splitStr[0], splitStr[1]}, strArr[1:]...)
	}

	instructionSTR := strings.ToUpper(strings.ToLower(strArr[0]))
	instructionType, ok := InstructionToOpType[instructionSTR]
	if !ok {
		return token, errors.New("Unknown instruction type: " + strArr[0])
	}

	parent.children = append(parent.children, NewToken(instruction, instructionSTR))

	var err error

	switch instructionType {
	case R:
		err = LexRType(strArr, parent)
		break
	case I:
		err = LexIType(strArr, parent)
		break
	case S:
		err = LexSType(strArr, parent)
		break
	case B:
		err = LexBType(strArr, parent)
		break
	case U:
		err = LexUType(strArr, parent)
		break
	case J:
		err = LexJType(strArr, parent)
		break
	}

	if err != nil {
		return nil, err
	}
	return parent.children[len(parent.children)-1], nil
}

func ParseRegisters(strArr []string, parent *Token) (int, error) {
	var numInst = 0
	for _, str := range strArr {
		if len(str) == 0 {
			continue
		}
		_, ok := Register[strings.TrimSpace(str)]
		if !ok {
			return 0, errors.New("Unknown register type: " + strArr[0])
		}
		parent.children = append(parent.children, NewToken(register, strings.TrimSpace(str)))
		numInst++
	}
	return numInst, nil
}

func LexRType(strArr []string, parent *Token) error {
	numInst, err := ParseRegisters(strArr, parent)
	if err != nil {
		return err
	}
	if numInst != 3 {
		return errors.New("Wrong number of arguments. Expected 3, got " + strconv.Itoa(numInst))
	}
	return nil
}

func LexIType(strArr []string, parent *Token) error {
	numInst, err := ParseRegisters(strArr[:len(strArr)-1], parent)
	if err != nil {
		return err
	}

	_, err = strconv.Atoi(strArr[len(strArr)-1])
	if err != nil {
		return errors.New("Last arg is not of type int " + strArr[len(strArr)-1])
	}
	numInst++

	if numInst != 3 {
		return errors.New("Wrong number of arguments. Expected 3, got " + strconv.Itoa(numInst))
	}

	parent.children = append(parent.children, NewToken(literal, strArr[len(strArr)-1]))

	return nil
}

func LexSType(strArr []string, parent *Token) error {
	numInst, err := ParseRegisters(strArr[:len(strArr)-1], parent)
	if err != nil {
		return err
	}
	var vals = strings.Split(strings.TrimSpace(strArr[len(strArr)-1]), "(")
	if len(vals) != 2 {
		return errors.New("Last arg is not of format 'offset(register)' " + strArr[len(strArr)-1])
	}

	_, err = strconv.Atoi(vals[0])
	if err != nil {
		return errors.New("Offset is not a number " + vals[0])
	}

	var _, ok = Register[strings.TrimSpace(vals[1][:len(vals[1])-1])]
	if !ok {
		return errors.New("Unknown register type: " + vals[1])
	}

	parent.children = append(
		parent.children,
		&Token{complexValue, strings.TrimSpace(strArr[len(strArr)-1]), []*Token{
			NewToken(literal, vals[0]),
			NewToken(register, strings.TrimSpace(vals[1][:len(vals[1])-1])),
		}})
	numInst++

	if numInst != 3 {
		return errors.New("Wrong number of arguments. Expected 2, got " + strconv.Itoa(numInst))
	}
	return nil
}

func LexBType(strArr []string, parent *Token) error {
	return LexIType(strArr, parent)
}

func LexUType(strArr []string, parent *Token) error {
	numInst, err := ParseRegisters(strArr[:len(strArr)-1], parent)
	if err != nil {
		return err
	}

	_, err = strconv.Atoi(strArr[len(strArr)-1])
	if err != nil {
		return errors.New("Last arg is not of type int " + strArr[len(strArr)-1])
	}
	numInst++

	if numInst != 2 {
		return errors.New("Wrong number of arguments. Expected 2, got " + strconv.Itoa(numInst))
	}

	parent.children = append(parent.children, NewToken(literal, strArr[len(strArr)-1]))

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

func copybytes(source []byte, dest []byte, amt int, sourceoffset int, destoffset int) {
	for i := 0; i < amt; i++ {
		dest[destoffset+i] = source[sourceoffset+i]
	}
}

func GenerateELFHeaders(e_entry []byte, e_phof []byte, e_shoff []byte, e_phnum []byte, e_shnum []byte, e_shstrndx []byte) []byte {
	var elfheader []byte
	//Magic Number
	elfheader[0x0] = 0x7F
	elfheader[0x1] = 0x45  //E
	elfheader[0x2] = 0x4c  //L
	elfheader[0x3] = 0x46  //F
	elfheader[0x4] = 0x01  //32bit
	elfheader[0x5] = 0x01  //LE
	elfheader[0x6] = 0x01  //ELF Version
	elfheader[0x7] = 0x03  //Linux ABI
	elfheader[0x10] = 0x02 //Executable
	elfheader[0x12] = 0xF3 //RISC-V
	elfheader[0x14] = 0x01
	//EntryPoint Address
	copybytes(e_entry, elfheader, 2, 0, 0x18)
	//Program Header Address
	copybytes(e_phof, elfheader, 2, 0, 0x1C)
	//Section Header Address
	copybytes(e_shoff, elfheader, 2, 0, 0x20)
	elfheader[0x24] = 0x00 //???
	elfheader[0x28] = 0x34 //Len = 52 Bytes
	elfheader[0x2A] = 0x20 //32bits
	//Amount of Entries in Program Header
	copybytes(e_phnum, elfheader, 2, 0, 0x2C)
	elfheader[0x2E] = 0x28 //Header entry size
	//Amount of Entries in Section Header
	copybytes(e_shnum, elfheader, 2, 0, 0x30)
	//Index of SHT entry containing names
	copybytes(e_shstrndx, elfheader, 2, 0, 0x32)
	return elfheader
}

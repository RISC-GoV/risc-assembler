package main

import (
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
	_ = TokenizeString(string(data))
	err = os.WriteFile(outputFile, []byte(""), 011)
	check(err)
}

func TokenizeString(str string) *Token {
	var token *Token = NewToken(none, str)
	for cursor := 0; cursor < len(str); {
		//Start extracting section for tokenizing
		for sectionCursor := cursor; sectionCursor < len(str); sectionCursor++ {
			if str[sectionCursor] == '\n' && sectionCursor != len(str)-1 {
				if str[sectionCursor] == '.' || sectionCursor == len(str)-1 { //End of previous section was reached

					currentSection := str[cursor:sectionCursor]
					sectionToken, err := TokenizeSection(currentSection)
					check(err)
					cursor = sectionCursor //Move cursor to start of next section
					token.children = append(token.children, sectionToken)
					break
				}
			}
		}
	}

	return token
}

func TokenizeSection(section string) (*Token, error) {
	return NewToken(none, ""), nil
}

func TokenizeLine(line string, parent *Token) (*Token, error) { //todo
	var token *Token
	if len(line) == 0 {
		return token, nil
	}
	if strings.TrimSpace(line) == "ret" {
		token = NewToken(ret, line)
		parent.children = append(parent.children, token)
		return token, nil
	}
	if line[0] == '.' {
		token = NewToken(labelSection, line[1:])
		parent.children = append(parent.children, token)
		return token, nil
	}
	if line[len(line)-1] == ':' {
		_, err := strconv.Atoi(line[1 : len(line)-1])
		if err != nil {
			token = NewToken(globalLabel, line[:len(line)-1])
			parent.children = append(parent.children, token)
			return token, nil
		}
		token = NewToken(localLabel, line[1:len(line)-1])
		parent.children = append(parent.children, token)
		return token, nil
	}

	var strArr []string = strings.Split(line, ",")
	var splitStr = strings.Split(strArr[0], " ")
	if len(splitStr) > 1 {
		strArr = append([]string{splitStr[1]}, strArr[1:]...)
	}
	for index, str := range strArr {
		strArr[index] = strings.TrimSpace(str)
	}
	instructionSTR := strings.ToLower(strings.ToLower(splitStr[0]))
	instructionType, ok := InstructionToOpType[instructionSTR]
	if !ok {
		return token, errors.New("Unknown instruction type: '" + splitStr[0] + "'")
	}

	parent.children = append(parent.children, NewToken(instruction, instructionSTR))

	var err error

	switch instructionType.opType {
	case R:
		err = ParseRegisters(strArr, parent)
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
	case CI:
		err = LexJType(strArr, parent)
		break
	case CSS:
		err = LexJType(strArr, parent)
		break
	case CL:
		err = LexJType(strArr, parent)
		break
	case CJ:
		err = LexJType(strArr, parent)
		break
	case CR:
		err = ParseRegisters(strArr, parent)
		break
	case CB:
		err = LexJType(strArr, parent)
		break
	case CIW:
		err = LexJType(strArr, parent)
		break
	case CS:
		err = LexJType(strArr, parent)
		break
	default:
		panic("unhandled default case")
	}

	if err != nil {
		return nil, err
	}
	return parent.children[len(parent.children)-1], nil
}

func ParseRegisters(strArr []string, parent *Token) error { //todo check for errors
	var numInst = 0
	for _, str := range strArr {
		if len(str) == 0 {
			continue
		}
		parent.children = append(parent.children, NewToken(register, strings.TrimSpace(str)))
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
		return errors.New("I TYPE: Last arg is not of type int: '" + strArr[len(strArr)-1] + "'")
	}
	parent.children = append(parent.children, NewToken(literal, strArr[len(strArr)-1]))

	return nil
}

func LexSType(strArr []string, parent *Token) error {
	err := ParseRegisters(strArr[:len(strArr)-1], parent)
	if err != nil {
		return err
	}
	var vals = strings.Split(strArr[len(strArr)-1], "(")
	if len(vals) != 2 {
		return errors.New("S TYPE: Last arg is not of format 'offset(register)' " + strArr[len(strArr)-1])
	}

	_, err = strconv.Atoi(vals[0])
	if err != nil {
		return errors.New("S TYPE: Offset is not a number " + vals[0])
	}

	parent.children = append(
		parent.children,
		&Token{complexValue, strArr[len(strArr)-1], []*Token{
			NewToken(literal, vals[0]),
			NewToken(register, strings.TrimSpace(vals[1][:len(vals[1])-1])),
		}})
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

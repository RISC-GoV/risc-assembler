package main

import (
	"errors"
	"strconv"
	"unsafe"
)

func invertBits(integer uint32) uint32 {
	var result uint32 = 0

	for range 32 {
		result |= integer & 1
		integer >>= 1
		result <<= 1
	}

	return result
}

func TranslateRType(opcode int, rd int, func3 int, rs1 int, rs2 int, func7 int) uint32 {
	//masks
	opcode &= 0b1111111
	rd &= 0b11111
	func3 &= 0b111
	rs1 &= 0b11111
	rs2 &= 0b11111
	func7 &= 0b1111111

	//shifts
	res := uint32(func7)
	res <<= 5
	res |= uint32(rs2)
	res <<= 5
	res |= uint32(rs1)
	res <<= 3
	res |= uint32(func3)
	res <<= 5
	res |= uint32(rd)
	res <<= 7
	res |= uint32(opcode)

	return invertBits(res)
}

func TranslateIType(opcode int, rd int, func3 int, rs1 int, imm int) uint32 {
	//masks
	opcode &= 0b1111111
	rd &= 0b11111
	func3 &= 0b111
	rs1 &= 0b11111
	imm &= 0b111111111111

	//shifts
	res := uint32(imm)
	res <<= 5
	res |= uint32(rs1)
	res <<= 3
	res |= uint32(func3)
	res <<= 5
	res |= uint32(rd)
	res <<= 7
	res |= uint32(opcode)

	return invertBits(res)
}

func TranslateSType(opcode int, func3 int, rs1 int, rs2 int, imm int) uint32 {
	//masks
	opcode &= 0b1111111
	imm4_0 := imm & 0b11111
	func3 &= 0b111
	rs1 &= 0b11111
	rs2 &= 0b11111
	imm11_5 := imm
	imm11_5 >>= 5
	imm11_5 &= 0b1111111

	//shifts
	res := uint32(imm11_5)
	res <<= 5
	res |= uint32(rs2)
	res <<= 5
	res |= uint32(rs1)
	res <<= 3
	res |= uint32(func3)
	res <<= 5
	res |= uint32(imm4_0)
	res <<= 7
	res |= uint32(opcode)

	return invertBits(res)
}

func TranslateBType(opcode int, func3 int, rs1 int, rs2 int, imm int) uint32 {
	//masks
	opcode &= 0b1111111
	imm11 := imm >> 11 & 0b1
	imm4_1 := imm >> 1 & 0b1111
	func3 &= 0b111
	rs1 &= 0b11111
	rs2 &= 0b11111
	imm10_5 := imm >> 5 & 0b111111
	imm12 := imm >> 12 & 0b1

	//shifts
	res := uint32(imm12)
	res <<= 6
	res |= uint32(imm10_5)
	res <<= 5
	res |= uint32(rs2)
	res <<= 5
	res |= uint32(rs1)
	res <<= 3
	res |= uint32(func3)
	res <<= 4
	res |= uint32(imm4_1)
	res <<= 1
	res |= uint32(imm11)
	res <<= 7
	res |= uint32(opcode)

	return invertBits(res)
}

func TranslateUType(opcode int, rd int, imm int) uint32 {
	//masks
	opcode &= 0b1111111
	rd &= 0b11111
	imm &= 0b11111111111111111111

	//shifts
	res := uint32(imm)
	res <<= 5
	res |= uint32(rd)
	res <<= 7
	res |= uint32(opcode)

	return invertBits(res)
}
func TranslateJType(opcode int, rd int, imm int) uint32 {
	//masks
	opcode &= 0b1111111
	rd &= 0b11111
	imm19_12 := imm >> 12 & 0b11111111
	imm11 := imm >> 11 & 0b1
	imm10_1 := imm >> 1 & 0b1111111111
	imm20 := imm >> 20 & 0b1

	//shifts
	res := uint32(imm20)
	res <<= 10
	res |= uint32(imm10_1)
	res <<= 1
	res |= uint32(imm11)
	res <<= 8
	res |= uint32(imm19_12)
	res <<= 5
	res |= uint32(rd)
	res <<= 7
	res |= uint32(opcode)

	return invertBits(res)
}

func (p *Program) InstructionToBinary(t *Token) (uint32, error) {
	if t.tokenType != instruction {
		return 0, errors.New("expected instruction")
	}
	switch t.opPair.opType {
	case R: // add x0, x0, x0
		if len(t.children) != 3 {
			return 0, errors.New(t.value + " is not a valid instruction")
		}
		opcode := int(t.opPair.opByte[0])
		func3 := int(t.opPair.opByte[1])
		func7 := int(t.opPair.opByte[2])
		rd, err := t.children[0].getRegisterFromABI()
		if err != nil {
			return 0, err
		}
		rs1, err := t.children[1].getRegisterFromABI()
		if err != nil {
			return 0, err
		}
		rs2, err := t.children[2].getRegisterFromABI()
		if err != nil {
			return 0, err
		}
		return TranslateRType(opcode, rd, func3, rs1, rs2, func7), nil
	case I: // lw x0, 0(x0)
		opcode := int(t.opPair.opByte[0])
		func3 := int(t.opPair.opByte[1])
		rd, err := t.children[0].getRegisterFromABI()
		if err != nil {
			return 0, err
		}
		rs1, imm, err := p.parseComplexValue(t.children[1])
		if err != nil {
			return 0, err
		}
		return TranslateIType(opcode, rd, func3, rs1, imm), nil
	case S: //sw x0, 0(x0)
		opcode := int(t.opPair.opByte[0])
		func3 := int(t.opPair.opByte[1])
		rs1, err := t.children[0].getRegisterFromABI()
		if err != nil {
			return 0, err
		}
		rs2, imm, err := p.parseComplexValue(t.children[1])
		if err != nil {
			return 0, err
		}
		return TranslateSType(opcode, func3, rs1, rs2, imm), nil
	case B: // beq x0, x0, 0
		opcode := int(t.opPair.opByte[0])
		func3 := int(t.opPair.opByte[1])
		rs1, err := t.children[0].getRegisterFromABI()
		if err != nil {
			return 0, err
		}
		rs2, err := t.children[1].getRegisterFromABI()
		imm, err := parseLabelOrLiteral(t.children[2])
		if err != nil {
			return 0, err
		}
		return TranslateBType(opcode, func3, rs1, rs2, imm), nil
	case U: // lui x0, 0
		opcode := int(t.opPair.opByte[0])
		rd, err := t.children[0].getRegisterFromABI()
		if err != nil {
			return 0, err
		}
		imm, err := parseLabelOrLiteral(t.children[1])
		if err != nil {
			return 0, err
		}
		return TranslateUType(opcode, rd, imm), nil
	case J: // jal x0, 0
		opcode := int(t.opPair.opByte[0])
		rd, err := t.children[0].getRegisterFromABI()
		if err != nil {
			return 0, err
		}
		imm, err := parseLabelOrLiteral(t.children[1])
		if err != nil {
			return 0, err
		}
		return TranslateJType(opcode, rd, imm), nil
	case CI:
		return 0, errors.New("todo: compressed instructions") //todo
	case CSS:
		return 0, errors.New("todo: compressed instructions") //todo
	case CL:
		return 0, errors.New("todo: compressed instructions") //todo
	case CJ:
		return 0, errors.New("todo: compressed instructions") //todo
	case CR:
		return 0, errors.New("todo: compressed instructions") //todo
	case CB:
		return 0, errors.New("todo: compressed instructions") //todo
	case CIW:
		return 0, errors.New("todo: compressed instructions") //todo
	case CS:
		return 0, errors.New("todo: compressed instructions") //todo
	default:
		return 0, errors.New("unhandled default case")
	}
}

func (p *Program) parseComplexValue(tok *Token) (int, int, error) {
	switch tok.tokenType { //0(x0) OU 0(.LC1)
	case complexValue:
		if tok.children[0].tokenType == literal {
			if tok.children[1].tokenType == register {
				reg, err := matchTokenValid(tok.children[1].value)
				if err != nil {
					return 0, 0, err
				}
				imm, err := strconv.Atoi(tok.children[0].value)
				if err != nil {
					return 0, 0, err
				}
				return reg, imm, nil
			} else { //constant
				con, ok := labelPositions[tok.children[1].value]
				if !ok {
					return 0, 0, errors.New(tok.children[1].value + " not found")
				}
				imm, err := strconv.Atoi(tok.children[0].value)
				if err != nil {
					return 0, 0, err
				}
				return con, imm, nil
			}
		} else { //modifier
			if tok.children[1].tokenType == register {
				reg, err := matchTokenValid(tok.children[1].value)
				if err != nil {
					return 0, 0, err
				}
				reg, err = handleModifier(tok.children[0].value, reg)
				if err != nil {
					return 0, 0, err
				}
				return reg, 0, nil
			} else { //constant
				con, ok := labelPositions[tok.children[1].value]
				if !ok {
					return 0, 0, errors.New(tok.children[1].value + " not found")
				}
				con, err := handleModifier(tok.children[0].value, con)
				if err != nil {
					return 0, 0, err
				}
				return con, 0, nil
			}
		}
	case varLabel:
		fallthrough
	case constantValue:
		fallthrough
	case literal:
		val, err := parseLabelOrLiteral(tok)
		if err != nil {
			return 0, 0, err
		}
		return 0, val, nil

	}
	return 0, 0, nil
}

func handleModifier(mod string, val int) (int, error) {
	switch mod {
	case "%lo":
		return int(uint8(val)), nil
	case "%hi":
		for range unsafe.Sizeof(val) - 1 {
			val >>= 8
		}
		return int(uint8(val)), nil
	}
	return 0, errors.New("modifier not found")
}

func parseLabelOrLiteral(tok *Token) (int, error) {
	switch tok.tokenType {
	case varLabel:
		fallthrough
	case constantValue:
		imm, ok := labelPositions[tok.children[0].value]
		if !ok {
			return 0, errors.New(tok.children[0].value + "not found")
		}
		return imm, nil
	case literal:
		imm, err := strconv.Atoi(tok.children[0].value)
		if err != nil {
			return 0, err
		}
		return imm, nil
	}
	return 0, errors.New("wrong token type")
}

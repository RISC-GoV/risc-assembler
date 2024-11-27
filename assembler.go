package main

import "errors"

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

	return res
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

	return res
}

func TranslateSType(opcode int, imm4_0 int, func3 int, rs1 int, rs2 int, imm11_5 int) uint32 {
	//masks
	opcode &= 0b1111111
	imm4_0 &= 0b11111
	func3 &= 0b111
	rs1 &= 0b11111
	rs2 &= 0b11111
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

	return res
}

func TranslateBType(opcode int, imm11 int, imm4_1 int, func3 int, rs1 int, rs2 int, imm10_5 int, imm12 int) uint32 {
	//masks
	opcode &= 0b1111111
	imm11 &= 0b1
	imm4_1 &= 0b1111
	func3 &= 0b111
	rs1 &= 0b11111
	rs2 &= 0b11111
	imm10_5 &= 0b111111
	imm12 &= 0b1

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

	return res
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

	return res
}
func TranslateJType(opcode int, rd int, imm19_12 int, imm11 int, imm10_1 int, imm20 int) uint32 {
	//masks
	opcode &= 0b1111111
	rd &= 0b11111
	imm19_12 &= 0b11111111
	imm11 &= 0b1
	imm10_1 &= 0b1111111111
	imm20 &= 0b1

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

	return res
}

func InstructionToBinary(t Token) (uint32, error) {
	if t.tokenType != instruction {
		return 0, errors.New("expected instruction")
	}
	opPair, ok := InstructionToOpType[t.value]
	if !ok {
		return 0, errors.New(t.value + " is not a valid instruction")
	}
	switch opPair.opType {
	case R:
		opcode := int(opPair.opByte[0])
		func3 := int(opPair.opByte[1])
		func7 := int(opPair.opByte[2])
		rd := 0  //todo
		rs1 := 0 //todo
		rs2 := 0 //todo
		return TranslateRType(opcode, rd, func3, rs1, rs2, func7), nil
	case I:
		opcode := int(opPair.opByte[0])
		func3 := int(opPair.opByte[1])
		rd := 0  //todo
		rs1 := 0 //todo
		imm := 0 //todo
		return TranslateIType(opcode, rd, func3, rs1, imm), nil
	case S:
		opcode := int(opPair.opByte[0])
		imm4_0 := 0
		func3 := int(opPair.opByte[1])
		rs1 := 0
		rs2 := 0
		imm11_5 := 0
		return TranslateSType(opcode, imm4_0, func3, rs1, rs2, imm11_5), nil
	case B:
		opcode := int(opPair.opByte[0])
		imm11 := 0  //todo
		imm4_1 := 0 //todo
		func3 := int(opPair.opByte[1])
		rs1 := 0     //todo
		rs2 := 0     //todo
		imm10_5 := 0 //todo
		imm12 := 0   //todo
		return TranslateBType(opcode, imm11, imm4_1, func3, rs1, rs2, imm10_5, imm12), nil
	case U:
		opcode := int(opPair.opByte[0])
		rd := 0  //todo
		imm := 0 //todo
		return TranslateUType(opcode, rd, imm), nil
	case J:
		opcode := int(opPair.opByte[0])
		rd := 0       //todo
		imm19_12 := 0 //todo
		imm11 := 0    //todo
		imm10_1 := 0  //todo
		imm20 := 0    //todo
		return TranslateJType(opcode, rd, imm19_12, imm11, imm10_1, imm20), nil
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

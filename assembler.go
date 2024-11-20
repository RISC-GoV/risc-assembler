package main

func TranslateRType(opcode int, rd int, func3 int, rs1 int, rs2 int, func7 int) int {
	//masks
	opcode &= 0b1111111
	rd &= 0b11111
	func3 &= 0b111
	rs1 &= 0b11111
	rs2 &= 0b11111
	func7 &= 0b1111111

	//shifts
	res := func7
	res <<= 5
	res |= rs2
	res <<= 5
	res |= rs1
	res <<= 3
	res |= func3
	res <<= 5
	res |= rd
	res <<= 7
	res |= opcode

	return res
}

func TranslateIType(opcode int, rd int, func3 int, rs1 int, imm int) int {
	//masks
	opcode &= 0b1111111
	rd &= 0b11111
	func3 &= 0b111
	rs1 &= 0b11111
	imm &= 0b111111111111

	//shifts
	res := imm
	res <<= 5
	res |= rs1
	res <<= 3
	res |= func3
	res <<= 5
	res |= rd
	res <<= 7
	res |= opcode

	return res
}

func TranslateSType(opcode int, imm4_0 int, func3 int, rs1 int, rs2 int, imm11_5 int) int {
	//masks
	opcode &= 0b1111111
	imm4_0 &= 0b11111
	func3 &= 0b111
	rs1 &= 0b11111
	rs2 &= 0b11111
	imm11_5 &= 0b1111111

	//shifts
	res := imm11_5
	res <<= 5
	res |= rs2
	res <<= 5
	res |= rs1
	res <<= 3
	res |= func3
	res <<= 5
	res |= imm4_0
	res <<= 7
	res |= opcode

	return res
}

func TranslateBType(opcode int, imm11 int, imm4_1 int, func3 int, rs1 int, rs2 int, imm10_5 int, imm12 int) int {
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
	res := imm12
	res <<= 6
	res |= imm10_5
	res <<= 5
	res |= rs2
	res <<= 5
	res |= rs1
	res <<= 3
	res |= func3
	res <<= 4
	res |= imm4_1
	res <<= 1
	res |= imm11
	res <<= 7
	res |= opcode

	return res
}

func TranslateUType(opcode int, rd int, imm int) int {
	//masks
	opcode &= 0b1111111
	rd &= 0b11111
	imm &= 0b11111111111111111111

	//shifts
	res := imm
	res <<= 5
	res |= rd
	res <<= 7
	res |= opcode

	return res
}
func TranslateJType(opcode int, rd int, imm19_12 int, imm11 int, imm10_1 int, imm20 int) int {
	//masks
	opcode &= 0b1111111
	rd &= 0b11111
	imm19_12 &= 0b11111111
	imm11 &= 0b1
	imm10_1 &= 0b1111111111
	imm20 &= 0b1

	//shifts
	res := imm20
	res <<= 10
	res |= imm10_1
	res <<= 1
	res |= imm11
	res <<= 8
	res |= imm19_12
	res <<= 5
	res |= rd
	res <<= 7
	res |= opcode

	return res
}

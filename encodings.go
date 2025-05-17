package assembler

type Instruction struct {
	opcode    byte
	rd        int
	rs1       int
	rs2       int
	funct3    byte
	funct7    byte
	immediate int32
}

type Assembler struct {
	labels     map[string]int
	lineNumber int
	Token      *Token
	output     []uint32
	currentPC  int
}

func (a *Assembler) encodeRType(inst *Instruction) uint32 {
	result := uint32(inst.opcode)
	result |= uint32(inst.rd) << 7
	result |= uint32(inst.funct3) << 12
	result |= uint32(inst.rs1) << 15
	result |= uint32(inst.rs2) << 20
	result |= uint32(inst.funct7) << 25
	return result
}

func (a *Assembler) encodeIType(inst *Instruction) uint32 {
	result := uint32(inst.opcode)
	result |= uint32(inst.rd) << 7
	result |= uint32(inst.funct3) << 12
	result |= uint32(inst.rs1) << 15
	result |= (uint32(inst.immediate) & 0xFFF) << 20
	return result
}

func (a *Assembler) encodeSType(inst *Instruction) uint32 {
	imm := uint32(inst.immediate)
	result := uint32(inst.opcode)
	result |= (imm & 0x1F) << 7 // imm[4:0]
	result |= uint32(inst.funct3) << 12
	result |= uint32(inst.rs1) << 15
	result |= uint32(inst.rs2) << 20
	result |= (imm & 0xFE0) << 20 // imm[11:5]
	return result
}

func (a *Assembler) encodeBType(inst *Instruction) uint32 {
	imm := uint32(inst.immediate)
	result := uint32(inst.opcode)
	result |= ((imm >> 11) & 0x1) << 7 // imm[11]
	result |= ((imm >> 1) & 0xF) << 8  // imm[4:1]
	result |= uint32(inst.funct3) << 12
	result |= uint32(inst.rs1) << 15
	result |= uint32(inst.rs2) << 20
	result |= ((imm >> 5) & 0x3F) << 25 // imm[10:5]
	result |= ((imm >> 12) & 0x1) << 31 // imm[12]
	return result
}

func (a *Assembler) encodeUType(inst *Instruction) uint32 {
	result := uint32(inst.opcode)
	result |= uint32(inst.rd) << 7
	result |= (uint32(inst.immediate) & 0xFFFFF000)
	return result
}

func (a *Assembler) encodeJType(inst *Instruction) uint32 {
	imm := uint32(inst.immediate)
	result := uint32(inst.opcode)
	result |= uint32(inst.rd) << 7
	result |= ((imm >> 12) & 0xFF) << 12 // imm[19:12]
	result |= ((imm >> 11) & 0x1) << 20  // imm[11]
	result |= ((imm >> 1) & 0x3FF) << 21 // imm[10:1]
	result |= ((imm >> 20) & 0x1) << 31  // imm[20]
	return result
}

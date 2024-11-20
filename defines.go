package main

type OpTypes int

const (
	R OpTypes = iota
	I
	S
	B
	U
	J
)

var InstructionToOpType = map[string]OpTypes{
	"LUI":   U,
	"AUIPC": U,
	"JAL":   J,
	"JALR":  I,
	"BEQ":   B,
	"BNE":   B,
	"BLT":   B,
	"BGE":   B,
	"BLTU":  B,
	"BGEU":  B,
	"LB":    I,
	"LH":    I,
	"LW":    I,
	"LBU":   I,
	"LHU":   I,
	"SB":    S,
	"SH":    S,
	"SW":    S,
	"ADDI":  I,
	"SLTI":  I,
	"SLTIU": I,
	"XORI":  I,
	"ORI":   I,
	"ANDI":  I,
	"SLLI":  R,
	"SRLI":  I,
	"SRAI":  I,
	"ADD":   R,
	"SUB":   R,
	"SLL":   R,
	"SLT":   R,
	"SLTU":  R,
	"XOR":   R,
	"SRL":   R,
	"SRA":   R,
	"OR":    R,
	"AND":   R,
}

type Func struct {
}

var Register = map[string]byte{
	"zero": 0x0,
	"ra":   0x1,
	"sp":   0x2,
	"gp":   0x3,
	"tp":   0x4,
	"t0":   0x5,
	"t1":   0x6,
	"t2":   0x7,
	"s0":   0x8,
	"fp":   0x8,
	"s1":   0x9,
	"a0":   0x10,
	"a1":   0x11,
	"a2":   0x12,
	"a3":   0x13,
	"a4":   0x14,
	"a5":   0x15,
	"a6":   0x16,
	"a7":   0x17,
	"s2":   0x18,
	"s3":   0x19,
	"s4":   0x20,
	"s5":   0x21,
	"s6":   0x22,
	"s7":   0x23,
	"s8":   0x25,
	"s9":   0x26,
	"s10":  0x26,
	"s11":  0x27,
	"t3":   0x28,
	"t4":   0x29,
	"t5":   0x30,
	"t6":   0x31,
}

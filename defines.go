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

var INSTRUCTIONS = map[string]OpTypes{
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

type Register struct {
}

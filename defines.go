package main

type OpCode int

const (
	R OpCode = iota
	I
	S
	B
	U
	J
	CI  // Compressed I-type
	CSS // Compressed S-type
	CL  // Compressed Load
	CJ  // Compressed Jump
	CR  // Compressed R-type
	CB  // Compressed Branch
	CIW // Compressed Immediate Word
	CS  // Compressed Store
)

// String method to return the name of OpCode instead of its numeric value
func (op OpCode) String() string {
	names := []string{
		"R", "I", "S", "B", "U", "J",
		"CI", "CSS", "CL", "CJ", "CR", "CB", "CIW", "CS",
	}
	if op >= 0 && int(op) < len(names) {
		return names[op]
	}
	return "Unknown"
}

type OpPair struct {
	opType OpCode
	opByte []byte
}

var PseudoToInstruction = map[string]func([]string) []string{
	"mv":  handleMV,
	"j":   handleJ,
	"jr":  handleJR,
	"ble": handleBLE,
	"li":  handleLI,
}

var InstructionToOpType = map[string]OpPair{

	// U TYPE
	"lui":   {U, []byte{0b0110111}},
	"auipc": {U, []byte{0b0010111}},

	// J TYPE
	"jal": {J, []byte{0b1101111}},

	// B TYPE
	"beq":  {B, []byte{0b1100011, 0x0}},
	"bne":  {B, []byte{0b1100011, 0x1}},
	"blt":  {B, []byte{0b1100011, 0x4}},
	"bge":  {B, []byte{0b1100011, 0x5}},
	"bltu": {B, []byte{0b1100011, 0x6}},
	"bgeu": {B, []byte{0b1100011, 0x7}},

	// I TYPE: opbyte = opcode, func3, imm[11:5]
	"jalr":   {I, []byte{0b1100111, 0x0, 0x0}},
	"lb":     {I, []byte{0b0000011, 0x0, 0x0}},
	"lh":     {I, []byte{0b0000011, 0x1, 0x0}},
	"lw":     {I, []byte{0b0000011, 0x2, 0x0}},
	"lbu":    {I, []byte{0b0000011, 0x4, 0x0}},
	"lhu":    {I, []byte{0b0000011, 0x5, 0x0}},
	"sb":     {S, []byte{0b0100011, 0x0, 0x0}},
	"sh":     {S, []byte{0b0100011, 0x1, 0x0}},
	"sw":     {S, []byte{0b0100011, 0x2, 0x0}},
	"addi":   {I, []byte{0b0010011, 0x0, 0x0}},
	"slti":   {I, []byte{0b0010011, 0x2, 0x0}},
	"sltiu":  {I, []byte{0b0010011, 0x3, 0x0}},
	"xori":   {I, []byte{0b0010011, 0x4, 0x0}},
	"ori":    {I, []byte{0b0010011, 0x6, 0x0}},
	"andi":   {I, []byte{0b0010011, 0x7, 0x0}},
	"slli":   {I, []byte{0b0010011, 0x1, 0x00}},
	"srli":   {I, []byte{0b0010011, 0x5, 0x00}},
	"srai":   {I, []byte{0b0010011, 0x5, 0x32}},
	"ebreak": {I, []byte{0b1110011, 0x0, 0x1}}, // Adjust as needed
	"ecall":  {I, []byte{0b1110011, 0x0, 0x0}}, // Adjust as needed
	"call":   {I, []byte{}},                    // Adjust as needed

	// R TYPE
	"add":  {R, []byte{0b0110011, 0x0, 0x00}},
	"sub":  {R, []byte{0b0110011, 0x0, 0x20}},
	"sll":  {R, []byte{0b0110011, 0x1, 0x00}},
	"slt":  {R, []byte{0b0110011, 0x2, 0x00}},
	"sltu": {R, []byte{0b0110011, 0x3, 0x00}},
	"xor":  {R, []byte{0b0110011, 0x4, 0x00}},
	"srl":  {R, []byte{0b0110011, 0x5, 0x00}},
	"sra":  {R, []byte{0b0110011, 0x5, 0x20}},
	"or":   {R, []byte{0b0110011, 0x6, 0x00}},
	"and":  {R, []byte{0b0110011, 0x7, 0x00}},
	"nop":  {R, []byte{0b0000000, 0x0, 0x0}}, // Representing NOP, adjust as necessary

	// C TYPE (Compressed Instructions)
	"c.add": {CR, []byte{0b000000, 0b000, 0b000, 0b000, 0b000, 0b0110011}}, // C-type add
	"c.sub": {CR, []byte{0b000000, 0b000, 0b000, 0b000, 0b001, 0b0110011}}, // C-type sub
	"c.nop": {CR, []byte{0b000000, 0b000, 0b000, 0b000, 0b001, 0b0000000}}, // C-type nop
	"c.j":   {CJ, []byte{0b1100011}},                                       // Compressed jump
	"c.jal": {CJ, []byte{0b1101111}},                                       // Compressed jump and link
	"c.beq": {CB, []byte{0b000000}},                                        // Compressed branch if equal
	"c.bne": {CB, []byte{0b000001}},                                        // Compressed branch if not equal
	"c.lw":  {CL, []byte{0b000001}},                                        // Compressed load word
	"c.sw":  {CS, []byte{0b000001}},                                        // Compressed store word
}

func getTType(ti TokenType) string {
	switch ti {
	case global:
		return "global"
	case section:
		return "section"
	case data:
		return "data"
	case sign:
		return "sign"
	case modifier:
		return "modifier"
	case symbol:
		return "symbol"
	case localLabel:
		return "localLabel"
	case globalLabel:
		return "globalLabel"
	case complexValue:
		return "complexValue"
	case instruction:
		return "instruction"
	case register:
		return "register"
	case literal:
		return "literal"
	case entrypoint:
		return "entrypoint"
	case constant:
		return "constant"
	case constantValue:
		return "constantValue"
	case varValue:
		return "varValue"
	case varLabel:
		return "varLabel"
	case varSize:
		return "varSize"
	case constLabel:
		return "constLabel"
	case comment:
		return "comment"
	case ret:
		return "ret"
	case none:
		return "none"
	}
	return ""
}

func getOpCode(opc OpCode) string {
	switch opc {
	case R:
		return "R"
	case I:
		return "I"
	case S:
		return "S"
	case B:
		return "B"
	case U:
		return "U"
	case J:
		return "J"
	case CI:
		return ""
	case CSS:
		return ""
	case CL:
		return ""
	case CJ:
		return ""
	case CR:
		return ""
	case CB:
		return ""
	case CIW:
		return ""
	case CS:
		return ""
	}
	return ""
}

// return false if string
func getVarSize(vT string) (int, bool) {
	switch vT {
	case ".string":
		fallthrough
	case ".asciz":
		return 0, false
	case ".byte":
		return 8, true
	case ".hword":
		return 16, true
	case ".word":
		return 32, true
	case ".dword":
		return 64, true
	}
	return 0, true
}

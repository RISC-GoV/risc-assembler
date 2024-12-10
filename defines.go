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

type OpPair struct {
	opType OpCode
	opByte []byte
}

var InstructionToOpType = map[string]OpPair{

	// U TYPE
	"lui":   {U, []byte{0b0110111}},
	"auipc": {U, []byte{0b0010111}},

	// J TYPE
	"jal":  {J, []byte{0b1101111}},
	"jalr": {I, []byte{0b1100111, 0x0}},

	// B TYPE
	"beq":  {B, []byte{0b1100011, 0x0}},
	"bne":  {B, []byte{0b1100011, 0x1}},
	"blt":  {B, []byte{0b1100011, 0x4}},
	"bge":  {B, []byte{0b1100011, 0x5}},
	"bltu": {B, []byte{0b1100011, 0x6}},
	"bgeu": {B, []byte{0b1100011, 0x7}},

	// I TYPE
	"lb":     {I, []byte{0b0000011, 0x0}},
	"lh":     {I, []byte{0b0000011, 0x1}},
	"lw":     {I, []byte{0b0000011, 0x2}},
	"lbu":    {I, []byte{0b0000011, 0x4}},
	"lhu":    {I, []byte{0b0000011, 0x5}},
	"sb":     {S, []byte{0b0100011, 0x0}},
	"sh":     {S, []byte{0b0100011, 0x1}},
	"sw":     {S, []byte{0b0100011, 0x2}},
	"addi":   {I, []byte{0b0010011, 0x0}},
	"slti":   {I, []byte{0b0010011, 0x2}},
	"sltiu":  {I, []byte{0b0010011, 0x3}},
	"xori":   {I, []byte{0b0010011, 0x4}},
	"ori":    {I, []byte{0b0010011, 0x6}},
	"andi":   {I, []byte{0b0010011, 0x7}},
	"slli":   {I, []byte{0b0010011, 0x1, 0x00}},
	"srli":   {I, []byte{0b0010011, 0x5, 0x00}},
	"srai":   {I, []byte{0b0010011, 0x5, 0x20}},
	"ebreak": {I, []byte{0b0000000, 0x0, 0x0}}, // Adjust as needed
	"ecall":  {I, []byte{0b0000000, 0x0, 0x1}}, // Adjust as needed

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
	"c.li":  {CI, []byte{0b000000, 0b000}},                                 // Compressed load immediate
}

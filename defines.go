package main

type OpCode int

const (
	R OpCode = iota
	I
	S
	B
	U
	J
	CI
	CSS
	CL
	CJ
	CR
	CB
	CIW
	CS
)

type OpPair struct {
	opType OpCode
	opByte []byte
}

var InstructionToOpType = map[string]OpPair{

	// I TYPE

	"lui":    {U, []byte{0b0110111}},
	"auipc":  {U, []byte{0b0110011}},
	"jal":    {J, []byte{0b1101111}},
	"jalr":   {I, []byte{0b1101111, 0x0}},
	"beq":    {B, []byte{0b1100011, 0x0}},
	"bne":    {B, []byte{0b1100011, 0x1}},
	"blt":    {B, []byte{0b1100011, 0x4}},
	"bge":    {B, []byte{0b1100011, 0x5}},
	"bltu":   {B, []byte{0b1100011, 0x6}},
	"bgeu":   {B, []byte{0b1100011, 0x7}},
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
	"add":    {R, []byte{0b0110011, 0x0, 0x00}},
	"sub":    {R, []byte{0b0110011, 0x0, 0x20}},
	"sll":    {R, []byte{0b0110011, 0x1, 0x00}},
	"slt":    {R, []byte{0b0110011, 0x2, 0x00}},
	"sltu":   {R, []byte{0b0110011, 0x3, 0x00}},
	"xor":    {R, []byte{0b0110011, 0x4, 0x00}},
	"srl":    {R, []byte{0b0110011, 0x5, 0x00}},
	"sra":    {R, []byte{0b0110011, 0x5, 0x20}},
	"or":     {R, []byte{0b0110011, 0x6, 0x00}},
	"and":    {R, []byte{0b0110011, 0x7, 0x00}},
	"nop":    {R, []byte{0b0110011, 0x7, 0x00}}, //todo check code
	"ebreak": {I, []byte{0b0110011, 0x7, 0x00}}, //todo check code
	"ecall":  {I, []byte{0b0110011, 0x7, 0x00}}, //todo check code

	// C TYPE
	//TODO change binary and hex codes
	"fld":  {CL, []byte{0b0110011, 0x7, 0x00}},
	"flw":  {CI, []byte{0b0110011, 0x7, 0x00}},
	"fsd":  {CSS, []byte{0b0110011, 0x7, 0x00}},
	"fsw":  {CSS, []byte{0b0110011, 0x7, 0x00}},
	"li":   {CI, []byte{0b0110011, 0x7, 0x00}},
	"j":    {CJ, []byte{0b0110011, 0x7, 0x00}},
	"beqz": {CB, []byte{0b0110011, 0x7, 0x00}},
	"jr":   {CR, []byte{0b0110011, 0x7, 0x00}},
	"mv":   {CR, []byte{0b0110011, 0x7, 0x00}},
}

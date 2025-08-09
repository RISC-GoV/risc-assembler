package assembler

import (
	"testing"
)

var labelPositionsMockup map[string]int

func init() {
	// Global variables needed for testing
	labelPositionsMockup = map[string]int{
		// Assembly section labels
		".loop": 100,
		".exit": 200,
		".data": 300,
		".text": 400,
		".LC0":  500,
		".LC1":  200, // Updated value from second map

		// Program labels
		"main":  1000,
		"start": 2000,
		"end":   3000,

		// Data labels
		"data1":      4000,
		"data2":      5000,
		"test_label": 100,

		// Additional labels from third map
		"label1":    100,
		"label2":    200,
		"constant1": 300,
		"variable1": 400,
	}
}

func TestTranslateRType(t *testing.T) {
	type args struct {
		opcode int
		rd     int
		func3  int
		rs1    int
		rs2    int
		func7  int
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "add x1, x2, x3",
			args: args{
				opcode: 0b0110011, // 51 (R-type opcode)
				rd:     1,
				func3:  0b000,
				rs1:    2,
				rs2:    3,
				func7:  0b0000000,
			},
			want: 0b00000000001100010000000010110011,
		},
		{
			name: "sub x10, x20, x15",
			args: args{
				opcode: 0b0110011, // 51 (R-type opcode)
				rd:     10,
				func3:  0b000,
				rs1:    20,
				rs2:    15,
				func7:  0b0100000,
			},
			want: 0b01000000111110100000010100110011,
		},
		{
			name: "xor x5, x10, x31",
			args: args{
				opcode: 0b0110011, // 51 (R-type opcode)
				rd:     5,
				func3:  0b100,
				rs1:    10,
				rs2:    31,
				func7:  0b0000000,
			},
			want: 0b00000001111101010100001010110011,
		},
		{
			name: "or x15, x3, x4",
			args: args{
				opcode: 0b0110011, // 51 (R-type opcode)
				rd:     15,
				func3:  0b110,
				rs1:    3,
				rs2:    4,
				func7:  0b0000000,
			},
			want: 0b00000000010000011110011110110011,
		},
		{
			name: "and x31, x0, x5",
			args: args{
				opcode: 0b0110011, // 51 (R-type opcode)
				rd:     31,
				func3:  0b111,
				rs1:    0,
				rs2:    5,
				func7:  0b0000000,
			},
			want: 0b00000000010100000111111110110011,
		},
		{
			name: "sll x4, x5, x6",
			args: args{
				opcode: 0b0110011, // 51 (R-type opcode)
				rd:     4,
				func3:  0b001,
				rs1:    5,
				rs2:    6,
				func7:  0b0000000,
			},
			want: 0b00000000011000101001001000110011,
		},
		{
			name: "srl x2, x3, x4",
			args: args{
				opcode: 0b0110011, // 51 (R-type opcode)
				rd:     2,
				func3:  0b101,
				rs1:    3,
				rs2:    4,
				func7:  0b0000000,
			},
			want: 0b00000000010000011101000100110011,
		},
		{
			name: "sra x3, x4, x5",
			args: args{
				opcode: 0b0110011, // 51 (R-type opcode)
				rd:     3,
				func3:  0b101,
				rs1:    4,
				rs2:    5,
				func7:  0b0100000,
			},
			want: 0b01000000010100100101000110110011,
		},
		{
			name: "slt x1, x2, x3",
			args: args{
				opcode: 0b0110011, // 51 (R-type opcode)
				rd:     1,
				func3:  0b010,
				rs1:    2,
				rs2:    3,
				func7:  0b0000000,
			},
			want: 0b00000000001100010010000010110011,
		},
		{
			name: "sltu x10, x20, x30",
			args: args{
				opcode: 0b0110011, // 51 (R-type opcode)
				rd:     10,
				func3:  0b011,
				rs1:    20,
				rs2:    30,
				func7:  0b0000000,
			},
			want: 0b00000001111010100011010100110011,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TranslateRType(tt.args.opcode, tt.args.rd, tt.args.func3, tt.args.rs1, tt.args.rs2, tt.args.func7); got != tt.want {
				t.Errorf("TranslateRType() = 0x%08X, want 0x%08X", got, tt.want)
			}
		})
	}
}

func TestTranslateIType(t *testing.T) {
	type args struct {
		opcode int
		rd     int
		func3  int
		rs1    int
		imm    int
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "addi x1, x2, 10",
			args: args{
				opcode: 0b0010011, // 19 (I-type opcode for immediate operations)
				rd:     1,
				func3:  0b000,
				rs1:    2,
				imm:    10,
			},
			want: 0b00000000101000010000000010010011,
		},
		{
			name: "lw x10, 4(x15)",
			args: args{
				opcode: 0b0000011, // 3 (I-type opcode for loads)
				rd:     10,
				func3:  0b010,
				rs1:    15,
				imm:    4,
			},
			want: 0b00000000010001111010010100000011,
		},
		{
			name: "slti x5, x10, -100",
			args: args{
				opcode: 0b0010011, // 19 (I-type opcode for immediate operations)
				rd:     5,
				func3:  0b010,
				rs1:    10,
				imm:    -100,
			},
			want: 0b11111001110001010010001010010011,
		},
		{
			name: "sltiu x15, x20, 4095",
			args: args{
				opcode: 0b0010011, // 19 (I-type opcode for immediate operations)
				rd:     15,
				func3:  0b011,
				rs1:    20,
				imm:    4095, // Maximum 12-bit immediate value (0xFFF)
			},
			want: 0b11111111111110100011011110010011,
		},
		{
			name: "xori x31, x0, -1",
			args: args{
				opcode: 0b0010011, // 19 (I-type opcode for immediate operations)
				rd:     31,
				func3:  0b100,
				rs1:    0,
				imm:    -1,
			},
			want: 0b11111111111100000100111110010011,
		},
		{
			name: "ori x4, x5, 42",
			args: args{
				opcode: 0b0010011, // 19 (I-type opcode for immediate operations)
				rd:     4,
				func3:  0b110,
				rs1:    5,
				imm:    42,
			},
			want: 0b00000010101000101110001000010011,
		},
		{
			name: "andi x2, x3, 0xFF",
			args: args{
				opcode: 0b0010011, // 19 (I-type opcode for immediate operations)
				rd:     2,
				func3:  0b111,
				rs1:    3,
				imm:    0xFF,
			},
			want: 0b00001111111100011111000100010011,
		},
		{
			name: "jalr x1, 20(x3)",
			args: args{
				opcode: 0b1100111, // 103 (I-type opcode for jump and link register)
				rd:     1,
				func3:  0b000,
				rs1:    3,
				imm:    20,
			},
			want: 0b00000001010000011000000011100111,
		},
		{
			name: "lb x10, 0(x20)",
			args: args{
				opcode: 0b0000011, // 3 (I-type opcode for loads)
				rd:     10,
				func3:  0b000,
				rs1:    20,
				imm:    0,
			},
			want: 0b00000000000010100000010100000011,
		},
		{
			name: "lhu x5, -200(x10)",
			args: args{
				opcode: 0b0000011, // 3 (I-type opcode for loads)
				rd:     5,
				func3:  0b101,
				rs1:    10,
				imm:    -200,
			},
			want: 0b11110011100001010101001010000011,
		},
		{
			name: "Test with overflow values that should be masked",
			args: args{
				opcode: 0b11111111, // Should be masked to 0b1111111
				rd:     0b111111,   // Should be masked to 0b11111
				func3:  0b1111,     // Should be masked to 0b111
				rs1:    0b111111,   // Should be masked to 0b11111
				imm:    0xFFFFF,    // Should be masked to 0xFFF (12 bits)
			},
			want: 0b11111111111111111111111111111111,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TranslateIType(tt.args.opcode, tt.args.rd, tt.args.func3, tt.args.rs1, tt.args.imm); got != tt.want {
				t.Errorf("TranslateIType() = 0x%08X, want 0x%08X", got, tt.want)
			}
		})
	}
}

func TestTranslateSType(t *testing.T) {
	type args struct {
		opcode int
		func3  int
		rs1    int
		rs2    int
		imm    int
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "sw x1, 0(x2)",
			args: args{
				opcode: 0b0100011, // 35 (S-type opcode for stores)
				func3:  0b010,
				rs1:    2,
				rs2:    1,
				imm:    0,
			},
			want: 0b00000000000100010010000000100011,
		},
		{
			name: "sb x10, 4(x15)",
			args: args{
				opcode: 0b0100011, // 35 (S-type opcode for stores)
				func3:  0b000,
				rs1:    15,
				rs2:    10,
				imm:    4,
			},
			want: 0b00000000101001111000001000100011,
		},
		{
			name: "sh x20, -8(x5)",
			args: args{
				opcode: 0b0100011, // 35 (S-type opcode for stores)
				func3:  0b001,
				rs1:    5,
				rs2:    20,
				imm:    -8,
			},
			want: 0b11111111010000101001110000100011,
		},
		{
			name: "sw x31, 2047(x0)",
			args: args{
				opcode: 0b0100011, // 35 (S-type opcode for stores)
				func3:  0b010,
				rs1:    0,
				rs2:    31,
				imm:    2047, // Maximum positive 12-bit immediate
			},
			want: 0b01111111111100000010111110100011,
		},
		{
			name: "sw x5, -2048(x10)",
			args: args{
				opcode: 0b0100011, // 35 (S-type opcode for stores)
				func3:  0b010,
				rs1:    10,
				rs2:    5,
				imm:    -2048, // Minimum negative 12-bit immediate
			},
			want: 0b10000000010101010010000000100011,
		},
		{
			name: "sb x0, 100(x31)",
			args: args{
				opcode: 0b0100011, // 35 (S-type opcode for stores)
				func3:  0b000,
				rs1:    31,
				rs2:    0,
				imm:    100,
			},
			want: 0b00000110000011111000001000100011,
		},
		{
			name: "sh x15, -100(x20)",
			args: args{
				opcode: 0b0100011, // 35 (S-type opcode for stores)
				func3:  0b001,
				rs1:    20,
				rs2:    15,
				imm:    -100,
			},
			want: 0b11111000111110100001111000100011,
		},
		{
			name: "sw x10, 8(x5)",
			args: args{
				opcode: 0b0100011, // 35 (S-type opcode for stores)
				func3:  0b010,
				rs1:    5,
				rs2:    10,
				imm:    8,
			},
			want: 0b00000000101000101010010000100011,
		},
		{
			name: "sw x20, -4(x15)",
			args: args{
				opcode: 0b0100011, // 35 (S-type opcode for stores)
				func3:  0b010,
				rs1:    15,
				rs2:    20,
				imm:    -4,
			},
			want: 0b11111111010001111010111000100011,
		},
		{
			name: "sb x2, 1023(x3)",
			args: args{
				opcode: 0b0100011, // 35 (S-type opcode for stores)
				func3:  0b000,
				rs1:    3,
				rs2:    2,
				imm:    1023,
			},
			want: 0b00111110001000011000111110100011,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TranslateSType(tt.args.opcode, tt.args.func3, tt.args.rs1, tt.args.rs2, tt.args.imm); got != tt.want {
				t.Errorf("TranslateSType() = 0x%08X, want 0x%08X", got, tt.want)
			}
		})
	}
}

func TestTranslateBType(t *testing.T) {
	type args struct {
		opcode int
		func3  int
		rs1    int
		rs2    int
		imm    int
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "beq x1, x2, 16",
			args: args{
				opcode: 0b1100011, // 99 (B-type opcode for branches)
				func3:  0b000,
				rs1:    1,
				rs2:    2,
				imm:    16,
			},
			want: 0b00000000001000001000100001100011,
		},
		{
			name: "bne x10, x15, -16",
			args: args{
				opcode: 0b1100011, // 99 (B-type opcode for branches)
				func3:  0b001,
				rs1:    10,
				rs2:    15,
				imm:    -16,
			},
			want: 0b11111110111101010001100011100011,
		},
		{
			name: "blt x20, x5, 1024",
			args: args{
				opcode: 0b1100011, // 99 (B-type opcode for branches)
				func3:  0b100,
				rs1:    20,
				rs2:    5,
				imm:    1024,
			},
			want: 0b01000000010110100100000001100011,
		},
		{
			name: "bge x5, x10, -1024",
			args: args{
				opcode: 0b1100011, // 99 (B-type opcode for branches)
				func3:  0b101,
				rs1:    5,
				rs2:    10,
				imm:    -1024,
			},
			want: 0b11000000101000101101000011100011,
		},
		{
			name: "bltu x15, x20, 2",
			args: args{
				opcode: 0b1100011, // 99 (B-type opcode for branches)
				func3:  0b110,
				rs1:    15,
				rs2:    20,
				imm:    2,
			},
			want: 0b00000001010001111110000101100011,
		},
		{
			name: "bgeu x0, x31, -2",
			args: args{
				opcode: 0b1100011, // 99 (B-type opcode for branches)
				func3:  0b111,
				rs1:    0,
				rs2:    31,
				imm:    -2,
			},
			want: 0b11111111111100000111111111100011,
		},
		{
			name: "beq x3, x4, 2048",
			args: args{
				opcode: 0b1100011, // 99 (B-type opcode for branches)
				func3:  0b000,
				rs1:    3,
				rs2:    4,
				imm:    2048, // This will get masked because it's beyond 13-bit signed range
			},
			want: 0b00000000010000011000000011100011,
		},
		{
			name: "bne x25, x26, -2048",
			args: args{
				opcode: 0b1100011, // 99 (B-type opcode for branches)
				func3:  0b001,
				rs1:    25,
				rs2:    26,
				imm:    -2048,
			},
			want: 0b10000001101011001001000011100011,
		},
		{
			name: "Test with overflow values that should be masked",
			args: args{
				opcode: 0b11111111, // Should be masked to 0b1111111
				func3:  0b1111,     // Should be masked to 0b111
				rs1:    0b111111,   // Should be masked to 0b11111
				rs2:    0b111111,   // Should be masked to 0b11111
				imm:    0xFFFFF,    // Should be appropriately handled for B-type
			},
			want: 0b11111111111111111111111111111111, // Expected value after masking
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TranslateBType(tt.args.opcode, tt.args.func3, tt.args.rs1, tt.args.rs2, tt.args.imm); got != tt.want {
				t.Errorf("TranslateBType() = 0x%08X, want 0x%08X", got, tt.want)
			}
		})
	}
}

func TestTranslateUType(t *testing.T) {
	type args struct {
		opcode int
		rd     int
		imm    int
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "lui x1, 1",
			args: args{
				opcode: 0b0110111, // 55 (U-type opcode for load upper immediate)
				rd:     1,
				imm:    1,
			},
			want: 0b00000000000000000001000010110111,
		},
		{
			name: "auipc x10, 0xFFFFF",
			args: args{
				opcode: 0b0010111, // 23 (U-type opcode for add upper immediate to pc)
				rd:     10,
				imm:    0xFFFFF,
			},
			want: 0b11111111111111111111010100010111,
		},
		{
			name: "lui x15, 0x12345",
			args: args{
				opcode: 0b0110111, // 55 (U-type opcode for load upper immediate)
				rd:     15,
				imm:    0x12345,
			},
			want: 0b00010010001101000101011110110111,
		},
		{
			name: "auipc x20, 0xABCDE",
			args: args{
				opcode: 0b0010111, // 23 (U-type opcode for add upper immediate to pc)
				rd:     20,
				imm:    0xABCDE,
			},
			want: 0b10101011110011011110101000010111,
		},
		{
			name: "lui x0, 0",
			args: args{
				opcode: 0b0110111, // 55 (U-type opcode for load upper immediate)
				rd:     0,
				imm:    0,
			},
			want: 0b00000000000000000000000000110111,
		},
		{
			name: "auipc x31, 0x80000",
			args: args{
				opcode: 0b0010111, // 23 (U-type opcode for add upper immediate to pc)
				rd:     31,
				imm:    0x80000, // Will be masked to 20 bits
			},
			want: 0b10000000000000000000111110010111,
		},
		{
			name: "lui x5, -1",
			args: args{
				opcode: 0b0110111, // 55 (U-type opcode for load upper immediate)
				rd:     5,
				imm:    -1,
			},
			want: 0b11111111111111111111001010110111,
		},
		{
			name: "auipc x25, -4096",
			args: args{
				opcode: 0b0010111, // 23 (U-type opcode for add upper immediate to pc)
				rd:     25,
				imm:    -4096,
			},
			want: 0b11111111000000000000110010010111,
		},
		{
			name: "lui x10, 0x12",
			args: args{
				opcode: 0b0110111, // 55 (U-type opcode for load upper immediate)
				rd:     10,
				imm:    0x12,
			},
			want: 0b00000000000000010010010100110111,
		},
		{
			name: "auipc x15, 0xFFFFF",
			args: args{
				opcode: 0b0010111, // 23 (U-type opcode for add upper immediate to pc)
				rd:     15,
				imm:    0xFFFFF,
			},
			want: 0b11111111111111111111011110010111,
		},
		{
			name: "Test with overflow values that should be masked",
			args: args{
				opcode: 0b11111111, // Should be masked to 0b1111111
				rd:     0b111111,   // Should be masked to 0b11111
				imm:    0xFFFFFFFF, // Should be masked to 20 bits
			},
			want: 0b11111111111111111111111111111111, // Expected after masking
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TranslateUType(tt.args.opcode, tt.args.rd, tt.args.imm); got != tt.want {
				t.Errorf("TranslateUType() = 0x%08X, want 0x%08X", got, tt.want)
			}
		})
	}
}

func TestTranslateJType(t *testing.T) {
	type args struct {
		opcode int
		rd     int
		imm    int
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "Basic JAL instruction",
			args: args{
				opcode: 0b1101111, // JAL opcode
				rd:     1,         // x1 (ra)
				imm:    1024,      // Jump forward 1024 bytes
			},
			want: 0b01000000000000000000000011101111, // Expected binary encoding
		},
		{
			name: "JAL with maximum positive immediate",
			args: args{
				opcode: 0b1101111,
				rd:     10,      // x10 (a0)
				imm:    1048574, // Max positive jump distance
			},
			want: 0b01111111111111111111010101101111, // Expected binary encoding
		},
		{
			name: "JAL with negative immediate",
			args: args{
				opcode: 0b1101111,
				rd:     5,     // x5 (t0)
				imm:    -1024, // Jump backward 1024 bytes
			},
			want: 0b11000000000111111111001011101111, // Expected binary encoding
		},
		{
			name: "JAL with all instruction parts at max values",
			args: args{
				opcode: 0b1101111, // Will be masked to 0b1111111
				rd:     31,        // x31 (t6)
				imm:    0xFFFFF,   // Large immediate value
			},
			want: 0b01111111111111111111111111101111, // Expected binary encoding after masking
		},
		{
			name: "JAL with zero immediate",
			args: args{
				opcode: 0b1101111,
				rd:     0, // x0 (zero)
				imm:    0, // No jump distance
			},
			want: 0b00000000000000000000000001101111, // Expected binary encoding
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TranslateJType(tt.args.opcode, tt.args.rd, tt.args.imm); got != tt.want {
				t.Errorf("TranslateJType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgram_InstructionToBinary(t *testing.T) {
	tests := []struct {
		name               string
		token              *Token
		relativeInstrCount int
		want               uint32
		wantErr            bool
	}{
		{
			name: "R-Type Instruction - ADD",
			token: &Token{
				tokenType: instruction,
				value:     "add",
				opPair: &OpPair{
					opType: R,
					opByte: []byte{0b0110011, 0b000, 0b0000000}, // opcode, func3, func7 for ADD
				},
				children: []*Token{
					{tokenType: register, value: "x1"}, // rd
					{tokenType: register, value: "x2"}, // rs1
					{tokenType: register, value: "x3"}, // rs2
				},
			},
			relativeInstrCount: 0,
			want:               0b00000000001100010000000010110011, // Encoded ADD x1, x2, x3
			wantErr:            false,
		},
		{
			name: "I-Type Instruction - ADDI",
			token: &Token{
				tokenType: instruction,
				value:     "addi",
				opPair: &OpPair{
					opType: I,
					opByte: []byte{0b0010011, 0b000, 0b0000000}, // opcode, func3, unused for I-type
				},
				children: []*Token{
					{tokenType: register, value: "x1"}, // rd
					{tokenType: register, value: "x2"}, // rs1
					{tokenType: literal, value: "10"},  // immediate
				},
			},
			relativeInstrCount: 0,
			want:               0b00000000101000010000000010010011, // Encoded ADDI x1, x2, 10
			wantErr:            false,
		},
		{
			name: "I-Type Load Instruction - LW",
			token: &Token{
				tokenType: instruction,
				value:     "lw",
				opPair: &OpPair{
					opType: I,
					opByte: []byte{0b0000011, 0b010, 0b0000000}, // opcode, func3, unused for I-type
				},
				children: []*Token{
					{tokenType: register, value: "x1"}, // rd
					{tokenType: complexValue, value: "", children: []*Token{
						{tokenType: literal, value: "4"},   // offset
						{tokenType: register, value: "x2"}, // base register
					}},
				},
			},
			relativeInstrCount: 0,
			want:               0b00000000010000010010000010000011, // Encoded LW x1, 4(x2)
			wantErr:            false,
		},
		{
			name: "S-Type Instruction - SW",
			token: &Token{
				tokenType: instruction,
				value:     "sw",
				opPair: &OpPair{
					opType: S,
					opByte: []byte{0b0100011, 0b010, 0b0000000}, // opcode, func3, unused for S-type
				},
				children: []*Token{
					{tokenType: register, value: "x2"}, // rs1 (base)
					{tokenType: complexValue, value: "", children: []*Token{
						{tokenType: literal, value: "8"},   // offset
						{tokenType: register, value: "x1"}, // rs2 (source)
					}},
				},
			},
			relativeInstrCount: 0,
			want:               0b00000000000100010010010000100011, // Encoded SW x2, 8(x1)
			wantErr:            false,
		},
		{
			name: "B-Type Instruction - BEQ",
			token: &Token{
				tokenType: instruction,
				value:     "beq",
				opPair: &OpPair{
					opType: B,
					opByte: []byte{0b1100011, 0b000, 0b0000000}, // opcode, func3, unused for B-type
				},
				children: []*Token{
					{tokenType: register, value: "x1"}, // rs1
					{tokenType: register, value: "x2"}, // rs2
					{tokenType: literal, value: "16"},  // offset
				},
			},
			relativeInstrCount: 0,
			want:               0b00000000001000001000100001100011, // Encoded BEQ x1, x2, 16
			wantErr:            false,
		},
		{
			name: "U-Type Instruction - LUI",
			token: &Token{
				tokenType: instruction,
				value:     "lui",
				opPair: &OpPair{
					opType: U,
					opByte: []byte{0b0110111, 0b000, 0b0000000}, // opcode, unused, unused for U-type
				},
				children: []*Token{
					{tokenType: register, value: "x3"},     // rd
					{tokenType: literal, value: "1048576"}, // immediate (0x100000)
				},
			},
			relativeInstrCount: 0,
			want:               0b00000000000000000000000110110111, // Encoded LUI x3, 0x100000
			wantErr:            false,
		},
		{
			name: "J-Type Instruction - JAL",
			token: &Token{
				tokenType: instruction,
				value:     "jal",
				opPair: &OpPair{
					opType: J,
					opByte: []byte{0b1101111, 0b000, 0b0000000}, // opcode, unused, unused for J-type
				},
				children: []*Token{
					{tokenType: register, value: "x1"},  // rd
					{tokenType: literal, value: "1024"}, // immediate
				},
			},
			relativeInstrCount: 0,
			want:               0b01000000000000000000000011101111, // Encoded JAL x1, 1024
			wantErr:            false,
		},
		{
			name: "Special Instruction - ECALL",
			token: &Token{
				tokenType: instruction,
				value:     "ecall",
				opPair: &OpPair{
					opType: I,
					opByte: []byte{0b1110011, 0b000, 0b000000000000}, // opcode, func3, func12 for ECALL
				},
			},
			relativeInstrCount: 0,
			want:               0b00000000000000000000000001110011, // Encoded ECALL
			wantErr:            false,
		},
		{
			name: "Special Instruction - EBREAK",
			token: &Token{
				tokenType: instruction,
				value:     "ebreak",
				opPair: &OpPair{
					opType: I,
					opByte: []byte{0b1110011, 0b000, 0b000000000001}, // opcode, func3, func12 for EBREAK
				},
			},
			relativeInstrCount: 0,
			want:               0b00000000000100000000000001110011, // Encoded EBREAK
			wantErr:            false,
		},
		{
			name: "I-Type Instruction with Label - ADDI",
			token: &Token{
				tokenType: instruction,
				value:     "addi",
				opPair: &OpPair{
					opType: I,
					opByte: []byte{0b0010011, 0b000, 0b0000000}, // opcode, func3, unused for I-type
				},
				children: []*Token{
					{tokenType: register, value: "x1"},    // rd
					{tokenType: register, value: "x2"},    // rs1
					{tokenType: varLabel, value: ".loop"}, // label reference
				},
			},
			relativeInstrCount: 4,
			want:               0b00000110000000010000000010010011, // Assuming loop is at position 8
			wantErr:            false,
		},
		{
			name: "Character Literal in Immediate - ADDI",
			token: &Token{
				tokenType: instruction,
				value:     "addi",
				opPair: &OpPair{
					opType: I,
					opByte: []byte{0b0010011, 0b000, 0b0000000}, // opcode, func3, unused for I-type
				},
				children: []*Token{
					{tokenType: register, value: "x1"}, // rd
					{tokenType: register, value: "x0"}, // rs1
					{tokenType: literal, value: "'A'"}, // ASCII value of 'A' is 65
				},
			},
			relativeInstrCount: 0,
			want:               0b00000100000100000000000010010011, // Using wrong value for illustration
			wantErr:            false,
		},
		{
			name: "Invalid Instruction Type",
			token: &Token{
				tokenType: literal, // Not an instruction token
				value:     "not_an_instruction",
			},
			relativeInstrCount: 0,
			want:               0,
			wantErr:            true,
		},
		{
			name: "R-Type with Invalid Register Count",
			token: &Token{
				tokenType: instruction,
				value:     "add",
				opPair: &OpPair{
					opType: R,
					opByte: []byte{0b0110011, 0b000, 0b0000000}, // opcode, func3, func7 for ADD
				},
				children: []*Token{
					{tokenType: register, value: "x1"}, // rd
					{tokenType: register, value: "x2"}, // rs1
					// Missing rs2
				},
			},
			relativeInstrCount: 0,
			want:               0,
			wantErr:            true,
		},
		{
			name: "Hex Literal in Immediate - ADDI",
			token: &Token{
				tokenType: instruction,
				value:     "addi",
				opPair: &OpPair{
					opType: I,
					opByte: []byte{0b0010011, 0b000, 0b0000000}, // opcode, func3, unused for I-type
				},
				children: []*Token{
					{tokenType: register, value: "x1"},  // rd
					{tokenType: register, value: "x0"},  // rs1
					{tokenType: literal, value: "0xFF"}, // Hex value 255
				},
			},
			relativeInstrCount: 0,
			want:               0b00001111111100000000000010010011, // Encoded ADDI x1, x0, 255
			wantErr:            false,
		},
		{
			name: "Complex Value with Modifier - LUI with %hi",
			token: &Token{
				tokenType: instruction,
				value:     "lui",
				opPair: &OpPair{
					opType: U,
					opByte: []byte{0b0110111, 0b000, 0b0000000}, // opcode, unused, unused for U-type
				},
				children: []*Token{
					{tokenType: register, value: "x1"}, // rd
					{tokenType: complexValue, value: "", children: []*Token{
						{tokenType: literal, value: "%hi"},   // Modifier
						{tokenType: varLabel, value: "data"}, // Label
					}},
				},
			},
			relativeInstrCount: 8,
			want:               0b00000000000000000000000010110111, // Using a placeholder value
			wantErr:            true,                               // This would actually error as there's no labelPositions defined in the test
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { //reset labelPositions
			p := &Program{
				machinecode: []byte{0},
				variables:   []byte{0},
				constants:   []byte{0},
				strings:     []byte{0},
				entrypoint:  [4]byte{0},
				compilationVariables: &Compilation{
					labelPositions:              make(map[string]int),
					compilationEntryPoint:       "",
					instructionCount:            0,
					instructionCountCompilation: 0,
					variableCount:               0,
					constantCount:               0,
					stringCount:                 0,
					callbackInstructions:        nil,
				},
			}
			got, err := p.InstructionToBinary(tt.token, tt.relativeInstrCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstructionToBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("InstructionToBinary() = %032b, want %032b", got, tt.want)
				t.Errorf("InstructionToBinary() = 0x%08X, want 0x%08X", got, tt.want)
			}
		})
	}
}

func Test_parseIntValue(t *testing.T) {
	type args struct {
		valueStr string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Parse decimal integer",
			args:    args{valueStr: "123"},
			want:    123,
			wantErr: false,
		},
		{
			name:    "Parse negative decimal integer",
			args:    args{valueStr: "-456"},
			want:    -456,
			wantErr: false,
		},
		{
			name:    "Parse hexadecimal integer with 0x prefix",
			args:    args{valueStr: "0x1A"},
			want:    26,
			wantErr: false,
		},
		{
			name:    "Parse hexadecimal integer with 0X prefix",
			args:    args{valueStr: "0XFF"},
			want:    255,
			wantErr: false,
		},
		{
			name:    "Parse negative hexadecimal integer",
			args:    args{valueStr: "-0xFF"},
			want:    -255,
			wantErr: false,
		},
		{
			name:    "Parse character literal",
			args:    args{valueStr: "'A'"},
			want:    65, // ASCII for 'A'
			wantErr: false,
		},
		{
			name:    "Parse newline escape character",
			args:    args{valueStr: "'\\n'"},
			want:    10, // ASCII for newline
			wantErr: false,
		},
		{
			name:    "Parse tab escape character",
			args:    args{valueStr: "'\\t'"},
			want:    9, // ASCII for tab
			wantErr: false,
		},
		{
			name:    "Parse return escape character",
			args:    args{valueStr: "'\\r'"},
			want:    13, // ASCII for carriage return
			wantErr: false,
		},
		{
			name:    "Parse backslash escape character",
			args:    args{valueStr: "'\\\\'"},
			want:    92, // ASCII for backslash
			wantErr: false,
		},
		{
			name:    "Parse single quote escape character",
			args:    args{valueStr: "'\\''"},
			want:    39, // ASCII for single quote
			wantErr: false,
		},
		{
			name:    "Parse null character",
			args:    args{valueStr: "'\\0'"},
			want:    0, // ASCII for null
			wantErr: false,
		},
		{
			name:    "Parse with whitespace",
			args:    args{valueStr: " 123 "},
			want:    123,
			wantErr: false,
		},
		{
			name:    "Parse hex with whitespace",
			args:    args{valueStr: " 0xFF "},
			want:    255,
			wantErr: false,
		},
		{
			name:    "Invalid format",
			args:    args{valueStr: "abc"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Invalid hex format",
			args:    args{valueStr: "0xZZ"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Value too large for 32-bit",
			args:    args{valueStr: "9223372036854775807"}, // MaxInt64
			want:    0,
			wantErr: true,
		},
		{
			name:    "Invalid character literal - too many characters",
			args:    args{valueStr: "'ABC'"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Invalid escape sequence",
			args:    args{valueStr: "'\\z'"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Maximum 32-bit integer",
			args:    args{valueStr: "2147483647"}, // MaxInt32
			want:    2147483647,
			wantErr: false,
		},
		{
			name:    "Minimum 32-bit integer",
			args:    args{valueStr: "-2147483648"}, // MinInt32
			want:    -2147483648,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIntValue(tt.args.valueStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIntValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseIntValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgram_parseComplexValue(t *testing.T) {
	// Helper function to create tokens for testing
	createToken := func(tokenType TokenType, value string) *Token {
		return &Token{
			tokenType: tokenType,
			value:     value,
		}
	}

	// Helper function to create complex value tokens
	createComplexToken := func(firstValue string, firstType TokenType, secondValue string, secondType TokenType) *Token {
		return &Token{
			tokenType: complexValue,
			children: []*Token{
				{
					tokenType: firstType,
					value:     firstValue,
				},
				{
					tokenType: secondType,
					value:     secondValue,
				},
			},
		}
	}

	type fields struct {
		machinecode []byte
		variables   []byte
		constants   []byte
		strings     []byte
		entrypoint  [4]byte
	}
	type args struct {
		tok                *Token
		relativeInstrCount int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		want1   int
		wantErr bool
	}{
		{
			name: "Parse offset(register) - positive offset",
			fields: fields{
				machinecode: []byte{},
				variables:   []byte{},
				constants:   []byte{},
				strings:     []byte{},
				entrypoint:  [4]byte{},
			},
			args: args{
				tok:                createComplexToken("8", literal, "x1", register),
				relativeInstrCount: 0,
			},
			want:    1, // Register x1
			want1:   8, // Offset 8
			wantErr: false,
		},
		{
			name: "Parse offset(register) - negative offset",
			fields: fields{
				machinecode: []byte{},
				variables:   []byte{},
				constants:   []byte{},
				strings:     []byte{},
				entrypoint:  [4]byte{},
			},
			args: args{
				tok:                createComplexToken("-12", literal, "x2", register),
				relativeInstrCount: 0,
			},
			want:    2,   // Register x2
			want1:   -12, // Offset -12
			wantErr: false,
		},
		{
			name: "Parse offset(label) - constant reference",
			fields: fields{
				machinecode: []byte{},
				variables:   []byte{},
				constants:   []byte{},
				strings:     []byte{},
				entrypoint:  [4]byte{},
			},
			args: args{
				tok:                createComplexToken("4", literal, ".LC1", constantValue),
				relativeInstrCount: 0,
			},
			want:    labelPositionsMockup[".LC1"], // Label position
			want1:   4,                            // Offset 4
			wantErr: false,
		},
		{
			name: "Parse modifier - %lo",
			fields: fields{
				machinecode: []byte{},
				variables:   []byte{},
				constants:   []byte{},
				strings:     []byte{},
				entrypoint:  [4]byte{},
			},
			args: args{
				tok: &Token{
					tokenType: complexValue,
					children: []*Token{
						{
							tokenType: modifier,
							value:     "%lo",
						},
						{
							tokenType: varLabel,
							value:     "test_label",
						},
					},
				},
				relativeInstrCount: 0,
			},
			want:    0,
			want1:   labelPositionsMockup["test_label"] & 0xFFF, // Lower 12 bits of test_label position
			wantErr: false,
		},
		{
			name: "Parse modifier - %hi",
			fields: fields{
				machinecode: []byte{},
				variables:   []byte{},
				constants:   []byte{},
				strings:     []byte{},
				entrypoint:  [4]byte{},
			},
			args: args{
				tok: &Token{
					tokenType: complexValue,
					children: []*Token{
						{
							tokenType: modifier,
							value:     "%hi",
						},
						{
							tokenType: varLabel,
							value:     "test_label",
						},
					},
				},
				relativeInstrCount: 0,
			},
			want:    0,
			want1:   ((labelPositionsMockup["test_label"]) >> 12) & 0xFFFFF, // Upper 20 bits of test_label position
			wantErr: false,
		},
		{
			name: "Parse register directly",
			fields: fields{
				machinecode: []byte{},
				variables:   []byte{},
				constants:   []byte{},
				strings:     []byte{},
				entrypoint:  [4]byte{},
			},
			args: args{
				tok:                createToken(register, "x5"),
				relativeInstrCount: 0,
			},
			want:    0,
			want1:   5, // Register x5 value
			wantErr: false,
		},
		{
			name: "Parse literal directly",
			fields: fields{
				machinecode: []byte{},
				variables:   []byte{},
				constants:   []byte{},
				strings:     []byte{},
				entrypoint:  [4]byte{},
			},
			args: args{
				tok:                createToken(literal, "42"),
				relativeInstrCount: 0,
			},
			want:    0,
			want1:   42,
			wantErr: false,
		},
		{
			name: "Parse label directly",
			fields: fields{
				machinecode: []byte{},
				variables:   []byte{},
				constants:   []byte{},
				strings:     []byte{},
				entrypoint:  [4]byte{},
			},
			args: args{
				tok:                createToken(varLabel, "test_label"),
				relativeInstrCount: 4, // Assume we're at position 4
			},
			want:    0,
			want1:   96, // Label position (100) - relative position (4)
			wantErr: false,
		},
		{
			name: "Invalid label",
			fields: fields{
				machinecode: []byte{},
				variables:   []byte{},
				constants:   []byte{},
				strings:     []byte{},
				entrypoint:  [4]byte{},
			},
			args: args{
				tok:                createToken(varLabel, "nonexistent_label"),
				relativeInstrCount: 0,
			},
			want:    0,
			want1:   0,
			wantErr: true,
		},
		{
			name: "Invalid register",
			fields: fields{
				machinecode: []byte{},
				variables:   []byte{},
				constants:   []byte{},
				strings:     []byte{},
				entrypoint:  [4]byte{},
			},
			args: args{
				tok:                createToken(register, "x99"), // Invalid register number
				relativeInstrCount: 0,
			},
			want:    0,
			want1:   0,
			wantErr: true,
		},
		{
			name: "Parse label with relative position",
			fields: fields{
				machinecode: []byte{},
				variables:   []byte{},
				constants:   []byte{},
				strings:     []byte{},
				entrypoint:  [4]byte{},
			},
			args: args{
				tok:                createToken(varLabel, "test_label"),
				relativeInstrCount: 50, // We're at position 50
			},
			want:    0,
			want1:   50, // Label position (100) - relative position (50)
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Program{
				machinecode: tt.fields.machinecode,
				variables:   tt.fields.variables,
				constants:   tt.fields.constants,
				strings:     tt.fields.strings,
				entrypoint:  tt.fields.entrypoint,
			}
			got, got1, err := p.parseComplexValue(tt.args.tok, tt.args.relativeInstrCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Program.parseComplexValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Program.parseComplexValue() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Program.parseComplexValue() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_handleModifier(t *testing.T) {
	type args struct {
		mod                string
		val                int
		relativeInstrCount int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// Basic %lo test cases
		{
			name: "%lo with zero values",
			args: args{
				mod:                "%lo",
				val:                0,
				relativeInstrCount: 0,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "%lo with positive value",
			args: args{
				mod:                "%lo",
				val:                0x1234,
				relativeInstrCount: 0,
			},
			want:    0x234, // Only lower 12 bits
			wantErr: false,
		},
		{
			name: "%lo with positive value and relative instruction count",
			args: args{
				mod:                "%lo",
				val:                0x1234,
				relativeInstrCount: 0x100,
			},
			want:    0x334, // (0x1234 + 0x100) & 0xFFF
			wantErr: false,
		},
		{
			name: "%lo with value that wraps around 12 bits",
			args: args{
				mod:                "%lo",
				val:                0xFFF,
				relativeInstrCount: 1,
			},
			want:    0x000, // (0xFFF + 1) & 0xFFF = 0x1000 & 0xFFF = 0x000
			wantErr: false,
		},
		{
			name: "%lo with negative value",
			args: args{
				mod:                "%lo",
				val:                -0x1234,
				relativeInstrCount: 0,
			},
			want:    0xDCC, // -0x1234 & 0xFFF (preserves lowest 12 bits)
			wantErr: false,
		},

		// Basic %hi test cases
		{
			name: "%hi with zero values",
			args: args{
				mod:                "%hi",
				val:                0,
				relativeInstrCount: 0,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "%hi with positive value",
			args: args{
				mod:                "%hi",
				val:                0x12345,
				relativeInstrCount: 0,
			},
			want:    0x12, // (0x12345 >> 12) & 0xFFFFF
			wantErr: false,
		},
		{
			name: "%hi with positive value and relative instruction count",
			args: args{
				mod:                "%hi",
				val:                0x12000,
				relativeInstrCount: 0x1000,
			},
			want:    ((0x12000 + 0x1000) >> 12) & 0xFFFFF,
			wantErr: false,
		},
		{
			name: "%hi with negative value",
			args: args{
				mod:                "%hi",
				val:                -0x12345,
				relativeInstrCount: 0,
			},
			want:    (-0x12345 >> 12) & 0xFFFFF,
			wantErr: false,
		},

		// Basic %pcrel_lo test cases
		{
			name: "%pcrel_lo with zero value",
			args: args{
				mod:                "%pcrel_lo",
				val:                0,
				relativeInstrCount: 0,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "%pcrel_lo with positive value",
			args: args{
				mod:                "%pcrel_lo",
				val:                0x1234,
				relativeInstrCount: 0,
			},
			want:    0x234,
			wantErr: false,
		},
		{
			name: "%pcrel_lo with negative value",
			args: args{
				mod:                "%pcrel_lo",
				val:                -0x1234,
				relativeInstrCount: 0,
			},
			want:    0xDCC, // -0x1234 & 0xFFF
			wantErr: false,
		},
		{
			name: "%pcrel_lo with large value",
			args: args{
				mod:                "%pcrel_lo",
				val:                0x123456,
				relativeInstrCount: 0,
			},
			want:    0x456,
			wantErr: false,
		},

		// Basic %pcrel_hi test cases
		{
			name: "%pcrel_hi with zero value",
			args: args{
				mod:                "%pcrel_hi",
				val:                0,
				relativeInstrCount: 0,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "%pcrel_hi with positive value",
			args: args{
				mod:                "%pcrel_hi",
				val:                0x12345,
				relativeInstrCount: 0,
			},
			want:    0x12, // (0x12345 >> 12) & 0xFFFFF
			wantErr: false,
		},
		{
			name: "%pcrel_hi with negative value",
			args: args{
				mod:                "%pcrel_hi",
				val:                -0x12345,
				relativeInstrCount: 0,
			},
			want:    (-0x12345 >> 12) & 0xFFFFF,
			wantErr: false,
		},
		{
			name: "%pcrel_hi with large value",
			args: args{
				mod:                "%pcrel_hi",
				val:                0x876543210,
				relativeInstrCount: 0,
			},
			want:    0x76543, // (0x876543210 >> 12) & 0xFFFFF
			wantErr: false,
		},

		// Error cases
		{
			name: "unknown modifier",
			args: args{
				mod:                "%unknown",
				val:                0x1234,
				relativeInstrCount: 0,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "empty modifier",
			args: args{
				mod:                "",
				val:                0x1234,
				relativeInstrCount: 0,
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleModifier(tt.args.mod, tt.args.val, tt.args.relativeInstrCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleModifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handleModifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseLabelOrLiteral(t *testing.T) {
	type args struct {
		tok                    *Token
		instructionRelativePos int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// varLabel tests
		{
			name: "varLabel token - existing label",
			args: args{
				tok: &Token{
					tokenType: varLabel,
					value:     "label1",
				},
				instructionRelativePos: 10,
			},
			want:    labelPositionsMockup["label1"] - 10, // 100 - 10
			wantErr: false,
		},
		{
			name: "varLabel token - non-existing label",
			args: args{
				tok: &Token{
					tokenType: varLabel,
					value:     "nonexistent_label",
				},
				instructionRelativePos: 10,
			},
			want:    0,
			wantErr: true,
		},

		// varValue tests
		{
			name: "varValue token - existing variable",
			args: args{
				tok: &Token{
					tokenType: varValue,
					value:     "variable1",
				},
				instructionRelativePos: 100,
			},
			want:    labelPositionsMockup["variable1"] - 100,
			wantErr: false,
		},
		{
			name: "varValue token - non-existing variable",
			args: args{
				tok: &Token{
					tokenType: varValue,
					value:     "nonexistent_var",
				},
				instructionRelativePos: 100,
			},
			want:    0,
			wantErr: true,
		},

		// constantValue tests
		{
			name: "constantValue token - existing constant",
			args: args{
				tok: &Token{
					tokenType: constantValue,
					value:     "constant1",
				},
				instructionRelativePos: 50,
			},
			want:    labelPositionsMockup["constant1"] - 50,
			wantErr: false,
		},
		{
			name: "constantValue token - non-existing constant",
			args: args{
				tok: &Token{
					tokenType: constantValue,
					value:     "nonexistent_constant",
				},
				instructionRelativePos: 50,
			},
			want:    0,
			wantErr: true,
		},

		// literal tests
		{
			name: "literal token - decimal value",
			args: args{
				tok: &Token{
					tokenType: literal,
					value:     "42",
				},
				instructionRelativePos: 0,
			},
			want:    42,
			wantErr: false,
		},
		{
			name: "literal token - negative decimal value",
			args: args{
				tok: &Token{
					tokenType: literal,
					value:     "-42",
				},
				instructionRelativePos: 0,
			},
			want:    -42,
			wantErr: false,
		},
		{
			name: "literal token - hex value",
			args: args{
				tok: &Token{
					tokenType: literal,
					value:     "0xABC",
				},
				instructionRelativePos: 0,
			},
			want:    0xABC,
			wantErr: false,
		},
		{
			name: "literal token - negative hex value",
			args: args{
				tok: &Token{
					tokenType: literal,
					value:     "-0xABC",
				},
				instructionRelativePos: 0,
			},
			want:    -0xABC,
			wantErr: false,
		},
		{
			name: "literal token - invalid value",
			args: args{
				tok: &Token{
					tokenType: literal,
					value:     "not_a_number",
				},
				instructionRelativePos: 0,
			},
			want:    0,
			wantErr: true,
		},

		// register tests
		{
			name: "register token - valid register x0",
			args: args{
				tok: &Token{
					tokenType: register,
					value:     "x0",
				},
				instructionRelativePos: 0,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "register token - valid register x5",
			args: args{
				tok: &Token{
					tokenType: register,
					value:     "x5",
				},
				instructionRelativePos: 0,
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "register token - valid alias sp for x2",
			args: args{
				tok: &Token{
					tokenType: register,
					value:     "sp",
				},
				instructionRelativePos: 0,
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "register token - invalid register",
			args: args{
				tok: &Token{
					tokenType: register,
					value:     "x99",
				},
				instructionRelativePos: 0,
			},
			want:    0,
			wantErr: true,
		},

		// invalid token type test
		{
			name: "invalid token type",
			args: args{
				tok: &Token{
					tokenType: 999, // Invalid token type
					value:     "anything",
				},
				instructionRelativePos: 0,
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Program{}
			got, err := p.parseLabelOrLiteral(tt.args.tok, tt.args.instructionRelativePos)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLabelOrLiteral() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseLabelOrLiteral() = %v, want %v", got, tt.want)
			}
		})
	}
}

package assembler

import (
	"os"
	"reflect"
	"testing"
)

func fileFromString(content string) *os.File {
	tmpfile, _ := os.CreateTemp("", "testfile")
	tmpfile.WriteString(content)
	tmpfile.Seek(0, 0)
	return tmpfile
}

func TestPreprocess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"mv valid", "mv x1 x2", []string{"addi x1, x2, 0"}},
		{"mv invalid", "mv x1", []string{"invalid mv instruction"}},
		{"j valid", "j label", []string{"jal x0, label"}},
		{"j invalid", "j", []string{"invalid j instruction"}},
		{"jal with 2 args", "jal x1, label", []string{"jal x1, label"}},
		{"jal with 1 arg", "jal label", []string{"jal x1, label"}},
		{"jr valid", "jr x5", []string{"jalr x0, x5, 0"}},
		{"jr invalid", "jr", []string{"invalid jr instruction"}},
		{"add with 3 args", "add x1 x2 x3", []string{"add x1 x2 x3"}},
		{"add with 2 args", "add x1 x2", []string{"add x1 x1 x2"}},
		{"sub with 3 args", "sub x1 x2 x3", []string{"sub x1 x2 x3"}},
		{"sub with 2 args", "sub x1 x2", []string{"sub x1 x1 x2"}},
		{"ble valid", "ble x1 x2 label", []string{"bge x2, x1, label"}},
		{"ble invalid", "ble x1", []string{"invalid ble instruction"}},
		{"li small imm", "li x1 42", []string{"addi x1, x0, 42"}},
		{"li large imm", "li x1 5000", []string{
			"lui x1, %hi(5000)",
			"addi x1, x1, %lo(5000)",
		}},
		{"li symbol", "li x1 symbol", []string{
			"lui x1, %hi(symbol)",
			"addi x1, x1, %lo(symbol)",
		}},
		{"la valid", "la x1 symbol", []string{
			"auipc x1, %pcrel_hi(symbol)",
			"addi x1, x1, %pcrel_lo(symbol)",
		}},
		{"la invalid", "la x1", []string{"invalid la instruction"}},
		{"ret", "ret", []string{"jalr x0, 0(x1)"}},
		{"nop", "nop", []string{"addi x0, x0, 0"}},
		{"comments", "add x1 x2 x3 # this is a comment", []string{"add x1 x2 x3"}},
		{"empty line", "", []string{}},
		{"tab characters", "add\tx1\tx2\tx3", []string{"add x1 x2 x3"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := fileFromString(tt.input)
			defer os.Remove(file.Name())
			got := Preprocess(file)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Preprocess() = %v, want %v", got, tt.expected)
			}
		})
	}
}

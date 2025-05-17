package assembler

import "testing"

func TestAssembler_encodeRType(t *testing.T) {
	type fields struct {
		labels     map[string]int
		lineNumber int
		Token      *Token
		output     []uint32
		currentPC  int
	}
	type args struct {
		inst *Instruction
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Assembler{
				labels:     tt.fields.labels,
				lineNumber: tt.fields.lineNumber,
				Token:      tt.fields.Token,
				output:     tt.fields.output,
				currentPC:  tt.fields.currentPC,
			}
			if got := a.encodeRType(tt.args.inst); got != tt.want {
				t.Errorf("Assembler.encodeRType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssembler_encodeIType(t *testing.T) {
	type fields struct {
		labels     map[string]int
		lineNumber int
		Token      *Token
		output     []uint32
		currentPC  int
	}
	type args struct {
		inst *Instruction
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Assembler{
				labels:     tt.fields.labels,
				lineNumber: tt.fields.lineNumber,
				Token:      tt.fields.Token,
				output:     tt.fields.output,
				currentPC:  tt.fields.currentPC,
			}
			if got := a.encodeIType(tt.args.inst); got != tt.want {
				t.Errorf("Assembler.encodeIType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssembler_encodeSType(t *testing.T) {
	type fields struct {
		labels     map[string]int
		lineNumber int
		Token      *Token
		output     []uint32
		currentPC  int
	}
	type args struct {
		inst *Instruction
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Assembler{
				labels:     tt.fields.labels,
				lineNumber: tt.fields.lineNumber,
				Token:      tt.fields.Token,
				output:     tt.fields.output,
				currentPC:  tt.fields.currentPC,
			}
			if got := a.encodeSType(tt.args.inst); got != tt.want {
				t.Errorf("Assembler.encodeSType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssembler_encodeBType(t *testing.T) {
	type fields struct {
		labels     map[string]int
		lineNumber int
		Token      *Token
		output     []uint32
		currentPC  int
	}
	type args struct {
		inst *Instruction
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Assembler{
				labels:     tt.fields.labels,
				lineNumber: tt.fields.lineNumber,
				Token:      tt.fields.Token,
				output:     tt.fields.output,
				currentPC:  tt.fields.currentPC,
			}
			if got := a.encodeBType(tt.args.inst); got != tt.want {
				t.Errorf("Assembler.encodeBType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssembler_encodeUType(t *testing.T) {
	type fields struct {
		labels     map[string]int
		lineNumber int
		Token      *Token
		output     []uint32
		currentPC  int
	}
	type args struct {
		inst *Instruction
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Assembler{
				labels:     tt.fields.labels,
				lineNumber: tt.fields.lineNumber,
				Token:      tt.fields.Token,
				output:     tt.fields.output,
				currentPC:  tt.fields.currentPC,
			}
			if got := a.encodeUType(tt.args.inst); got != tt.want {
				t.Errorf("Assembler.encodeUType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssembler_encodeJType(t *testing.T) {
	type fields struct {
		labels     map[string]int
		lineNumber int
		Token      *Token
		output     []uint32
		currentPC  int
	}
	type args struct {
		inst *Instruction
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Assembler{
				labels:     tt.fields.labels,
				lineNumber: tt.fields.lineNumber,
				Token:      tt.fields.Token,
				output:     tt.fields.output,
				currentPC:  tt.fields.currentPC,
			}
			if got := a.encodeJType(tt.args.inst); got != tt.want {
				t.Errorf("Assembler.encodeJType() = %v, want %v", got, tt.want)
			}
		})
	}
}

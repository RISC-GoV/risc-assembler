package main

import (
	"testing"
)

func TestNewToken(t *testing.T) {
	str := "testing value"
	expected := &Token{
		tokenType: data,
		value:     str,
		children:  make([]Token, 0),
	}
	result := NewToken(data, str)
	if result.value != expected.value {
		t.Errorf("Expected %s, got %s", expected.value, result.value)
		return
	}
	if result.tokenType != expected.tokenType {
		t.Errorf("Expected %s, got %s", expected.tokenType, result.tokenType)
		return
	}
	if len(result.children) != len(expected.children) {
		t.Errorf("Expected %d, got %d", len(expected.children), len(result.children))
	}
	for i, child := range result.children {
		if child.tokenType != expected.children[i].tokenType {
			t.Errorf("Expected %s, got %s", expected.children[i].tokenType, child.tokenType)
		}
	}
}

func TestGetValue(t *testing.T) {
	token := &Token{
		tokenType: data,
		value:     "test",
		children:  make([]Token, 0),
	}
	expected := token.value
	result := token.getValue()
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestGetValueFromLiteral(t *testing.T) {
	token := &Token{
		tokenType: literal,
		value:     "123",
	}
	expected := 123
	result, err := token.getValueFromLiteral()
	if err != nil {
		t.Error(err)
	}
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestGetValueFromLiteralBinary(t *testing.T) {
	token := &Token{
		tokenType: literal,
		value:     "0b1101",
	}
	expected := 0b1101
	result, err := token.getValueFromLiteral()
	if err != nil {
		t.Error(err)
	}
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestGetValueFromLiteralHex(t *testing.T) {
	token := &Token{
		tokenType: literal,
		value:     "0xAA",
	}
	expected := 0xAA
	result, err := token.getValueFromLiteral()
	if err != nil {
		t.Error(err)
	}
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestGetValueFromLiteralErrorNotLiteral(t *testing.T) {
	token := &Token{
		tokenType: data,
		value:     "123",
	}
	_, err := token.getValueFromLiteral()
	if err == nil {
		t.Error("Expected error when getting literal value of non literal")
	}
}

func TestGetValueFromLiteralErrorNotNumber(t *testing.T) {
	token := &Token{
		tokenType: literal,
		value:     "abc",
	}
	_, err := token.getValueFromLiteral()
	if err == nil {
		t.Error("Expected error when getting literal value of non number")
	}
}

func TestGetRegisterFromABI(t1 *testing.T) {
	type fields struct {
		tokenType TokenType
		value     string
		children  []*Token
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name:    "Valid Register x0",
			fields:  fields{tokenType: register, value: "x0"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Valid Register x1",
			fields:  fields{tokenType: register, value: "x1"},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Valid Register x2 (sp)",
			fields:  fields{tokenType: register, value: "x2"},
			want:    2,
			wantErr: false,
		},
		{
			name:    "Valid Register sp",
			fields:  fields{tokenType: register, value: "sp"},
			want:    2,
			wantErr: false,
		},
		{
			name:    "Valid Register a0",
			fields:  fields{tokenType: register, value: "a0"},
			want:    10,
			wantErr: false,
		},
		{
			name:    "Valid Register t6",
			fields:  fields{tokenType: register, value: "t6"},
			want:    31,
			wantErr: false,
		},
		{
			name:    "Valid Register fp",
			fields:  fields{tokenType: register, value: "fp"},
			want:    8,
			wantErr: false,
		},
		{
			name:    "Valid Register s1",
			fields:  fields{tokenType: register, value: "s1"},
			want:    9,
			wantErr: false,
		},

		// Special case for "pc"
		{
			name:    "Valid Register pc",
			fields:  fields{tokenType: register, value: "pc"},
			want:    -1,
			wantErr: false, // pc is valid but does not belong to general registers
		},

		// Invalid register tests
		{
			name:    "Invalid Register invalid_register",
			fields:  fields{tokenType: register, value: "invalid_register"},
			want:    -1,
			wantErr: true,
		},
		{
			name:    "Invalid Register x32",
			fields:  fields{tokenType: register, value: "x32"},
			want:    -1,
			wantErr: true,
		},
		{
			name:    "Invalid Register out_of_range_register",
			fields:  fields{tokenType: register, value: "out_of_range_register"},
			want:    -1,
			wantErr: true,
		},

		// Non-register tokens
		{
			name:    "Non-Register Token",
			fields:  fields{tokenType: other, value: "non_register"},
			want:    -1,
			wantErr: true,
		},

		// Valid alternative names for registers
		{
			name:    "Valid alternative name for register x0 (zero)",
			fields:  fields{tokenType: register, value: "zero"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Valid alternative name for register x1 (ra)",
			fields:  fields{tokenType: register, value: "ra"},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Valid alternative name for register x2 (sp)",
			fields:  fields{tokenType: register, value: "sp"},
			want:    2,
			wantErr: false,
		},
		{
			name:    "Valid alternative name for register x3 (gp)",
			fields:  fields{tokenType: register, value: "gp"},
			want:    3,
			wantErr: false,
		},

		// Tests for higher registers
		{
			name:    "Valid Register x28 (t3)",
			fields:  fields{tokenType: register, value: "x28"},
			want:    28,
			wantErr: false,
		},
		{
			name:    "Valid Register x29 (t4)",
			fields:  fields{tokenType: register, value: "x29"},
			want:    29,
			wantErr: false,
		},
		{
			name:    "Valid Register x30 (t5)",
			fields:  fields{tokenType: register, value: "x30"},
			want:    30,
			wantErr: false,
		},
		{
			name:    "Valid Register x31 (t6)",
			fields:  fields{tokenType: register, value: "x31"},
			want:    31,
			wantErr: false,
		},

		// Edge cases with values like '0', 'x', etc.
		{
			name:    "Edge Case: x (single character)",
			fields:  fields{tokenType: register, value: "x"},
			want:    -1,
			wantErr: true, // Invalid register
		},
		{
			name:    "Edge Case: empty value",
			fields:  fields{tokenType: register, value: ""},
			want:    -1,
			wantErr: true, // Empty string is invalid
		},
		{
			name:    "Edge Case: number as register name",
			fields:  fields{tokenType: register, value: "123"},
			want:    -1,
			wantErr: true, // Numbers as register names should be invalid
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Token{
				tokenType: tt.fields.tokenType,
				value:     tt.fields.value,
				children:  tt.fields.children,
			}
			got, err := t.getRegisterFromABI()
			if (err != nil) != tt.wantErr {
				t1.Errorf("getRegisterFromABI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("getRegisterFromABI() got = %v, want %v", got, tt.want)
			}
		})
	}
}

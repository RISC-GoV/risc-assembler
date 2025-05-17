package assembler

import (
	"reflect"
	"testing"
)

func TestNewToken(t *testing.T) {
	parent := &Token{tokenType: global, value: "parent"}
	op := &OpPair{} // Define OpPair if needed elsewhere

	tests := []struct {
		name string
		args struct {
			tokenType     TokenType
			value         string
			parent        *Token
			pair_optional []*OpPair
		}
		want *Token
	}{
		{
			name: "Without OpPair",
			args: struct {
				tokenType     TokenType
				value         string
				parent        *Token
				pair_optional []*OpPair
			}{
				tokenType:     register,
				value:         "x1",
				parent:        parent,
				pair_optional: nil,
			},
			want: &Token{tokenType: register, value: "x1", parent: parent, children: []*Token{}},
		},
		{
			name: "With OpPair",
			args: struct {
				tokenType     TokenType
				value         string
				parent        *Token
				pair_optional []*OpPair
			}{
				tokenType:     register,
				value:         "x2",
				parent:        parent,
				pair_optional: []*OpPair{op},
			},
			want: &Token{tokenType: register, value: "x2", parent: parent, children: []*Token{}, opPair: op},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewToken(tt.args.tokenType, tt.args.value, tt.args.parent, tt.args.pair_optional...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_getValue(t *testing.T) {
	tok := &Token{value: "test"}
	if tok.getValue() != "test" {
		t.Errorf("getValue() = %v, want %v", tok.getValue(), "test")
	}
}

func TestToken_getRegisterNumericValue(t *testing.T) {
	tests := []struct {
		name    string
		token   *Token
		want    int
		wantErr bool
	}{
		{"Valid ABI name", &Token{tokenType: register, value: "a0"}, 10, false},
		{"Valid literal", &Token{tokenType: register, value: "5"}, 5, false},
		{"Invalid literal", &Token{tokenType: register, value: "abc"}, 0, true},
		{"Invalid type", &Token{tokenType: literal, value: "x1"}, -1, true},
		{"Out of range literal", &Token{tokenType: register, value: "40"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.token.getRegisterNumericValue()
			if (err != nil) != tt.wantErr {
				t.Errorf("getRegisterNumericValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("getRegisterNumericValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_getValueFromLiteral(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    int
		wantErr bool
	}{
		{"Valid literal", "12", 12, false},
		{"Invalid literal", "abc", 0, true},
		{"Out of range", "32", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := &Token{value: tt.value}
			got, err := tok.getValueFromLiteral()
			if (err != nil) != tt.wantErr {
				t.Errorf("getValueFromLiteral() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("getValueFromLiteral() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToken_getRegisterFromABI(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    int
		wantErr bool
	}{
		{"Valid name a1", "a1", 11, false},
		{"Valid name x2", "x2", 2, false},
		{"Invalid name", "xyz", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := &Token{value: tt.value}
			got, err := tok.getRegisterFromABI()
			if (err != nil) != tt.wantErr {
				t.Errorf("getRegisterFromABI() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("getRegisterFromABI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchTokenValid(t *testing.T) {
	tests := []struct {
		val     string
		want    int
		wantErr bool
	}{
		{"zero", 0, false},
		{"x1", 1, false},
		{"ra", 1, false},
		{"t6", 31, false},
		{"invalid", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.val, func(t *testing.T) {
			got, err := matchTokenValid(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchTokenValid() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("matchTokenValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

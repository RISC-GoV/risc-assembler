package assembler

import "testing"

func TestOpCode_String(t *testing.T) {
	tests := []struct {
		name string
		op   OpCode
		want string
	}{
		{"R type", R, "R"},
		{"I type", I, "I"},
		{"S type", S, "S"},
		{"B type", B, "B"},
		{"U type", U, "U"},
		{"J type", J, "J"},
		{"CI type", CI, "CI"},
		{"CSS type", CSS, "CSS"},
		{"CL type", CL, "CL"},
		{"CJ type", CJ, "CJ"},
		{"CR type", CR, "CR"},
		{"CB type", CB, "CB"},
		{"CIW type", CIW, "CIW"},
		{"CS type", CS, "CS"},
		{"Invalid type", OpCode(99), "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.op.String(); got != tt.want {
				t.Errorf("OpCode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTType(t *testing.T) {
	type args struct {
		ti TokenType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"global", args{global}, "global"},
		{"section", args{section}, "section"},
		{"modifier", args{modifier}, "modifier"},
		{"symbol", args{symbol}, "symbol"},
		{"localLabel", args{localLabel}, "localLabel"},
		{"globalLabel", args{globalLabel}, "globalLabel"},
		{"complexValue", args{complexValue}, "complexValue"},
		{"instruction", args{instruction}, "instruction"},
		{"register", args{register}, "register"},
		{"literal", args{literal}, "literal"},
		{"entrypoint", args{entrypoint}, "entrypoint"},
		{"constant", args{constant}, "constant"},
		{"constantValue", args{constantValue}, "constantValue"},
		{"varValue", args{varValue}, "varValue"},
		{"varLabel", args{varLabel}, "varLabel"},
		{"varSize", args{varSize}, "varSize"},
		{"unknown", args{TokenType(99)}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTType(tt.args.ti); got != tt.want {
				t.Errorf("getTType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOpCode(t *testing.T) {
	type args struct {
		opc OpCode
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"R type", args{R}, "R"},
		{"I type", args{I}, "I"},
		{"S type", args{S}, "S"},
		{"B type", args{B}, "B"},
		{"U type", args{U}, "U"},
		{"J type", args{J}, "J"},
		{"CI type", args{CI}, ""},
		{"CSS type", args{CSS}, ""},
		{"CL type", args{CL}, ""},
		{"CJ type", args{CJ}, ""},
		{"CR type", args{CR}, ""},
		{"CB type", args{CB}, ""},
		{"CIW type", args{CIW}, ""},
		{"CS type", args{CS}, ""},
		{"Invalid type", args{OpCode(99)}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOpCode(tt.args.opc); got != tt.want {
				t.Errorf("getOpCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getVarSize(t *testing.T) {
	type args struct {
		vT string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{"string", args{".string"}, 0, false},
		{"asciz", args{".asciz"}, 0, false},
		{"byte", args{".byte"}, 8, true},
		{"hword", args{".hword"}, 16, true},
		{"word", args{".word"}, 32, true},
		{"dword", args{".dword"}, 64, true},
		{"unknown", args{".unknown"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getVarSize(tt.args.vT)
			if got != tt.want {
				t.Errorf("getVarSize() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getVarSize() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

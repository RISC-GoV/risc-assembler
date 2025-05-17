package assembler

import (
	"os"
	"reflect"
	"testing"
)

func TestAssembler_Assemble(t *testing.T) {
	// Create a test file
	testContent := `.section .text
.globl main
main:
	addi sp, sp, -16
	sw ra, 12(sp)
	li a0, 42
	jal ra, print_int
	lw ra, 12(sp)
	addi sp, sp, 16
	ret

print_int:
	mv a1, a0
	li a0, 1
	ecall
	ret

.section .data
message: .string "Hello, World!"
value: .word 42`

	// Create folder if it doesn't exist
	if _, err := os.Stat("./tresult"); os.IsNotExist(err) {
		err = os.Mkdir("./tresult", 0755)
		if err != nil {
			t.Fatalf("Failed to create tresult directory: %v", err)
		}
	}

	// Schedule cleanup of the test directory
	defer func() {
		os.RemoveAll("./tresult")
	}()

	err := os.WriteFile("./tresult/test_input.s", []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	type fields struct {
		labels     map[string]int
		lineNumber int
		Token      *Token
		output     []uint32
		currentPC  int
	}
	type args struct {
		filename     string
		outputFolder string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Basic assembly test",
			fields: fields{
				labels:     make(map[string]int),
				lineNumber: 0,
				Token:      nil,
				output:     []uint32{},
				currentPC:  0,
			},
			args: args{
				filename:     "./tresult/test_input.s",
				outputFolder: "./tresult",
			},
			wantErr: false,
		},
		{
			name: "Non-existent file test",
			fields: fields{
				labels:     make(map[string]int),
				lineNumber: 0,
				Token:      nil,
				output:     []uint32{},
				currentPC:  0,
			},
			args: args{
				filename:     "./tresult/nonexistent_file.s",
				outputFolder: "./tresult",
			},
			wantErr: true,
		},
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
			err := a.Assemble(tt.args.filename, tt.args.outputFolder)
			if (err != nil) != tt.wantErr {
				t.Errorf("Assembler.Assemble() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_printTokenTree(t *testing.T) {
	// Create sample token tree for testing
	root := NewToken(global, "global", nil)
	section := NewToken(section, ".text", root)
	root.children = append(root.children, section)

	label := NewToken(globalLabel, "main:", section)
	section.children = append(section.children, label)

	op := &OpPair{opType: R, opByte: []byte{0x33}}
	instr := NewToken(instruction, "add", label, op)
	label.children = append(label.children, instr)

	instr.children = append(instr.children, NewToken(register, "x1", instr))
	instr.children = append(instr.children, NewToken(register, "x2", instr))
	instr.children = append(instr.children, NewToken(register, "x3", instr))

	type args struct {
		t     *Token
		depth int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test printTokenTree with depth 0",
			args: args{
				t:     root,
				depth: 0,
			},
			want: "global global\n\tsection .text\n\t\tglobalLabel main:\n\t\t\tinstruction add R\n\t\t\t\tregister x1\n\t\t\t\tregister x2\n\t\t\t\tregister x3\n",
		},
		{
			name: "Test printTokenTree with depth 1",
			args: args{
				t:     section,
				depth: 1,
			},
			want: "\tsection .text\n\t\tglobalLabel main:\n\t\t\tinstruction add R\n\t\t\t\tregister x1\n\t\t\t\tregister x2\n\t\t\t\tregister x3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := printTokenTree(tt.args.t, tt.args.depth); got != tt.want {
				t.Errorf("printTokenTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeEmptyStrings(t *testing.T) {
	type args struct {
		arr []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "No empty strings",
			args: args{
				arr: []string{"hello", "world", "test"},
			},
			want: []string{"hello", "world", "test"},
		},
		{
			name: "With empty strings",
			args: args{
				arr: []string{"hello", "", "world", "", "test"},
			},
			want: []string{"hello", "world", "test"},
		},
		{
			name: "Only empty strings",
			args: args{
				arr: []string{"", ""},
			},
			want: []string{},
		},
		{
			name: "With commas",
			args: args{
				arr: []string{"hello", ",", "world", ",", "test"},
			},
			want: []string{"hello", "world", "test"},
		},
		{
			name: "Empty array",
			args: args{
				arr: []string{},
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeEmptyStrings(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeEmptyStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanupStr(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "String with spaces",
			args: args{
				str: "  hello world  ",
			},
			want: "hello world",
		},
		{
			name: "String with commas",
			args: args{
				str: "hello,world",
			},
			want: "helloworld",
		},
		{
			name: "String with spaces and commas",
			args: args{
				str: "  hello, world  ",
			},
			want: "hello world",
		},
		{
			name: "Empty string",
			args: args{
				str: "",
			},
			want: "",
		},
		{
			name: "Only spaces",
			args: args{
				str: "    ",
			},
			want: "",
		},
		{
			name: "Only commas",
			args: args{
				str: ",,,",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanupStr(tt.args.str); got != tt.want {
				t.Errorf("cleanupStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssembler_Parse(t *testing.T) {
	// Create basic test environment
	root := NewToken(global, "global", nil)
	section := NewToken(section, ".text", root)
	root.children = append(root.children, section)

	type fields struct {
		labels     map[string]int
		lineNumber int
		Token      *Token
		output     []uint32
		currentPC  int
	}
	type args struct {
		lineParts []string
		parent    *Token
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Token
		wantErr bool
	}{
		{
			name: "Parse section directive",
			fields: fields{
				labels:     make(map[string]int),
				lineNumber: 0,
				Token:      root,
				output:     []uint32{},
				currentPC:  0,
			},
			args: args{
				lineParts: []string{".section", ".data"},
				parent:    root,
			},
			want:    root,
			wantErr: false,
		},
		{
			name: "Parse global label",
			fields: fields{
				labels:     make(map[string]int),
				lineNumber: 0,
				Token:      root,
				output:     []uint32{},
				currentPC:  0,
			},
			args: args{
				lineParts: []string{"main:"},
				parent:    section,
			},
			want:    section,
			wantErr: false,
		},
		{
			name: "Parse instruction with registers",
			fields: fields{
				labels:     make(map[string]int),
				lineNumber: 0,
				Token:      root,
				output:     []uint32{},
				currentPC:  0,
			},
			args: args{
				lineParts: []string{"add", "x1", "x2", "x3"},
				parent:    section,
			},
			want:    section,
			wantErr: false,
		},
		{
			name: "Parse unknown instruction",
			fields: fields{
				labels:     make(map[string]int),
				lineNumber: 0,
				Token:      root,
				output:     []uint32{},
				currentPC:  0,
			},
			args: args{
				lineParts: []string{"unknowninstr", "x1", "x2"},
				parent:    section,
			},
			want:    nil, // Will error before returning
			wantErr: true,
		},
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
			got, err := a.Parse(tt.args.lineParts, tt.args.parent)
			if (err != nil) != tt.wantErr {
				t.Errorf("Assembler.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("Assembler.Parse() returned nil, expected non-nil")
			}
		})
	}
}

func TestParseRegisters(t *testing.T) {
	// Create parent token for testing
	parent := NewToken(instruction, "add", nil)

	type args struct {
		strArr []string
		parent *Token
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantLen int
	}{
		{
			name: "Parse multiple registers",
			args: args{
				strArr: []string{"x1", "x2", "x3"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 3,
		},
		{
			name: "Parse with empty strings",
			args: args{
				strArr: []string{"x1", "", "x3"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name: "Parse empty array",
			args: args{
				strArr: []string{},
				parent: parent,
			},
			wantErr: false,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset parent children
			tt.args.parent.children = []*Token{}

			err := ParseRegisters(tt.args.strArr, tt.args.parent)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRegisters() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.args.parent.children) != tt.wantLen {
				t.Errorf("ParseRegisters() added %v children, want %v", len(tt.args.parent.children), tt.wantLen)
			}
		})
	}
}

func TestLexIType(t *testing.T) {
	// Sample instruction tokens for I-type instructions
	parent := NewToken(instruction, "addi", nil)

	type args struct {
		strArr []string
		parent *Token
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantLen int
	}{
		{
			name: "Simple I-type instruction",
			args: args{
				strArr: []string{"x1", "x2", "42"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 3, // 2 registers + 1 immediate
		},
		{
			name: "I-type instruction with offset and register",
			args: args{
				strArr: []string{"x1", "8(x2)"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 2, // 1 register + 1 complex value
		},
		{
			name: "I-type instruction with variable reference",
			args: args{
				strArr: []string{"x1", "x2", "variable_name"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 3, // 2 registers + 1 var reference
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset parent children
			tt.args.parent.children = []*Token{}

			err := LexIType(tt.args.strArr, tt.args.parent)
			if (err != nil) != tt.wantErr {
				t.Errorf("LexIType() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.args.parent.children) != tt.wantLen {
				t.Errorf("LexIType() added %v children, want %v", len(tt.args.parent.children), tt.wantLen)
			}
		})
	}
}

func TestLexSType(t *testing.T) {
	// Sample instruction tokens for S-type instructions
	parent := NewToken(instruction, "sw", nil)

	type args struct {
		strArr []string
		parent *Token
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantLen int
	}{
		{
			name: "Simple S-type instruction",
			args: args{
				strArr: []string{"x1", "8(x2)"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 2, // 1 register + 1 complex value
		},
		{
			name: "S-type instruction with variable reference",
			args: args{
				strArr: []string{"x1", "var_offset(x2)"},
				parent: parent,
			},
			wantErr: true, // Should error on non-numeric offset
			wantLen: 1,    // Only the register will be added
		},
		{
			name: "S-type instruction with modifier",
			args: args{
				strArr: []string{"x1", "%hi(symbol)(x2)"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 2, // 1 register + 1 complex value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset parent children
			tt.args.parent.children = []*Token{}

			err := LexSType(tt.args.strArr, tt.args.parent)
			if (err != nil) != tt.wantErr {
				t.Errorf("LexSType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLexBType(t *testing.T) {
	// Sample instruction tokens for B-type instructions
	parent := NewToken(instruction, "beq", nil)

	type args struct {
		strArr []string
		parent *Token
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantLen int
	}{
		{
			name: "Simple B-type instruction",
			args: args{
				strArr: []string{"x1", "x2", "label"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 3, // 2 registers + 1 label
		},
		{
			name: "B-type instruction with immediate",
			args: args{
				strArr: []string{"x1", "x2", "42"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 3, // 2 registers + 1 immediate
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset parent children
			tt.args.parent.children = []*Token{}

			err := LexBType(tt.args.strArr, tt.args.parent)
			if (err != nil) != tt.wantErr {
				t.Errorf("LexBType() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.args.parent.children) != tt.wantLen {
				t.Errorf("LexBType() added %v children, want %v", len(tt.args.parent.children), tt.wantLen)
			}
		})
	}
}

func TestLexUType(t *testing.T) {
	// Sample instruction tokens for U-type instructions
	parent := NewToken(instruction, "lui", nil)

	type args struct {
		strArr []string
		parent *Token
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantLen int
	}{
		{
			name: "Simple U-type instruction",
			args: args{
				strArr: []string{"x1", "4096"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 2, // 1 register + 1 immediate
		},
		{
			name: "U-type instruction with variable reference",
			args: args{
				strArr: []string{"x1", "variable_name"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 2, // 1 register + 1 var reference
		},
		{
			name: "U-type instruction with hi modifier",
			args: args{
				strArr: []string{"x1", "%hi(symbol)"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 2, // 1 register + 1 var reference with modifier
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset parent children
			tt.args.parent.children = []*Token{}

			err := LexUType(tt.args.strArr, tt.args.parent)
			if (err != nil) != tt.wantErr {
				t.Errorf("LexUType() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.args.parent.children) != tt.wantLen {
				t.Errorf("LexUType() added %v children, want %v", len(tt.args.parent.children), tt.wantLen)
			}
		})
	}
}

func TestLexJType(t *testing.T) {
	// Sample instruction tokens for J-type instructions
	parent := NewToken(instruction, "jal", nil)

	type args struct {
		strArr []string
		parent *Token
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantLen int
	}{
		{
			name: "Simple J-type instruction",
			args: args{
				strArr: []string{"x1", "label"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 2, // 1 register + 1 label
		},
		{
			name: "J-type instruction with immediate",
			args: args{
				strArr: []string{"x1", "42"},
				parent: parent,
			},
			wantErr: false,
			wantLen: 2, // 1 register + 1 immediate
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset parent children
			tt.args.parent.children = []*Token{}

			err := LexJType(tt.args.strArr, tt.args.parent)
			if (err != nil) != tt.wantErr {
				t.Errorf("LexJType() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.args.parent.children) != tt.wantLen {
				t.Errorf("LexJType() added %v children, want %v", len(tt.args.parent.children), tt.wantLen)
			}
		})
	}
}

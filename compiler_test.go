package assembler

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSplitValues(t *testing.T) {
	tests := []struct {
		name     string
		valueStr string
		want     []string
	}{
		{
			name:     "Single value",
			valueStr: "10",
			want:     []string{"10"},
		},
		{
			name:     "Multiple values",
			valueStr: "10, 20, 30",
			want:     []string{"10", "20", "30"},
		},
		{
			name:     "Values with whitespace",
			valueStr: " 10,  20 , 30 ",
			want:     []string{"10", "20", "30"},
		},
		{
			name:     "Empty string",
			valueStr: "",
			want:     []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitValues(tt.valueStr)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to create a temporary assembly file for testing
func createTempAssemblyFile(content string) (string, error) {
	tmpDir, err := ioutil.TempDir("", "asm_test")
	if err != nil {
		return "", err
	}

	tmpFile := filepath.Join(tmpDir, "test.s")
	if err := ioutil.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	return tmpFile, nil
}

// Helper function to clean up temporary test files
func cleanupTempFiles(path string) {
	os.RemoveAll(filepath.Dir(path))
	os.RemoveAll("./output.parser")
}

func TestProgramHandleString(t *testing.T) {
	// Test cases with assembly source code
	tests := []struct {
		name            string
		assemblySource  string
		expectedLabel   string
		expectedContent string
	}{
		{
			name: "Basic string handling",
			assemblySource: `
.data
string_label: .string "hello"
`,
			expectedLabel:   "string_label",
			expectedContent: "hello",
		},
		{
			name: "Empty string",
			assemblySource: `
.data
empty_string: .asciz ""
`,
			expectedLabel:   "empty_string",
			expectedContent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global variables for each test
			c := Compilation{}
			c.labelPositions = map[string]int{}
			c.stringCount = 8

			// Create temporary assembly file
			tempFile, err := createTempAssemblyFile(tt.assemblySource)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer cleanupTempFiles(tempFile)

			// Use the assembler to parse the file
			asm := Assembler{}
			asm.Token = NewToken(global, "", nil)

			file, err := os.Open(tempFile)
			if err != nil {
				t.Fatalf("Failed to open temp file: %v", err)
			}
			defer file.Close()

			lines := Preprocess(file)

			actualParent := asm.Token
			for _, line := range lines {
				lineParts := splitLine(line)
				if len(lineParts) == 0 {
					continue
				}
				actualParent, err = asm.Parse(lineParts, actualParent)
				if err != nil {
					t.Fatalf("Parse error: %v", err)
				}
			}

			// Find the string label token in the AST
			var stringToken *Token
			findStringLabel := func(token *Token) {
				if token.tokenType == varLabel && token.value == tt.expectedLabel+":" {
					stringToken = token
				}
				for _, child := range token.children {
					if child.tokenType == varLabel && child.value == tt.expectedLabel+":" {
						stringToken = child
					}
				}
			}

			// Traverse the AST to find our string token
			traverseToken(asm.Token, findStringLabel)

			if stringToken == nil {
				t.Fatalf("String label %s not found in token tree", tt.expectedLabel)
			}

			// Create a new program and handle the string token
			p := &Program{
				strings: []byte{0}, // Initial empty strings section
			}
			p.compilationVariables = &Compilation{}
			p.compilationVariables.labelPositions = map[string]int{}
			p.compilationVariables.stringCount = 8
			p.handleString(stringToken)

			// Check if the label was added correctly
			_, exists := p.compilationVariables.labelPositions[tt.expectedLabel]
			if !exists {
				t.Errorf("Label %s was not added to labelPositions", tt.expectedLabel)
			}

			// Extract the string content from p.strings
			// We need to skip the initial 0 byte and stop at the terminating 0 byte
			if len(p.strings) <= 1 {
				if tt.expectedContent != "" {
					t.Errorf("String content not found in program strings")
				}
			} else {
				// Start from position 1 (after initial 0)
				content := p.strings[1:]
				// Remove the terminating null byte if present
				if len(content) > 0 && content[len(content)-1] == 0 {
					content = content[:len(content)-1]
				}

				if string(content) != tt.expectedContent {
					t.Errorf("String content = %q, want %q", string(content), tt.expectedContent)
				}
			}
		})
	}
}

// Helper function to traverse token tree
func traverseToken(token *Token, visitor func(*Token)) {
	visitor(token)
	for _, child := range token.children {
		traverseToken(child, visitor)
	}
}

// Helper function to split a line into parts
func splitLine(line string) []string {
	// Split by whitespace and remove empty strings
	var parts []string
	for _, part := range strings.Fields(line) {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

func TestProgramRecursiveCompilation(t *testing.T) {
	tests := []struct {
		name              string
		assemblySource    string
		expectedVarsLen   int
		expectedConstsLen int
		expectedLabels    []string
	}{
		{
			name: "Variable label with .byte directive",
			assemblySource: `
.data
var_byte: .byte 10, 20, 30
`,
			expectedVarsLen:   3, // 3 bytes: 10, 20, 30
			expectedConstsLen: 0,
			expectedLabels:    []string{"var_byte"},
		},
		{
			name: "Constant with .hword directive",
			assemblySource: `
.const_hword: .hword 10, 20
`,
			expectedVarsLen:   0,
			expectedConstsLen: 4, // 2 hwords (2 bytes each): 10, 20
			expectedLabels:    []string{".const_hword"},
		},
		{
			name: "Global label",
			assemblySource: `
.text
global_label:
`,
			expectedVarsLen:   0,
			expectedConstsLen: 0,
			expectedLabels:    []string{"global_label"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Create temporary assembly file
			tempFile, err := createTempAssemblyFile(tt.assemblySource)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer cleanupTempFiles(tempFile)

			// Use the assembler to parse the file
			asm := Assembler{}
			asm.Token = NewToken(global, "", nil)

			file, err := os.Open(tempFile)
			if err != nil {
				t.Fatalf("Failed to open temp file: %v", err)
			}
			defer file.Close()

			lines := Preprocess(file)
			for _, line := range lines {
				lineParts := splitLine(line)
				if len(lineParts) == 0 {
					continue
				}
				_, err = asm.Parse(lineParts, asm.Token)
				if err != nil {
					t.Fatalf("Parse error: %v", err)
				}
			}

			// Create a new program and perform recursive compilation
			p := &Program{
				compilationVariables: &Compilation{
					labelPositions: make(map[string]int),
				},
				strings: []byte{0}, // Initial empty strings section
			}

			// Perform recursive compilation on the root token
			err = p.recursiveCompilation(asm.Token)
			if err != nil {
				t.Errorf("%s", err)
			}

			// Check lengths
			if len(p.variables) != tt.expectedVarsLen {
				t.Errorf("recursiveCompilation() variables length = %d, want %d", len(p.variables), tt.expectedVarsLen)
			}
			if len(p.constants) != tt.expectedConstsLen {
				t.Errorf("recursiveCompilation() constants length = %d, want %d", len(p.constants), tt.expectedConstsLen)
			}

			// Check labels
			for _, label := range tt.expectedLabels {
				_, exists := p.compilationVariables.labelPositions[label]
				if !exists {
					t.Errorf("Label %s was not added to labelPositions", label)
				}
			}
		})
	}
}

func TestCompile(t *testing.T) {
	tests := []struct {
		name               string
		assemblySource     string
		expectedCodeLen    int
		expectedVarsLen    int
		expectedConstsLen  int
		expectedEntrypoint string
	}{
		{
			name: "Simple program with main label",
			assemblySource: `
.text
main:
  add x0, x0, x0
`,
			expectedCodeLen:    4, // One instruction
			expectedVarsLen:    0,
			expectedConstsLen:  0,
			expectedEntrypoint: "main",
		},
		{
			name: "Program with explicit entry point",
			assemblySource: `
.text
.globl start
start:
  mv x0, x1
`,
			expectedCodeLen:    4, // One instruction
			expectedVarsLen:    0,
			expectedConstsLen:  0,
			expectedEntrypoint: "start",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary assembly file
			tempFile, err := createTempAssemblyFile(tt.assemblySource)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer cleanupTempFiles(tempFile)

			// Use the assembler to parse the file
			asm := Assembler{}
			err = asm.Assemble(tempFile, "")
			if err != nil {
				t.Fatalf("Assemble error: %v", err)
			}
			c := Compilation{}
			c.labelPositions = map[string]int{}
			c.stringCount = 8
			// Compile the program
			prog, err := c.compile(asm.Token)
			if err != nil {
				t.Fatalf("Compile error: %v", err)
			}

			// Check expected results
			if len(prog.machinecode) != tt.expectedCodeLen {
				t.Errorf("compile() machinecode length = %d, want %d", len(prog.machinecode), tt.expectedCodeLen)
			}
			if len(prog.variables) != tt.expectedVarsLen {
				t.Errorf("compile() variables length = %d, want %d", len(prog.variables), tt.expectedVarsLen)
			}
			if len(prog.constants) != tt.expectedConstsLen {
				t.Errorf("compile() constants length = %d, want %d", len(prog.constants), tt.expectedConstsLen)
			}
			compilationEntryPoint := prog.compilationVariables.compilationEntryPoint
			// Check if the entrypoint matches
			if compilationEntryPoint != "" && compilationEntryPoint != tt.expectedEntrypoint {
				t.Errorf("compile() entrypoint = %s, want %s", compilationEntryPoint, tt.expectedEntrypoint)
			}
		})
	}
}

func TestProgramCallDescendants(t *testing.T) {
	// Create an assembly source with a simple structure
	assemblySource := `
.text
label1:
  add x0, x0, x0
  sub x1, x1, x1
label2:
  addi x2, x2, 80
`

	// Create temporary assembly file
	tempFile, err := createTempAssemblyFile(assemblySource)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer cleanupTempFiles(tempFile)

	// Use the assembler to parse the file
	asm := Assembler{}
	asm.Token = NewToken(global, "", nil)

	file, err := os.Open(tempFile)
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	lines := Preprocess(file)

	actualParent := asm.Token
	for _, line := range lines {
		lineParts := splitLine(line)
		if len(lineParts) == 0 {
			continue
		}
		actualParent, err = asm.Parse(lineParts, actualParent)
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
	}

	// Find the .text section token
	var textSectionToken *Token
	findTextSection := func(token *Token) {
		if token.tokenType == section && token.value == ".text" {
			textSectionToken = token
		}
	}

	traverseToken(asm.Token, findTextSection)

	if textSectionToken == nil {
		t.Fatalf(".text section not found in token tree")
	}

	// Test callDescendants function
	calledTokens := []*Token{}

	// Inject tracking function
	mockRecursive := func(token *Token) error {
		calledTokens = append(calledTokens, token)
		return nil
	}

	p := &Program{}
	p.callDescendants(textSectionToken, mockRecursive)

	// Check if all children of the text section were processed
	if len(calledTokens) != len(textSectionToken.children) {
		t.Errorf("Expected %d calls, got %d", len(textSectionToken.children), len(calledTokens))
	}

	// Verify that each child was processed
	for i, child := range textSectionToken.children {
		if i < len(calledTokens) && calledTokens[i] != child {
			t.Errorf("Expected token %v at index %d, got %v", child, i, calledTokens[i])
		}
	}
}

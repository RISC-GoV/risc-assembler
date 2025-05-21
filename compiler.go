package assembler

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Compilation struct {
	labelPositions              map[string]int //= make(map[string]int)
	compilationEntryPoint       string
	instructionCount            int
	instructionCountCompilation int
	variableCount               int
	constantCount               int
	stringCount                 int //= 8
	callbackInstructions        [][2]interface{}
}

func (c *Compilation) compile(token *Token) (Program, error) {
	prog := Program{}
	prog.compilationVariables = c
	prog.strings = append(prog.strings, uint8(00))
	err := prog.recursiveCompilation(token)
	if err != nil {
		return Program{}, err
	}
	fmt.Print("final instructions size (should match instructions): ")
	fmt.Println(prog.compilationVariables.instructionCountCompilation)

	for _, fun := range prog.compilationVariables.callbackInstructions {
		err := fun[0].(func(int) error)(fun[1].(int))
		if err != nil {
			return Program{}, err
		}
	}
	if len(prog.strings) == 1 {
		prog.strings = nil
	}
	if prog.compilationVariables.compilationEntryPoint != "" {
		val, _ := prog.compilationVariables.labelPositions[prog.compilationVariables.compilationEntryPoint]
		var bts = make([]byte, 4)
		binary.LittleEndian.PutUint32(bts, uint32(val))
		prog.entrypoint = [4]byte(bts)
	} else {
		val, _ := prog.compilationVariables.labelPositions["main"]
		var bts = make([]byte, 4)
		binary.LittleEndian.PutUint32(bts, uint32(val))
		prog.entrypoint = [4]byte(bts)
	}
	return prog, nil
}

func (p *Program) recursiveCompilation(token *Token) error {
	switch token.tokenType {
	case varLabel:
		var varValue []byte
		switch token.children[0].value {
		case ".string":
			fallthrough
		case ".asciz":
			p.handleString(token)
			goto endGoTo
		case ".byte":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					return err
				}
				varValue = append(varValue, byte(val))
			}
		case ".hword":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					return err
				}
				varValue = binary.LittleEndian.AppendUint16(varValue, uint16(val))
			}
		case ".word":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					return err
				}
				varValue = binary.LittleEndian.AppendUint32(varValue, uint32(val))
			}
		case ".dword":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					return err
				}
				varValue = binary.LittleEndian.AppendUint64(varValue, uint64(val))
			}
		}

		p.compilationVariables.labelPositions[strings.ReplaceAll(token.value, ":", "")] = p.compilationVariables.instructionCountCompilation + len(p.variables)

		p.variables = append(p.variables, varValue...)
	case constant:
		var varValue []byte
		switch token.children[0].value {
		case ".string":
			fallthrough
		case ".asciz":
			p.handleString(token)
			goto endGoTo
		case ".byte":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					return err
				}
				varValue = append(varValue, byte(val))
			}
		case ".hword":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					return err
				}
				varValue = binary.LittleEndian.AppendUint16(varValue, uint16(val))
			}
		case ".word":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					return err
				}
				varValue = binary.LittleEndian.AppendUint32(varValue, uint32(val))
			}
		case ".dword":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					return err
				}
				varValue = binary.LittleEndian.AppendUint64(varValue, uint64(val))
			}
		default:
			if token.children[0].tokenType == instruction {
				p.compilationVariables.labelPositions[strings.ReplaceAll(token.value, ":", "")] = p.compilationVariables.instructionCountCompilation + len(p.constants)
				p.callDescendants(token, p.recursiveCompilation)
				goto endGoTo
			}
		}
		p.compilationVariables.labelPositions[strings.ReplaceAll(token.value, ":", "")] = p.compilationVariables.instructionCountCompilation

		p.constants = append(p.constants, varValue...)

	// case constant:
	// 	p.compilationVariables.labelPositions[strings.Replace(token.value, ":", "", 1)] = p.compilationVariables.instructionCountCompilation + p.compilationVariables.variableCount + len(p.constants)
	// 	p.callDescendants(token)
	case globalLabel:
		p.compilationVariables.labelPositions[strings.Replace(token.value, ":", "", 1)] = p.compilationVariables.instructionCountCompilation
		fallthrough
	case section:
		fallthrough
	case global:
		err := p.callDescendants(token, p.recursiveCompilation)
		if err != nil {
			return err
		}
	case instruction:
		p.compilationVariables.callbackInstructions = append(p.compilationVariables.callbackInstructions,
			[2]interface{}{func(relativeInstrCount int) error {
				val, err := p.InstructionToBinary(token, relativeInstrCount)
				if err != nil {
					return err
				}
				p.machinecode = binary.LittleEndian.AppendUint32(p.machinecode, val)
				return nil
			},
				p.compilationVariables.instructionCountCompilation})
		p.compilationVariables.instructionCountCompilation += 4
	case entrypoint:
		if token.value == ".globl" {
			p.compilationVariables.compilationEntryPoint = token.children[0].value
		}

	}
endGoTo:
	return nil
}

func (p *Program) callDescendants(token *Token, recursionFn func(*Token) error) error {
	for _, child := range token.children {
		err := recursionFn(child)
		if err != nil {
			return err
		}
	}
	return nil
}

// Helper function to split comma-separated values and trim whitespace since we repeat it in all vars
func splitValues(valueStr string) []string {
	values := strings.Split(valueStr, ",")
	for i, v := range values {
		values[i] = strings.TrimSpace(v)
	}
	return values
}

func (p *Program) handleString(token *Token) {
	p.compilationVariables.labelPositions[strings.ReplaceAll(token.value, ":", "")] = p.compilationVariables.instructionCount + len(p.strings)
	for _, ch := range token.children[1].value {
		p.strings = append(p.strings, byte(ch))
	}
	p.strings = append(p.strings, uint8(0))
}

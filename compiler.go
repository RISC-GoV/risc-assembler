package assembler

import (
	"encoding/binary"
	"fmt"
	"strings"
)

func compile(token *Token) Program {
	prog := Program{}
	prog.strings = append(prog.strings, uint8(00))
	prog.recursiveCompilation(token)
	fmt.Print("final instructions size (should match instructions): ")
	fmt.Println(instructionCountCompilation)

	for _, fun := range callbackInstructions {
		fun[0].(func(int))(fun[1].(int))
	}
	if len(prog.strings) == 1 {
		prog.strings = nil
	}
	if compilationEntryPoint != "" {
		val, _ := labelPositions[compilationEntryPoint]
		var bts = make([]byte, 4)
		binary.LittleEndian.PutUint32(bts, uint32(val))
		prog.entrypoint = [4]byte(bts)
	} else {
		val, _ := labelPositions["main"]
		var bts = make([]byte, 4)
		binary.LittleEndian.PutUint32(bts, uint32(val))
		prog.entrypoint = [4]byte(bts)
	}
	return prog
}

var (
	labelPositions              map[string]int = make(map[string]int)
	compilationEntryPoint       string
	instructionCount            int
	instructionCountCompilation int
	variableCount               int
	constantCount               int
	stringCount                 int = 8
	callbackInstructions        [][2]interface{}
)

func (p *Program) recursiveCompilation(token *Token) {
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
					panic(err)
				}
				varValue = append(varValue, byte(val))
			}
		case ".hword":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					panic(err)
				}
				varValue = binary.LittleEndian.AppendUint16(varValue, uint16(val))
			}
		case ".word":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					panic(err)
				}
				varValue = binary.LittleEndian.AppendUint32(varValue, uint32(val))
			}
		case ".dword":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					panic(err)
				}
				varValue = binary.LittleEndian.AppendUint64(varValue, uint64(val))
			}
		}

		labelPositions[strings.ReplaceAll(token.value, ":", "")] = instructionCountCompilation + len(p.variables)

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
					panic(err)
				}
				varValue = append(varValue, byte(val))
			}
		case ".hword":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					panic(err)
				}
				varValue = binary.LittleEndian.AppendUint16(varValue, uint16(val))
			}
		case ".word":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					panic(err)
				}
				varValue = binary.LittleEndian.AppendUint32(varValue, uint32(val))
			}
		case ".dword":
			// Handle multiple comma-separated values
			values := splitValues(token.children[1].value)
			for _, valueStr := range values {
				val, err := parseIntValue(valueStr)
				if err != nil {
					panic(err)
				}
				varValue = binary.LittleEndian.AppendUint64(varValue, uint64(val))
			}
		default:
			if token.children[0].tokenType == instruction {
				labelPositions[strings.ReplaceAll(token.value, ":", "")] = instructionCountCompilation + len(p.constants)
				p.callDescendants(token)
				goto endGoTo
			}
		}
		labelPositions[strings.ReplaceAll(token.value, ":", "")] = instructionCountCompilation

		p.constants = append(p.constants, varValue...)

	// case constant:
	// 	labelPositions[strings.Replace(token.value, ":", "", 1)] = instructionCountCompilation + variableCount + len(p.constants)
	// 	p.callDescendants(token)
	case globalLabel:
		labelPositions[strings.Replace(token.value, ":", "", 1)] = instructionCountCompilation
		fallthrough
	case section:
		fallthrough
	case global:
		p.callDescendants(token)
	case instruction:
		callbackInstructions = append(callbackInstructions,
			[2]interface{}{func(relativeInstrCount int) {
				val, err := p.InstructionToBinary(token, relativeInstrCount)
				if err != nil {
					panic(err)
				}
				p.machinecode = binary.LittleEndian.AppendUint32(p.machinecode, val)
			},
				instructionCountCompilation})
		instructionCountCompilation += 4
	case entrypoint:
		if token.value == ".globl" {
			compilationEntryPoint = token.children[0].value
		}

	}
endGoTo:
}

func (p *Program) callDescendants(token *Token) {
	for _, tk := range token.children {
		p.recursiveCompilation(tk)
	}
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
	labelPositions[strings.ReplaceAll(token.value, ":", "")] = instructionCount + len(p.strings)
	for _, ch := range token.children[1].value {
		p.strings = append(p.strings, byte(ch))
	}
	p.strings = append(p.strings, uint8(0))
}

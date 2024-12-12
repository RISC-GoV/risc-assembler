package main

import (
	"encoding/binary"
	"strconv"
	"strings"
)

func compile(token *Token) Program {
	prog := Program{}
	prog.strings = append(prog.strings, uint8(00))
	prog.recursiveCompilation(token)

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
	labelPositions        map[string]int = make(map[string]int)
	compilationEntryPoint string
	instructionCount      int
	variableCount         int
	constantCount         int
	stringCount           int = 8
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
			val, err := strconv.Atoi(token.children[1].value)
			if err != nil {
				panic(err)
			}
			varValue = append(varValue, byte(val))
		case ".hword":
			val, err := strconv.Atoi(token.children[1].value)
			if err != nil {
				panic(err)
			}
			varValue = binary.LittleEndian.AppendUint16(varValue, uint16(val))
		case ".word":
			val, err := strconv.Atoi(token.children[1].value)
			if err != nil {
				panic(err)
			}
			varValue = binary.LittleEndian.AppendUint32(varValue, uint32(val))
		case ".dword":
			val, err := strconv.Atoi(token.children[1].value)
			if err != nil {
				panic(err)
			}
			varValue = binary.LittleEndian.AppendUint64(varValue, uint64(val))
		}

		labelPositions[strings.ReplaceAll(token.value, ":", "")] = instructionCount + len(p.variables)

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
			val, err := strconv.Atoi(token.children[1].value)
			if err != nil {
				panic(err)
			}
			varValue = append(varValue, byte(val))
		case ".hword":
			val, err := strconv.Atoi(token.children[1].value)
			if err != nil {
				panic(err)
			}
			varValue = binary.LittleEndian.AppendUint16(varValue, uint16(val))
		case ".word":
			val, err := strconv.Atoi(token.children[1].value)
			if err != nil {
				panic(err)
			}
			varValue = binary.LittleEndian.AppendUint32(varValue, uint32(val))
		case ".dword":
			val, err := strconv.Atoi(token.children[1].value)
			if err != nil {
				panic(err)
			}
			varValue = binary.LittleEndian.AppendUint64(varValue, uint64(val))
		}
		p.constants = append(p.constants, varValue...)
	// case constant:
	// 	labelPositions[strings.Replace(token.value, ":", "", 1)] = instructionCount + variableCount + len(p.constants)
	// 	p.callDescendants(token)
	case globalLabel:
		labelPositions[strings.Replace(token.value, ":", "", 1)] = len(p.machinecode)
		fallthrough
	case section:
		fallthrough
	case global:
		p.callDescendants(token)
	case instruction:
		val, err := p.InstructionToBinary(token)
		if err != nil {
			panic(err)
		}
		p.machinecode = binary.LittleEndian.AppendUint32(p.machinecode, val)
	case modifier:
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

func (p *Program) handleString(token *Token) {
	for _, ch := range token.children[1].value {
		p.strings = append(p.strings, byte(ch))
	}
	p.strings = append(p.strings, uint8(0))
}

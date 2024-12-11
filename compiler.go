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
	return prog
}

type labelMention struct {
	name            string
	mentionPosition string
}

var (
	labelMentions         []labelMention
	labelPositions        map[string]int
	compilationEntryPoint string
)

func (p *Program) recursiveCompilation(token *Token) {
	switch token.tokenType {
	case global:
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
			binary.LittleEndian.PutUint16(varValue, uint16(val))
		case ".word":
			val, err := strconv.Atoi(token.children[1].value)
			if err != nil {
				panic(err)
			}
			binary.LittleEndian.PutUint32(varValue, uint32(val))
		case ".dword":
			val, err := strconv.Atoi(token.children[1].value)
			if err != nil {
				panic(err)
			}
			binary.LittleEndian.PutUint64(varValue, uint64(val))
		}

		labelPositions[token.children[0].value] = len(p.variables)
		p.variables = append(p.variables, varValue...)
	case constantValue:
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
			binary.LittleEndian.PutUint16(varValue, uint16(val))
		case ".word":
			val, err := strconv.Atoi(token.children[1].value)
			if err != nil {
				panic(err)
			}
			binary.LittleEndian.PutUint32(varValue, uint32(val))
		case ".dword":
			val, err := strconv.Atoi(token.children[1].value)
			if err != nil {
				panic(err)
			}
			binary.LittleEndian.PutUint64(varValue, uint64(val))
		}

		labelPositions[token.children[0].value] = len(p.constants)
		p.constants = append(p.constants, varValue...)
	case constant:
		labelPositions[strings.Replace(strings.Replace(token.value, ".", "", 1), ":", "", 1)] = len(p.machinecode)
		p.callDescendants(token)
	case globalLabel:
		labelPositions[strings.Replace(token.value, ".", "", 1)] = len(p.machinecode)
		fallthrough
	case section:
		p.callDescendants(token)
	case instruction:
		val, err := InstructionToBinary(token)
		if err != nil {
			panic(err)
		}
		binary.LittleEndian.PutUint32(p.machinecode, val)
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

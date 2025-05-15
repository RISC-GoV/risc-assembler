package main

import (
	"errors"
	"fmt"
	"strconv"
)

type TokenType int

const (
	global = iota
	section
	symbol
	localLabel
	globalLabel
	modifier
	complexValue
	instruction
	register
	literal
	entrypoint
	constant
	constantValue
	varValue
	varLabel
	varSize
)

type Token struct {
	tokenType TokenType
	value     string
	opPair    *OpPair
	children  []*Token
	parent    *Token
}

func NewToken(tokenType TokenType, value string, parent *Token, pair_optional ...*OpPair) *Token {

	var pair *OpPair
	if len(pair_optional) > 0 {
		pair = pair_optional[0]
	}
	return &Token{tokenType: tokenType, value: value, children: make([]*Token, 0), parent: parent, opPair: pair}
}

func (t *Token) getValue() string {
	return t.value
}

func (t *Token) getValueFromLiteral() (int, error) {
	if t.tokenType != literal {
		return 0, errors.New("token is not of type literal")
	}
	ret, err := strconv.Atoi(t.value)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (t *Token) getRegisterFromABI() (int, error) {
	if t.tokenType != register {
		return -1, errors.New("token is not of type register")
	}
	res, err := matchTokenValid(t.value)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func matchTokenValid(val string) (int, error) {
	switch val {
	case "x0", "zero":
		return 0, nil
	case "x1", "ra":
		return 1, nil
	case "x2", "sp":
		return 2, nil
	case "x3", "gp":
		return 3, nil
	case "x4", "tp":
		return 4, nil
	case "x5", "t0":
		return 5, nil
	case "x6", "t1":
		return 6, nil
	case "x7", "t2":
		return 7, nil
	case "x8", "s0", "fp":
		return 8, nil
	case "x9", "s1":
		return 9, nil
	case "x10", "a0":
		return 10, nil
	case "x11", "a1":
		return 11, nil
	case "x12", "a2":
		return 12, nil
	case "x13", "a3":
		return 13, nil
	case "x14", "a4":
		return 14, nil
	case "x15", "a5":
		return 15, nil
	case "x16", "a6":
		return 16, nil
	case "x17", "a7":
		return 17, nil
	case "x18", "s2":
		return 18, nil
	case "x19", "s3":
		return 19, nil
	case "x20", "s4":
		return 20, nil
	case "x21", "s5":
		return 21, nil
	case "x22", "s6":
		return 22, nil
	case "x23", "s7":
		return 23, nil
	case "x24", "s8":
		return 24, nil
	case "x25", "s9":
		return 25, nil
	case "x26", "s10":
		return 26, nil
	case "x27", "s11":
		return 27, nil
	case "x28", "t3":
		return 28, nil
	case "x29", "t4":
		return 29, nil
	case "x30", "t5":
		return 30, nil
	case "x31", "t6":
		return 31, nil
	case "pc":
		return -1, nil
	default:
		return -1, fmt.Errorf("invalid register: %d", register)
	}
}

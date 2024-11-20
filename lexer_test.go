package main

import "testing"

func TestLexLine(t *testing.T) {
	asm := "add x6, x7, x28"
	var _ = []byte{
		byte(0b00000001),
		byte(0b11000011),
		byte(0b10000011),
		byte(0b00110011)}
	_, err := LexLine(asm, nil) //TODO add parent
	if err != nil {
		t.Error(err)
		return
	}
	//TODO add tests
}

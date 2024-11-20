package main

import "testing"

func TestLexLine(t *testing.T) {
	asm := "add x6, x7, x28"
	var expected = []byte{
		byte(0b00000001),
		byte(0b11000011),
		byte(0b10000011),
		byte(0b00110011)}
	result, err := LexLine(asm, nil)
	if err != nil {
		t.Error(err)
		return
	}
	for i, b := range result {
		if b != expected[i] {
			t.Errorf("%d: expected %s, got %s", i, b, result[i])
			return
		}
	}
}

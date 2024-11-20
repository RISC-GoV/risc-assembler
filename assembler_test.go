package main

import (
	"testing"
)

func TestRType(t *testing.T) {
	opcode := 0b0110011
	rd := 0b00010
	func3 := 0b000
	rs1 := 0b00011
	rs2 := 0b00100
	func7 := 0b0000000

	expected := 0b00000000010000011000000100110011
	result := RType(opcode, rd, func3, rs1, rs2, func7)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
		return
	}
}

func TestIType(t *testing.T) {
	opcode := 0b0010011
	rd := 0b00001
	func3 := 0b000
	rs1 := 0b00010
	imm := 0b000000000101

	expected := 0b00000000010100010000000010010011
	result := TranslateIType(opcode, rd, func3, rs1, imm)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
		return
	}
}

func TestSType(t *testing.T) {
	opcode := 0b0100011
	ihm4_0 := 0b00101
	func3 := 0b000
	rs1 := 0b00000
	rs2 := 0b01100
	imm11_5 := 0b0000000

	expected := 0b00000000110000000000001010100011
	result := TranslateSType(opcode, ihm4_0, func3, rs1, rs2, imm11_5)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
		return
	}
}

func TestBType(t *testing.T) {
	opcode := 0b1100011
	ihm_11 := 0b0
	ihm4_1 := 0b0110
	func3 := 0b000
	rs1 := 0b10111
	rs2 := 0b01111
	imm10_5 := 0b000000
	imm12 := 0b0

	expected := 0b00000000111110111000011001100011
	result := TranslateBType(opcode, ihm_11, ihm4_1, func3, rs1, rs2, imm10_5, imm12)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
		return
	}
}

func TestUType(t *testing.T) {
	opcode := 0b0110111
	rd := 0b10001
	ihm := 0b00000000000000010011

	expected := 0b00000000000000010011100010110111
	result := TranslateUType(opcode, rd, ihm)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
		return
	}
}

func TestJType(t *testing.T) {
	opcode := 0b1101111
	rd := 0b10011
	ihm19_12 := 0b00000000
	ihm11 := 0b0
	ihm10_1 := 0b0000000111
	ihm20 := 0b0

	expected := 0b00000000111000000000100111101111
	result := TranslateJType(opcode, rd, ihm19_12, ihm11, ihm10_1, ihm20)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
		return
	}
}

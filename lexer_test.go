package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestLexFile(t *testing.T) {
	LexFile("./test-files/TestLexFileText.asm", "./test-files/TestLexFileRESULT.bin")
}

func TestTokenizeString(t *testing.T) {
	const str = `
add(int, int):
        addi    sp,sp,-32
        sw      ra,28(sp)
        sw      s0,24(sp)
        addi    s0,sp,32
        sw      a0,-20(s0)
        sw      a1,-24(s0)
        lw      a4,-20(s0)
        lw      a5,-24(s0)
        add     a5,a4,a5
        mv      a0,a5
        lw      ra,28(sp)
        lw      s0,24(sp)
        addi    sp,sp,32
        jr      ra
.LC0:
        .string "Hello World"
.LC1:
        .string "%d"
main:
        addi    sp,sp,-16
        sw      ra,12(sp)
        sw      s0,8(sp)
        addi    s0,sp,16
        lui     a5,%hi(.LC0)
        addi    a0,a5,%lo(.LC0)
        call    printf
        li      a1,1
        li      a0,1
        call    add(int, int)
        mv      a5,a0
        mv      a1,a5
        lui     a5,%hi(.LC1)
        addi    a0,a5,%lo(.LC1)
        call    printf
        li      a5,0
        mv      a0,a5
        lw      ra,12(sp)
        lw      s0,8(sp)
        addi    sp,sp,16
        jr      ra		
`
	var token *Token = TokenizeString(str)
	PrintTokenTree(token, 0)
}

func PrintTokenTree(token *Token, depth int) {
	fmt.Print("|"+repeatChar(" ", depth), strings.Split(token.value, " ")[0], "   OF TYPE", token.tokenType, "|  ")
	for _, tok := range token.children {
		PrintTokenTree(tok, depth+1)
	}
}

func repeatChar(char string, depth int) string {
	return strings.Repeat(char, depth)
}

func TestLexLine(t *testing.T) {
	asms := []string{
		"add t1, t2, t3",
		"addi t1, t2, 100",
		"jal t1, 100",
		"lui t1, 100",
		"beq t1, t2, 100",
		"sw t1, 100(t2)",
	}
	var token *Token = &Token{global, "", []*Token{}}
	for index, asm := range asms {
		_, err := TokenizeLine(asm, token) //TODO add parent
		if err != nil {
			t.Error(err, "FOREACH RAN UNTIL LINE (0-5):", index)
			return
		}
	}
	//TODO add tests
}

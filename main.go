package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: assembler <input.s>")
		os.Exit(1)
	}
	asm := Assembler{}
	asm.Token = NewToken(global, "", nil)
	if err := asm.Assemble(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func copybytes(source []byte, dest []byte, amt int, sourceoffset int, destoffset int) {
	for i := 0; i < amt; i++ {
		dest[destoffset+i] = source[sourceoffset+i]
	}
}
func GenerateELFSectionHeaders() []byte {
	return []byte{}
}

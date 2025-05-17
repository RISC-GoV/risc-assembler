# risc-assembler


## example usage
`
go get github.com/RISC-GoV/risc-assembler@dirty
`
```go
package main

import (
	"fmt"
	"github.com/RISC-GoV/risc-assembler"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: assembler <input.s>")
		os.Exit(1)
	}
	asm := assembler.Assembler{}
	asm.Token = assembler.NewToken(global, "", nil)
	if err := asm.Assemble(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

```
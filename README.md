# risc-assembler


## example usage
`
go get github.com/RISC-GoV/risc-assembler@dirty
`
if the latest commit was not taken then go on branch dirty,
take the latest commit hash then
`
go get github.com/RISC-GoV/risc-assembler@[HASH]
`

```go
package main

import (
	"fmt"
	"os"

	assembler "github.com/RISC-GoV/risc-assembler"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <input.asm> <outputFolde?r>")
		os.Exit(1)
	}
	if len(os.Args) == 3 {
		asm := assembler.Assembler{}
		if err := asm.Assemble(os.Args[1], os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}else {
		asm := assembler.Assembler{}
		if err := asm.Assemble(os.Args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}
```
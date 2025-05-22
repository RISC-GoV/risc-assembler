# risc-assembler

A RISC-V assembler written in Go, designed to convert RISC-V assembly code into executable ELF binary files.

## Features

- Support for RISC-V instructions
- ELF file generation
- Integrated preprocessor
- Instruction encoding
- Unit tests for main components

## Installation

To install the library in your project:

```bash
go get github.com/RISC-GoV/risc-assembler
```

## Usage

After installation, you can use the assembler programmatically in Go:

### API Functions

### - `assembler.Assembler.Assemble(filename string, outputFolder string) error`
Assembles the RISC-V assembly file at `filename` and writes an ELF binary to `outputFolder` or returns an error.

### - `assembler.Assembler.AssembleLine(line string) ([]byte, error)`
Assembles the RISC-V assembly string `line` and returns its corresponding byte code or an error.

### - `assembler.Preprocess(file *os.File) []string`
Processes macros and directives in the `source` file and returns cleaned instructions.

### - `assembler.PreprocessLine(line string) []string`
Processes macros and directives in `line` and returns cleaned instructions.

### - `assembler.BuildELFFile(program Program) *[]byte`
Wraps machine code into a valid ELF binary format.

**Example: Assembling a RISC-V file**

```go
package main

import (
    "github.com/RISC-GoV/risc-assembler/assembler"
    "os"
)

func main() {
    err := assembler.Assemble("input.asm", "output.elf")
    if err != nil {
        panic(err)
    }
}
```

## Project Structure

- `assembler.go`: Core logic of the assembler
- `compiler.go`: Compiles instructions into machine code
- `elfbuilder.go`: Builds ELF output files
- `preprocesser.go`: Handles macros and directives
- `tokens.go`: Lexical analysis of instructions
- `encodings.go`: Encodes RISC-V instructions
- `defines.go`: Constants and shared structures
- `type_encoders.go`: Type-specific instruction encoders
- `*.test.go`: Unit tests for each component
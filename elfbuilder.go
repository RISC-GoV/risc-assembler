package main

import "encoding/binary"

func GenerateELFHeaders(e_entry [4]byte, e_phnum [2]byte) *[0x34]byte {
	var elfHeader [0x34]byte
	// Magic Number
	elfHeader[0x0] = 0x7F
	elfHeader[0x1] = 0x45  // E
	elfHeader[0x2] = 0x4c  // L
	elfHeader[0x3] = 0x46  // F
	elfHeader[0x4] = 0x01  // 32bit
	elfHeader[0x5] = 0x01  // LE
	elfHeader[0x6] = 0x01  // ELF Version
	elfHeader[0x7] = 0x03  // Linux ABI
	elfHeader[0x10] = 0x02 // Executable
	elfHeader[0x12] = 0xF3 // RISC-V
	elfHeader[0x14] = 0x01
	// EntryPoint Address
	for i := 0; i < 4; i++ {
		elfHeader[0x18+i] = e_entry[i]
	}
	elfHeader[0x1C] = 0x34 // Program Header Address
	elfHeader[0x28] = 0x34 // Len = 52 Bytes
	elfHeader[0x2A] = 0x20 // 32bits
	// Amount of Entries in Program Header
	elfHeader[0x2C] = e_phnum[0]
	elfHeader[0x2D] = e_phnum[1]
	elfHeader[0x2E] = 0x20 // Header entry size
	return &elfHeader
}

func GenerateSingleELFProgramHeader(htype byte, offset [4]byte, size [4]byte) *[0x20]byte {
	var programHeader [0x20]byte
	programHeader[0x00] = 0x01 // PT_LOAD
	for i := 0; i < 4; i++ {
		programHeader[i+0x04] = offset[i] // File offset
		programHeader[i+0x08] = offset[i] // Memory offset
	}
	for i := 0; i < 4; i++ {
		programHeader[i+0x10] = size[i] // Size
		programHeader[i+0x14] = size[i] // Size
	}
	programHeader[0x18] = htype // RWX
	return &programHeader
}

func BuildELFFile(program Program) *[]byte {
	var headerAmount uint16 = 1

	if program.constants != nil {
		headerAmount++
	}
	if program.variables != nil {
		headerAmount++
	}
	if program.strings != nil {
		headerAmount++
	}

	finalOffset := uint32(headerAmount*0x20) + 0x34
	offset := make([]byte, 4)
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(offset, finalOffset)
	binary.LittleEndian.PutUint32(size, uint32(len(program.machinecode)))

	file := GenerateSingleELFProgramHeader(0x05, *(*[4]byte)(offset), *(*[4]byte)(size))[:]
	finalOffset += uint32(len(program.machinecode))

	if program.variables != nil {
		binary.LittleEndian.PutUint32(offset, finalOffset)
		binary.LittleEndian.PutUint32(size, uint32(len(program.variables)))
		file = append(file, GenerateSingleELFProgramHeader(0x06, *(*[4]byte)(offset), *(*[4]byte)(size))[:]...)
		finalOffset += uint32(len(program.variables))
	}

	if program.constants != nil {
		binary.LittleEndian.PutUint32(offset, finalOffset)
		binary.LittleEndian.PutUint32(size, uint32(len(program.constants)))
		file = append(file, GenerateSingleELFProgramHeader(0x04, *(*[4]byte)(offset), *(*[4]byte)(size))[:]...)
		finalOffset += uint32(len(program.constants))
	}

	if program.strings != nil {
		binary.LittleEndian.PutUint32(offset, finalOffset)
		binary.LittleEndian.PutUint32(size, uint32(len(program.strings)))
		file = append(file, GenerateSingleELFProgramHeader(0x04, *(*[4]byte)(offset), *(*[4]byte)(size))[:]...)
	}

	hamt := make([]byte, 2)
	binary.LittleEndian.PutUint16(hamt, headerAmount)

	file = append(GenerateELFHeaders(program.entrypoint, *(*[2]byte)(hamt))[:], file...)
	file = append(file, program.machinecode...)
	file = append(file, program.variables...)
	file = append(file, program.constants...)
	file = append(file, program.strings...)
	return &file
}

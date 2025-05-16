# Comprehensive RISC-V Assembly Example
# Demonstrating all instruction types from the provided map

.section .text
.globl _start

_start:
    # ========== U-TYPE INSTRUCTIONS ==========
    # lui - Load Upper Immediate: Loads the immediate value into the upper 20 bits of the destination register
    lui x1, 0xABCDE     # Load 0xABCDE000 into x1
    lui x2, 0x12345     # Load 0x12345000 into x2

    # auipc - Add Upper Immediate to PC: Adds the immediate value to the PC and stores the result in the destination register
    auipc x3, 0x1000    # x3 = PC + 0x1000000
    auipc x4, 0x2       # x4 = PC + 0x2000

    # ========== J-TYPE INSTRUCTIONS ==========
    # jal - Jump and Link: Jumps to the target address and stores the address of the next instruction in the destination register
    jal x5, function1    # Jump to function1 and store return address in x5
    jal x0, skip1        # Jump to skip1 (equivalent to 'j skip1')

function1:
    # Function body
    addi x10, x0, 1      # x10 = 1
    jalr x0, 0(x5)      # Return to caller

skip1:
    # ========== B-TYPE INSTRUCTIONS ==========
    # Branch instructions
    addi x6, x0, 10     # x6 = 10
    addi x7, x0, 20     # x7 = 20

    # beq - Branch if Equal
    beq x6, x6, beq_target  # Always branches since x6 == x6
    addi x8, x0, 1      # Will be skipped
beq_target:

    # bne - Branch if Not Equal
    bne x6, x7, bne_target  # Branches since x6 != x7
    addi x8, x0, 2      # Will be skipped
bne_target:

    # blt - Branch if Less Than
    blt x6, x7, blt_target  # Branches since x6 < x7
    addi x8, x0, 3      # Will be skipped
blt_target:

    # bge - Branch if Greater or Equal
    bge x7, x6, bge_target  # Branches since x7 > x6
    addi x8, x0, 4      # Will be skipped
bge_target:

    # bltu - Branch if Less Than Unsigned
    addi x9, x0, -1     # x9 = -1 (0xFFFFFFFF in two's complement)
    bltu x6, x9, bltu_target  # Branches since x6 < x9 (unsigned: 10 < 0xFFFFFFFF)
    addi x8, x0, 5      # Will be skipped
bltu_target:

    # bgeu - Branch if Greater or Equal Unsigned
    bgeu x9, x7, bgeu_target  # Branches since x9 > x7 (unsigned: 0xFFFFFFFF > 20)
    addi x8, x0, 6      # Will be skipped
bgeu_target:

    # ========== I-TYPE INSTRUCTIONS ==========
    # jalr - Jump and Link Register
    lui x11, 0x10000    # x11 = 0x10000000
    jalr x12, 8(x11)    # Jump to address in x11 + 8, store return address in x12

    # Load instructions
    lui x13, 0x10000    # Base address for memory operations

    # lb - Load Byte
    lb x14, 0(x13)      # Load signed byte from memory
    lb x15, 4(x13)      # Load signed byte from memory with offset 4

    # lh - Load Halfword
    lh x16, 0(x13)      # Load signed halfword from memory
    lh x17, 2(x13)      # Load signed halfword from memory with offset 2

    # lw - Load Word
    lw x18, 0(x13)      # Load word from memory
    lw x19, 4(x13)      # Load word from memory with offset 4

    # lbu - Load Byte Unsigned
    lbu x20, 0(x13)     # Load unsigned byte from memory
    lbu x21, 4(x13)     # Load unsigned byte from memory with offset 4

    # lhu - Load Halfword Unsigned
    lhu x22, 0(x13)     # Load unsigned halfword from memory
    lhu x23, 2(x13)     # Load unsigned halfword from memory with offset 2

    # Immediate arithmetic instructions
    addi x24, x0, 100   # x24 = 0 + 100
    addi x25, x24, -50  # x25 = x24 - 50

    # slti - Set Less Than Immediate
    slti x26, x24, 200  # x26 = (x24 < 200) ? 1 : 0
    slti x27, x24, 50   # x27 = (x24 < 50) ? 1 : 0

    # sltiu - Set Less Than Immediate Unsigned
    sltiu x28, x24, 200 # x28 = (x24 < 200 unsigned) ? 1 : 0
    sltiu x29, x24, -1  # x29 = (x24 < 0xFFFFFFFF unsigned) ? 1 : 0

    # xori - XOR Immediate
    xori x30, x24, 0xFF # x30 = x24 ^ 0xFF
    xori x31, x0, -1    # x31 = 0 ^ 0xFFFFFFFF = 0xFFFFFFFF

    # ori - OR Immediate
    ori x1, x24, 0x0F   # x1 = x24 | 0x0F
    ori x2, x0, 0xFF    # x2 = 0 | 0xFF = 0xFF

    # andi - AND Immediate
    andi x3, x24, 0x0F  # x3 = x24 & 0x0F
    andi x4, x31, 0xFF  # x4 = 0xFFFFFFFF & 0xFF = 0xFF

    # Shift instructions
    addi x5, x0, 1      # x5 = 1

    # slli - Shift Left Logical Immediate
    slli x6, x5, 3      # x6 = x5 << 3 = 8
    slli x7, x5, 31     # x7 = x5 << 31 = 0x80000000

    # srli - Shift Right Logical Immediate
    srli x8, x7, 4      # x8 = x7 >> 4 (logical) = 0x08000000
    srli x9, x7, 31     # x9 = x7 >> 31 (logical) = 1

    # srai - Shift Right Arithmetic Immediate
    srai x10, x7, 4     # x10 = x7 >> 4 (arithmetic) = 0xF8000000
    srai x11, x7, 31    # x11 = x7 >> 31 (arithmetic) = 0xFFFFFFFF

    # Special I-type instructions
    ecall               # Environment call
    ebreak              # Environment breakpoint

    # ========== S-TYPE INSTRUCTIONS ==========
    # Store instructions
    lui x12, 0x10000    # Base address for memory operations
    addi x13, x0, 0x42  # Value to store

    # sb - Store Byte
    sb x13, 0(x12)      # Store lowest byte of x13 to memory
    sb x13, 4(x12)      # Store lowest byte of x13 to memory with offset 4

    # sh - Store Halfword
    sh x13, 0(x12)      # Store lowest halfword of x13 to memory
    sh x13, 2(x12)      # Store lowest halfword of x13 to memory with offset 2

    # sw - Store Word
    sw x13, 0(x12)      # Store x13 to memory
    sw x13, 4(x12)      # Store x13 to memory with offset 4

    # ========== R-TYPE INSTRUCTIONS ==========
    addi x14, x0, 100   # x14 = 100
    addi x15, x0, 50    # x15 = 50

    # add - Add
    add x16, x14, x15   # x16 = x14 + x15 = 150
    add x17, x14, x0    # x17 = x14 + 0 = x14 = 100

    # sub - Subtract
    sub x18, x14, x15   # x18 = x14 - x15 = 50
    sub x19, x15, x14   # x19 = x15 - x14 = -50

    # sll - Shift Left Logical
    addi x20, x0, 3     # x20 = 3
    sll x21, x14, x20   # x21 = x14 << x20 = 100 << 3 = 800

    # slt - Set Less Than
    slt x22, x15, x14   # x22 = (x15 < x14) ? 1 : 0 = 1
    slt x23, x14, x15   # x23 = (x14 < x15) ? 1 : 0 = 0

    # sltu - Set Less Than Unsigned
    addi x24, x0, -1    # x24 = -1 (0xFFFFFFFF)
    sltu x25, x14, x24  # x25 = (x14 < x24 unsigned) ? 1 : 0 = 1
    sltu x26, x24, x14  # x26 = (x24 < x14 unsigned) ? 1 : 0 = 0

    # xor - Bitwise XOR
    xor x27, x14, x15   # x27 = x14 ^ x15
    xor x28, x14, x14   # x28 = x14 ^ x14 = 0

    # srl - Shift Right Logical
    addi x20, x0, 2     # x20 = 2
    srl x29, x14, x20   # x29 = x14 >> x20 (logical) = 100 >> 2 = 25

    # sra - Shift Right Arithmetic
    addi x24, x0, -100  # x24 = -100
    sra x30, x24, x20   # x30 = x24 >> x20 (arithmetic) = -100 >> 2 = -25

    # or - Bitwise OR
    or x31, x14, x15    # x31 = x14 | x15
    or x1, x14, x0      # x1 = x14 | 0 = x14

    # and - Bitwise AND
    and x2, x14, x15    # x2 = x14 & x15
    and x3, x14, -1     # x3 = x14 & 0xFFFFFFFF = x14

    # nop - No Operation
    nop                 # No operation (addi x0, x0, 0)

    # ========== C-TYPE INSTRUCTIONS (Compressed) ==========
    # Note: Compressed instructions are 16-bit variants of regular instructions

    # add - Compressed Add
    add x4, x5        # x4 = x4 + x5

    # sub - Compressed Subtract
    sub x6, x7        # x6 = x6 - x7

    # nop - Compressed No Operation
    nop               # No operation

    # j - Compressed Jump
    j compressed_target  # Jump to compressed_target

    # jal - Compressed Jump and Link
    jal compressed_func  # Jump to compressed_func and save return address

compressed_func:
    # Function body
    add x10, x11
    ret                 # Return

compressed_target:
    # beq - Compressed Branch if Equal
    beq x8, x9, comp_beq_target

comp_beq_target:
    # bne - Compressed Branch if Not Equal
    bne x10, x11, comp_bne_target

comp_bne_target:
    # lw - Compressed Load Word
    lw x12, 4(x13)    # Load word from x13+4 into x12

    # sw - Compressed Store Word
    sw x14, 4(x15)    # Store word from x14 to x15+4

    # ========== DATA SECTION REFERENCE EXAMPLES ==========
    # Load from various data types to demonstrate they work
    la x10, word_data      # Load address of word_data
    lw x11, 0(x10)         # Load first word (0x12345678)

    la x12, string_data    # Load address of string_data
    lb x13, 0(x12)         # Load first byte of string ('T')

    la x14, asciz_data     # Load address of asciz_data
    lb x15, 0(x14)         # Load first byte of string ('T')

    la x16, byte_data      # Load address of byte_data
    lb x17, 0(x16)         # Load first byte (0x12)
    lb x18, 6(x16)         # Load 7th byte (0xEF)

    la x19, hword_data     # Load address of hword_data
    lh x20, 0(x19)         # Load first halfword (0x1122)
    lh x21, 2(x19)         # Load second halfword (0x3344)

    la x22, dword_data     # Load address of dword_data
    lw x23, 0(x22)         # Load first 4 bytes of dword (lower part)
    lw x24, 4(x22)         # Load second 4 bytes of dword (upper part)

    la x25, aligned_data   # Load address of aligned_data
    lw x26, 0(x25)         # Load aligned word (0x01234567)

    # End program with exit code
end:
    addi a7, x0, 93         # flag to exit
    addi a2, x0, 0          # flag to confirm successful exit
    ecall

# Data section for memory operations
.section .data
# .word - 32-bit (4-byte) data items
word_data:
    .word 0x12345678, 0xABCDEF01, 0x87654321, 0x10203040

# .string - Null-terminated string
string_data:
    .string "This is a null-terminated string"

.LC0:
    .byte 0x12, 0x34, 0x56, 0x78, 0xAB, 0xCD, 0xEF

# .asciz - Same as .string, also null-terminated
asciz_data:
    .asciz "This is another null-terminated string using .asciz"

# .byte - 8-bit (1-byte) data items
byte_data:
    .byte 0x12, 0x34, 0x56, 0x78, 0xAB, 0xCD, 0xEF
    .byte 'H', 'e', 'l', 'l', 'o', 0   # Character values as bytes

# .hword - 16-bit (2-byte) data items (halfword)
hword_data:
    .hword 0x1122, 0x3344, 0x5566, 0x7788, 0x99AA, 0xBBCC

# .dword - 64-bit (8-byte) data items (double word)
dword_data:
    .dword 0x1122334455667788, 0x99AABBCCDDEEFF00
    .dword 0xFEDCBA9876543210, 0x0123456789ABCDEF

# Aligned data examples
.align 4  # Align to a 2^4 = 16-byte boundary
aligned_data:
    .word 0x01234567   # Word-aligned data
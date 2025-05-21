    .section .text
    .globl   _start

_start:
# --------------------
# U-TYPE Instructions
# --------------------
    lui      x1, 0xABCDE                                           # Load upper immediate into x1
    lui      x2, 0x12345                                           # Load upper immediate into x2
    auipc    x3, 0x1000                                            # Add upper immediate to PC and store in x3
    auipc    x4, 0x2                                               # Similar operation with smaller immediate

# --------------------
# J-TYPE Instructions
# --------------------
    jal      x5, function1                                         # Jump to function1 and save return address in x5
    nop                                                            # Filler instruction

after_function1:
    jal      x0, skip1                                             # Unconditional jump to skip1 (x0 = discard return)

function1:
    addi     x10, x0, 1                                            # Set x10 = 1
    jalr     x0, 0(x5)                                             # Return to caller using x5 (JALR)

skip1:
# --------------------
# B-TYPE Instructions (Conditional branches)
# --------------------
    addi     x6, x0, 10                                            # x6 = 10
    addi     x7, x0, 20                                            # x7 = 20

    beq      x6, x6, beq_target                                    # x6 == x6 → branch taken
    nop
beq_target:

    bne      x6, x7, bne_target                                    # x6 != x7 → branch taken
    nop
bne_target:

    blt      x6, x7, blt_target                                    # x6 < x7 → branch taken
    nop
blt_target:

    bge      x7, x6, bge_target                                    # x7 >= x6 → branch taken
    nop
bge_target:

    addi     x9, x0, -1                                            # x9 = -1 (0xFFFFFFFF)
    bltu     x6, x9, bltu_target                                   # Unsigned comparison → 10 < 0xFFFFFFFF
    nop
bltu_target:

    bgeu     x9, x7, bgeu_target                                   # Unsigned → 0xFFFFFFFF >= 20 → branch taken
    nop
bgeu_target:

# --------------------
# I-TYPE Instructions (Immediate and Loads)
# --------------------

    lui      x11, 0x1                                              # x11 = 4096
    addi     x11, x11, -2000  
    addi     x11, x11, -1996  
    jalr     x12, 12(x11)                                           # x11 + 8 = 108

after_jalr:
    addi      x13, x0, %lo(word_data)                                            # x13 = address for memory access
    lb       x14, 0(x13)                                           # Load byte (signed)
    lb       x15, 4(x13)                                           # Load another byte
    lh       x16, 0(x13)                                           # Load halfword (signed)
    lh       x17, 2(x13)                                           # Load next halfword
    lw       x18, 0(x13)                                           # Load word (32-bit)
    lw       x19, 4(x13)                                           # Load next word
    lbu      x20, 0(x13)                                           # Load byte (unsigned)
    lbu      x21, 4(x13)
    lhu      x22, 0(x13)                                           # Load halfword (unsigned)
    lhu      x23, 2(x13)

    addi     x24, x0, 100                                          # x24 = 100
    addi     x25, x24, -50                                         # x25 = 50
    slti     x26, x24, 200                                         # x26 = 1 (100 < 200)
    slti     x27, x24, 50                                          # x27 = 0
    sltiu    x28, x24, 200                                         # Unsigned less than → 1
    sltiu    x29, x24, -1                                          # 100 < 0xFFFFFFFF → 1
    xori     x30, x24, 0xFF                                        # Bitwise XOR
    xori     x31, x0, -1                                           # x31 = 0xFFFFFFFF
    ori      x1, x24, 0x0F                                         # Bitwise OR
    ori      x2, x0, 0xFF                                          # x2 = 0xFF
    andi     x3, x24, 0x0F                                         # Bitwise AND
    andi     x4, x31, 0xFF                                         # x4 = 0xFF

    addi     x5, x0, 1
    slli     x6, x5, 3                                             # Shift left logical by 3 → x6 = 8
    slli     x7, x5, 31                                            # Large shift
    srli     x8, x7, 4                                             # Logical right shift
    srli     x9, x7, 31                                            # Shift to LSB
    srai     x10, x7, 4                                            # Arithmetic shift
    srai     x11, x7, 31                                           # Sign extend

# --------------------
# S-TYPE Instructions (Stores)
# --------------------
    addi     x12, x0, %lo(word_data)                                         # x12 = base address "Hello"
    addi     x13, x0, %lo(byte_data)                                      # x13 = data to store "using"

    sb       x13, 0(x12)                                           # Store byte
    sb       x13, 4(x12)
    sh       x13, 0(x12)                                           # Store halfword
    sh       x13, 2(x12)
    sw       x13, 0(x12)                                           # Store word
    sw       x13, 4(x12)

# --------------------
# R-TYPE Instructions (Register-to-register ops)
# --------------------
    addi     x14, x0, 100
    addi     x15, x0, 50

    add      x16, x14, x15                                         # 100 + 50 = 150
    add      x17, x14, x0                                          # x17 = x14

    sub      x18, x14, x15                                         # 100 - 50 = 50
    sub      x19, x15, x14                                         # 50 - 100 = -50

    addi     x20, x0, 3
    sll      x21, x14, x20                                         # Shift left by 3
    slt      x22, x15, x14                                         # x15 < x14 → 1
    slt      x23, x14, x15                                         # x14 < x15 → 0

    addi     x24, x0, -1
    sltu     x25, x14, x24                                         # Unsigned compare
    sltu     x26, x24, x14                                         # Unsigned compare

    xor      x27, x14, x15                                         # Bitwise XOR
    xor      x28, x14, x14                                         # Result = 0

    addi     x20, x0, 2
    srl      x29, x14, x20                                         # Logical right
    addi     x24, x0, -100
    sra      x30, x24, x20                                         # Arithmetic right

    or       x31, x14, x15                                         # Bitwise OR
    or       x1, x14, x0                                           # OR with zero (copy)
    and      x2, x14, x15                                          # Bitwise AND
    and      x3, x14, -1                                           # AND with all 1s (copy)
    nop

# --------------------
# Simulated Compressed Instructions
# --------------------
    add      x4, x4, x5                                            # Simulated C.ADD
    sub      x6, x6, x7                                            # Simulated C.SUB
    nop

    jal      x5, compressed_target                                 # Jump to compressed_target

compressed_func:
    add      x10, x10, x11                                         # Simulated function body
    ret                                                            # Return

compressed_target:
    beq      x8, x9, comp_beq_target
    nop
comp_beq_target:
    bne      x10, x11, comp_bne_target
    nop
comp_bne_target:
    lw       x12, 4(x13)                                           # Load word
    sw       x14, 4(x15)                                           # Store word

# --------------------
# Data Section Access
# --------------------
    la       x10, word_data
    lw       x11, 0(x10)                                           # Load from word array

    la       x12, string_data
    lb       x13, 0(x12)                                           # Load character

    la       x14, asciz_data
    lb       x15, 0(x14)

    la       x16, byte_data
    lb       x17, 0(x16)
    lb       x18, 6(x16)

    la       x19, hword_data
    lh       x20, 0(x19)
    lh       x21, 2(x19)

    la       x22, dword_data
    lw       x23, 0(x22)
    lw       x24, 4(x22)

    la       x25, aligned_data
    lw       x26, 0(x25)

# --------------------
# Program Exit
# --------------------
end:
    li       a7, 93                                                # syscall: exit
    li       a0, 0                                                 # exit code 0
    ecall                                                          # make syscall

# --------------------
# Data Section
# --------------------
    .section .data

word_data:
    .word    0x12345678, 0xABCDEF01, 0x87654321, 0x10203040

string_data:
    .string  "This is a null-terminated string"

asciz_data:
    .asciz   "This is another null-terminated string using .asciz"

byte_data:
    .byte    0x12, 0x34, 0x56, 0x78, 0xAB, 0xCD, 0xEF
    .byte    'H', 'e', 'l', 'l', 'o', 0

hword_data:
    .hword   0x1122, 0x3344, 0x5566, 0x7788, 0x99AA, 0xBBCC

dword_data:
    .dword   0x1122334455667788, 0x99AABBCCDDEEFF00
    .dword   0xFEDCBA9876543210, 0x0123456789ABCDEF

    .align   4
aligned_data:
    .word    0x01234567

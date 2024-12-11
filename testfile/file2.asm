.data
.LC0:
    .string "Hello, World!"

.text
.globl _start
_start:
    # Load the address of the string into a register
    la a0, .LC0

    # Print the string
    li a7, 4            # System call number for printing string
    li a1, 1            # File descriptor (stdout)
    li a2, 13           # Length of the string
    ecall

    # Exit the program
    li a7, 10           # Exit system call code
    ecall
.data
num1: .word 5
num2: .word 7
result: .word 0

.text
.globl _start
_start:
    # Load the first number into register t0
    lw t0, num1
    
    # Load the second number into register t1
    lw t1, num2

    add t2, t0, t1

    sw t2, result
    
    # Exit the program
    li a7, 10         # Exit system call code
    ecall
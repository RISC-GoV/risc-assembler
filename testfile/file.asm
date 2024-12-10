.data
    num1: .word 10        # First number
    num2: .word 20        # Second number
    result: .word 0       # To store the result

.LC0:
    vr: .word 10

.LC1:
    lw   t0, num1        # Load num1 into t0
    lw   t1, num2        # Load num2 into t1

.text
    .globl _start

_start:
    # Load numbers into registers
    lw   t0, num1        # Load num1 into t0
    lw   t1, num2        # Load num2 into t1

    # Perform addition
    add  t2, t0, t1      # t2 = num1 + num2

    # Perform subtraction
    sub  t3, t1, t0      # t3 = num2 - num1

    # Store results
    sw   t2, result      # Store addition result
    sw   t3, result + 4  # Store subtraction result

    # Conditional operations
    slt  t4, t0, t1      # t4 = (num1 < num2) ? 1 : 0
    beq  t4, x0, skip    # If t4 == 0, skip

    # Perform XOR operation
    xor  t5, t0, t1      # t5 = num1 ^ num2

    # Store XOR result
    sw   t5, result + 8  # Store XOR result

skip:
    # Logical operations
    and  t6, t0, t1      # t6 = num1 & num2
    or   t7, t0, t1      # t7 = num1 | num2

    # Store logical results
    sw   t6, result + 12 # Store AND result
    sw   t7, result + 16 # Store OR result

    # Shift operations
    slli t8, t0, 2       # t8 = num1 << 2
    srli t9, t1, 1       # t9 = num2 >> 1

    # Store shift results
    sw   t8, result + 20 # Store shift left result
    sw   t9, result + 24 # Store shift right result

    # End of program (for simulation purposes)
    ebreak                # Trigger debugger break (or replace with ecall for OS)
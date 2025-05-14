.section .data
dt:     .word 50       # Example data
de:     .word 0        # Space for stored value

.LC0:
    .word 420          # Read-only data

.section .text
.globl _start

    # ===========================
    # Loop to increment x1 until it equals x2
    # ===========================
loop_start:
    beq x1, x2, loop_end  # If x1 == x2, exit loop
    addi x2, x2, 1        # Increment x1
    jal x0, loop_start    # Jump back to start of loop

loop_end:
    # x1 is now equal to x2
    
    # Store the final value
    sw x1, 4(x11)         # Store the final value of x1 into de

    # Option 1: Infinite loop once x1 == x2
infinite_loop:
    jal x0, infinite_loop # Infinite loop - program stays here forever

    # Option 2: Exit program (uncomment these lines if you want to exit instead)
    # li a7, 10           # Exit syscall number (10)
    # ecall               # Make syscall to exit

# ===========================
# Branch targets (not used in main logic)
# ===========================
branch_equal:
    # Code for when x1 == x2
    jal x0, continue

branch_notequal:
    # Code for when x1 != x2
    jal x0, continue

branch_lt:
    # Code for when x1 < x2
    jal x0, continue

branch_ge:
    # Code for when x2 >= x1
    jal x0, continue

continue:
    # Continue execution
    jal x0, loop_start

# ===========================
# Main Program (_start)
# ===========================
_start:
    # Register initialization
    li x1, 42             # x1 = 42 (also acts as ra later)
    li x2, 100            # x2 = 100

    # Arithmetic operations
    add x3, x1, x2        # x3 = x1 + x2
    sub x4, x2, x1        # x4 = x2 - x1

    # Logical operations
    and x5, x1, x2        # x5 = x1 & x2
    or  x6, x1, x2        # x6 = x1 | x2
    xor x7, x1, x2        # x7 = x1 ^ x2

    # Shift operations
    slli x8, x1, 2        # x8 = x1 << 2
    srli x9, x2, 1        # x9 = x2 >> 1 (logical)
    srai x10, x2, 1       # x10 = x2 >> 1 (arithmetic)

    # Memory operations (assumes x11 holds address)
    la x11, dt            # Load address of dt into x11
    lw x12, 0(x11)        # Load word from dt into x12
    sw x1, 4(x11)         # Store x1 into de (4 bytes after dt)

    # Branches (for demonstration, not used in loop)
    beq x1, x2, branch_equal
    bne x1, x2, branch_notequal
    blt x1, x2, branch_lt
    bge x2, x1, branch_ge

    # Set up loop
    mv x14, x1            # x14 = x1 (starting point)

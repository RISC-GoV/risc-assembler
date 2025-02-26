.data
hello_msg:
.asciz "Hello, RISC-V!\n"
prompt_msg:
   .asciz "Enter a number: "
result_msg:
   .asciz "Result: "

    .section .bss
    .lcomm num, 4  # Reserve space for input integer

    .text
    .globl _start

_start:
    # Load address of hello_msg into a0 (without la)
    lui a0, %hi(hello_msg)  # Load upper 20 bits of hello_msg
    addi a0, a0, %lo(hello_msg)  # Add lower 12 bits

    # Print "Hello, RISC-V!"
    li a7, 4          # ecall code for print string
    ecall             # System call

    # Debug breakpoint
    ebreak            # Trigger debugger breakpoint

    # Load address of prompt_msg into a0 (without la)
    lui a0, %hi(prompt_msg)
    addi a0, a0, %lo(prompt_msg)

    # Print "Enter a number: "
    li a7, 4
    ecall

    # Read integer from user
    li a7, 5          # ecall code for reading integer
    ecall
    sw a0, num        # Store the input integer

    # Load input and add 10
    lw t0, num        # Load input number into t0
    addi t0, t0, 10   # Add 10
    sw t0, num        # Store result

    # Load address of result_msg into a0 (without la)
    lui a0, %hi(result_msg)
    addi a0, a0, %lo(result_msg)

    # Print "Result: "
    li a7, 4
    ecall

    # Print result integer
    li a7, 1          # ecall code for print integer
    lw a0, num        # Load result
    ecall

    # Exit program
    li a7, 10         # ecall code for exit
    ecall
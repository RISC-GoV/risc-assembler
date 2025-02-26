.data
    dt: .word 50  # Example data
    de: .word 0           # Space for stored value


.LC0:
    .word 420  # Read-only data

.LC1:
    .string "Hello world"
  
function:
    # Example function
    addi x15, x1, 5   # Add immediate
    
    # Comparisons
    slt x16, x1, x2   # Set if less than
    slti x17, x1, 50  # Set if less than immediate
    
    
    # Environment call
    ecall             # System call
    
    # More complex operations
    auipc x19, 4      # Add upper immediate to PC 
skip:
    # Some operations to skip
    add x13, x1, x2
    
continue:
    # Control transfer
    jal x14, function # Jump and link
    jalr x0, x14, 0   # Return
    



.section .text
.globl _start
_start:
    # Register initialization
    addi x0, 0          # x0 is hardwired to zero
    addi x1, 42         # x1 = 42
    addi x2, 100        # x2 = 100
    
    # Arithmetic operations
    add x3, x1, x2    # x3 = x1 + x2
    sub x4, x2, x1    # x4 = x2 - x1
    
    # Logical operations
    and x5, x1, x2    # Bitwise AND
    or x6, x1, x2     # Bitwise OR
    xor x7, x1, x2    # Bitwise XOR
    
    # Shifts
    slli x8, x1, 2    # Shift left logical immediate
    srli x9, x2, 1    # Shift right logical immediate
    srai x10, x2, 1   # Shift right arithmetic immediate
    
    # Load and store
    lw x12, 0(x11)    # Load word
    sw x1, 4(x11)     # Store word
    
    # Branches
    beq x1, x2, skip  # Branch if equal
    bne x1, x2, continue # Branch if not equal
    blt x1, x2, continue # Branch if less than
    bge x2, x1, continue # Branch if greater or equal
 
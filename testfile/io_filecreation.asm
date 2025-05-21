main:
    addi a7, x0, 56         # syscall OPENAT (filesystem)
    addi a0, x0, -100       # path relative to current directory
    la a1, .path       # path to use
    addi a2, x0, 256        # create file if not exist O_CREATE
    addi a3, x0, 511        # authorization level 0x777
    ecall
    addi a7, x0, 64         # syscall WRITE
    addi a0, x0, 0          # First file opened by program so id = 0
    la a1, .content    # content to write
    addi a2, x0, 4         # size of string in byte
    ecall
    addi a7, x0, 93         # flag to exit
    addi a2, x0, 0          # flag to confirm successful exit
    ecall
.data
.path:
    .asciz "./louis.txt"
.content:
    .asciz "0522"
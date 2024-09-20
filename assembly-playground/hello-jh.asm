; John Hammond Video

global _start

section .text:
_start:
    mov eax, 0x4            ; write syscall
    mov ebx, 1              ; use stdout as the fd
    mov ecx, message        ; use the message as the buffer
    mov edx, message_len   ; and supply the message length
    int 0x80

    mov eax, 0x1            ; exit syscall
    mov ebx, 0             ; error code (success)
    int 0x80


section .data:
    message: db "Hello World", 0xA
    message_len equ $-message
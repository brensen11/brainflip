	global main
	extern printf
    extern ExitProcess
	section .text

main: 
    lea rcx, [rel msg]
    sub rsp, 8+16          ; 8 for alignment, 16 - shadow space
    call printf
    add rsp, 8+16

    xor rcx, rcx            ; set rcx to 0 for Exit code 0
    sub rsp, 8+16           ; 8 for alignment, 16 - shadow space
    call ExitProcess

msg:
	db "Hello World", 10, 0
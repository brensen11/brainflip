package generator

import (
	lp "brainflip-go/lexparse"
	"brainflip-go/utils"
	"strconv"
	"strings"
)

// mask_1: 0xFFFFFFFF
// mask_2: ... (101010101010...)
// mask_4:
// mask_8:

var asm_win64_template = `bits 64
default rel
segment .text

extern ExitProcess
extern my_putchar
extern my_getchar
extern calloc

; Index Vector masks
zeroes db 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
vmovdqu  ymm0, [zeroes] ; load 0s
mask_1 dw 0xFFFF FFFF ; 11111111
mask_2 dw 0xAAAA AAAA ; 10101010
mask_3 dw 0x8888 8888 ; 10001000
mask_4 dw 0x8080 8080 ; 10000000

global main
main:
    push    rbp
    mov     rbp, rsp
    sub     rsp, 32
    mov     rcx, 1024 * 1024 * 4
    mov     rdx, 1
    call    calloc
    mov     rdi, rax
	add     rdi, 1024 * 1024 * 2
	xor     rcx, rcx

{MAIN_CODE}

    xor     rcx, rcx
    call    ExitProcess`

func Generate(instructions *[]lp.Instruction) string {
	var asm_b utils.Builderf
	const TAPE_PTR string = "rdi"

	bracket_pairs := lp.Locate_Brackets(*instructions)

	// main run function
	for i := 0; i < len(*instructions); i++ {
		instruction := (*instructions)[i]

		switch i_t := instruction.(type) {
		case lp.Move_right:
			asm_b.Add_instr("inc     %s", TAPE_PTR)
		case lp.Move_left:
			asm_b.Add_instr("dec     %s", TAPE_PTR)
		case lp.Inc:
			asm_b.Add_instr("inc     BYTE [%s]", TAPE_PTR)
		case lp.Dec:
			asm_b.Add_instr("dec     BYTE [%s]", TAPE_PTR)
		case lp.Output:
			asm_b.Add_instr("mov     cl, BYTE [%s]", TAPE_PTR)
			asm_b.Add_instr("call    my_putchar")
		case lp.Input:
			asm_b.Add_instr("call    my_getchar")
			asm_b.Add_instr("mov     BYTE [%s], al", TAPE_PTR)
		case lp.Left_loop:
			asm_b.Add_instr("cmp     BYTE [%s], 0", TAPE_PTR)
			asm_b.Add_instr("je      right_%s", strconv.Itoa(bracket_pairs[i]))
			asm_b.Add_label("left_%s", strconv.Itoa(i))
		case lp.Right_loop:
			asm_b.Add_instr("cmp     BYTE [%s], 0", TAPE_PTR)
			asm_b.Add_instr("jne     left_%s", strconv.Itoa(bracket_pairs[i]))
			asm_b.Add_label("right_%s", strconv.Itoa(i))
		case lp.Add:
			asm_b.Add_instr("add     %s, %s", i_t.Op1.ToAsm(), i_t.Op2.ToAsm())
		case lp.Sub:
			asm_b.Add_instr("sub     %s, %s", i_t.Op1.ToAsm(), i_t.Op2.ToAsm())
		case lp.Mul:
			asm_b.Add_instr("imul    %s, %s", i_t.Op1.ToAsm(), i_t.Op2.ToAsm())
		case lp.Store:
			asm_b.Add_instr("mov     %s, %s", i_t.Op1.ToAsm(), i_t.Op2.ToAsm())
		case lp.Raw:
			asm_b.Add_instr(i_t.Raw)
		default:
			// do nothing
		}
	}
	return strings.Replace(asm_win64_template, "{MAIN_CODE}", asm_b.String(), 1)
}

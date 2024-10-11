package generator

import (
	lp "brainflip-go/lexparse"
	"brainflip-go/utils"
	"strconv"
	"strings"
)

var asm_win64_template = `bits 64
default rel
segment .text

extern ExitProcess
extern my_putchar
extern my_getchar
extern calloc

; Index Vector masks
zeroes db 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
ones dq 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF
indices_1 db 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31;
indices_2 db 0, -1, 2, -1, 4, -1, 6, -1, 8, -1, 10, -1, 12, -1, 14, -1;
indices_4 db 0, -1, -1, -1, 4, -1, -1, -1, 8, -1, -1, -1, 12, -1, -1, -1;
indices_8 db 0, -1, -1, -1, -1, -1, -1, -1, 8, -1, -1, -1, -1, -1, -1, -1;

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

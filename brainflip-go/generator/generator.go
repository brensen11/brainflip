package generator

import (
	lp "brainflip-go/lexparse"
	"fmt"
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
zeroes db 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
ones dq 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF
indices_1 db 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15;
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

type builder struct {
	strings.Builder
}

func (asm_b *builder) add_instr(instr string, args ...any) {
	asm_b.WriteString(fmt.Sprintf("\t"+instr+"\n", args...))
}

func (asm_b *builder) add_label(instr string, args ...any) {
	asm_b.WriteString(fmt.Sprintf(instr+":\n", args...))
}

func Generate(instructions *[]lp.Instruction) string {
	var asm_b builder
	const TAPE_PTR string = "rdi"

	bracket_pairs := lp.Locate_Brackets(*instructions)

	// main run function
	for i := 0; i < len(*instructions); i++ {
		instruction := (*instructions)[i]

		switch i_t := instruction.(type) {
		case lp.Move_right:
			asm_b.add_instr("inc     %s", TAPE_PTR)
		case lp.Move_left:
			asm_b.add_instr("dec     %s", TAPE_PTR)
		case lp.Inc:
			asm_b.add_instr("inc     BYTE [%s]", TAPE_PTR)
		case lp.Dec:
			asm_b.add_instr("dec     BYTE [%s]", TAPE_PTR)
		case lp.Output:
			asm_b.add_instr("mov     cl, BYTE [%s]", TAPE_PTR)
			asm_b.add_instr("call    my_putchar")
		case lp.Input:
			asm_b.add_instr("call    my_getchar")
			asm_b.add_instr("mov     BYTE [%s], al", TAPE_PTR)
		case lp.Left_loop:
			asm_b.add_instr("cmp     BYTE [%s], 0", TAPE_PTR)
			asm_b.add_instr("je      right_%s", strconv.Itoa(bracket_pairs[i]))
			asm_b.add_label("left_%s", strconv.Itoa(i))
		case lp.Right_loop:
			asm_b.add_instr("cmp     BYTE [%s], 0", TAPE_PTR)
			asm_b.add_instr("jne     left_%s", strconv.Itoa(bracket_pairs[i]))
			asm_b.add_label("right_%s", strconv.Itoa(i))
		case lp.Add:
			asm_b.add_instr("add     %s, %s", i_t.Op1.ToAsm(), i_t.Op2.ToAsm())
		case lp.Sub:
			asm_b.add_instr("sub     %s, %s", i_t.Op1.ToAsm(), i_t.Op2.ToAsm())
		case lp.Mul:
			asm_b.add_instr("imul    %s, %s", i_t.Op1.ToAsm(), i_t.Op2.ToAsm())
		case lp.Store:
			asm_b.add_instr("mov     %s, %s", i_t.Op1.ToAsm(), i_t.Op2.ToAsm())
		case lp.Raw:
			asm_b.add_instr(i_t.Raw)
		default:
			// do nothing
		}
	}
	return strings.Replace(asm_win64_template, "{MAIN_CODE}", asm_b.String(), 1)
}

package generator

import (
	"brainflip-go/lexparse"
	"fmt"
	"strconv"
	"strings"
)

var asm_win64_template = `bits 64
default rel
segment .text

extern ExitProcess
extern putchar
extern getchar
extern calloc

global main
main:
    push    rbp
    mov     rbp, rsp
    sub     rsp, 32
    mov     rcx, 1024 * 1024 * 4
    mov     rdx, 1
    call    calloc
    mov     rdi, rax

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

func Generate(program *lexparse.Program) string {
	var asm_b builder
	const TAPE_PTR string = "rdi"

	// main run function
	for i := 0; i < len(program.Instructions); i++ {
		instruction := program.Instructions[i]

		switch instruction.(type) {
		case lexparse.Move_right:
			asm_b.add_instr("inc     %s", TAPE_PTR)
		case lexparse.Move_left:
			asm_b.add_instr("dec     %s", TAPE_PTR)
		case lexparse.Inc:
			asm_b.add_instr("inc     BYTE [%s]", TAPE_PTR)
		case lexparse.Dec:
			asm_b.add_instr("dec     BYTE [%s]", TAPE_PTR)
		case lexparse.Output:
			asm_b.add_instr("mov     rcx, [%s]", TAPE_PTR)
			asm_b.add_instr("call    putchar")
		case lexparse.Input:
			asm_b.add_instr("call    getchar")
			asm_b.add_instr("mov     [%s], rax", TAPE_PTR)
		case lexparse.Left_loop:
			asm_b.add_instr("cmp     BYTE [%s], 0", TAPE_PTR)
			asm_b.add_instr("je      right_%s", strconv.Itoa(program.BracketPairs[i]))
			asm_b.add_label("left_%s", strconv.Itoa(i))
		case lexparse.Right_loop:
			asm_b.add_instr("cmp     BYTE [%s], 0", TAPE_PTR)
			asm_b.add_instr("jne      left_%s", strconv.Itoa(program.BracketPairs[i]))
			asm_b.add_label("right_%s", strconv.Itoa(i))
		default:
			// do nothing
		}
	}
	return strings.Replace(asm_win64_template, "{MAIN_CODE}", asm_b.String(), 1)
}

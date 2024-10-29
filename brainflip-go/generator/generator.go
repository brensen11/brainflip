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
extern my_print
extern calloc

; Index Vector masks
zeroes db 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
movdqu  xmm0, [zeroes] ; load 0s
mask_1 dd 0xFFFFFFFF ; 11111111
mask_2 dd 0xAAAAAAAA ; 10101010
mask_3 dd 0x88888888 ; 10001000
mask_4 dd 0x80808080 ; 10000000
{OUTPUT}

global main
main:
    push    rbp
    mov     rbp, rsp
    sub     rsp, 32
{SETUP_TAPE_CODE}
	add     rdi, 1024 * 1024 * 2
	xor     rcx, rcx

{MAIN_CODE}

    xor     rcx, rcx
    call    ExitProcess`

func setup(TAPE *[]byte, POINTER int) string {
	var asm_b utils.Builderf
	if TAPE != nil && POINTER > -1 {
		asm_b.Add_instr("mov     rcx, %d", len(*TAPE))
		asm_b.Add_instr("mov     rdx, 1")
		asm_b.Add_instr("call    calloc")
		asm_b.Add_instr("mov     rdi, rax")
		asm_b.Add_instr("add     rdi, %d", POINTER)
		asm_b.Add_instr("lea     rcx, [my_string]")
		asm_b.Add_instr("call    my_print")
	} else {
		asm_b.Add_instr("mov     rcx, 1024 * 1024 * 4")
		asm_b.Add_instr("mov     rdx, 1")
		asm_b.Add_instr("call    calloc")
		asm_b.Add_instr("mov     rdi, rax")
	}

	return asm_b.String()
}

func to_string_def_instr(output string) string {
	var string_def strings.Builder
	string_open := false
	for _, char := range output {
		switch char {
		case '\t':
			if string_open {
				string_def.WriteString("\", ")
				string_open = false
			}
			string_def.WriteString("0x09, ")
		case '\n':
			if string_open {
				string_def.WriteString("\", ")
				string_open = false
			}
			string_def.WriteString("0x0A, ")
		case '\r':
			break
		default:
			if !string_open {
				string_def.WriteRune('"')
				string_open = true
			}
			string_def.WriteRune(char)
		}
	}
	if string_open {
		string_def.WriteString("\", ")
	}
	return "my_string db " + string_def.String() + "0"
}

func Generate(instructions *[]lp.Instruction, TAPE *[]byte, POINTER int, output string) string {
	var asm_b utils.Builderf
	const TAPE_PTR string = "rdi"
	var assembly_string string = asm_win64_template

	assembly_string = strings.Replace(assembly_string, "{OUTPUT}", to_string_def_instr(output), 1)
	assembly_string = strings.Replace(assembly_string, "{SETUP_TAPE_CODE}", setup(TAPE, POINTER), 1)

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
	return strings.Replace(assembly_string, "{MAIN_CODE}", asm_b.String(), 1)
}

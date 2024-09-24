package main

import (
	"brainflip-go/utils"
	"fmt"
	"os"
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
    mov     rcx, 4096
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

func compile(program string) string {
	var asm_b builder
	const TAPE_PTR string = "rdi"

	// var POINTER int = 0

	bracketPairs := utils.Locate_Brackets(program)
	// main run function
	for i := 0; i < len(program); i++ {
		cmd := program[i]
		switch cmd {
		case '>':
			asm_b.add_instr("inc     %s", TAPE_PTR)
		case '<':
			asm_b.add_instr("dec     %s", TAPE_PTR)
		case '+':
			asm_b.add_instr("inc     BYTE [%s]", TAPE_PTR)
		case '-':
			asm_b.add_instr("dec     BYTE [%s]", TAPE_PTR)
		case '.':
			asm_b.add_instr("mov     rcx, [%s]", TAPE_PTR)
			asm_b.add_instr("call    putchar")
		case ',':
			asm_b.add_instr("call    getchar")
			asm_b.add_instr("mov     [%s], rax", TAPE_PTR)
		case '[':
			asm_b.add_instr("cmp     BYTE [%s], 0", TAPE_PTR)
			asm_b.add_instr("je      right_%s", strconv.Itoa(bracketPairs[i]))
			asm_b.add_label("left_%s", strconv.Itoa(i))
		case ']':
			asm_b.add_instr("cmp     BYTE [%s], 0", TAPE_PTR)
			asm_b.add_instr("jne      left_%s", strconv.Itoa(bracketPairs[i]))
			asm_b.add_label("right_%s", strconv.Itoa(i))
		default:
			// do nothing
		}
	}
	return asm_b.String()
}

func compile_file(filename string) string {
	bf_data, prog_err := os.ReadFile(filename)
	if prog_err != nil {
		fmt.Println("Usage: ./brainflip <brainflip.b>")
		panic("There was an error reading: " + filename)
	}

	bf_prog := string(bf_data)
	bf_asm := compile(bf_prog)
	asm := strings.Replace(asm_win64_template, "{MAIN_CODE}", bf_asm, 1)
	return asm
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./brainflip <brainflip.b>")
		return
	}

	assembly := compile_file(os.Args[1])

	paths := strings.Split(os.Args[1], "/")
	bf_filename := paths[len(paths)-1]
	program_name := strings.Split(bf_filename, ".")[0]

	file, err := os.Create(fmt.Sprintf("%s-win.asm", program_name))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write([]byte(assembly))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

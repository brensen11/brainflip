package main

import (
	"brainflip-go/utils"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

// TODO put as like driver in parent package or something
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./brainflip <brainflip.b> ", os.Args)
	}

	prog, prog_err := os.ReadFile(os.Args[1])
	if prog_err != nil {
		fmt.Println("Usage: ./brainflip <brainflip.b> ", os.Args)
		return
	}

	program := string(prog)
	assembly := compile(program)

	tmpl, tmpl_err := os.ReadFile("compiler/win.tmpl")
	if tmpl_err != nil {
		fmt.Println("Could not find template assembly file")
		return
	}

	template := string(tmpl)
	assembly_program := strings.Replace(template, "{MAIN_CODE}", assembly, 1)

	paths := strings.Split(os.Args[1], "/")
	filename := paths[len(paths)-1]
	program_name := strings.Split(filename, ".")[0]

	fmt.Println("arg 1: ", os.Args[1])
	fmt.Println("prog name: ", program_name)

	file, err := os.Create(fmt.Sprintf("%s-win.asm", program_name))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write([]byte(assembly_program))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

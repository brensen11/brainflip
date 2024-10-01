package compiler

import (
	"brainflip-go/generator"
	"brainflip-go/lexparse"
	"brainflip-go/optimize"
	"brainflip-go/utils"
	"fmt"
)

func Compile(filename string, outfile string, loop_optimize bool) {
	// ---------- get program name ----------
	// .........

	// ---------- readfile ----------
	bf_prog := utils.Readfile(filename)

	// ---------- lexparse ----------
	program := lexparse.Lexparse(bf_prog)

	// ---------- optimize ----------
	if loop_optimize {
		optimize.Optimize_simple_loops(program)
		fmt.Print(lexparse.Instructions_string(program.Instructions))
		return
	}

	// ---------- generator ----------
	assembly := generator.Generate(program)

	// ---------- writefile ----------
	utils.Writefile(assembly, outfile) // , filename)
}

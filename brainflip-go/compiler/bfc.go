package compiler

import (
	"brainflip-go/generator"
	"brainflip-go/lexparse"
	"brainflip-go/optimize"
	"brainflip-go/utils"
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
		optimize.Optimize_scans(program)
		// fmt.Println("Optimizer On!!!!!")
	}

	// ---------- generator ----------
	assembly := generator.Generate(program)

	// ---------- writefile ----------
	utils.Writefile(assembly, outfile) // , filename)
}

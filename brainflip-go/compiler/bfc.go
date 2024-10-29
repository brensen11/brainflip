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
	instructions := lexparse.Lexparse(bf_prog)

	// ---------- optimize ----------
	var TAPE *[]byte = nil
	var POINTER int = -1
	var output string = ""
	if loop_optimize {
		TAPE, POINTER, output = optimize.Optimize_partialeval(instructions)
		optimize.Optimize_simple_loops(instructions)
		optimize.Optimize_scans(instructions)
	}

	// ---------- generator ----------
	assembly := generator.Generate(instructions, TAPE, POINTER, output)

	// ---------- writefile ----------
	utils.Writefile(assembly, outfile) // , filename)
}

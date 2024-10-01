package compiler

import (
	"brainflip-go/generator"
	"brainflip-go/lexparse"
	"brainflip-go/utils"
)

func Compile(filename string, outfile string) {
	// ---------- get program name ----------
	// .........

	// ---------- readfile ----------
	bf_prog := utils.Readfile(filename)

	// ---------- lexparse ----------
	program := lexparse.Lexparse(bf_prog)

	// ---------- optimize ----------

	// ---------- generator ----------
	assembly := generator.Generate(program)

	// ---------- writefile ----------
	utils.Writefile(assembly, outfile) // , filename)
}

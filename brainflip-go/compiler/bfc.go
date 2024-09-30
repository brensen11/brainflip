package compiler

import (
	"brainflip-go/generator"
	"brainflip-go/lexparse"
	"fmt"
	"os"
)

func Compile(filename string) {
	// ---------- get program name ----------

	// ---------- readfile ----------
	bf_data, prog_err := os.ReadFile(filename)
	if prog_err != nil {
		panic("There was an error reading: " + filename)
	}
	bf_prog := string(bf_data)

	// ---------- lexparse ----------
	program := lexparse.Lexparse(bf_prog)

	// ---------- optimize ----------

	// ---------- generator ----------
	assembly := generator.Generate(program)

	// ---------- writefile ----------
	file, err := os.Create("out-win.asm")
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

package main

import (
	"brainflip-go/compiler"
	"brainflip-go/interpreter"
	"flag"
)

// TODO fill in with driver code, (calls bfc, bfi, respectively)
func main() {
	// ---------- Handle flags ----------
	out := flag.String("o", "out-win.asm", "Name of output file")

	interpret := flag.Bool("i", false, "Interpret Mode On")
	profile := flag.Bool("p", false, "Profile Mode On")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		panic("Must provide .b file")
	}

	filename := args[0]

	// ---------- Interpret with profiler if -p ----------
	if *profile {
		interpreter.Interpret_profiler(filename)
		return
	}

	// ---------- Interpret if -i ----------
	if *interpret {
		interpreter.Interpret(filename)
		return
	}

	compiler.Compile(filename, *out)
}

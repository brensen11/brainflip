package main

import (
	"brainflip-go/compiler"
	"brainflip-go/interpreter"
	"flag"
	"fmt"
	"os"
	"strings"
)

// TODO fill in with driver code, (calls bfc, bfi, respectively)
func main() {
	// ---------- Handle input file ----------
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./brainflip <brainflip.b>")
		return
	}
	if !strings.HasSuffix(os.Args[1], ".b") {
		fmt.Println("Usage: ./brainflip <brainflip.b>")
		fmt.Println("first argument must be a brainflip file!")
		return
	}
	filename := os.Args[1]

	// ---------- Handle flags ----------
	var out string
	flag.StringVar(&out, "o", "out-win.asm", "Name of output file")

	var interpret bool
	flag.BoolVar(&interpret, "i", false, "Interpret Mode On")

	var profile bool
	flag.BoolVar(&profile, "p", false, "Profile Mode On")

	flag.Parse()

	// ---------- Interpret with profiler if -p ----------
	if profile {
		interpreter.Interpret_profiler(filename)
		return
	}

	// ---------- Interpret if -i ----------
	if interpret {
		interpreter.Interpret(filename)
		return
	}

	compiler.Compile(filename)
}

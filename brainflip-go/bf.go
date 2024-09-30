package main

import (
	"brainflip-go/compiler"
	"flag"
	"fmt"
	"os"
	"strings"
)

// TODO fill in with driver code, (calls bfc, bfi, respectively)
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./brainflip <brainflip.b>")
		return
	}
	if !strings.HasSuffix(os.Args[1], ".b") {
		fmt.Println("Usage: ./brainflip <brainflip.b>")
		fmt.Println("first argument must be a brainflip file!")
		return
	}

	var out string
	flag.StringVar(&out, "o", "out-win.asm", "Name of output file")

	var interpret bool
	flag.BoolVar(&interpret, "i", false, "Interpret Mode On")

	compiler.Compile(os.Args[1])
}

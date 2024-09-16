package main

import "fmt"

func run(program string) {
	const TAPE_SIZE = 1024 * 4
	var TAPE [TAPE_SIZE]byte
	var POINTER int
	var PC int

	m := make(map[int]int)

	for i, v := range program {
		if v == '[' {

		} else if v == ']' {

		}
	}
}

func main() {
	fmt.Println("hello world")
}

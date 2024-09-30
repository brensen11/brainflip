package main

import (
	"brainflip-go/lexparse"
	"fmt"
	"os"
	"sort"
	"strings"
)

func interpret(program string) {
	const TAPE_SIZE = 1024 * 4
	var TAPE [TAPE_SIZE]byte
	var POINTER int = 0

	bracketPairs := lexparse.Locate_Brackets(program)
	// main run function
	for PC := 0; PC < len(program); PC++ {
		cmd := program[PC]
		// fmt.Println("Checking for: ", cmd)
		switch cmd {
		case '>':
			if POINTER == TAPE_SIZE {
				panic("Tape Pointer Out of Bounds!!")
			}
			POINTER++
		case '<':
			if POINTER == 0 {
				panic("Tape Pointer Out of Bounds!!")
			}
			POINTER--
		case '+':
			TAPE[POINTER]++
		case '-':
			TAPE[POINTER]--
		case '.':
			fmt.Printf("%c", TAPE[POINTER])
		case ',':
			// do nothing
		case '[':
			if TAPE[POINTER] == 0 {
				PC = bracketPairs[PC]
				continue
			}
		case ']':
			if TAPE[POINTER] != 0 {
				PC = bracketPairs[PC]
				continue
			}
		default:
			// do nothing
		}
	}
}

func interpret_profiler(program string) {
	var cmd_count [8]int

	const TAPE_SIZE = 1024 * 4
	var TAPE [TAPE_SIZE]byte
	var POINTER int = 0
	bracketPairs := lexparse.Locate_Brackets(program)
	rightBrackCount := make(map[int]int)

	// main run function
	for PC := 0; PC < len(program); PC++ {
		cmd := program[PC]
		// fmt.Println("Checking for: ", cmd)
		switch cmd {
		case '>':
			if POINTER == TAPE_SIZE {
				panic("Tape Pointer Out of Bounds!!")
			}
			POINTER++
			cmd_count[0]++
		case '<':
			if POINTER == 0 {
				panic("Tape Pointer Out of Bounds!!")
			}
			POINTER--
			cmd_count[1]++
		case '+':
			TAPE[POINTER]++
			cmd_count[2]++
		case '-':
			TAPE[POINTER]--
			cmd_count[3]++
		case '.':
			fmt.Printf("%c", TAPE[POINTER])
			cmd_count[4]++
		case ',':
			cmd_count[5]++
			// do nothing
		case '[':
			cmd_count[6]++
			if TAPE[POINTER] == 0 {
				PC = bracketPairs[PC]
				continue
			}
		case ']':
			rightBrackCount[PC]++
			cmd_count[7]++
			if TAPE[POINTER] != 0 {
				PC = bracketPairs[PC]
				continue
			}
		default:
			// do nothing
		}
	}

	println()
	fmt.Println("Instruction proc count:")
	fmt.Println(">: ", cmd_count[0])
	fmt.Println("<: ", cmd_count[1])
	fmt.Println("+: ", cmd_count[2])
	fmt.Println("-: ", cmd_count[3])
	fmt.Println(".: ", cmd_count[4])
	fmt.Println(",: ", cmd_count[5])
	fmt.Println("[: ", cmd_count[6])
	fmt.Println("]: ", cmd_count[7])

	sort_func := func(slice []int) func(int, int) bool {
		return func(i, j int) bool {
			return rightBrackCount[slice[i]] > rightBrackCount[slice[j]]
		}
	}

	simple, complex := lexparse.Categorize_Brackets(program)
	sort.Slice(simple, sort_func(simple))
	sort.Slice(complex, sort_func(complex))

	format_stmt := func(idx int) string {
		loop_stmt := program[bracketPairs[idx] : idx+1]
		loop_stmt = strings.Replace(loop_stmt, "\r", "", -1)
		loop_stmt = strings.Replace(loop_stmt, "\n", "", -1)
		return fmt.Sprintf("Loop: %s at [%d:%d] occured %d times", loop_stmt, bracketPairs[idx], idx+1, rightBrackCount[idx])
	}

	println()
	fmt.Println("Simple Innermost Loops")
	for _, idx := range simple {
		fmt.Println(format_stmt(idx))
	}

	println()
	fmt.Println("\nComplex Innermost Loops")
	for _, idx := range complex {
		fmt.Println(format_stmt(idx))
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./brainflip <brainflip.b> ", os.Args)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Usage: ./brainflip <brainflip.b> ", os.Args)
	}
	program := string(data)

	if len(os.Args) > 2 && os.Args[2] == "-p" {
		interpret_profiler(program)
	} else {
		interpret(program)
	}

}

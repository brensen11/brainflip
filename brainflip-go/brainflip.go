package main

import (
	"fmt"
	"os"
	"sort"
)

type stack []int

func (s *stack) Push(v int) {
	*s = append(*s, v)
}

func (s *stack) Pop() int {
	l := len(*s)
	if l == 0 {
		panic("Tried to pop empty stack")
	}
	val := (*s)[l-1]
	*s = (*s)[:l-1]
	return val
}

func interpret(program string) {
	const TAPE_SIZE = 1024 * 4
	var TAPE [TAPE_SIZE]byte
	var POINTER int = 0

	bracketPairs := locate_brackets(program)
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

func categorize_brackets(program string, bracketPairs map[int]int) ([]int, []int) {
	const (
		SIMPLE int = iota
		COMPLEX
		CLOSED
	)
	simples := make([]int, 0, len(program)/2)
	complexes := make([]int, 0, len(program)/2)

	// a simple loop is one that...
	// - contains no i/o,
	// - has 0 net change to the pointer
	// - and that changes p[0] by either +1 or -1 after each loop iteration, aka by the end of the loop [0] has changed by 1
	// (the cell that the pointer references when the loop body starts executing)
	ptr_rel_loc := 0
	p0_changes := 0
	state := CLOSED
	for i, v := range program {
		switch v {
		case '[':
			state = SIMPLE
			ptr_rel_loc = 0
			p0_changes = 0
		case ']':
			if state == SIMPLE {
				if ptr_rel_loc == 0 && (p0_changes == 1 || p0_changes == -1) {
					simples = append(simples, bracketPairs[i])
				} else {
					complexes = append(complexes, bracketPairs[i])
				}
			} else if state == COMPLEX {
				complexes = append(complexes, bracketPairs[i])
			}
			state = CLOSED
		case '.':
			state = COMPLEX
		case ',':
			state = COMPLEX
		case '>':
			ptr_rel_loc++
		case '<':
			ptr_rel_loc--
		case '+':
			if ptr_rel_loc == 0 {
				p0_changes++
			}
		case '-':
			if ptr_rel_loc == 0 {
				p0_changes--
			}
		}
	}
	return simples, complexes
}

func interpret_profiler(program string) {
	var cmd_count [8]int

	const TAPE_SIZE = 1024 * 4
	var TAPE [TAPE_SIZE]byte
	var POINTER int = 0
	bracketPairs := locate_brackets(program)
	leftBrackCount := make(map[int]int)

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
			if TAPE[POINTER] == 0 {
				PC = bracketPairs[PC]
				continue
			}
			cmd_count[6]++
			leftBrackCount[PC]++
		case ']':
			if TAPE[POINTER] != 0 {
				PC = bracketPairs[PC]
				continue
			}
			cmd_count[7]++
		default:
			// do nothing
		}
	}

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
			return leftBrackCount[slice[i]] > leftBrackCount[slice[j]]
		}
	}

	simple, complex := categorize_brackets(program, bracketPairs)
	sort.Slice(simple, sort_func(simple))
	sort.Slice(complex, sort_func(complex))

	fmt.Println("Simple Innermost Loops")
	for _, v := range simple {
		fmt.Println(v, leftBrackCount[v])
	}

	fmt.Println("Complex Innermost Loops")
	for _, v := range complex {
		fmt.Println(v, leftBrackCount[v])
	}

}

func locate_brackets(program string) map[int]int {
	bracketPairs := make(map[int]int)
	stack := make(stack, 0, len(program)/2)

	for i, v := range program {
		if v == '[' {
			stack.Push(i)
		} else if v == ']' {
			l_loc := stack.Pop()
			bracketPairs[l_loc] = i
			bracketPairs[i] = l_loc
		}
	}
	if len(stack) != 0 {
		panic("Mismatching [ & ]")
	}

	return bracketPairs
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

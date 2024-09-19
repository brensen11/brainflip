package main

import (
	"fmt"
	"os"
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

func run(program string) {
	const TAPE_SIZE = 1024 * 4
	var TAPE [TAPE_SIZE]byte
	var POINTER int = 0

	bracketPairs := make(map[int]int)
	stack := make(stack, len(program)/2)
	stack = stack[:0]

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

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Incorrect usage: ", os.Args)
	}
	program := string(data)
	run(program)
}

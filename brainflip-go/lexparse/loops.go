package lexparse

import (
	"brainflip-go/utils"
)

func Categorize_Brackets(instructions []Instruction) ([]int, []int) {
	const (
		SIMPLE int = iota
		COMPLEX
		CLOSED
	)
	simples := make([]int, 0, len(instructions)/2)
	complexes := make([]int, 0, len(instructions)/2)

	// a simple loop is one that...
	// - contains no i/o,
	// - has 0 net change to the pointer
	// - and that changes p[0] by either +1 or -1 after each loop iteration, aka by the end of the loop [0] has changed by 1
	// (the cell that the pointer references when the loop body starts executing)
	ptr_rel_loc := 0
	p0_changes := 0
	state := CLOSED
	for i, v := range instructions {
		switch v.(type) {
		case Left_loop:
			state = SIMPLE
			ptr_rel_loc = 0
			p0_changes = 0
		case Right_loop:
			if state == SIMPLE {
				if ptr_rel_loc == 0 && (p0_changes == 1 || p0_changes == -1) {
					simples = append(simples, i)
				} else {
					complexes = append(complexes, i)
				}
			} else if state == COMPLEX {
				complexes = append(complexes, i)
			}
			state = CLOSED
		case Output:
			state = COMPLEX
		case Input:
			state = COMPLEX
		case Move_right:
			ptr_rel_loc++
		case Move_left:
			ptr_rel_loc--
		case Inc:
			if ptr_rel_loc == 0 {
				p0_changes++
			}
		case Dec:
			if ptr_rel_loc == 0 {
				p0_changes--
			}
		}
	}
	return simples, complexes
}

// Key Value pair of indices for matching '[' and ']' characters
func Locate_Brackets(instructions []Instruction) map[int]int {
	bracketPairs := make(map[int]int)
	stack := make(utils.Stack, 0, len(instructions)/2)

	for i, v := range instructions {
		switch v.(type) {
		case Left_loop:
			stack.Push(i)
		case Right_loop:
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

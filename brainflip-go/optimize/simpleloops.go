package optimize

import (
	lp "brainflip-go/lexparse"
	"fmt"
)

func Optimize_simple_loops(instructions *[]lp.Instruction) {
	bracketPairs := lp.Locate_Brackets(*instructions)
	simple_loops, _ := lp.Categorize_Brackets(*instructions)

	for i := len(simple_loops) - 1; i >= 0; i-- { // for each simple loop
		right_loop_index := simple_loops[i]
		left_loop_index := bracketPairs[right_loop_index]
		loop_instructions := (*instructions)[left_loop_index : right_loop_index+1] // get instruction '[' -> ']' inclusive

		// Doing analysis to see which parts of the loop are incremented by what
		rel_cell_change := make(map[int]int)
		REL_PTR := 0
		for i := range loop_instructions {

			instruction := loop_instructions[i]
			switch instruction.(type) {
			case lp.Move_right:
				REL_PTR++
			case lp.Move_left:
				REL_PTR--
			case lp.Inc:
				rel_cell_change[REL_PTR] += 1
			case lp.Dec:
				rel_cell_change[REL_PTR] -= 1
			}
		}

		var new_instructions = make([]lp.Instruction, 0)

		p0 := rel_cell_change[0] // p[0] represents the pointers relative position as you enter the loop
		if p0 != 1 && p0 != -1 {
			fmt.Println(lp.Instructions_string(loop_instructions))
			fmt.Println(p0)
			panic("Something went wrong with the count of the p[0] relative change")
		}

		if len(rel_cell_change) > 1 { // if there are other changes in the loop then we store p[0] for later use
			if p0 == -1 {
				new_instructions = append(new_instructions, lp.Store{lp.R0, lp.Offset(0)}) // R0 = p[0]
			} else {
				new_instructions = append(new_instructions, lp.Store{lp.R0, lp.Imm(0)})  //  R0 = 256
				new_instructions = append(new_instructions, lp.Sub{lp.R0, lp.Offset(0)}) // 	R0 = R1 - p[0]
			}
		}

		for offset, change := range rel_cell_change {
			if offset == 0 {
				continue
			} // skip for p[0] as we know it changes by 1 until 0

			add := change > 0
			for i := 0; i < abs(change); i++ {
				if add {
					new_instructions = append(new_instructions, lp.Add{lp.Offset(offset), lp.R0})
				} else {
					new_instructions = append(new_instructions, lp.Sub{lp.Offset(offset), lp.R0})
				}
			}
		}

		// Set p[0] every time
		new_instructions = append(new_instructions, lp.Store{lp.Offset(0), lp.Imm(0)}) // p[0] = 0

		*instructions = lp.Instructions_replace(*instructions, left_loop_index, right_loop_index+1, new_instructions)
	}
}

// what a dumb thing
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// psuedo := `{0: -1, 1: 3, -4: 5}
//  R0 = p[0]
// 	- - - - - // for each kv
//  p[1] = p[1] + R0
//  p[1] = p[1] + R0
//  p[1] = p[1] + R0
// 	- - - - - // for each kv
// 	p[3] = p[3] - R0
// 	p[3] = p[3] - R0
// 	p[3] = p[3] - R0
// 	p[3] = p[3] - R0
// 	- - - - - // counter var set to 0
// 	p[0] = 0
// `
// psuedo_2 := `{0: 1, 1: 3, 4: 5}
//  R0 = 256
// 	R0 = R0 - p[0]
// 	- - - - -
// 	R1 = R0
// 	R1 = R1 * ((255 - 3) + 1) // inclusiveness for 255 itself
// 	p[1] = p[1] + R1
// 	- - - - -
// 	p[0] = 0
// `

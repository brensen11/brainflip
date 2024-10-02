package optimize

import (
	lp "brainflip-go/lexparse"
	"fmt"
)

func Optimize_simple_loops(program *lp.Program) {

	for i := len(*program.Simple_loops) - 1; i >= 0; i-- {
		right_loop_index := (*program.Simple_loops)[i]
		left_loop_index := (*program.BracketPairs)[right_loop_index]
		loop_instructions := (*program.Instructions)[left_loop_index : right_loop_index+1]

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

		var new_instructions []lp.Instruction
		loop_increment := rel_cell_change[0]

		if loop_increment != 1 && loop_increment != -1 { // TODO replace assert
			fmt.Println(lp.Instructions_string(loop_instructions))
			fmt.Println(loop_increment)
			panic("Something went wrong with the count of the p[0] relative change")
		}

		// psuedo := `{0: -1, 1: 3, -4: 5}
		// 	R1 = p[0] // init counter var
		// 	- - - - - // for each kv
		// 	R2 = R1
		// 	R2 = R2 * 3
		// 	p[1] = p[1] + R2
		// 	- - - - - // for each kv
		// 	R2 = R1
		// 	R2 = R2 * 5
		// 	p[-4] = p[-4] + R2
		// 	- - - - - // counter var set to 0
		// 	p[0] = 0
		// `
		// psuedo_2 := `{0: 1, 1: 3, 4: 5}
		// 	R1 = p[0]
		// 	- - - - -
		// 	R2 = R1
		// 	R2 = R2 * ((255 - 3) + 1) // inclusiveness for 255 itself
		// 	p[1] = p[1] + R2
		// 	- - - - -
		// 	p[0] = 0
		// `
		reverse := loop_increment == 1
		if len(rel_cell_change) > 1 {
			new_instructions = append(new_instructions, lp.Store{lp.R0, lp.Offset(0)}) // R0 = p[0]
		}
		for offset, change := range rel_cell_change {
			if offset == 0 {
				continue
			}

			// new_instructions = append(new_instructions, lp.)
			new_instructions = append(new_instructions, lp.Store{lp.R1, lp.R0}) // R1 = R0
			if reverse {
				new_instructions = append(new_instructions, lp.Store{lp.R2, lp.Imm(256)}) // R2 = 256
				new_instructions = append(new_instructions, lp.Sub{lp.R2, lp.R1})         // R2 = R2 - R1
				new_instructions = append(new_instructions, lp.Store{lp.R1, lp.R2})       // R1 = R2
			}
			new_instructions = append(new_instructions, lp.Mul{lp.R1, lp.Imm(change)})    // R1 = R1 * IMM
			new_instructions = append(new_instructions, lp.Add{lp.Offset(offset), lp.R1}) // p[offset] = p[offset] + R1
		}
		new_instructions = append(new_instructions, lp.Store{lp.Offset(0), lp.Imm(0)}) // p[0] = 0

		*program.Instructions = lp.Instructions_replace(*program.Instructions, left_loop_index, right_loop_index+1, new_instructions)
		*program.Simple_loops, *program.Complex_loops = lp.Categorize_Brackets(*program.Instructions)
		*program.BracketPairs = lp.Locate_Brackets(*program.Instructions)
	}
}

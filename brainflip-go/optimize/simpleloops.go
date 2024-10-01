package optimize

import (
	"brainflip-go/lexparse"
	"fmt"
	"slices"
)

func Optimize_simple_loops(program lexparse.Program) {

	for i := len(program.Simple_loops) - 1; i >= 0; i-- {
		right_loop_index := program.Simple_loops[i]
		left_loop_index := program.BracketPairs[right_loop_index]
		loop_instructions := program.Instructions[left_loop_index : right_loop_index+1]

		rel_cell_change := make(map[int]int)
		REL_PTR := 0
		for i := range loop_instructions {

			instruction := loop_instructions[i]
			switch instruction.(type) {
			case lexparse.Move_right:
				REL_PTR++
			case lexparse.Move_left:
				REL_PTR--
			case lexparse.Inc:
				rel_cell_change[REL_PTR] += 1
			case lexparse.Dec:
				rel_cell_change[REL_PTR] -= 1
			}
		}

		var new_instructions []lexparse.Instruction
		loop_increment := rel_cell_change[0]

		if loop_increment != 1 && loop_increment != -1 { // TODO replace assert
			fmt.Println(lexparse.Instructions_string(loop_instructions))
			fmt.Println(loop_increment)
			panic("Something went wrong with the count of the p[0] relative change")
		}

		// psuedo := `{0: -1, 1: 3, -4: 5}
		// 	TMP_1 = p[0] // init counter var
		// 	- - - - - // for each kv
		// 	TMP_2 = TMP_1
		// 	TMP_2 = TMP_2 * 3
		// 	p[1] = p[1] + TMP_2
		// 	- - - - - // for each kv
		// 	TMP_2 = TMP_1
		// 	TMP_2 = TMP_2 * 5
		// 	p[-4] = p[-4] + TMP_2
		// 	- - - - - // counter var set to 0
		// 	p[0] = 0
		// `
		// psuedo_2 := `{0: 1, 1: 3, 4: 5}
		// 	TMP_1 = p[0]
		// 	- - - - -
		// 	TMP_2 = TMP_1
		// 	TMP_2 = TMP_2 * ((255 - 3) + 1) // inclusiveness for 255 itself
		// 	p[1] = p[1] + TMP_2
		// 	- - - - -
		// 	p[0] = 0
		// `
		reverse := loop_increment == 1
		new_instructions = append(new_instructions, lexparse.Store_Reg_Offset{lexparse.R1, 0}) // TMP_0 = p[0]
		for offset, change := range rel_cell_change {
			// if offset == 0 {
			// 	continue
			// }

			// new_instructions = append(new_instructions, lexparse.)
			new_instructions = append(new_instructions, lexparse.Store_Reg_Reg{lexparse.R2, lexparse.R1}) // TMP_1 = TMP_0
			times := change
			if reverse {
				times = 256 - change
			}
			new_instructions = append(new_instructions, lexparse.Mul_Reg_Imm{lexparse.R2, times})     // TMP_1 = TMP_1 * IMM
			new_instructions = append(new_instructions, lexparse.Add_Offset_Reg{offset, lexparse.R2}) // p[offset] = p[offset] + TMP_1
		}
		new_instructions = append(new_instructions, lexparse.Set_Offset_Imm{0, 0}) // p[0] = 0

		program.Instructions = slices.Replace(program.Instructions, left_loop_index, right_loop_index+1, new_instructions...)
	}
}

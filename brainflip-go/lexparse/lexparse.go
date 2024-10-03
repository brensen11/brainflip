package lexparse

func Lexparse(program string) *[]Instruction {
	var instructions []Instruction

	// main run function
	for i := 0; i < len(program); i++ {
		cmd := program[i]
		switch cmd {
		case '>':
			instructions = append(instructions, Move_right{})
		case '<':
			instructions = append(instructions, Move_left{})
		case '+':
			instructions = append(instructions, Inc{})
		case '-':
			instructions = append(instructions, Dec{})
		case '.':
			instructions = append(instructions, Output{})
		case ',':
			instructions = append(instructions, Input{})
		case '[':
			instructions = append(instructions, Left_loop{})
		case ']':
			instructions = append(instructions, Right_loop{})
		default:
			// do nothing
		}
	}

	return &instructions
}

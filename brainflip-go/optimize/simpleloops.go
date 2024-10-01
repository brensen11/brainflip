package optimize

import (
	"brainflip-go/lexparse"
)

func optimize_simple_loops(program lexparse.Program) {

	// for i, v := range program.Simple_loops {

	// }

	for i := 0; i < len(program.Instructions); i++ {
		instruction := program.Instructions[i]

		switch instruction.(type) {
		case lexparse.Move_right:
		case lexparse.Move_left:
		case lexparse.Inc:
		case lexparse.Dec:
		case lexparse.Output:
		case lexparse.Input:
		case lexparse.Left_loop:
		case lexparse.Right_loop:
		case lexparse.Set:
		case lexparse.Add:
		case lexparse.Sub:
		}
	}
}

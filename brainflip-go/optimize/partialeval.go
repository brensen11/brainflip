package optimize

import (
	lp "brainflip-go/lexparse"
	"os"
	"strings"
)

func Optimize_partialeval(instructions *[]lp.Instruction) (*[]byte, int, string) {
	const TAPE_SIZE = 1024 * 1024 * 4
	var TAPE []byte = make([]byte, TAPE_SIZE)
	var POINTER int = TAPE_SIZE / 2

	bracketPairs := lp.Locate_Brackets(*instructions)
	// main run function
	var PC int
	var end_loop bool = false
	var out_builder strings.Builder
	for PC = 0; PC < len(*instructions); PC++ {
		cmd := (*instructions)[PC]
		// fmt.Println("Checking for: ", cmd)
		switch cmd.(type) {
		case lp.Move_right:
			if POINTER == TAPE_SIZE {
				panic("Tape Pointer Out of Bounds!!")
			}
			POINTER++
		case lp.Move_left:
			if POINTER == 0 {
				panic("Tape Pointer Out of Bounds!!")
			}
			POINTER--
		case lp.Inc:
			TAPE[POINTER]++
		case lp.Dec:
			TAPE[POINTER]--
		case lp.Output:
			os.Stdout.Write([]byte{TAPE[POINTER]})
			out_builder.WriteByte(TAPE[POINTER])
		case lp.Input:
			end_loop = true
			// var input [1]byte
			// _, err := os.Stdin.Read(input[:])
			// if err != nil && err == io.EOF {
			// 	TAPE[POINTER] = 255
			// } else {
			// 	TAPE[POINTER] = input[0]
			// }
		case lp.Left_loop:
			left_count := 1
			for loop_inspect_ptr := PC; ; loop_inspect_ptr++ {
				loop_cmd := (*instructions)[loop_inspect_ptr]
				switch loop_cmd.(type) {
				case lp.Right_loop:
					left_count--
				case lp.Input:
					end_loop = true
				}

				if left_count == 0 || end_loop {
					break
				}
			}

			if TAPE[POINTER] == 0 {
				PC = bracketPairs[PC]
				continue
			}
		case lp.Right_loop:
			if TAPE[POINTER] != 0 {
				PC = bracketPairs[PC]
				continue
			}
		}
		if end_loop {
			break
		}
	}
	*instructions = (*instructions)[PC:]        // delete part of the program that is unevaluated
	return &TAPE, POINTER, out_builder.String() // return the tape and tape head
}

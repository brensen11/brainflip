package interpreter

import (
	"brainflip-go/lexparse"
	"brainflip-go/utils"
	"io"
	"os"
)

func Interpret(filename string) {
	code := utils.Readfile(filename)
	program := lexparse.Lexparse(code)

	const TAPE_SIZE = 1024 * 1024 * 64
	var TAPE [TAPE_SIZE]byte
	var POINTER int = TAPE_SIZE / 2

	bracketPairs := lexparse.Locate_Brackets(*program.Instructions)
	// main run function
	for PC := 0; PC < len(*program.Instructions); PC++ {
		cmd := (*program.Instructions)[PC]
		// fmt.Println("Checking for: ", cmd)
		switch cmd.(type) {
		case lexparse.Move_right:
			if POINTER == TAPE_SIZE {
				panic("Tape Pointer Out of Bounds!!")
			}
			POINTER++
		case lexparse.Move_left:
			if POINTER == 0 {
				panic("Tape Pointer Out of Bounds!!")
			}
			POINTER--
		case lexparse.Inc:
			TAPE[POINTER]++
		case lexparse.Dec:
			TAPE[POINTER]--
		case lexparse.Output:
			os.Stdout.Write([]byte{TAPE[POINTER]})
		case lexparse.Input:
			var input [1]byte
			_, err := os.Stdin.Read(input[:])
			if err != nil && err == io.EOF {
				TAPE[POINTER] = 255
			} else {
				TAPE[POINTER] = input[0]
			}
		case lexparse.Left_loop:
			if TAPE[POINTER] == 0 {
				PC = bracketPairs[PC]
				continue
			}
		case lexparse.Right_loop:
			if TAPE[POINTER] != 0 {
				PC = bracketPairs[PC]
				continue
			}
		default:
			// do nothing
		}
	}
}

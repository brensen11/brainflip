package utils

import (
	"fmt"
	"os"
)

// Stack Data Structure
type Stack []int

func (s *Stack) Push(v int) {
	*s = append(*s, v)
}

func (s *Stack) Pop() int {
	l := len(*s)
	if l == 0 {
		panic("Tried to pop empty stack")
	}
	val := (*s)[l-1]
	*s = (*s)[:l-1]
	return val
}

func Readfile(filename string) string {
	bf_data, prog_err := os.ReadFile(filename)
	if prog_err != nil {
		panic("There was an error reading: " + filename)
	}
	return string(bf_data)
}

func Writefile(content string) {
	file, err := os.Create("out-win.asm")
	if err != nil {
		panic(fmt.Sprint("Error creating file:", err))
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		panic(fmt.Sprint("Error writing to file:", err))
	}
}

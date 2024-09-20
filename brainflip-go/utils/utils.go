package utils

// Stack Data Structure
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

// Key Value pair of indices for matching '[' and ']' characters
func Locate_Brackets(program string) map[int]int {
	bracketPairs := make(map[int]int)
	stack := make(stack, 0, len(program)/2)

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

	return bracketPairs
}

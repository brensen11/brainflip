package lexparse

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// powers of two supported by ... me! that is
func isPowerOfTwo(n int) bool {
	np := abs(n)
	return np == 1 || np == 2 || np == 4 || np == 8
}

type Scan struct {
	L     int
	R     int
	Moves int
}

func Locate_Scans(instructions *[]Instruction) []Scan {
	var valid bool = false
	var L int
	var moves = 0
	var scans []Scan
	for i, v := range *instructions {
		switch v.(type) {
		case Left_loop:
			L = i
			valid = true
			moves = 0
		case Right_loop:
			if valid && isPowerOfTwo(moves) {
				scans = append(scans, Scan{L, i, moves})
			}
			moves = 0
			valid = false
		case Move_right:
			moves++
		case Move_left:
			moves--
		default:
			moves = 0
			valid = false
		}
	}
	return scans
}

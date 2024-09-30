package lexparse

type Instruction interface {
	isInstruction()
}

// ----- default 8 instruction -----
type Move_right struct{}
type Move_left struct{}
type Inc struct{}
type Dec struct{}
type Output struct{}
type Input struct{}
type Left_loop struct{}
type Right_loop struct{}

// ----- implementation specific IR -----

type Set struct {
	value  int
	offset int
}
type Add struct {
	value  uint
	offset int
}
type Sub struct {
	value  uint
	offset int
}

// types extend Instruction

func (Move_right) isInstruction() {}
func (Move_left) isInstruction()  {}
func (Inc) isInstruction()        {}
func (Dec) isInstruction()        {}
func (Output) isInstruction()     {}
func (Input) isInstruction()      {}
func (Left_loop) isInstruction()  {}
func (Right_loop) isInstruction() {}
func (Set) isInstruction()        {}
func (Add) isInstruction()        {}
func (Sub) isInstruction()        {}

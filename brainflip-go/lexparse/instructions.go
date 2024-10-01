package lexparse

import (
	"fmt"
	"strings"
)

type Instruction interface {
	isInstruction()
	String() string
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
	offset int
	value  int
}
type Add struct {
	offset int
	value  uint
}
type Sub struct {
	offset int
	value  uint
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

// toString stmts

func (Move_right) String() string { return ">" }
func (Move_left) String() string  { return "<" }
func (Inc) String() string        { return "+" }
func (Dec) String() string        { return "-" }
func (Output) String() string     { return "." }
func (Input) String() string      { return "," }
func (Left_loop) String() string  { return "[" }
func (Right_loop) String() string { return "]" }
func (set Set) String() string {
	return fmt.Sprint("set(offset:", set.offset, "value:", set.value, ")")
}
func (add Add) String() string {
	return fmt.Sprint("add(offset:", add.offset, "value:", add.value, ")")
}
func (sub Sub) String() string {
	return fmt.Sprint("sub(offset:", sub.offset, "value:", sub.value, ")")
}

func Instructions_string(instrs []Instruction) string {
	var b strings.Builder
	for _, v := range instrs {
		b.WriteString(v.String())
	}
	return b.String()
}

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
type Add struct {
	Op1 Operand
	Op2 Operand
}
type Sub struct {
	Op1 Operand
	Op2 Operand
}
type Mul struct {
	Op1 Operand
	Op2 Operand
}
type Store struct {
	Op1 Operand
	Op2 Operand
}
type Raw struct {
	Raw string
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

func (Add) isInstruction()   {}
func (Sub) isInstruction()   {}
func (Mul) isInstruction()   {}
func (Store) isInstruction() {}
func (Raw) isInstruction()   {}

// toString stmts
func (Move_right) String() string { return ">" }
func (Move_left) String() string  { return "<" }
func (Inc) String() string        { return "+" }
func (Dec) String() string        { return "-" }
func (Output) String() string     { return "." }
func (Input) String() string      { return "," }
func (Left_loop) String() string  { return "[" }
func (Right_loop) String() string { return "]" }

func (i Add) String() string {
	return fmt.Sprint("\n", i.Op1, " = ", i.Op1, " + ", i.Op2, "\n")
}
func (i Sub) String() string {
	return fmt.Sprint("\n", i.Op1, " = ", i.Op1, " - ", i.Op2, "\n")
}
func (i Mul) String() string {
	return fmt.Sprint("\n", i.Op1, " = ", i.Op1, " * ", i.Op2, "\n")
}
func (i Store) String() string {
	return fmt.Sprint("\n", i.Op1, " = ", i.Op2, "\n")
}
func (raw Raw) String() string { return raw.Raw }

func Instructions_string(instrs []Instruction) string {
	var b strings.Builder
	for _, v := range instrs {
		b.WriteString(v.String())
	}
	final := b.String()
	final = strings.Replace(final, "\n\n", "\n", -1)
	return final
}

func Instructions_replace(instrs []Instruction, i int, j int, replace_instrs []Instruction) []Instruction {
	new_instrs := make([]Instruction, 0)
	new_instrs = append(new_instrs, instrs[0:i]...)
	new_instrs = append(new_instrs, replace_instrs[:]...)
	new_instrs = append(new_instrs, instrs[j:]...)
	return new_instrs
}

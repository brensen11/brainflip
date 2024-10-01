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
type Tmp_Reg int

const (
	R1 Tmp_Reg = iota
	R2
	R3
	R4
)

type Store_Reg_Offset struct { // R = p[Offset]
	Reg    Tmp_Reg
	Offset int
}
type Store_Reg_Reg struct { // R1 = R2
	Reg_1 Tmp_Reg
	Reg_2 Tmp_Reg
}
type Set_Offset_Imm struct {
	Offset int
	Imm    int
}
type Add_Offset_Reg struct {
	Offset_1 int
	Reg      Tmp_Reg
}
type Sub struct {
	Offset_1 int
	Offset_2 int
}
type Mul_Reg_Imm struct {
	Reg Tmp_Reg
	Imm int
}
type Raw struct {
	raw string
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

func (Store_Reg_Offset) isInstruction() {}
func (Store_Reg_Reg) isInstruction()    {}
func (Set_Offset_Imm) isInstruction()   {}
func (Add_Offset_Reg) isInstruction()   {}
func (Sub) isInstruction()              {}
func (Mul_Reg_Imm) isInstruction()      {}
func (Raw) isInstruction()              {}

// toString stmts

func (Move_right) String() string { return ">" }
func (Move_left) String() string  { return "<" }
func (Inc) String() string        { return "+" }
func (Dec) String() string        { return "-" }
func (Output) String() string     { return "." }
func (Input) String() string      { return "," }
func (Left_loop) String() string  { return "[" }
func (Right_loop) String() string { return "]" }

func (store Store_Reg_Offset) String() string {
	return fmt.Sprintf("\nR%d = p[%d]\n", store.Reg, store.Offset)
}
func (store Store_Reg_Reg) String() string {
	return fmt.Sprintf("\nR%d = R%d\n", store.Reg_1, store.Reg_2)
}

func (set Set_Offset_Imm) String() string {
	return fmt.Sprintf("\np[%d] = %d\n", set.Offset, set.Imm)
}
func (add Add_Offset_Reg) String() string {
	return fmt.Sprintf("\np[%d] = p[%d] + R%d\n", add.Offset_1, add.Offset_1, add.Reg)
}
func (sub Sub) String() string {
	return fmt.Sprintf("\np[%d] = p[%d] + p[%d]\n", sub.Offset_1, sub.Offset_1, sub.Offset_2)
}
func (mul Mul_Reg_Imm) String() string {
	return fmt.Sprintf("\nR%d = R%d * %d\n", mul.Reg, mul.Reg, mul.Imm)
}
func (raw Raw) String() string {
	return raw.raw
}

func Instructions_string(instrs []Instruction) string {
	var b strings.Builder
	for _, v := range instrs {
		b.WriteString(v.String())
	}
	return b.String()
}

// func Instruction_replace(i int, j int, instrs []Instruction, new_instrs []Instruction) {

// }

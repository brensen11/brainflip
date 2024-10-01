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
	R0 Tmp_Reg = iota
	R1
	R2
	R3
)

// type Reg struct{}
// type Offset struct{}
// type Imm struct{}

// type Mul struct {
// 	Op1, Op2
// }

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
type Sub_Imm_Reg struct {
	Imm int
	Reg Tmp_Reg
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
func (Sub_Imm_Reg) isInstruction()      {}
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
func (sub Sub_Imm_Reg) String() string {
	return fmt.Sprintf("\nR%d = %d - R%d\n", sub.Reg, sub.Imm, sub.Reg)
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
	final := b.String()
	final = strings.Replace(final, "\n\n", "\n", -1)
	return final
}

func Instructions_replace(instrs []Instruction, i int, j int, replace_instrs []Instruction) []Instruction {
	new_instrs := append(instrs[0:i], replace_instrs[:]...)
	new_instrs = append(new_instrs, instrs[j:]...)
	return new_instrs
}

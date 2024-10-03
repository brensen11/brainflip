package lexparse

import (
	"fmt"
)

type Operand interface {
	isOperand()
	String() string
	ToAsm() string
}

type reg int
type Offset int
type Imm int

func (reg) isOperand()    {}
func (Offset) isOperand() {}
func (Imm) isOperand()    {}

func (reg reg) String() string {
	return fmt.Sprintf("R%d", reg)
}
func (off Offset) String() string {
	return fmt.Sprintf("p[%d]", off)
}
func (imm Imm) String() string {
	return fmt.Sprintf("%d", imm)
}

// TODO this doesn't belong here move to generator??????

func (reg reg) ToAsm() string {
	regname := reg + 12
	return fmt.Sprintf("r%db", regname)
}
func (off Offset) ToAsm() string {
	return fmt.Sprintf("BYTE [rdi + %d]", off)
}
func (imm Imm) ToAsm() string {
	return imm.String()
}

const (
	R0 reg = iota
	R1
	R2
	R3
)

package stack

import (
	"chip-8/register"
	"chip-8/util"
)

type Stack struct {
	Stack   [80]register.Register
	pointer uint8

}

func New() *Stack{
    return &Stack{}
}

func (s *Stack) Push(val uint16) {
	if s.pointer < 80 {
		s.pointer++

		_, low := util.DecimalToBinary16(val)

		s.Stack[s.pointer].Write(low)
	}

}

func (s *Stack) Pop() (low [8]bool) {
	low = s.Stack[s.pointer].Read()
	if s.pointer > 0 {
		s.pointer--
	}

	return low
}

func (s *Stack) Peek() (low [8]bool){
	low = s.Stack[s.pointer].Read()

	return low
}

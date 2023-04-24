package stack

import (
	"chip-8/register"
	"chip-8/util"
)

type Stack struct {
	Stack   [80][2]register.Register
	pointer uint16

}

func New() *Stack{
    return &Stack{pointer: 0}
}

func (s *Stack) Push(val uint16) {
	if s.pointer < 80 {
		s.pointer++

		high, low := util.DecimalToBinary16(val)

		s.Stack[s.pointer][0].Write(high)
		s.Stack[s.pointer][1].Write(low)
	}

}

func (s *Stack) Pop() (high, low [8]bool) {

	high = s.Stack[s.pointer][0].Read()
	low = s.Stack[s.pointer][1].Read()
	if s.pointer > 0 {
		s.pointer--
	}

	return high, low
}

func (s *Stack) Peek() (high, low [8]bool){
	high = s.Stack[s.pointer][0].Read()
	low = s.Stack[s.pointer][1].Read()

	return high, low
}

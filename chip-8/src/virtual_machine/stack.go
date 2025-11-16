package chip8

import (
	"errors"
)

/*
	The stack is an array of 16 16-bit values, used to store the address
	that the interpreter shoud return to when finished with a subroutine.
	Chip-8 allows for up to 16 levels of nested subroutines.
*/

const STACK_SIZE = 16

type STACK struct {
	stack [STACK_SIZE]uint16
	SP    int8 //Stack Pointer, -1 = level main (currently not in a function execution)
}

func (s *STACK) Init() {
	s.SP = -1
}

func (s *STACK) Push(addr uint16) error {
	if s.SP == -1 || s.SP+1 < STACK_SIZE {
		s.SP++               //increase stack pointer
		s.stack[s.SP] = addr //insert data
		return nil
	} else { //STACK OVERFLOW
		return errors.New("STACK OVERFLOW")
	}
}

func (s *STACK) Pop() (uint16, error) {
	if s.SP >= 0 {
		s.SP--
		return s.stack[s.SP+1], nil
	} else {
		return uint16(0), errors.New("can not pop fron an empty stack")
	}
}

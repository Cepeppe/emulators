package main

import (
	"errors"
	"fmt"
	"strings"
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
	if s.SP+1 < STACK_SIZE {
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

func (st *STACK) Dump() {
	const innerWidth = 76

	printSep := func() {
		fmt.Printf("\t+%s+\n", strings.Repeat("-", innerWidth+2))
	}
	printLine := func(text string) {
		fmt.Printf("\t| %-*s |\n", innerWidth, text)
	}

	fmt.Printf("\n\tSTACK STATE DUMP\n")
	printSep()
	printLine("STACK CONTENT [index] [bit view] (dec)")
	printSep()

	// Two stack entries per line
	for i := 0; i < STACK_SIZE; i += 2 {
		v1 := st.stack[i]

		var line string
		if i+1 < STACK_SIZE {
			v2 := st.stack[i+1]
			line = fmt.Sprintf(
				"[%02d]: %s (%5d)  [%02d]: %s (%5d)",
				i, formatBits16(v1), v1,
				i+1, formatBits16(v2), v2,
			)
		} else {
			line = fmt.Sprintf(
				"[%02d]: %s (%5d)",
				i, formatBits16(v1), v1,
			)
		}

		printLine(line)
	}

	printSep()
	// SP is int8, shown as 8-bit value with padding like other 8-bit registers
	printLine(fmt.Sprintf("SP (stack pointer): %s (%3d)", formatBits8(uint8(st.SP)), st.SP))
	printSep()
}

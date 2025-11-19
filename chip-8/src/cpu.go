package main

import (
	"fmt"
	"strings"
)

/*
	Chip-8 has 16 general purpose 8-bit registers,
	usually referred to as Vx, where x is a hexadecimal digit (0 through F).
		VF register should not be used by any program,
		as it is used as a flag by some instructions.

	There is also a 16-bit register called I.
	This register is generally used to store memory addresses,
	so only the lowest (rightmost) 12 bits are usually used.

	There are two special 8-bit registers used as timers, the DT and the ST.
		-The delay timer DT is automatically decremented with a frequency of 60 Hz
		(60 times per second) whenever its value is greater than zero (> 0).
		That’s all it does.
		Its value can be read into a registry and written with FX07 and FX15 respectively.
		-The sound timer ST is also automatically decremented with a frequency of 60 Hz
		when its value is greater than zero.
		Additionally, when this happens (ST > 0) the system sounds the buzzer to produce a beep.
		So, for example, if you want your program to sound the buzzer for one second,
		you need to write 0x3C to ST. It can be set with FX18 but it can’t be read.

	The program counter (PC) should be 16-bit,
	and is used to store the currently executing address

	The stack pointer (SP) can be 8-bit,
	it is used to point to the topmost level of the stack
*/

const (
	V0 = 0
	V1 = 1
	V2 = 2
	V3 = 3
	V4 = 4
	V5 = 5
	V6 = 6
	V7 = 7
	V8 = 8
	V9 = 9
	VA = 10
	VB = 11
	VC = 12
	VD = 13
	VE = 14
	VF = 15
)

const cpuDumpInnerWidth = 55

type CPU struct {

	//General purpose registers
	// VF register should not be used by any program
	V_registers [16]uint8

	//I register. only the lowest (rightmost) 12 bits are usually used
	I uint16

	//The sound timer
	ST uint8

	//The delay timer
	DT uint8

	//Program Counter
	PC uint16

	//Current opcode
	curr_opcode uint16
}

func (cpu *CPU) Init() {
	for i := 0; i < 16; i++ {
		cpu.V_registers[i] = uint8(0)
	}
	cpu.I = uint16(0)
	cpu.ST = uint8(0)
	cpu.DT = uint8(0)
	cpu.PC = PROGRAM_MEMORY_START
	cpu.curr_opcode = 0
}

func (cpu *CPU) Dump() {
	// Local helpers for aligned box drawing
	printSep := func() {
		// if you are asking why the +2 in the following line, don't. It just works and it's a pretty print so it's not an issue
		fmt.Printf("\t+%s+\n", strings.Repeat("-", cpuDumpInnerWidth+2))
	}
	printLine := func(text string) {
		fmt.Printf("\t| %-*s |\n", cpuDumpInnerWidth, text)
	}

	fmt.Printf("\n\tCPU STATE DUMP\n")
	printSep()
	printLine("V REGISTERS (V0-VF)   [bit view] (dec)")
	printSep()

	// V0–VF: two registers per line
	for i := 0; i < 16; i += 2 {
		v1 := cpu.V_registers[i]
		v2 := cpu.V_registers[i+1]

		line := fmt.Sprintf(
			"V%X: %08b (%3d)  V%X: %08b (%3d)",
			i, v1, v1,
			i+1, v2, v2,
		)
		printLine(line)
	}

	printSep()
	printLine("SPECIAL REGISTERS [bit view] (dec)")
	printSep()

	// All lines share the same internal width; 8-bit registers are left-padded with spaces instead of bits
	printLine(fmt.Sprintf("I          : %s (%5d)", formatBits16(cpu.I), cpu.I))
	printLine(fmt.Sprintf("PC         : %s (%5d)", formatBits16(cpu.PC), cpu.PC))
	printLine(fmt.Sprintf("DT (Delay) : %s (%5d)", formatBits8(cpu.DT), cpu.DT))
	printLine(fmt.Sprintf("ST (Sound) : %s (%5d)", formatBits8(cpu.ST), cpu.ST))
	printLine(fmt.Sprintf("curr_opcode: %s (%5d)", formatBits16(cpu.curr_opcode), cpu.curr_opcode))

	printSep()
}

package main

import (
	"fmt"
)

//n: number of most significative bits to isolate
//For example:
//   instr = 11110000 00000000  k = 3
//   head  = 00000000 00000111
//   tail  = 00010000 00000000
func isolateBits(instr uint16, k uint8) (uint16, uint16) {

	// First 16-k bit more significative, moved towards right
	head := instr >> (16 - k)

	// Last 16-k bit less significative
	tail := instr & ((1 << (16 - k)) - 1)

	return head, tail
}

// 16-bit value formatted as "xxxxxxxx xxxxxxxx"
func formatBits16(v uint16) string {
	bits := fmt.Sprintf("%016b", v) // 16 bits
	return bits[:8] + " " + bits[8:]
}

// 8-bit value formatted in a 16-bit field where the high byte is spaces:
// "        00000000" (8 spaces, space separator, 8 bits)
func formatBits8(v uint8) string {
	bits := fmt.Sprintf("%08b", v) // 8 bits
	return fmt.Sprintf("%8s %s", "", bits)
}

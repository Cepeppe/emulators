package chip8

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

package main

import (
	"math/rand"
)

func (vm *VM) exec_clear_screen() {
	for i := 0; i < DISPLAY_HEIGHT; i++ {
		for j := 0; j < DISPLAY_LENGTH/8; j++ {
			vm.display.pixels[i][j] = uint8(0)
		}
	}
	vm.display.dirty = true
}

// (00EE)
// Return from a subroutine
// pop the last address from the stack and set the PC to it.
func (vm *VM) exec_subroutine_return() {
	addr, err := vm.stack.Pop()
	if err != nil {
		// Handle stack underflow (return without CALL)
		panic(err)
	}
	vm.cpu.PC = addr
}

// (1NNN)
// jump to the address in NNN. Sets the PC to NNN.
func (vm *VM) exec_jump(nnn uint16) {
	vm.cpu.PC = nnn
}

// (2NNN)
// Call the subroutine at address NNN.
// Save the address of the next instruction on the stack (PC is already
// advanced by fetch) and then jump to NNN.
func (vm *VM) exec_call_nnn(nnn uint16) {
	// At this point PC already points to the next instruction (because of fetch())
	if err := vm.stack.Push(vm.cpu.PC); err != nil {
		// Handle stack overflow
		panic(err)
	}
	vm.cpu.PC = nnn
}

// (3XNN)
// SE VX, NN : Skip the next instruction if VX == NN
// Compare the value of register VX with NN and if they are equal, increment PC by two
//
//	if VX == NN: PC += 2
func (vm *VM) exec_se_vxnn(x uint16, nn uint16) {
	if vm.cpu.V_registers[x] == uint8(nn) { //EQUALS
		vm.cpu.PC += 2
	}
}

// (4XNN)
// SNE VX, NN: Skip the next instruction if VX != NN.
// Compare the value of register VX with NN and if they are not equal, increment PC by two.
//
// if VX != NN: PC += 2
func (vm *VM) exec_sne_vxnn(x uint16, nn uint16) {
	if vm.cpu.V_registers[x] != uint8(nn) { //NOT EQUALS
		vm.cpu.PC += 2
	}
}

// (5XY0)
// SE VX, VY: Skip the next instruction if VX == VY.
// Compare the value of register VX with the value of VY and if they are equal, increment PC by two.
//
// if VX == VY: PC += 2
func (vm *VM) exec_se_vxvy(x uint16, y uint16) {
	if vm.cpu.V_registers[x] == vm.cpu.V_registers[y] {
		vm.cpu.PC += 2
	}
}

// (6XNN)
// Load the value NN into the register VX.
//
// VX = NN
func (vm *VM) exec_load_vxnn(x uint16, nn uint16) {
	vm.cpu.V_registers[x] = uint8(nn)
}

// (7XNN)
// Add the value NN to the value of register VX and store the result in VX.
//
// VX := VX + NN
func (vm *VM) exec_add_vxnn(x uint16, nn uint16) {
	vm.cpu.V_registers[x] += uint8(nn)
}

// (8XY0)
// Put the value of register VY into VX.
//
// VX = VY
func (vm *VM) exec_load_vxvy(x uint16, y uint16) {
	vm.cpu.V_registers[x] = vm.cpu.V_registers[y]
}

// (8XY1)
// Perform a bitwise OR between the values of VX and VY and store the result in VX.
//
// VX = VX | VY
func (vm *VM) exec_or_vxvy(x uint16, y uint16) {
	vm.cpu.V_registers[x] |= vm.cpu.V_registers[y]
}

// (8XY2)
// Perform a bitwise AND between the values of VX and VY and store the result in VX.
//
// VX = VX | VY
func (vm *VM) exec_and_vxvy(x uint16, y uint16) {
	vm.cpu.V_registers[x] &= vm.cpu.V_registers[y]
}

// (8XY3)
// Perform a bitwise xor between the values of VX and VY and store the result in VX.
//
// VX = VX ^ VY
func (vm *VM) exec_xor_vxvy(x uint16, y uint16) {
	vm.cpu.V_registers[x] ^= vm.cpu.V_registers[y]
}

// (8XY4)
// Add the values of VX and VY and store the result in VX. Put the carry bit in VF (if there is overflow, set VF to 1, otherwise 0).
//
// if VX + VY > 0xFF: VF := 1
// else: VF := 0  VX := VX + VY
func (vm *VM) exec_add_vxvy(x uint16, y uint16) {

	var add_result uint16 = uint16(vm.cpu.V_registers[x]) + uint16(vm.cpu.V_registers[y])

	if add_result > 0xFF { // overflow
		vm.cpu.V_registers[0xF] = uint8(1)
	} else { // no overflow
		vm.cpu.V_registers[0xF] = uint8(0)
	}

	vm.cpu.V_registers[x] = uint8(add_result)
}

// (8XY5)
// Subtract the value of VY from VX and store the result in VX. Put the borrow in VF (if there is borrow, VX > VY, set VF to 1, otherwise 0).
//
// if VX > VY: VF := 1;
// else: VF := 0;
// VX := VX - VY;
func (vm *VM) exec_sub_vxvy(x uint16, y uint16) {
	if vm.cpu.V_registers[x] > vm.cpu.V_registers[y] {
		vm.cpu.V_registers[0xF] = uint8(1)
	} else {
		vm.cpu.V_registers[0xF] = uint8(0)
	}

	vm.cpu.V_registers[x] = vm.cpu.V_registers[x] - vm.cpu.V_registers[y]
}

// (8XY6)
// Shift right, or divide VX by two. Store the least significant bit of VX in VF, and then divide VX and store its value in VX
//
// VF := VX & 0x01 ; VX := VX / 2
func (vm *VM) exec_shr_vxvy(x uint16) {
	vm.cpu.V_registers[0xF] = vm.cpu.V_registers[x] & 0x01
	vm.cpu.V_registers[x] >>= 1
}

// (8XY7)
// Subtract the value of VY from VX and store the result in VX. Set VF to 1 if there is no borrow, to 0 otherwise.
//
// if VY > VX: VF := 1;
// else: VF := 0;
// VX := VY - VX
func (vm *VM) exec_subn_vxvy(x uint16, y uint16) {
	if vm.cpu.V_registers[y] > vm.cpu.V_registers[x] {
		vm.cpu.V_registers[0xF] = uint8(1)
	} else {
		vm.cpu.V_registers[0xF] = uint8(0)
	}

	vm.cpu.V_registers[x] = vm.cpu.V_registers[y] - vm.cpu.V_registers[x]
}

// (8XYE)
// Shift left, or multiply VX by two. Store the most significant bit of VX in VF, and then multiply VX and store its value in VX
//
// VF := MSB(VX); VX <<= 1
func (vm *VM) exec_shl_vxvy(x uint16) {
	vm.cpu.V_registers[0xF] = vm.cpu.V_registers[x] >> 7
	vm.cpu.V_registers[x] <<= 1
}

// (9XY0)
// Skip the next instruction if the values of VX and VY are not equal.
//
// if VX != VY: PC := PC + 2
func (vm *VM) exec_sne_vxvy(x uint16, y uint16) {
	if vm.cpu.V_registers[x] != vm.cpu.V_registers[y] {
		vm.cpu.PC += 2
	}
}

// (ANNN)
// Set the value of I to the address NNN.
//
// I := NNN
func (vm *VM) exec_ldi_nnn(nnn uint16) {
	vm.cpu.I = nnn
}

// (BNNN)
// Jump to the location NNN + V0.
//
// PC := V0 + NNN
func (vm *VM) exec_jmp_v0nnn(nnn uint16) {
	vm.cpu.PC = uint16(vm.cpu.V_registers[0]) + nnn
}

// (CXNN)
// Generate a random byte (from 0 to 255), do a bitwise AND with NN and store the result to VX.
//
// VX := random() & NN
func (vm *VM) exec_rnd_vxnn(x uint16, nn uint16) {
	vm.cpu.V_registers[x] = uint8(rand.Intn(256)) & uint8(nn)
}

// (EX9E)
// Skip the next instruction if the key with the value of VX (the value inside register VX) is currently pressed. Basically, increase PC by two if the key corresponding to the value in VX is pressed.
//
// if keys[VX] == 1: PC := PC + 2

func (vm *VM) exec_skp_vx(x uint16) {
	key_index := vm.cpu.V_registers[x]

	if vm.keypad.keys_state[key_index] {
		vm.cpu.PC += 2
	}
}

// (EXA1)
// Skip the next instruction if the key with the value of VX (the value inside register VX) is currently not pressed. Basically, increase PC by two if the key corresponding to the value in VX is not pressed
//
// if keys[VX] == 0: PC := PC + 2
func (vm *VM) exec_sknp_vx(x uint16) {
	key_index := vm.cpu.V_registers[x]

	if !vm.keypad.keys_state[key_index] {
		vm.cpu.PC += 2
	}
}

// (FX07)
// Read the delay timer register value into VX.
//
// VX := DT
func (vm *VM) exec_load_vx_dt(x uint16) {
	vm.cpu.V_registers[x] = vm.timers.DT
}

// (FX0A)
// Wait for a key press, and then store the value of the key to VX.
//
// K := wait_input() ; VX := K
func (vm *VM) exec_load_vx_k(x uint16) {
	keyPressed := false

	for i := 0; i < 16; i++ {
		if vm.keypad.keys_state[i] {
			vm.cpu.V_registers[x] = uint8(i)
			keyPressed = true
			break
		}
	}

	// If no key is pressed, step back the PC so that this same instruction
	// is re-executed on the next cycle (emulating a blocking wait for input).
	if !keyPressed {
		vm.cpu.PC -= 2
	}
}

// (FX15)
// Load the value of VX into the delay timer DT.
//
// DT := VX
func (vm *VM) exec_load_dt_vx(x uint16) {
	vm.timers.SetDT(vm.cpu.V_registers[x])
}

// (FX18)
// Load the value of VX into the sound time ST.
//
// ST := VX
func (vm *VM) exec_load_st_vx(x uint16) {
	vm.timers.SetST(vm.cpu.V_registers[x])
}

// (FX1E)
// Add the values of I and VX, and store the result in I.
//
// I := I + VX
func (vm *VM) exec_add_i_vx(x uint16) {
	vm.cpu.I += uint16(vm.cpu.V_registers[x])
}

// (FX29)
// Set the location of the sprite for the digit VX to I.
// The font sprites start at address 0x050, and contain the
// hexadecimal digits from 1..F. Each font has a length of
// 0x05 bytes. The memory address for the value in VX is put in I.
//
// I := FONTS_BASE_ADDRESS + VX * 0x05		(0x05 = 5 bytes is number of sprites (bytes) in every font char)
func (vm *VM) exec_load_f_vx(x uint16) {
	vm.cpu.I = FONTS_BASE_ADDRESS + uint16(vm.cpu.V_registers[x])*0x05
}

// (FX33)
// Store the binary-coded decimal representation of VX in memory locations I, I+1 and I+2.
//
// VX is in the range 0..255. We take its decimal representation and split it into
// hundreds, tens and ones:
//
// h := VX / 100
// t := (VX - h * 100) / 10
// o := VX - h * 100 - t * 10
//
// Then we store them at:
// RAM[I]     := h
// RAM[I + 1] := t
// RAM[I + 2] := o
func (vm *VM) exec_load_b_vx(x uint16) {
	vx := vm.cpu.V_registers[x]

	h := vx / 100
	t := (vx - h*100) / 10
	o := vx - h*100 - t*10

	vm.mem.memory[vm.cpu.I] = h
	vm.mem.memory[vm.cpu.I+1] = t
	vm.mem.memory[vm.cpu.I+2] = o
}

// (FX55)
// Store registers from V0 to VX in the main memory,
// starting at location I. Note that X is the number
// of the register, so we can use it in the loop.
// In the following pseudo-code, V[i] allows for
// indexed register access, so that VX == V[X].
// NB: x register is inclusive
//
// for reg in 0..X: RAM[I + reg] := V[reg]
func (vm *VM) exec_load_i_vx(x uint16) {
	for idx := uint16(0); idx <= x; idx++ {
		vm.mem.memory[vm.cpu.I+idx] = vm.cpu.V_registers[idx]
	}
}

// (FX65)
// Load the memory data starting at address I into the registers V0 to VX
// NB: x register is inclusive
//
// for reg in 0..X: V[reg] := RAM[I + reg]
func (vm *VM) exec_load_vx_i(x uint16) {
	for idx := uint16(0); idx <= x; idx++ {
		vm.cpu.V_registers[idx] = vm.mem.memory[vm.cpu.I+idx]
	}
}

// (DXYN)
// DRW VX, VY, N
// Draws an N-byte tall sprite starting at memory location I
// at coordinates (VX, VY). Each sprite byte encodes a row of 8 pixels.
// Pixels are XORed onto the display. VF is set to 1 if any pixel
// is unset as a result of this operation (collision), otherwise 0.
// This implementation wraps around the screen edges.
func (vm *VM) exec_draw_vxvy_n(x, y, n uint16) {
	vx := vm.cpu.V_registers[x]
	vy := vm.cpu.V_registers[y]

	// Reset collision flag
	vm.cpu.V_registers[0xF] = 0

	for row := uint16(0); row < n; row++ {
		spriteByte := vm.mem.memory[vm.cpu.I+row]

		// Y position with vertical wrap
		yPos := (uint16(vy) + row) % DISPLAY_HEIGHT
		yIdx := int(yPos)

		for bit := uint16(0); bit < 8; bit++ {
			// Check current bit of the sprite row (MSB is leftmost)
			if (spriteByte & (0x80 >> bit)) == 0 {
				continue
			}

			// X position with horizontal wrap
			xPos := (uint16(vx) + bit) % DISPLAY_LENGTH
			xIdx := int(xPos)

			byteIdx := xIdx / 8              // which byte in the row
			bitOffset := uint8(xIdx % 8)     // which bit in that byte (0..7)
			mask := uint8(0x80 >> bitOffset) // MSB-first mapping

			// Collision: pixel was set and will be toggled off
			if vm.display.pixels[yIdx][byteIdx]&mask != 0 {
				vm.cpu.V_registers[0xF] = 1
			}

			// XOR the pixel
			vm.display.pixels[yIdx][byteIdx] ^= mask
		}
	}

	// Mark display as dirty so the framebuffer will be rebuilt
	vm.display.dirty = true
}

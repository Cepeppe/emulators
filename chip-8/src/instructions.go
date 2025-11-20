package main

import (
	"math/rand"
)

// (00EE)
// Return from a subroutine
// pop the last address from the stack and set the PC to it.
func (vm *VM) exec_subroutine_return() {
	vm.cpu.PC, _ = vm.stack.Pop()
}

// (1NNN)
// jump to the address in NNN. Sets the PC to NNN.
func (vm *VM) exec_jump(nnn uint16) {
	vm.cpu.PC = nnn
}

// (2NNN)
// Call the subroutine at address NNN. It increments SP, puts the current PC at the top of the stack and sets PC to the address NNN.
func (vm *VM) exec_call_nnn(nnn uint16) {
	vm.cpu.PC += 2
	vm.stack.Push(vm.cpu.PC)
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
// VF := VX & 0x80 ; VX := VX * 2
func (vm *VM) exec_shl_vxvy(x uint16) {
	vm.cpu.V_registers[0xF] = vm.cpu.V_registers[x] & 0x80
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
// Skip the next instruction if the key with the value of VX is currently pressed. Basically, increase PC by two if the key corresponding to the value in VX is pressed.
//
// if keys[VX] == 1: PC := PC + 2

func (vm *VM) exec_skp_vx(x uint16){
	if vm.keypad.keys_state[x]{
		vm.cpu.PC+=2
	}
}

// (EXA1)
// Skip the next instruction if the key with the value of VX is currently not pressed. Basically, increase PC by two if the key corresponding to the value in VX is not pressed
//
// if keys[VX] == 0: PC := PC + 2
func (vm *VM) exec_sknp_vx(x uint16){
	if !vm.keypad.keys_state[x]{
		vm.cpu.PC+=2
	}
}

// (FX07)
// Read the delay timer register value into VX.
//
// VX := DT
func (vm *VM) exec_load_vx_dt(x uint16){
	vm.cpu.V_registers[x]=vm.timers.DT
}

// (FX0A)
// Wait for a key press, and then store the value of the key to VX.
//
// K := wait_input() ; VX := K
func (vm *VM) exec_load_vx_k(x uint16){
	var key uint8
	var end bool = false
	for !end{
		for i:=0; i<16; i++{
			if vm.keypad.keys_state[i]{
				key = uint8(i)
				end = true
				break
			}
		} 
	}
	vm.cpu.V_registers[x]=key
}

// (FX15)
// Load the value of VX into the delay timer DT.
//
// DT := VX
func (vm *VM) exec_load_dt_vx(x uint16){
	vm.timers.DT = vm.cpu.V_registers[x]
}

// (FX18)
// Load the value of VX into the sound time ST.
//
// ST := VX
func (vm *VM) exec_load_st_vx(x uint16){
	vm.timers.ST = vm.cpu.V_registers[x]
}

// (FX1E)
// Add the values of I and VX, and store the result in I.
//
// I := I + VX
func (vm *VM) exec_add_i_vx(x uint16){
	vm.cpu.I += uint16(vm.cpu.V_registers[x])
}

///////////////////////

// DXYN TODO

func (vm *VM) exec_DXYN(X, Y, N uint) {
	/*
		// values of inside registers at indexes X, Y
		// I register contains needed sprite memory address

		// cpu register VF is used to detect collisions, we set it to 0 before starting to draw
		vm.cpu.V_registers[0xF] = uint8(0)

		x_coord := vm.cpu.V_registers[X]
		y_coord := vm.cpu.V_registers[Y]
		I := vm.cpu.I

		sprites_to_draw := make([]byte, N)

		//save sprites to be drawn in slice "sprites_to_draw"
		for curr_mem_addr := I; curr_mem_addr < I+uint16(N); curr_mem_addr++ {
			sprites_to_draw[curr_mem_addr] = vm.mem.memory[curr_mem_addr]
		}

		//TODO: IMPLEMENT
	*/
}

// calculate the new value for the screen pixel in that byte applying XOR
// returns new pixel values and a flag that is true if collisions happened
func apply_xor_verify_collision(screen uint8, sprite uint8) (screen_byte_post_xor uint8, collision_happened bool) {

	//screen: 			0 0 1 0 1 0 1 1
	//sprite: 			1 0 0 0 1 1 0 1
	//collision: 		0 0 0 0 1 0 0 1  if screen & sprite (bitwise and) is not 0 then collision, else no collision
	//screen post xor:  1 0 1 0 0 1 1 0

	flag := false

	if screen&sprite != 0 {
		flag = true
	}
	return screen ^ sprite, flag
}

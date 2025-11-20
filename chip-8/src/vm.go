package main

import (
	"errors"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// Conceptual reference:
// type Instruction uint16
// type Address uint16

// VM holds the state of the CHIP-8 virtual machine for the Ebiten loop.
type VM struct {
	mem     MEM
	cpu     CPU
	stack   STACK
	display DISPLAY
	keypad  KEYPAD
	timers  TIMERS
}

// Update is called approximately 60 times per second.
// CPU steps, timers, and input handling will be placed here.
func (vm *VM) Update() error {

	vm.keypad.updatePressedKeys()
	if vm.timers.shouldBeep() {
		PlayBeep()
	}

	return nil
}

// Draw is called every frame to render the current display state.
func (vm *VM) Draw(screen *ebiten.Image) {
	vm.display.Draw(screen)
}

// Layout defines the logical resolution of the window.
func (vm *VM) Layout(outsideWidth, outsideHeight int) (int, int) {
	return DISPLAY_LENGTH * DISPLAY_SCALE, DISPLAY_HEIGHT * DISPLAY_SCALE
}

func NewVM(romPath string) (*VM, error) {
	var vm VM

	fmt.Println("chip-8 emulator booting..")

	// Load sprites into memory
	vm.mem.LoadSprites()

	// Load game from ROM(file system) to RAM
	if err := vm.mem.LoadProgramFromRom(romPath); err != nil {
		return nil, err
	}

	// Initialize cpu
	vm.cpu.Init()

	// Initialize stack
	vm.stack.Init()

	// Initialize keypad
	vm.keypad.Init()

	// Initialize display to all black
	vm.display.Init()

	// Initialize Timers
	vm.timers.Init()

	vm.Dump()

	return &vm, nil
}

func (vm *VM) ExecuteNextInstruction() error {
	instr, err := vm.fetch()
	if err != nil {
		return nil
	}

	return vm.decode_execute(instr)
}

// read 2 bytes from memory at address PC; PC = PC + 2
func (vm *VM) fetch() (uint16, error) {

	var pc uint16 = vm.cpu.PC

	var instruction uint16 = uint16(vm.mem.memory[pc])<<8 | uint16(vm.mem.memory[pc+1])

	vm.cpu.PC = pc + 2

	return instruction, nil
}

/*
Decodes and executes the instruction (calling appropriate method based on the decoded opcode)
Before calling appropriate method, decode phase will resolve actual values that must be passed to the method
Commands are decoded using the following rules:

op_code: first 4 bits (first nibble). 4 bits is a nibble
X: The second nibble. Used to look up one of the 16 registers (VX) from V0 through VF.
Y: The third nibble. Also used to look up one of the 16 registers (VY) from V0 through VF.
N: The fourth nibble. A 4-bit number.
NN: The second byte (third and fourth nibbles). An 8-bit immediate number.
NNN: The second, third and fourth nibbles. A 12-bit immediate memory address.

ex:
ooooxxxx yyyynnnn
oooo---- nnnnnnnn
oooonnnn nnnnnnnn
*/
func (vm *VM) decode_execute(instr uint16) error {
	op_code, nnn := isolateBits(instr, 4)
	x, nn := isolateBits(nnn, 4+4)
	y, n := isolateBits(nn, 4+4+4)

	// we save opcode in vm.current_opcode for diagnostic purposes only, for example when dumping cpu
	// state in case of crash/failure cpu.opcode is unused except that for this purpose but it's a good
	// choice with minimal tradeoff to have memorized in a fixed place the opcode that possibly broke the house.
	// And it's just 2 bytes.
	vm.cpu.curr_opcode = op_code

	switch op_code {
	case 0x0: //Execute machine language routine: must not be implemented in our case
		if nnn == 0x0E0 { // 00E0: Clear screenPermalink: It should clear the display, turning all pixels off to 0.

		} else if nnn == 0x0EE { // (00EE) Return from a subroutine
			// pop the last address from the stack and set the PC to it.
			vm.exec_subroutine_return()
		} else {
			// 0NNN: Execute machine language routine, must not be implemented in our case

		}
	case 0x1: // JUMP (1NNN): Jump to the address in NNN. Sets the PC to NNN.
		vm.exec_jump(nnn)

	case 0x2: // CALL (2NNN) subroutine: Call the subroutine at address NNN. It increments SP, puts the current PC at the top of the stack and sets PC to the address NNN.
		vm.exec_call_nnn(nnn)

	case 0x3: // CONDITIONAL JUMP EQUALITY (3XNN): SE VX, NN Skip the next instruction if VX == NN
		vm.exec_se_vxnn(x, nn)

	case 0x4: // CONDITIONAL JUMP NON-EQUALITY (4XNN): SE VX, NN Skip the next instruction if VX == NN
		vm.exec_sne_vxnn(x, nn)

	case 0x5: // CONDITIONAL JUMP EQUALITY BETWEEN REGISTERS (5XY0): SE VX, VY Skip the next instruction if VX == VY. Compare the value of register VX with the value of VY and if they are equal, increment PC by two.
		vm.exec_se_vxvy(x, y)

	case 0x6: // LOAD (6XNN): LD VX, NN Load the value NN into the register VX.
		vm.exec_load_vxnn(x, nn)

	case 0x7: // ADD (7XNN): ADD VX, NN Add the value NN to the value of register VX and store the result in VX.
		vm.exec_add_vxnn(x, nn)

	case 0x8:
		if n == 0x0 { // LOAD (8XY0): LD VX, VY Put the value of register VY into VX.
			vm.exec_load_vxvy(x, y)
		} else if n == 0x1 { // OR (8XY1): OR VX, VY Perform a bitwise OR between the values of VX and VY and store the result in VX.
			vm.exec_or_vxvy(x, y)
		} else if n == 0x2 { // AND (8XY2): AND VX, VY Perform a bitwise AND between the values of VX and VY and store the result in VX.
			vm.exec_and_vxvy(x, y)
		} else if n == 0x3 { // XOR (8XY3): XOR VX, VY Perform a bitwise XOR between the values of VX and VY and store the result in VX.
			vm.exec_xor_vxvy(x, y)
		} else if n == 0x4 { // ADD (8XY4): ADD VX, VY Add the values of VX and VY and store the result in VX. Put the carry bit in VF (if there is overflow, set VF to 1, otherwise 0).
			vm.exec_add_vxvy(x, y)
		} else if n == 0x5 { // SUB (8XY5): SUB VX, VY Subtract the value of VY from VX and store the result in VX. Put the borrow in VF (if there is borrow, VX > VY, set VF to 1, otherwise 0).
			vm.exec_sub_vxvy(x, y)
		} else if n == 0x6 { // SHR (8XY6): SHR VX {, VY} Shift right, or divide VX by two. Store the least significant bit of VX in VF, and then divide VX and store its value in VX
			vm.exec_shr_vxvy(x)
		} else if n == 0x7 { // SUBN (8XY7): SUBN VX, VY Subtract the value of VY from VX and store the result in VX. Set VF to 1 if there is no borrow, to 0 otherwise.
			vm.exec_subn_vxvy(x, y)
		} else if n == 0xE { // SHL (8XYE): SHL VX {, VY} Shift left, or multiply VX by two. Store the most significant bit of VX in VF, and then multiply VX and store its value in VX
			vm.exec_shl_vxvy(x)
		}

	case 0x9: //SNE (9XY0): SNE VX, VY Skip the next instruction if the values of VX and VY are not equal.
		vm.exec_sne_vxvy(x, y)

	case 0xA: // LOAD-I (ANNN): LD I, NNN Set the value of I to the address NNN.
		vm.exec_ldi_nnn(nnn)

	case 0xB: // JMP (BNNN): JMP V0, NNN Jump to the location NNN + V0.
		vm.exec_jmp_v0nnn(nnn)

	case 0xC: // RND (CXNN): RND VX, NN Generate a random byte (from 0 to 255), do a bitwise AND with NN and store the result to VX.
		vm.exec_rnd_vxnn(x, nn)

	case 0xD: // DRAW (DXYN): DRW VX, VY, N Draw instruction. For algorithm documentation see the impl or the pdf about dxyn
		//TODO: IMPL

	case 0xE:
		if nn == 0x9E { // SKIP PRESSED (EX9E): SKP VX Skip the next instruction if the key with the value of VX is currently pressed. Basically, increase PC by two if the key corresponding to the value in VX is pressed.
			vm.exec_skp_vx(x)
		} else if nn == 0xA1 { // SKIP NOT PRESSED (EXA1): SKNP VX Skip the next instruction if the key with the value of VX is currently not pressed. Basically, increase PC by two if the key corresponding to the value in VX is not pressed
			vm.exec_sknp_vx(x)
		}

	case 0xF:
		if nn == 0x07 { // LOAD (FX07): LD VX, DT Read the delay timer register value into VX.
			vm.exec_load_vx_dt(x)
		} else if nn == 0x0A { // LOAD (FX0A): LD VX, K Wait for a key press, and then store the value of the key to VX.
			vm.exec_load_vx_k(x)
		} else if nn == 0x15 { // LOAD (FX15): LD DT, VX Load the value of VX into the delay timer DT.
			vm.exec_load_dt_vx(x)
		} else if nn == 0x18 { // LOAD (FX18): LD ST, VX Load the value of VX into the sound time ST.
			vm.exec_load_st_vx(x)
		} else if nn == 0x1E { // ADDI (FX1E): ADD I, VX Add the values of I and VX, and store the result in I.
			vm.exec_add_i_vx(x)
		} else if nn == 0x29 {

		} else if nn == 0x33 {

		} else if nn == 0x55 {

		}

	default:
		return errors.New("invalid opcode: " + fmt.Sprint(op_code))
	}

	return nil
}

func execDisplay(vm *VM) error {
	return nil
}

func (vm *VM) Dump() {
	vm.mem.Dump(0x00, 0x1000) //0, 4096
	vm.cpu.Dump()
	vm.stack.Dump()
	vm.keypad.Dump()
	vm.display.Dump()
	vm.timers.Dump()
}

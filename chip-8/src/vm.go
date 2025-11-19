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
}

// Update is called approximately 60 times per second.
// CPU steps, timers, and input handling will be placed here.
func (g *VM) Update() error {
	return nil
}

// Draw is called every frame to render the current display state.
func (g *VM) Draw(screen *ebiten.Image) {
	g.display.Draw(screen)
}

// Layout defines the logical resolution of the window.
func (g *VM) Layout(outsideWidth, outsideHeight int) (int, int) {
	return DISPLAY_LENGTH * DISPLAY_SCALE, DISPLAY_HEIGHT * DISPLAY_SCALE
}

func NewVM(romPath string) (*VM, error) {
	var vm VM

	fmt.Println("chip-8 emulator booting..")

	// Load font
	vm.mem.LoadFont()

	// Load game from ROM(file system) to RAM
	if err := vm.mem.LoadProgramFromRom(romPath); err != nil {
		return nil, err
	}

	vm.mem.Dump(0x000, 4096)

	// Initialize cpu
	vm.cpu.Init()

	// Initialize stack
	vm.stack.Init()

	// Initialize keypad
	vm.keypad.Init()

	// Initialize display to all black
	vm.display.Init()

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
	//x, nn := isolateBits(nnn, 4+4)
	//y, n := isolateBits(nn, 4+4+4)

	switch op_code {
	case 0x0: //Execute machine language routine: must not be implemented in our case
		if nnn == 0x0E0 { // 00E0: Clear screenPermalink: It should clear the display, turning all pixels off to 0.

		} else if nnn == 0x0EE { //00EE: Return from a subroutine
			//pop the last address from the stack and set the PC to it.

		} else {
			// 0NNN: Execute machine language routine, must not be implemented in our case

		}
	case 0x1:

	case 0xD:

	default:
		return errors.New("invalid opcode")
	}

	return nil
}

func execDisplay(vm *VM) error {
	return nil
}

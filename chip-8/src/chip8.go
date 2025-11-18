package main

import (
	"fmt"
	"log"
	"os"

	c8 "github.com/Cepeppe/chip8/chip8"
	"github.com/hajimehoshi/ebiten/v2"
)

// VM holds the state of the CHIP-8 virtual machine for the Ebiten loop.
type VM struct {
	memory  c8.MEM
	cpu     c8.CPU
	stack   c8.STACK
	display c8.DISPLAY
	keypad  c8.KEYPAD
}

func NewVM(romPath string) (*VM, error) {
	var vm VM

	fmt.Println("chip-8 emulator booting..")

	// Load font
	vm.memory.LoadFont()

	// Load game from ROM(file system) to RAM
	if err := vm.memory.LoadProgramFromRom(romPath); err != nil {
		return nil, err
	}

	vm.memory.Dump(0x000, 4096)

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
	return c8.DISPLAY_LENGTH * c8.DISPLAY_SCALE, c8.DISPLAY_HEIGHT * c8.DISPLAY_SCALE
}

func main() {
	romPath := "..\\rom_games\\RPS.ch8" // ROM path can be set here (e.g., "../rom_games/PONG.ch8").

	game, err := NewVM(romPath)
	if err != nil {
		fmt.Println("Error while trying to load program:", err.Error())
		os.Exit(1)
	}

	// Configure the window size and title.
	ebiten.SetWindowSize(c8.DISPLAY_LENGTH*c8.DISPLAY_SCALE, c8.DISPLAY_HEIGHT*c8.DISPLAY_SCALE)
	ebiten.SetWindowTitle("CHIP-8 Emulator")

	// Start the Ebiten game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

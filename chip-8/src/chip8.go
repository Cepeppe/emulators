package main

import (
	"fmt"
	"log"
	"os"

	c8 "github.com/Cepeppe/chip8/chip8"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game holds the state of the CHIP-8 virtual machine for the Ebiten loop.
type Game struct {
	memory  c8.MEM
	cpu     c8.CPU
	stack   c8.STACK
	display c8.DISPLAY
	keypad  c8.KEYPAD
}

func NewGame(romPath string) (*Game, error) {
	var g Game

	fmt.Println("chip-8 emulator booting..")

	// Load font
	g.memory.LoadFont()

	// Load game from ROM(file system) to RAM
	if err := g.memory.LoadProgramFromRom(romPath); err != nil {
		return nil, err
	}

	// Initialize cpu
	g.cpu.Init()

	// Initialize stack
	g.stack.Init()

	// Initialize keypad
	g.keypad.Init()

	// Initialize display to all black
	g.display.Init()

	return &g, nil
}

// Update is called approximately 60 times per second.
// CPU steps, timers, and input handling will be placed here.
func (g *Game) Update() error {
	return nil
}

// Draw is called every frame to render the current display state.
func (g *Game) Draw(screen *ebiten.Image) {
	g.display.Draw(screen)
}

// Layout defines the logical resolution of the window.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return c8.DISPLAY_LENGTH * c8.DISPLAY_SCALE, c8.DISPLAY_HEIGHT * c8.DISPLAY_SCALE
}

func main() {
	romPath := "..\\rom_games\\RPS.ch8" // ROM path can be set here (e.g., "../rom_games/PONG.ch8").

	game, err := NewGame(romPath)
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

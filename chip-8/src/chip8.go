package main

import (
	"fmt"
	"log"
	"os"

	c8 "github.com/Cepeppe/chip8/chip8"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	romPath := "..\\rom_games\\RPS.ch8" // ROM path can be set here (e.g., "../rom_games/PONG.ch8").

	game, err := c8.NewVM(romPath)
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

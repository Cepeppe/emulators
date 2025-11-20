package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	romPath := "..\\rom_games\\RPS.ch8" // ROM path can be set here (e.g., "../rom_games/PONG.ch8").

	vm, err := NewVM(romPath)
	if err != nil {
		fmt.Println("Error while trying to load program:", err.Error())
		os.Exit(1)
	}
	
	if err := InitAudio(); err != nil {
        fmt.Println("An error didn't allow to init audio: ", err)
    }

	// Configure the window size and title.
	ebiten.SetWindowSize(DISPLAY_LENGTH*DISPLAY_SCALE, DISPLAY_HEIGHT*DISPLAY_SCALE)
	ebiten.SetWindowTitle("CHIP-8 Emulator")

	// Start the Ebiten game loop.
	if err := ebiten.RunGame(vm); err != nil {
		log.Fatal(err)
	}

	
}

package main

import (
	"fmt"
	"os"
	"github.com/Cepeppe/chip8/virtual_machine"
)

func main() {
	fmt.Println("chip-8 emulator booting..")
	var memory chip8.MEM
	var cpu chip8.CPU
	var stack chip8.STACK
	var display chip8.DISPLAY
	var keypad chip8.KEYPAD

	var chosen_game string = ""

	//Load font
	memory.LoadFont()
	//Load game from ROM(file system) to RAM
	err := memory.LoadProgramFromRom(chosen_game)
	if err != nil {
		fmt.Println("Error while trying to load program", err.Error())
		os.Exit(1)
	}

	//Initialize cpu
	cpu.Init()

	//Initialize stack
	stack.Init()

	//Initialize keypad
	keypad.Init()

	//Initialize display to all black
	display.Init()

	//Draw
	display.Draw()

	for {

	}

}

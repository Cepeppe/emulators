package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	romPath, err := selectGame("../rom_games")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

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

// selectROM lists all ROM files in romDir, prints them on the console and
// lets the user choose one by entering its number. It returns the full path
// to the selected ROM file.
func selectGame(romDir string) (string, error) {
	entries, err := os.ReadDir(romDir)
	if err != nil {
		return "", fmt.Errorf("failed to read ROM directory %q: %w", romDir, err)
	}

	var roms []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()

		// Optional: skip hidden files
		if strings.HasPrefix(name, ".") {
			continue
		}

		// Optional: filter by extension, e.g. .ch8
		// if filepath.Ext(name) != ".ch8" {
		//     continue
		// }

		roms = append(roms, name)
	}

	if len(roms) == 0 {
		return "", fmt.Errorf("no ROM files found in %s", romDir)
	}

	fmt.Println("Available games:")
	for i, name := range roms {
		fmt.Printf("  %d) %s\n", i+1, name)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Select a game [1-%d]: ", len(roms))

		line, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read selection: %w", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		idx, err := strconv.Atoi(line)
		if err != nil || idx < 1 || idx > len(roms) {
			fmt.Println("Invalid selection, try again.")
			continue
		}

		chosenName := roms[idx-1]
		fullPath := filepath.Join(romDir, chosenName)
		return fullPath, nil
	}

}

func terminal_cls() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

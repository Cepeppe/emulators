package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

/*

	The original implementation of the Chip-8 language used a
	64x32-pixel monochrome display with this format:

	+-------------------+
	|(0,0)	     (63,0) |
	|                   |
	|(0,31)	     (63,31)|
    +-------------------+

	Chip-8 draws graphics on screen through the use of sprites.
	A sprite is a group of bytes which are a binary representation of the desired picture.
	Chip-8 sprites may be up to 15 bytes, for a possible sprite size of 8x15.

	Programs may also refer to a group of sprites representing the hexadecimal
	digits 0 through F. These sprites are 5 bytes long, or 8x5 pixels.
	The data should be stored in the interpreter area of Chip-8 memory (0x000 to 0x1FF).
	Below is a listing of each character's bytes

	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80  // F

*/

const (
	DISPLAY_LENGTH = 64
	DISPLAY_HEIGHT = 32
	DISPLAY_SCALE  = 10
)

type DISPLAY struct {
	// We use DISPLAY_LENGTH/8 because in each one of
	// those bytes a single bit decodes a column value for that row
	pixels [DISPLAY_HEIGHT][DISPLAY_LENGTH / 8]uint8

	scale int           // es. 10 -> 64x32 becomes 640x320
	fb    *ebiten.Image // internal scaled framebuffer
	dirty bool          // true when pixels change
}

// Initialize display to all black (0 everywhere)
func (display *DISPLAY) Init() {
	for row_idx, row := range display.pixels {
		for col_idx := range row {
			display.pixels[row_idx][col_idx] = uint8(0)
		}
	}
	display.scale = DISPLAY_SCALE
	width := DISPLAY_LENGTH * display.scale
	height := DISPLAY_HEIGHT * display.scale
	display.fb = ebiten.NewImage(width, height)

	display.fb.Fill(color.Black)

	display.dirty = true
}

// Draws display.pixels on the actual display.
// Should be called from the Game.Draw(screen) method.
func (display *DISPLAY) Draw(screen *ebiten.Image) {
	// If the framebuffer (fb) does not exist yet for any reason, create it now.
	if display.fb == nil {
		width := DISPLAY_LENGTH * display.scale
		height := DISPLAY_HEIGHT * display.scale
		// Create a new Ebiten image for the framebuffer with the scaled dimensions.
		display.fb = ebiten.NewImage(width, height)
		// Initialize the framebuffer with a black background.
		display.fb.Fill(color.Black)
		// Mark the display as dirty to ensure it's drawn in the current frame.
		display.dirty = true
	}

	// If the logical content of the pixels has changed, the framebuffer must be rebuilt.
	if display.dirty {
		// Set the background to black.
		display.fb.Fill(color.Black)

		// Scan the CHIP-8 pixels (which are bit-packed in bytes).
		for y := 0; y < DISPLAY_HEIGHT; y++ {
			for byteIdx := 0; byteIdx < DISPLAY_LENGTH/8; byteIdx++ {
				b := display.pixels[y][byteIdx]
				if b == 0 {
					continue // No lit pixels in this 8-pixel block.
				}

				// Each bit of the byte represents a pixel (from left to right).
				for bit := 0; bit < 8; bit++ {
					// Create a mask to check the current bit (0x80 = 10000000).
					mask := uint8(0x80 >> bit)
					if b&mask == 0 {
						continue // The current pixel bit is 0 (off).
					}

					// Calculate the logical x-coordinate (column 0..63).
					x := byteIdx*8 + bit

					// Draw the scaled pixel at (x, y) onto the framebuffer.
					// Calculate the starting coordinates for the scaled block.
					startX := x * display.scale
					startY := y * display.scale

					// Draw the pixel block (scaling loop).
					for dy := 0; dy < display.scale; dy++ {
						for dx := 0; dx < display.scale; dx++ {
							// Set the corresponding pixel in the framebuffer to white.
							display.fb.Set(
								startX+dx,
								startY+dy,
								color.White,
							)
						}
					}
				}
			}
		}

		// The framebuffer has been updated, so mark the display as clean.
		display.dirty = false
	}

	// Draw the framebuffer (which is already scaled) onto the Ebiten screen.
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(display.fb, op)
}

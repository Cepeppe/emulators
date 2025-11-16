package chip8

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
)

type DISPLAY struct {
	// We use DISPLAY_LENGTH/8 because in each one of
	// those bytes a single bit decodes a column value for that row
	screen [DISPLAY_HEIGHT][DISPLAY_LENGTH / 8]uint8
}

//Initialize display to all black (0 everywhere)
func (display *DISPLAY) Init() {
	for row_idx, row := range display.screen {
		for col_idx := range row {
			display.screen[row_idx][col_idx] = uint8(0)
		}
	}
}

//Draws display.screen on the actual display
func (display *DISPLAY) Draw() {

}

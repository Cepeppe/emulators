package chip_8

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
	Below is a listing of each character's bytes, in binary and hexadecimal:


*/

const (
	DISPLAY_LENGTH = 64
	DISPLAY_HEIGHT = 32
)

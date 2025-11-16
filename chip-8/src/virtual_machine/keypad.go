package chip8

/*
	The computers which originally used the Chip-8 Language had
	a 16-key hexadecimal keypad with the following layout:

	1	2	3	C
	4	5	6	D
	7	8	9	E
	A	0	B	F

	This layout must be mapped into various other configurations
	to fit the keyboards of today's platforms.

	HEX based keypad (0x0-0xF), it's possible to use
	an array to store the current state of the key.
*/

type KEYPAD struct {
	keys_state [16]uint8
}

func (k *KEYPAD) Init() {

}

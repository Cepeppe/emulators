package chip_8

/*
	The stack is an array of 16 16-bit values, used to store the address 
	that the interpreter shoud return to when finished with a subroutine. 
	Chip-8 allows for up to 16 levels of nested subroutines.
*/

var stack [16]uint16
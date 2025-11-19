package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

/*
	The computers which originally used the Chip-8 Language had
	a 16-key hexadecimal keypad with the following layout:

	   DEFAULT						MINE
	1	2	3	C				1	2	3	4
	4	5	6	D		==>		Q	W	E	R
	7	8	9	E		==>		A	S	D	F
	A	0	B	F				Z	X	C	V

	This layout must be mapped into various other configurations
	to fit the keyboards of today's platforms.

	HEX based keypad (0x0-0xF), it's possible to use
	an array to store the current state of the key.
*/

const (
	// THESE CONST MEMORIZES THE INDEX MAPPING INFO ABOUT EVERY
	// DEFAULT KEYPAD KEY THAT IS BEING PRESSED, ALLOWING TO USE
	// "MINE" KEYPAD SETUP

	// "MINE" keypad will be correctly mapped by updatePressedKeys

	KEY_1 = 0
	KEY_2 = 1
	KEY_3 = 2
	KEY_C = 3

	KEY_4 = 4
	KEY_5 = 5
	KEY_6 = 6
	KEY_D = 7

	KEY_7 = 8
	KEY_8 = 9
	KEY_9 = 10
	KEY_E = 11

	KEY_A = 12
	KEY_0 = 13
	KEY_B = 14
	KEY_F = 15
)

type KEYPAD struct {
	keys_state [16]bool
}

func (k *KEYPAD) Init() {
	for i := 0; i < 16; i++ {
		k.keys_state[i] = false
	}
}

func (k *KEYPAD) updatePressedKeys() {
	/*
		   DEFAULT					    MINE
		1	2	3	C				1	2	3	4
		4	5	6	D		==>		Q	W	E	R
		7	8	9	E		==>		A	S	D	F
		A	0	B	F				Z	X	C	V
	*/

	// key on default keypad		key on "MINE" keypad

	k.keys_state[KEY_1] = ebiten.IsKeyPressed(ebiten.Key1)
	k.keys_state[KEY_2] = ebiten.IsKeyPressed(ebiten.Key2)
	k.keys_state[KEY_3] = ebiten.IsKeyPressed(ebiten.Key3)
	k.keys_state[KEY_C] = ebiten.IsKeyPressed(ebiten.Key4)

	k.keys_state[KEY_4] = ebiten.IsKeyPressed(ebiten.KeyQ)
	k.keys_state[KEY_5] = ebiten.IsKeyPressed(ebiten.KeyW)
	k.keys_state[KEY_6] = ebiten.IsKeyPressed(ebiten.KeyE)
	k.keys_state[KEY_D] = ebiten.IsKeyPressed(ebiten.KeyR)

	k.keys_state[KEY_7] = ebiten.IsKeyPressed(ebiten.KeyA)
	k.keys_state[KEY_8] = ebiten.IsKeyPressed(ebiten.KeyS)
	k.keys_state[KEY_9] = ebiten.IsKeyPressed(ebiten.KeyD)
	k.keys_state[KEY_E] = ebiten.IsKeyPressed(ebiten.KeyF)

	k.keys_state[KEY_A] = ebiten.IsKeyPressed(ebiten.KeyZ)
	k.keys_state[KEY_0] = ebiten.IsKeyPressed(ebiten.KeyX)
	k.keys_state[KEY_B] = ebiten.IsKeyPressed(ebiten.KeyC)
	k.keys_state[KEY_F] = ebiten.IsKeyPressed(ebiten.KeyV)
}

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
	// CHIP-8 native keypad layout:
	// 1  2  3  C
	// 4  5  6  D
	// 7  8  9  E
	// A  0  B  F

	KEY_1 = 0x1
	KEY_2 = 0x2
	KEY_3 = 0x3
	KEY_C = 0xC

	KEY_4 = 0x4
	KEY_5 = 0x5
	KEY_6 = 0x6
	KEY_D = 0xD

	KEY_7 = 0x7
	KEY_8 = 0x8
	KEY_9 = 0x9
	KEY_E = 0xE

	KEY_A = 0xA
	KEY_0 = 0x0
	KEY_B = 0xB
	KEY_F = 0xF
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
	   CHIP-8 keypad          my keypad layout
	   1   2   3   C           1   2   3   4
	   4   5   6   D    ==>    Q   W   E   R
	   7   8   9   E           A   S   D   F
	   A   0   B   F           Z   X   C   V
	*/

	// key CHIP-8      ->      tasto fisico PC

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

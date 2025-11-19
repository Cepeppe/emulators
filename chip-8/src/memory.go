package main

import (
	"fmt"
	"os"
)

/*
	Memory Map:
+---------------+= 0xFFF (4095) End of Chip-8 RAM
|               |
|               |
| 0x200 to 0xFFF|
|     Chip-8    |
| Program / Data|
|     Space     |
|               |
+---------------+= 0x200 (512) Start of most Chip-8 programs
|               |= 0x1FF (511) End of reserved area
| 0x000 to 0x1FF|
| Reserved for  |
|  interpreter  |
+---------------+= 0x000 (0) Start of Chip-8 RAM

0x000-0x1FF - Chip 8 interpreter (contains font set in emu) -> not used
0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
0x200-0xFFF - Program ROM and work RAM

*/

const MEMORY_SPACE_BYTES uint16 = 4096 //4kb

const MEMORY_START uint16 = 0x000         //0
const RESERVED_MEMORY_END uint16 = 0x1FF  //511
const PROGRAM_MEMORY_START uint16 = 0x200 //512
const MEMORY_END uint16 = 0xFFF           //4095

type MEM struct {
	memory [MEMORY_SPACE_BYTES]byte
}

func (mem *MEM) LoadProgramFromRom(file string) error {
	bytes_arr, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	for idx := range bytes_arr {
		mem.memory[PROGRAM_MEMORY_START+uint16(idx)] = bytes_arr[idx]
	}
	return nil
}

// startAddress is inclusive, endAddress is exclusive
func (mem *MEM) Dump(startAddress uint16, endAddress uint16) {

	var effective_end = endAddress

	if effective_end >= MEMORY_SPACE_BYTES {
		effective_end = MEMORY_SPACE_BYTES
	}

	fmt.Printf("\n	MEM DUMP FROM 0x%x to 0x%x\n", startAddress, endAddress)
	fmt.Printf("\n\t+------------------------------------------------------------------------------------------------------------------------------------+\n")

	for i := startAddress; i < effective_end; i++ {
		if i%32 == 0 {
			fmt.Printf("\t| %04d -> 0x%03x |", i, i)
		}

		if i%2 == 0 {
			fmt.Printf(" ")
		}

		fmt.Printf("%02x", mem.memory[i])

		if i%32 == 31 {
			fmt.Printf(" | ")
			for j := i - (i % 32); j <= i; j++ {
				car := uint8(mem.memory[j])

				if car == 0 {
					fmt.Printf(" ")
				} else if car < 33 || car > 126 {
					fmt.Printf("·")
				} else {
					fmt.Printf("%s", string(car))
				}
			}
			fmt.Println(" |")
		}

	}
	fmt.Printf("\t+------------------------------------------------------------------------------------------------------------------------------------+\n")
}

// Loads sprites that are used to draw the screen
func (m *MEM) LoadSprites() {

	// 0
	m.memory[0x050] = 0xF0 // 11110000
	m.memory[0x051] = 0x90 // 10010000
	m.memory[0x052] = 0x90 // 10010000
	m.memory[0x053] = 0x90 // 10010000
	m.memory[0x054] = 0xF0 // 11110000

	// 1
	m.memory[0x055] = 0x20 // 00100000
	m.memory[0x056] = 0x60 // 01100000
	m.memory[0x057] = 0x20 // 00100000
	m.memory[0x058] = 0x20 // 00100000
	m.memory[0x059] = 0x70 // 01110000

	// 2
	m.memory[0x05A] = 0xF0 // 11110000
	m.memory[0x05B] = 0x10 // 00010000
	m.memory[0x05C] = 0xF0 // 11110000
	m.memory[0x05D] = 0x80 // 10000000
	m.memory[0x05E] = 0xF0 // 11110000

	// 3
	m.memory[0x05F] = 0xF0 // 11110000
	m.memory[0x060] = 0x10 // 00010000
	m.memory[0x061] = 0xF0 // 11110000
	m.memory[0x062] = 0x10 // 00010000
	m.memory[0x063] = 0xF0 // 11110000

	// 4
	m.memory[0x064] = 0x90 // 10010000
	m.memory[0x065] = 0x90 // 10010000
	m.memory[0x066] = 0xF0 // 11110000
	m.memory[0x067] = 0x10 // 00010000
	m.memory[0x068] = 0x10 // 00010000

	// 5
	m.memory[0x069] = 0xF0 // 11110000
	m.memory[0x06A] = 0x80 // 10000000
	m.memory[0x06B] = 0xF0 // 11110000
	m.memory[0x06C] = 0x10 // 00010000
	m.memory[0x06D] = 0xF0 // 11110000

	// 6
	m.memory[0x06E] = 0xF0 // 11110000
	m.memory[0x06F] = 0x80 // 10000000
	m.memory[0x070] = 0xF0 // 11110000
	m.memory[0x071] = 0x90 // 10010000
	m.memory[0x072] = 0xF0 // 11110000

	// 7
	m.memory[0x073] = 0xF0 // 11110000
	m.memory[0x074] = 0x10 // 00010000
	m.memory[0x075] = 0x20 // 00100000
	m.memory[0x076] = 0x40 // 01000000
	m.memory[0x077] = 0x40 // 01000000

	// 8
	m.memory[0x078] = 0xF0 // 11110000
	m.memory[0x079] = 0x90 // 10010000
	m.memory[0x07A] = 0xF0 // 11110000
	m.memory[0x07B] = 0x90 // 10010000
	m.memory[0x07C] = 0xF0 // 11110000

	// 9
	m.memory[0x07D] = 0xF0 // 11110000
	m.memory[0x07E] = 0x90 // 10010000
	m.memory[0x07F] = 0xF0 // 11110000
	m.memory[0x080] = 0x10 // 00010000
	m.memory[0x081] = 0xF0 // 11110000

	// A
	m.memory[0x082] = 0xF0 // 11110000
	m.memory[0x083] = 0x90 // 10010000
	m.memory[0x084] = 0xF0 // 11110000
	m.memory[0x085] = 0x90 // 10010000
	m.memory[0x086] = 0x90 // 10010000

	// B
	m.memory[0x087] = 0xE0 // 11100000
	m.memory[0x088] = 0x90 // 10010000
	m.memory[0x089] = 0xE0 // 11100000
	m.memory[0x08A] = 0x90 // 10010000
	m.memory[0x08B] = 0xE0 // 11100000

	// C
	m.memory[0x08C] = 0xF0 // 11110000
	m.memory[0x08D] = 0x80 // 10000000
	m.memory[0x08E] = 0x80 // 10000000
	m.memory[0x08F] = 0x80 // 10000000
	m.memory[0x090] = 0xF0 // 11110000

	// D
	m.memory[0x091] = 0xE0 // 11100000
	m.memory[0x092] = 0x90 // 10010000
	m.memory[0x093] = 0x90 // 10010000
	m.memory[0x094] = 0x90 // 10010000
	m.memory[0x095] = 0xE0 // 11100000

	// E
	m.memory[0x096] = 0xF0 // 11110000
	m.memory[0x097] = 0x80 // 10000000
	m.memory[0x098] = 0xF0 // 11110000
	m.memory[0x099] = 0x80 // 10000000
	m.memory[0x09A] = 0xF0 // 11110000

	// F
	m.memory[0x09B] = 0xF0 // 11110000
	m.memory[0x09C] = 0x80 // 10000000
	m.memory[0x09D] = 0xF0 // 11110000
	m.memory[0x09E] = 0x80 // 10000000
	m.memory[0x09F] = 0x80 // 10000000
}

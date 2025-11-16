package chip_8

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

*/

const MEMORY_SPACE_BYTES int16 = 4096 //4kb

const MEMORY_START int16 = 0x000         //0
const RESERVED_MEMORY_END int16 = 0x1FF  //511
const PROGRAM_MEMORY_START int16 = 0x200 //512
const MEMORY_END int16 = 0xFFF           //4095

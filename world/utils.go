package world

import "math"

// JMT: no bounds checking at this time
func Nibble(arr []byte, i int) byte {
	if i%2 == 0 {
		return arr[i/2] & 0x0F
	}
	return (arr[i/2] >> 4) & 0x0F
}

// JMT: no bounds checking at this time
func WriteNibble(arr []byte, i int, b byte) {
	if i%2 == 0 {
		arr[i/2] = (arr[i/2] & 0xF0) | (b & 0x0F)
	} else {
		arr[i/2] = ((b << 4) & 0xF0) | (arr[i/2] & 0x0F)
	}
}

func floor16(in int32) int32 {
	return int32(math.Floor(float64(in) / 16.0))
}

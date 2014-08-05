package world

import (
	"bytes"
	"testing"
)

var Nibble_tests = []struct {
	arr []byte
	i   int
	b   byte
}{
	{[]byte{0x12, 0x34, 0x56, 0x78}, 3, 0x03},
	{[]byte{0x12, 0x34, 0x56, 0x78}, 4, 0x06},
}

func Test_Nibble(t *testing.T) {
	for _, tt := range Nibble_tests {
		outb := Nibble(tt.arr, tt.i)
		if outb != tt.b {
			t.Errorf("Given %v and %v, expected %v, got %v", tt.arr, tt.i, tt.b, outb)
		}
	}
}

var WriteNibble_tests = []struct {
	arr []byte
	i   int
	b   byte
	out []byte
}{
	{[]byte{0x12, 0x34, 0x56, 0x78}, 3, 0x09, []byte{0x12, 0x94, 0x56, 0x78}},
	{[]byte{0x12, 0x34, 0x56, 0x78}, 4, 0x09, []byte{0x12, 0x34, 0x59, 0x78}},
}

func Test_WriteNibble(t *testing.T) {
	for _, tt := range WriteNibble_tests {
		out := tt.arr
		WriteNibble(out, tt.i, tt.b)
		if bytes.Compare(out, tt.out) != 0 {
			t.Errorf("Given %v, %v and %v, expected %v, got %v", tt.arr, tt.i, tt.b, tt.out, out)
		}
	}
}

var floor16_tests = []struct {
	in  int32
	out int32
}{
	{3, 0},
	{21, 1},
	{-3, -1},
}

func Test_floor16(t *testing.T) {
	for _, tt := range floor16_tests {
		out := floor16(tt.in)
		if out != tt.out {
			t.Errorf("Given %v, expected %v, got %v", tt.in, tt.out, out)
		}
	}
}

package world

import "testing"

var point_tests = []struct {
	X  int32
	Y  int32
	Z  int32
	CX int32
	CZ int32
}{
	{0, 0, 0, 0, 0},
	{0, 128, 0, 0, 0},
	{27, 0, -15, 1, -1},
}

func Test_MakePoint(t *testing.T) {
	for _, tt := range point_tests {
		mp := MakePoint(tt.X, tt.Y, tt.Z)
		if tt.X != mp.X || tt.Y != mp.Y || tt.Z != mp.Z {
			t.Errorf("Given %d, %d, %d, expected same, got %d, %d, %d", tt.X, tt.Y, tt.Z, mp.X, mp.Y, mp.Z)
		}
	}
}

func Test_ChunkXZ(t *testing.T) {
	for _, tt := range point_tests {
		p := MakePoint(tt.X, tt.Y, tt.Z)
		ttXZ := p.ChunkXZ()
		if ttXZ.X != tt.CX || ttXZ.Z != tt.CZ {
			t.Errorf("Given %d, %d, %d, wanted %d, %d, got %d, %d", tt.X, tt.Y, tt.Z, tt.CX, tt.CZ, ttXZ.X, ttXZ.Z)
		}
	}
}

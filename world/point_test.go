package world

import "testing"

var point_tests = []struct {
	X  int32
	Y  int32
	Z  int32
	CX int32
	CZ int32
	I  int
}{
	{0, 0, 0, 0, 0, 0},
	{0, 128, 0, 0, 0, 0},
	{27, 0, -15, 1, -1, -229},
}

func Test_MakePoint(t *testing.T) {
	for _, tt := range point_tests {
		p := MakePoint(tt.X, tt.Y, tt.Z)
		if tt.X != p.X || tt.Y != p.Y || tt.Z != p.Z {
			t.Errorf("Given %d, %d, %d, expected same, got %d, %d, %d", tt.X, tt.Y, tt.Z, p.X, p.Y, p.Z)
		}
	}
}

func Test_ChunkXZ(t *testing.T) {
	for _, tt := range point_tests {
		p := MakePoint(tt.X, tt.Y, tt.Z)
		ttXZ := p.ChunkXZ()
		if ttXZ.X != tt.CX || ttXZ.Z != tt.CZ {
			t.Errorf("Given %d, %d, %d, expected %d, %d, got %d, %d", tt.X, tt.Y, tt.Z, tt.CX, tt.CZ, ttXZ.X, ttXZ.Z)
		}
	}
}

func Test_Index(t *testing.T) {
	for _, tt := range point_tests {
		p := MakePoint(tt.X, tt.Y, tt.Z)
		i := p.Index()
		if i != tt.I {
			t.Errorf("Given %d, %d, %d, expected %d, got %d", tt.X, tt.Y, tt.Z, tt.I, i)
		}
	}
}

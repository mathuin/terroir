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

func Test_NewMakePoint(t *testing.T) {
	for _, tt := range point_tests {
		np := NewPoint(tt.X, tt.Y, tt.Z)
		mp := MakePoint(tt.X, tt.Y, tt.Z)
		if np.X != mp.X || np.Y != mp.Y || np.Z != mp.Z {
			t.Errorf("Given %d, %d, %d, expected %#+v to equal %#+v", tt.X, tt.Y, tt.Z, np, mp)
		}
	}
}

func Test_WhichChunk(t *testing.T) {
	for _, tt := range point_tests {
		p := MakePoint(tt.X, tt.Y, tt.Z)
		cx, cz := p.WhichChunk()
		if cx != tt.CX || cz != tt.CZ {
			t.Errorf("Given %d, %d, %d, wanted %d, %d, got %d, %d", tt.X, tt.Y, tt.Z, tt.CX, tt.CZ, cx, cz)
		}
	}
}

var location_tests = []struct {
	LX, LY, LZ float64
	PX, PY, PZ int32
}{
	{27.9561, 63.527, -148.113, 27, 63, -149},
}

func Test_ToPoint(t *testing.T) {
	for _, tt := range location_tests {
		l := MakeLocation(tt.LX, tt.LY, tt.LZ)
		p := l.ToPoint()
		if p.X != tt.PX || p.Y != tt.PY || p.Z != tt.PZ {
			t.Errorf("Given %f, %f, %f, wanted %d, %d, %d, got %d, %d, %d", tt.LX, tt.LY, tt.LZ, tt.PX, tt.PY, tt.PZ, p.X, p.Y, p.Z)
		}
	}
}

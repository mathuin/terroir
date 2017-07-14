package world

import "testing"

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

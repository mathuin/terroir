package idt

import (
	"log"
	"testing"
)

var IDT_tests = []struct {
	coords [][2]float64
	values []int
	base   [][2]int
	nnear  int
	outf   []int16
	outt   []int16
}{
	{
		[][2]float64{
			{0, 0}, {8, 0}, {16, 0}, {24, 0},
			{0, 8}, {8, 8}, {16, 8}, {24, 8},
			{0, 16}, {8, 16}, {16, 16}, {24, 16},
			{0, 24}, {8, 24}, {16, 24}, {24, 24},
		},
		[]int{
			0, 0, 1, 1,
			0, 1, 1, 1,
			0, 1, 2, 3,
			1, 2, 3, 3,
		},
		[][2]int{
			{2, 2}, {6, 2}, {10, 2}, {14, 2}, {18, 2}, {22, 2},
			{2, 6}, {6, 6}, {10, 6}, {14, 6}, {18, 6}, {22, 6},
			{2, 10}, {6, 10}, {10, 10}, {14, 10}, {18, 10}, {22, 10},
			{2, 14}, {6, 14}, {10, 14}, {14, 14}, {18, 14}, {22, 14},
			{2, 18}, {6, 18}, {10, 18}, {14, 18}, {18, 18}, {22, 18},
			{2, 22}, {6, 22}, {10, 22}, {14, 22}, {18, 22}, {22, 22},
		},
		4,
		[]int16{
			0, 0, 0, 1, 1, 1,
			0, 1, 1, 1, 1, 1,
			0, 1, 1, 1, 1, 1,
			0, 1, 1, 2, 2, 2,
			0, 1, 1, 2, 2, 3,
			1, 2, 2, 3, 3, 3,
		},
		[]int16{
			0, 0, 0, 1, 1, 1,
			0, 1, 1, 1, 1, 1,
			0, 1, 1, 1, 1, 1,
			0, 1, 1, 2, 2, 3,
			0, 1, 1, 2, 2, 3,
			1, 2, 2, 3, 3, 3,
		},
	},
}

func Test_IDT(t *testing.T) {
	for _, tt := range IDT_tests {
		Debug = true
		idt, err := NewIDT(tt.coords, tt.values)
		if err != nil {
			t.Fail()
		}
		outf, err := idt.Call(tt.base, tt.nnear, false)
		if err != nil {
			t.Fail()
		}
		outt, err := idt.Call(tt.base, tt.nnear, true)
		if err != nil {
			t.Fail()
		}
		Debug = false

		log.Printf("%+#v", outf)
		log.Printf("%+#v", outt)

		for i, v := range outf {
			if tt.outf[i] != v {
				t.Errorf("idt: given %+#v %+#v %+#v, expected %+#v, got %+#v ", tt.coords, tt.values, tt.base, tt.outf, outf)
				break
			}
		}
		for i, v := range outt {
			if tt.outt[i] != v {
				t.Errorf("majority: given %+#v %+#v %+#v, expected %+#v, got %+#v ", tt.coords, tt.values, tt.base, tt.outt, outt)
				break
			}
		}
	}
}

var dot_tests = []struct {
	x   []float64
	y   []int
	out int
}{
	{[]float64{1.0, 2.0, 3.0}, []int{4, -5, 6}, 12},
	{[]float64{1, 3, -5}, []int{4, -2, -1}, 3},
}

func Test_dot(t *testing.T) {
	for _, tt := range dot_tests {
		out, err := dot(tt.x, tt.y)
		if err != nil {
			t.Fail()
		}
		if out != tt.out {
			t.Errorf("given %+#v %+#v, expected %d, got %d", tt.x, tt.y, tt.out, out)
		}
	}
}

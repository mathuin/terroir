package idt

import "testing"

var IDT_tests = []struct {
	coords [][2]float64
	values []int
	base   [][2]int
	outf   []int16
	outt   []int16
}{
	{
		[][2]float64{{1, 1}, {1, 3}, {3, 1}, {3, 3}},
		[]int{11, 23, 41, 23},
		[][2]int{{1, 1}, {1, 3}, {3, 1}, {3, 3}},
		[]int16{11, 23, 41, 23},
		[]int16{11, 23, 41, 23},
	},
}

func Test_IDT(t *testing.T) {
	for _, tt := range IDT_tests {
		idt, err := NewIDT(tt.coords, tt.values)
		if err != nil {
			t.Fail()
		}
		outf, err := idt.Call(tt.base, 4, false)
		if err != nil {
			t.Fail()
		}
		outt, err := idt.Call(tt.base, 4, true)
		if err != nil {
			t.Fail()
		}

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

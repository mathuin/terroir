package idt

import "testing"

// tests:
// I HAVE NO IDEA
// make it run right
// put in sample values
// get sample output
// write them to test

var IDT_tests = []struct {
	coords [][2]float64
	values []int
	base   [][2]int
	outf   []int
	outt   []int
}{
	{
		[][2]float64{{0.0, 0.1}, {0.2, 0.3}, {0.4, 0.5}},
		[]int{11, 23, 41},
		[][2]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
		[]int{10, 26, 25, 28},
		[]int{41, 11, 11, 11},
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
				t.Errorf("idt: given %+#v %+#v %+#v %+#v, expected %+#v, got %+#v ", tt.coords, tt.values, tt.base, tt.outf, outf)
				break
			}
		}
		for i, v := range outt {
			if tt.outt[i] != v {
				t.Errorf("majority: given %+#v %+#v %+#v %+#v, expected %+#v, got %+#v ", tt.coords, tt.values, tt.base, tt.outt, outt)
				break
			}
		}
	}
}

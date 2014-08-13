package carto

import "testing"

var intknob_tests = []struct {
	name   string
	inval  int
	mymax  int
	mymin  int
	outval int
	msg    string
}{
	{"zero", 0, 3, 1, 1, "warning: zero 0 outside 1-3 range"},
	{"one", 1, 3, 1, 1, ""},
	{"two", 2, 3, 1, 2, ""},
	{"three", 3, 3, 1, 3, ""},
	{"four", 4, 3, 1, 3, "warning: four 4 outside 1-3 range"},
}

func Test_intknob(t *testing.T) {
	for _, tt := range intknob_tests {
		ik := IntKnob{name: tt.name, value: tt.inval}
		msg := ik.setValue(tt.mymax, tt.mymin)
		if ik.value != tt.outval || msg != tt.msg {
			t.Errorf("given name %s inval %d mymax %d mymin %d, expected %d %s, got %d %s", tt.name, tt.inval, tt.mymax, tt.mymin, tt.outval, tt.msg, ik.value, msg)
		}
	}
}

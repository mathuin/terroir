package world

import (
	"bytes"
	"testing"
)

var toHalf_tests = []struct {
	inlow  byte
	inhigh byte
	outtop byte
	outbot byte
}{
	{0x12, 0x34, 0x31, 0x42},
}

func Test_toHalf(t *testing.T) {
	for _, tt := range toHalf_tests {
		outtop, outbot := toHalf(tt.inlow, tt.inhigh)
		if outtop != tt.outtop || outbot != tt.outbot {
			t.Errorf("Given %x and %x, expected %x and %x, got %x and %x", tt.inlow, tt.inhigh, tt.outtop, tt.outbot, outtop, outbot)
		}
	}
}

func Test_toDouble(t *testing.T) {
	for _, tt := range toHalf_tests {
		inlow, inhigh := toDouble(tt.outtop, tt.outbot)
		if inlow != tt.inlow || inhigh != tt.inhigh {
			t.Errorf("Given %x and %x, expected %x and %x, got %x and %x", tt.outtop, tt.outbot, tt.inlow, tt.inhigh, inlow, inhigh)
		}
	}
}

func Test_Half(t *testing.T) {
	arrin := make([]byte, 4096)
	arrout := make([]byte, 2048)
	arrouttop := make([]byte, 2048)
	arroutbot := make([]byte, 2048)

	for i := range arrin {
		arrin[i] = 0x34
	}

	for i := range arrouttop {
		arrouttop[i] = 0x33
		arroutbot[i] = 0x44
	}

	arrout = Half(arrin, true)

	if !bytes.Equal(arrout, arrouttop) {
		t.Errorf("Given %x and true, expected %x, got %x", arrin[0], arrouttop[0], arrout[0])
	}

	arrout = Half(arrin, false)

	if !bytes.Equal(arrout, arroutbot) {
		t.Errorf("Given %x and false, expected %x, got %x", arrin[0], arroutbot[0], arrout[0])
	}
}

func Test_Halve(t *testing.T) {
	arrin := make([]byte, 4096)
	arrouttop := make([]byte, 2048)
	arroutbot := make([]byte, 2048)

	for i := range arrin {
		arrin[i] = 0x34
	}

	for i := range arrouttop {
		arrouttop[i] = 0x33
		arroutbot[i] = 0x44
	}

	arrtop, arrbot := Halve(arrin)

	if !bytes.Equal(arrtop, arrouttop) || !bytes.Equal(arrbot, arroutbot) {
		t.Errorf("Given %x, expected %x and %x, got %x and %x", arrin[0], arrouttop[0], arroutbot[0], arrtop[0], arrbot[0])
	}
}

func Test_Double(t *testing.T) {
	arrout := make([]byte, 4096)
	arroutcheck := make([]byte, 4096)
	arrintop := make([]byte, 2048)
	arrinbot := make([]byte, 2048)

	for i := range arrintop {
		arrintop[i] = 0x33
		arrinbot[i] = 0x44
	}

	for i := range arroutcheck {
		arroutcheck[i] = 0x34
	}

	arrout = Double(arrintop, arrinbot)

	if !bytes.Equal(arrout, arroutcheck) {
		t.Errorf("Given %x and %x, expected %x, got %x", arrintop[0], arrinbot[0], arroutcheck[0], arrout[0])
	}
}

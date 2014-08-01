package world

import "testing"

var half_tests = []struct {
	arrin  []byte
	topin  bool
	arrout []byte
}{
	{[]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xab}, true, []byte{0x31, 0x75, 0xa9}},
	{[]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xab}, false, []byte{0x42, 0x86, 0xb0}},
}

func Test_Half(t *testing.T) {
	var arrin [4096]byte
	var arrout [2048]byte
	var arroutcheck [2048]byte

	for i := range arrin {
		arrin[i] = 0x34
	}

	for i := range arroutcheck {
		arroutcheck[i] = 0x33
	}

	arrout = Half(arrin, true)

	if arrout != arroutcheck {
		t.Errorf("Half conversion on top failed")
	}

	for i := range arroutcheck {
		arroutcheck[i] = 0x44
	}

	arrout = Half(arrin, false)

	if arrout != arroutcheck {
		t.Errorf("Half conversion on bottom failed")
	}
}

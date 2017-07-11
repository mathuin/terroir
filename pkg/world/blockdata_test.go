package world

import "testing"

func Test_checkBlockData(t *testing.T) {
	b := blockNames["Stone"]
	if b.block != 1 || b.data != 0 {
		t.Fail()
	}
}

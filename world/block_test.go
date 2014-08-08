package world

import "testing"

func Test_checkBlockData(t *testing.T) {
	b := BlockNames["Stone"]
	if b.block != 1 || b.data != 0 {
		t.Fail()
	}
}

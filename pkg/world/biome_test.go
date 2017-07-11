package world

import "testing"

func Test_checkBiome(t *testing.T) {
	b := Biome["Desert"]
	if b != 2 {
		t.Fail()
	}
}

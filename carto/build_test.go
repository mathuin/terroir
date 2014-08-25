package carto

import (
	"testing"

	"github.com/mathuin/terroir/world"
)

var buildWorld_tests = []struct {
	name   string
	ll     FloatExtents
	elname string
	lcname string
}{
	{"BlockIsland", FloatExtents{-71.575, -71.576, 41.189, 41.191}, "elevation.tif", "landcover.tif"},
}

func Test_buildWorld(t *testing.T) {
	for _, tt := range buildWorld_tests {
		r := MakeRegion(tt.name, tt.ll, tt.elname, tt.lcname)
		// Debug = true
		r.buildMap()
		// Debug = false
		// Debug = true
		w, err := r.buildWorld()
		if err != nil {
			t.Fail()
		}
		// Debug = false

		w.SetSaveDir("/tmp")
		w.Write()

		nw, nwerr := world.ReadWorld("/tmp", tt.name, false)
		if nwerr != nil {
			t.Fail()
		}
		_ = nw
	}

}

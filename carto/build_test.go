package carto

import (
	"path"
	"testing"

	"github.com/mathuin/terroir/world"
)

var buildWorld_tests = []struct {
	name string
	ll   FloatExtents
}{
	{"BlockIsland", FloatExtents{-71.575, -71.576, 41.189, 41.191}},
}

func Test_buildWorld(t *testing.T) {
	for _, tt := range buildWorld_tests {
		r := MakeRegion(tt.name, tt.ll)
		r.vrts["elevation"] = path.Join(tt.name, "elevation.tif")
		r.vrts["landcover"] = path.Join(tt.name, "landcover.tif")
		Debug = true
		r.buildMap()
		w, err := r.buildWorld()
		if err != nil {
			t.Fail()
		}
		Debug = false

		w.SetSaveDir("/tmp")
		w.Write()

		nw, nwerr := world.ReadWorld("/tmp", tt.name, false)
		if nwerr != nil {
			t.Fail()
		}
		_ = nw
	}

}

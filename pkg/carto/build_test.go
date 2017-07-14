package carto

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mathuin/terroir/pkg/world"
)

var buildWorld_tests = []struct {
	name   string
	ll     FloatExtents
	elname string
	lcname string
}{
	{"BlockIsland", FloatExtents{-71.575, -71.576, 41.189, 41.191}, "elevation.tif", "landcover.tif"},
}

func Test_BuildWorld(t *testing.T) {
	td, nerr := ioutil.TempDir("", "")
	if nerr != nil {
		panic(nerr)
	}
	defer os.RemoveAll(td)

	for _, tt := range buildWorld_tests {
		r := MakeRegion(tt.name, tt.ll, tt.elname, tt.lcname)
		r.BuildMap()
		w, err := r.BuildWorld()
		if err != nil {
			t.Fail()
		}

		w.SetSaveDir(td)
		werr := w.Write()
		if werr != nil {
			panic(werr)
		}

		nw, nwerr := world.ReadWorld(td, tt.name, true)
		if nwerr != nil {
			panic(nwerr)
		}
		_ = nw
	}

}

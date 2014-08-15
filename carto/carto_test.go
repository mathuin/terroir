package carto

import (
	"path"
	"testing"

	"github.com/lukeroth/gdal"
)

var buildMap_tests = []struct {
	name     string
	ll       FloatExtents
	minmaxes []RasterInfo
}{
	{
		"BlockIsland",
		FloatExtents{-71.575, -71.576, 41.189, 41.191},
		[]RasterInfo{
			RasterInfo{"Int16", 11, 95},
			RasterInfo{"Int16", 62, 64},
			RasterInfo{"Int16", 0, 30},
			RasterInfo{"Int16", 1, 4},
			RasterInfo{"Int16", 0, 24},
		},
	},
}

func Test_buildMap(t *testing.T) {
	for _, tt := range buildMap_tests {
		r := MakeRegion(tt.name, tt.ll)
		r.tilesize = 16
		r.vrts["elevation"] = path.Join(tt.name, "elevation.tif")
		r.vrts["landcover"] = path.Join(tt.name, "landcover.tif")
		Debug = true
		r.buildMap()
		Debug = true

		// check the raster minmaxes
		ds, err := gdal.Open(r.mapfile, gdal.ReadOnly)
		if err != nil {
			t.Fail()
		}
		if Debug {
			datasetInfo(ds, "Test")
		}
		minmaxes := datasetMinMaxes(ds)

		if len(minmaxes) != len(tt.minmaxes) {
			t.Fatalf("len(minmaxes) %d != len(tt.minmaxes) %d", len(minmaxes), len(tt.minmaxes))
		}
		for i, v := range minmaxes {
			if tt.minmaxes[i] != v {
				t.Errorf("Raster #%d: expected \"%s\", got \"%s\"", i+1, tt.minmaxes[i], v)
			}
		}
	}
}

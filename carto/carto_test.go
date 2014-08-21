package carto

import (
	"path"
	"reflect"
	"testing"

	"github.com/mathuin/gdal"
)

var buildMap_tests = []struct {
	name     string
	ll       FloatExtents
	minmaxes []RasterInfo
	histos   []RasterHInfo
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
		[]RasterHInfo{
			RasterHInfo{"Int16", map[int]int{11: 49475, 21: 1989, 22: 812, 23: 291, 31: 2138, 41: 313, 42: 1051, 43: 2361, 90: 6706, 95: 400}},
			RasterHInfo{"Int16", map[int]int{62: 63631, 63: 1860, 64: 45}},
			RasterHInfo{"Int16", map[int]int{0: 16061, 1: 948, 2: 848, 3: 781, 4: 919, 5: 715, 6: 746, 7: 661, 8: 679, 9: 794, 10: 614, 11: 718, 12: 608, 13: 723, 14: 668, 15: 592, 16: 743, 17: 675, 18: 661, 19: 643, 20: 588, 21: 722, 22: 630, 23: 617, 24: 586, 25: 646, 26: 612, 27: 575, 28: 651, 29: 570, 30: 29542}},
			RasterHInfo{"Int16", map[int]int{1: 13163, 2: 39217, 3: 12367, 4: 789}},
			RasterHInfo{"Int16", map[int]int{0: 21061, 1: 9916, 2: 1515, 4: 3536, 24: 29508}},
		},
	},
}

func Test_buildMap(t *testing.T) {
	for _, tt := range buildMap_tests {
		r := MakeRegion(tt.name, tt.ll)
		r.tilesize = 16
		r.vrts["elevation"] = path.Join(tt.name, "elevation.tif")
		r.vrts["landcover"] = path.Join(tt.name, "landcover.tif")
		// Debug = true
		r.buildMap()
		// Debug = false

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

		histos := datasetHistograms(ds)
		if len(histos) != len(tt.histos) {
			t.Fatalf("len(histos) %d != len(tt.histos) %d", len(histos), len(tt.histos))
		}
		for i, v := range histos {
			if tt.histos[i].datatype != v.datatype {
				t.Errorf("Raster #%d: expected datatype \"%s\", got \"%s\"", i+1, tt.histos[i].datatype, v.datatype)
			}
			if !reflect.DeepEqual(tt.histos[i].buckets, v.buckets) {
				t.Errorf("Raster #%d: expected buckets \"%s\", got \"%s\"", i+1, tt.histos[i].buckets, v.buckets)
			}
		}
	}
}

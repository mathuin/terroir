package carto

import (
	"reflect"
	"testing"

	"github.com/mathuin/gdal"
)

var buildMap_tests = []struct {
	name   string
	elname string
	lcname string
	ll     FloatExtents
	histos []RasterHInfo
}{
	{
		"BlockIsland", "elevation.tif", "landcover.tif",
		FloatExtents{-71.575, -71.576, 41.189, 41.191},
		[]RasterHInfo{
			RasterHInfo{"Int16", map[int]int{31: 2083, 90: 6713, 43: 2322, 22: 757, 42: 1061, 11: 49617, 95: 358, 41: 316, 21: 2005, 23: 304}},
			RasterHInfo{"Int16", map[int]int{62: 63631, 63: 1860, 64: 45}},
			RasterHInfo{"Int16", map[int]int{7: 644, 6: 732, 5: 688, 15: 591, 27: 568, 13: 729, 12: 603, 10: 608, 3: 699, 0: 15919, 19: 641, 21: 728, 22: 637, 24: 580, 29: 567, 16: 741, 11: 720, 9: 791, 18: 660, 20: 583, 4: 896, 17: 675, 28: 655, 14: 666, 8: 678, 2: 743, 26: 607, 30: 29983, 1: 904, 23: 651, 25: 649}},
			RasterHInfo{"Int16", map[int]int{1: 13163, 2: 39217, 3: 12367, 4: 789}},
		},
	},
}

func Test_buildMap(t *testing.T) {
	for _, tt := range buildMap_tests {
		r := MakeRegion(tt.name, tt.ll, tt.elname, tt.lcname)
		r.tilesize = 16
		// Debug = true
		r.BuildMap()
		// Debug = false

		// check the raster minmaxes
		ds, err := gdal.Open(r.mapfile, gdal.ReadOnly)
		if err != nil {
			t.Fail()
		}
		if Debug {
			datasetInfo(ds, "Test")
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
				// JMT: crust raster is expected to vary
				if i != 3 {
					t.Errorf("Raster #%d: expected buckets \"%+#v\", got \"%+#v\"", i+1, tt.histos[i].buckets, v.buckets)
				}
			}
		}
	}
}

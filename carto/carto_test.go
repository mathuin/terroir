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
			RasterHInfo{"Int16", map[int]int{31: 2085, 90: 6701, 95: 357, 41: 312, 21: 2006, 42: 1060, 11: 49620, 43: 2331, 22: 764, 23: 300}},
			RasterHInfo{"Int16", map[int]int{62: 63631, 63: 1860, 64: 45}},
			RasterHInfo{"Int16", map[int]int{4: 896, 21: 728, 26: 607, 18: 660, 14: 666, 13: 729, 11: 720, 9: 791, 15: 591, 10: 608, 5: 687, 17: 675, 24: 580, 29: 567, 3: 700, 19: 641, 7: 645, 1: 905, 16: 742, 25: 649, 8: 680, 20: 583, 23: 651, 12: 603, 6: 731, 2: 742, 0: 15916, 22: 637, 27: 568, 28: 655, 30: 29983}},
			RasterHInfo{"Int16", map[int]int{1: 13163, 2: 39217, 3: 12367, 4: 789}},
		},
	},
}

func Test_buildMap(t *testing.T) {
	for _, tt := range buildMap_tests {
		r := MakeRegion(tt.name, tt.ll, tt.elname, tt.lcname)
		r.tilesize = 16
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

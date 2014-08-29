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
	histos []RasterInfo
}{
	{
		"BlockIsland", "elevation.tif", "landcover.tif",
		FloatExtents{-71.575, -71.576, 41.189, 41.191},
		[]RasterInfo{
			// landcover
			RasterInfo{"Int16", map[int]int{
				11: 49727,
				21: 1987,
				22: 734,
				23: 278,
				31: 2019,
				41: 298,
				43: 2317,
				42: 1067,
				90: 6767,
				95: 342,
			}},
			// elevation
			RasterInfo{"Int16", map[int]int{
				62: 63631,
				63: 1860,
				64: 45,
			}},
			// bathy
			RasterInfo{"Int16", map[int]int{
				0:  15809,
				1:  916,
				2:  729,
				3:  694,
				4:  893,
				5:  679,
				6:  747,
				7:  641,
				8:  689,
				9:  766,
				10: 595,
				11: 728,
				12: 595,
				13: 737,
				14: 654,
				15: 588,
				16: 750,
				17: 667,
				18: 671,
				19: 635,
				20: 580,
				21: 732,
				22: 638,
				23: 663,
				24: 578,
				25: 655,
				26: 599,
				27: 569,
				28: 651,
				29: 570,
				30: 30118,
			}},
			// crust -- may change
			RasterInfo{"Int16", map[int]int{
				1: 13163,
				2: 39217,
				3: 12367,
				4: 789,
			}},
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

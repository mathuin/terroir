package carto

import (
	"testing"

	"github.com/lukeroth/gdal"
)

var buildMap_tests = []struct {
	ll      FloatExtents
	elvrt   string
	lcvrt   string
	rasters [][]float64
}{
	{
		FloatExtents{-71.575, -71.576, 41.189, 41.191},
		"test_elevation.tif",
		"test_landcover.tif",
		[][]float64{[]float64{11, 95}, []float64{62, 64}, []float64{0, 30}, []float64{1, 4}},
	},
}

func Test_buildMap(t *testing.T) {
	for _, tt := range buildMap_tests {
		r := MakeRegion("Pie", tt.ll)
		r.tilesize = 16
		r.vrts["elevation"] = tt.elvrt
		r.vrts["landcover"] = tt.lcvrt
		Debug = true
		r.buildMap()
		Debug = false

		// check the raster minmaxes
		ds, err := gdal.Open(r.mapfile, gdal.ReadOnly)
		if err != nil {
			t.Fail()
		}
		bandCount := ds.RasterCount()
		for i := 0; i < bandCount; i++ {
			rbi := i + 1
			rb := ds.RasterBand(rbi)
			rbmin, minok := rb.GetMinimum()
			rbmax, maxok := rb.GetMaximum()
			if !minok || !maxok {
				rbmin, rbmax = rb.ComputeMinMax(0)
			}
			if rbmin != tt.rasters[i][0] || rbmax != tt.rasters[i][1] {
				t.Errorf("raster %d: expected (%f, %f) got (%f, %f)", rbi, tt.rasters[i][0], tt.rasters[i][1], rbmin, rbmax)
			}
		}

	}
}

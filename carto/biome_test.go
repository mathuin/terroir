package carto

import (
	"log"
	"strings"
	"testing"
)

var biome_tests = []struct {
	extents  FloatExtents
	gt       [6]float64
	x        int
	y        int
	maxdepth int
	lcarr    []int16
	elevarr  []int16
	bathyarr []int16
	outarr   []int16
}{
	{
		FloatExtents{-71.575, -71.576, 41.189, 41.191},
		[6]float64{2007975, 30, 0, 2282025, 0, -30},
		5, 7, 2,
		[]int16{
			11, 11, 11, 11, 11,
			11, 11, 11, 11, 11,
			11, 11, 71, 11, 11,
			11, 11, 31, 41, 11,
			11, 11, 11, 11, 11,
			11, 11, 11, 11, 11,
			11, 11, 11, 11, 11,
		},
		[]int16{
			62, 62, 62, 62, 62,
			62, 62, 62, 62, 62,
			62, 62, 63, 62, 62,
			62, 62, 63, 63, 62,
			62, 62, 62, 62, 62,
			62, 62, 62, 62, 62,
			62, 62, 62, 62, 62,
		},
		[]int16{
			2, 2, 2, 2, 2,
			2, 1, 1, 1, 2,
			2, 1, 0, 1, 1,
			2, 1, 0, 0, 1,
			2, 1, 1, 1, 1,
			2, 2, 2, 2, 2,
			2, 2, 2, 2, 2,
		},
		[]int16{
			24, 24, 24, 24, 24,
			24, 0, 0, 0, 24,
			24, 0, 1, 0, 0,
			24, 0, 2, 4, 0,
			24, 0, 0, 0, 0,
			24, 24, 24, 24, 24,
			24, 24, 24, 24, 24,
		},
	},
}

func Test_biome(t *testing.T) {
	for _, tt := range biome_tests {
		r := MakeRegion("Pie", FloatExtents{-71.575, -71.576, 41.189, 41.191})
		r.maxdepth = tt.maxdepth
		outarr := r.biome(tt.lcarr, tt.elevarr, tt.bathyarr)

		for i, v := range outarr {
			if v != tt.outarr[i] {
				t.Errorf("expected \n%s, got \n%s", printarr(tt.outarr, tt.x, tt.y), printarr(outarr, tt.x, tt.y))
				break
			}
		}
	}
}

func Test_newbiome(t *testing.T) {
	for _, tt := range biome_tests {
		for _, line := range strings.Split(printarr(tt.lcarr, tt.x, tt.y), "\n") {
			log.Print(line)
		}
		r := MakeRegion("Pie", tt.extents)
		r.maxdepth = tt.maxdepth
		// Debug = true
		outarr, err := r.newbiome(tt.x, tt.y, tt.gt, tt.lcarr, tt.elevarr, tt.bathyarr)
		// Debug = false
		if err != nil {
			t.Fail()
		}
		for i, v := range outarr {
			if v != tt.outarr[i] {
				t.Errorf("expected \n%s, got \n%s", printarr(tt.outarr, tt.x, tt.y), printarr(outarr, tt.x, tt.y))
				break
			}
		}
	}
}

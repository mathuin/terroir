package carto

import "testing"

var biome_tests = []struct {
	x        int
	y        int
	maxdepth int
	lcarr    []int16
	elevarr  []int16
	bathyarr []int16
	outarr   []int16
}{
	{
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

// func Test_newbiome(t *testing.T) {
// 	for _, tt := range biome_tests {
// 		r := MakeRegion("Pie", FloatExtents{-71.575, -71.576, 41.189, 41.191})
// 		r.maxdepth = tt.maxdepth
// 		outarr, err := r.newbiome(tt.x, tt.y, tt.lcarr, tt.elevarr, tt.bathyarr)
// 		if err != nil {
// 			t.Fail()
// 		}
// 		for i, v := range outarr {
// 			if v != tt.outarr[i] {
// 				t.Errorf("expected \n%s, got \n%s", printarr(tt.outarr, tt.x, tt.y), printarr(outarr, tt.x, tt.y))
// 				break
// 			}
// 		}
// 	}
// }

package carto

import (
	"fmt"
	"testing"
)

var bathy_tests = []struct {
	inarr  []int16
	inx    int
	iny    int
	outarr []int16
}{
	{
		[]int16{
			11, 11, 11, 11, 11,
			11, 11, 11, 11, 11,
			11, 11, 11, 11, 11,
			11, 11, 31, 32, 11,
			11, 11, 11, 11, 11,
			11, 11, 11, 11, 11,
			11, 11, 11, 11, 11,
		},
		5, 7,
		[]int16{
			2, 2, 2, 2, 2,
			2, 2, 2, 2, 2,
			2, 1, 1, 1, 1,
			2, 1, 0, 0, 1,
			2, 1, 1, 1, 1,
			2, 2, 2, 2, 2,
			2, 2, 2, 2, 2,
		}},
}

func Test_bathy(t *testing.T) {
	for _, tt := range bathy_tests {
		r := MakeRegion("Pie", FloatExtents{-71.575, -71.576, 41.189, 41.191})
		outarr := r.bathy(tt.inarr, tt.inx, tt.iny)

		for i, v := range outarr {
			if v != tt.outarr[i] {
				t.Errorf("given\n%s, wanted\n%s, got\n%s", printarr(tt.inarr, tt.inx, tt.iny), printarr(tt.outarr, tt.inx, tt.iny), printarr(outarr, tt.inx, tt.iny))
				break
			}
		}

	}
}

func printarr(arr []int16, x int, y int) string {
	retval := ""
	line := ""
	for i := 0; i < x*y; i++ {
		line = fmt.Sprintf("%s %d", line, arr[i])
		if i%x == x-1 {
			retval = fmt.Sprintf("%s%s\n", retval, line)
			line = ""
		}
	}
	return retval
}

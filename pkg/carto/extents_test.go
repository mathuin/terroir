package carto

import (
	"math"
	"testing"
)

var generateExtents_tests = []struct {
	ll     FloatExtents
	albers map[string]IntExtents
	wgs84  map[string]FloatExtents
}{
	{FloatExtents{-71.533, -71.62, 41.238, 41.142},
		map[string]IntExtents{
			"elevation": {2015232, 2002944, 2291712, 2273280},
			"landcover": {2015412, 2002764, 2291892, 2273100},
		},
		map[string]FloatExtents{
			"elevation": {-71.48245416519951, -71.68161580988412, 41.30762685177235, 41.12049902574309},
			"landcover": {-71.47980755995131, -71.68425418133639, 41.309590101757536, 41.11853504487958},
		},
	},
}

var epsilon = 1e-10

func Test_generateExtents(t *testing.T) {
	for _, tt := range generateExtents_tests {
		r := MakeRegion("Pie", tt.ll, "", "")
		r.tilesize = 1024
		r.generateExtents()
		for maptype, subarr := range tt.albers {
			for coord, value := range subarr {
				if r.albers[maptype][coord] != value {
					t.Errorf("albers %s: given %+#v, expected %+#v, got %+#v", maptype, tt.ll, tt.albers[maptype], r.albers[maptype])
					break
				}
			}
		}
		for maptype, subarr := range tt.wgs84 {
			for coord, value := range subarr {
				if math.Abs(r.wgs84[maptype][coord]-value)/value > epsilon {
					t.Errorf("wgs84 %s: given %+#v, expected %+#v, got %+#v", maptype, tt.ll, tt.wgs84[maptype], r.wgs84[maptype])
					break
				}
			}
		}
	}
}

var getCorners_tests = []struct {
	fromCS string
	toCS   string
	in     Extents
	out    FloatExtents
}{
	{"+proj=longlat +ellps=WGS84 +datum=WGS84 +no_defs", "+proj=aea +datum=NAD83 +lat_1=29.5 +lat_2=45.5 +lat_0=23 +lon_0=-96 +x_0=0 +y_0=0 +units=m", FloatExtents{-71.533, -71.62, 41.238, 41.142}, FloatExtents{2015094.9510428484, 2005360.7296084524, 2286126.7340978114, 2273891.6760678287}},
	{"+proj=aea +datum=NAD83 +lat_1=29.5 +lat_2=45.5 +lat_0=23 +lon_0=-96 +x_0=0 +y_0=0 +units=m", "+proj=longlat +ellps=WGS84 +datum=WGS84 +no_defs", IntExtents{2015412, 2002764, 2291892, 2273100}, FloatExtents{-71.47980755995131, -71.68425418133639, 41.309590101757536, 41.11853504487958}},
	{"+proj=aea +datum=NAD83 +lat_1=29.5 +lat_2=45.5 +lat_0=23 +lon_0=-96 +x_0=0 +y_0=0 +units=m", "+proj=longlat +ellps=WGS84 +datum=WGS84 +no_defs", IntExtents{2015232, 2002944, 2291712, 2273280}, FloatExtents{-71.48245416519951, -71.68161580988412, 41.30762685177235, 41.12049902574309}},
}

func Test_getCorners(t *testing.T) {
	for _, tt := range getCorners_tests {
		out := getCorners(tt.fromCS, tt.toCS, tt.in)
		for i, v := range out {
			if math.Abs(tt.out[i]-v)/v > epsilon {
				t.Errorf("Given %s %s %v, wanted %+#v, got %+#v", tt.fromCS, tt.toCS, tt.in, tt.out, out)
				break
			}
		}
	}
}

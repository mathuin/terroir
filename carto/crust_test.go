package carto

import (
	"log"
	"testing"
)

func Test_crust(t *testing.T) {
	r := MakeRegion("Pie", FloatExtents{-71.575, -71.576, 41.189, 41.191})
	crustBuffer := r.crust(100, 150)
	minwidth := int16(1)
	maxwidth := int16(5)
	for _, v := range crustBuffer {
		if v < minwidth {
			log.Panicf("too small - %d < %d", v, minwidth)
		}
		if v > maxwidth {
			log.Panicf("too large - %d > %d", v, maxwidth)
		}
	}
}

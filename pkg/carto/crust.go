package carto

import (
	"math/rand"

	"github.com/mathuin/terroir/pkg/idt"
)

func (r Region) crust(rXsize int, rYsize int) []int16 {
	minwidth := 1
	maxwidth := 5
	crustrange := maxwidth - minwidth
	coverage := 0.05

	bufferLen := rXsize * rYsize

	numcoords := int(float64(bufferLen) * coverage)
	crustCoords := make([][2]float64, numcoords)
	crustValues := make([]int, numcoords)
	for i := range crustCoords {
		crustCoords[i] = [2]float64{float64(rand.Intn(rXsize)), float64(rand.Intn(rYsize))}
		crustValues[i] = (rand.Int() % crustrange) + minwidth
	}

	crustBase := make([][2]int, bufferLen)
	for i := 0; i < bufferLen; i++ {
		crustBase[i] = [2]int{i % rXsize, i / rXsize}
	}
	crustIDT, err := idt.NewIDT(crustCoords, crustValues)
	if err != nil {
		panic(err)
	}
	crustBuffer, err := crustIDT.Call(crustBase, 31, false)
	if err != nil {
		panic(err)
	}
	return crustBuffer
}

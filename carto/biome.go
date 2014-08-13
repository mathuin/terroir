package carto

import (
	"log"

	"github.com/mathuin/terroir/world"
)

func (r *Region) biome(lcarr []int16, elevarr []int16, bathyarr []int16) (biomearr []int16) {
	bufferLen := len(lcarr)
	biomearr = make([]int16, bufferLen)
	for i := 0; i < bufferLen; i++ {
		lc := lcarr[i]
		// elev := elevarr[i]
		bathy := bathyarr[i]

		var biome string
		switch lc {
		case 11:
			if bathy > int16(r.maxdepth-1) {
				biome = "Deep Ocean"
			} else {
				biome = "Ocean"
			}
		default:
			biome = "Plains"
		}
		val, ok := world.Biome[biome]
		if !ok {
			log.Printf("%s is not a valid biome!", biome)
			val = -1
		}
		biomearr[i] = int16(val)
	}
	return biomearr
}

package carto

import (
	"fmt"
	"log"
	"math"

	"github.com/mathuin/gdal"
)

// Extents are arrays of four values:
// xmax, xmin, ymax, ymin
type FloatExtents [4]float64
type IntExtents [4]int

type Extents interface {
	ints() IntExtents
	floats() FloatExtents
}

func (i IntExtents) ints() IntExtents {
	return i
}

func (i IntExtents) floats() FloatExtents {
	f := FloatExtents{}
	for k, v := range i {
		f[k] = float64(v)
	}
	return f
}

func (f FloatExtents) ints() IntExtents {
	i := IntExtents{}
	for k, v := range i {
		i[k] = int(v)
	}
	return i
}

func (f FloatExtents) floats() FloatExtents {
	return f
}

const (
	xMax = iota
	xMin
	yMax
	yMin
)

func (r Region) generateExtents() {
	// get corners from wgs to albers
	marr := getCorners(wgs84_proj, albers_proj, r.ll)

	realsize := r.scale * r.tilesize

	var tiles IntExtents
	for i, v := range marr {
		// i % 2 == 0 for maxes
		if i%2 == 0 {
			tiles[i] = int(math.Ceil(v / float64(realsize)))
		} else {
			tiles[i] = int(math.Floor(v / float64(realsize)))
		}
	}

	var nae IntExtents
	for i, v := range tiles {
		nae[i] = v * realsize
	}
	r.albers["elevation"] = nae

	// landcover requires a maxdepth-sized border for calculating depth
	borderwidth := r.maxdepth * r.scale

	var nal IntExtents
	for i, v := range r.albers["elevation"] {
		// i % 2 == 0 for "maxes"
		if i%2 == 0 {
			nal[i] = v + borderwidth
		} else {
			nal[i] = v - borderwidth
		}
	}
	r.albers["landcover"] = nal

	// get corners from albers back to wgs
	for maptype := range r.albers {
		r.wgs84[maptype] = getCorners(albers_proj, wgs84_proj, r.albers[maptype])
	}

	return
}

func getCorners(fromCS string, toCS string, in Extents) FloatExtents {
	if Debug {
		log.Print("getCorners: ")
		log.Print(" fromCS: ", fromCS)
		log.Print(" toCS: ", toCS)
		log.Print(" in: ", in)
	}

	fromSR := gdal.CreateSpatialReference("")
	fromSR.FromProj4(fromCS)
	toSR := gdal.CreateSpatialReference("")
	toSR.FromProj4(toCS)

	fe := in.floats()
	xmax := fe[0]
	xmin := fe[1]
	ymax := fe[2]
	ymin := fe[3]

	corners := [][]float64{{xmin, ymin}, {xmin, ymax}, {xmax, ymin}, {xmax, ymax}}
	if Debug {
		log.Print("  corners: ", corners)
	}

	xfloat := Float64Arr{}
	yfloat := Float64Arr{}

	for _, corner := range corners {
		wkt := fmt.Sprintf("POINT (%f %f)", corner[0], corner[1])
		if Debug {
			log.Print("wkt: ", wkt)
		}
		point, err := gdal.CreateFromWKT(wkt, fromSR)
		if notnil(err) {
			panic(err)
		}
		if Debug {
			log.Print("point pre-transform:", point)
		}
		point.TransformTo(toSR)
		if Debug {
			log.Print("point post-transform:", point)
		}
		xfloat = append(xfloat, point.X(0))
		yfloat = append(yfloat, point.Y(0))
	}
	return FloatExtents{xfloat.max(), xfloat.min(), yfloat.max(), yfloat.min()}
}

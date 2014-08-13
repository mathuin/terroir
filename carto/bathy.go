package carto

import (
	"fmt"
	"strings"

	"github.com/mathuin/gdal"
)

func (r Region) bathy(inarr []int16, inx int, iny int) []int16 {
	inprod := inx * iny

	// mem driver!
	memdrv, err1 := gdal.GetDriverByName("MEM")
	if err1 != nil {
		panic(err1)
	}

	// create a source band from that array
	srcDS := memdrv.Create("src", inx, iny, 1, gdal.Int16, nil)
	srcBand := srcDS.RasterBand(1)
	err2 := srcBand.IO(gdal.Write, 0, 0, inx, iny, inarr, inx, iny, 0, 0)
	if err2 != nil && err2.Error() != "No Error" {
		panic(err2)
	}

	// create a target band
	destDS := memdrv.Create("dest", inx, iny, 1, gdal.Int16, nil)
	destBand := destDS.RasterBand(1)

	// someday water may not be so simple
	not11 := nomatch(inarr, []int16{11})

	// configure options
	options := []string{fmt.Sprintf("MAXDIST=%d", r.maxdepth), fmt.Sprintf("NODATA=%d", r.maxdepth), fmt.Sprintf("VALUES=%s", strings.Join(not11, ","))}

	// run computeproximity
	err3 := srcBand.ComputeProximity(destBand, options, gdal.ScaledProgress, nil)
	if err3 != nil && err3.Error() != "No Error" {
		panic(err3)
	}

	// get output
	outarr := make([]int16, inprod)
	err4 := destBand.IO(gdal.Read, 0, 0, inx, iny, outarr, inx, iny, 0, 0)
	if err4 != nil && err4.Error() != "No Error" {
		panic(err4)
	}

	return outarr
}

// a list of all the numbers in the incoming array
// that don't any of the specified values
func nomatch(arr []int16, ints []int16) []string {

	dmap := make(map[int16]bool)
	for _, v := range arr {
		if _, ok := dmap[v]; !ok {
			dmap[v] = true
		}
	}

	for _, v := range ints {
		delete(dmap, v)
	}

	retval := []string{}

	for k := range dmap {
		retval = append(retval, fmt.Sprintf("%d", k))
	}

	return retval
}

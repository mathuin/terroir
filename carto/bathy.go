package carto

import (
	"fmt"
	"log"
	"strings"

	"github.com/lukeroth/gdal"
)

func (r Region) bathy(darr []int16, depthx int, depthy int, gt [6]float64, proj string) []int16 {
	memdrv, err := gdal.GetDriverByName("MEM")
	if err != nil {
		panic(err)
	}
	depthDS := memdrv.Create("depth", depthx, depthy, 1, gdal.Int16, nil)
	depthDS.SetGeoTransform(gt)
	depthDS.SetProjection(proj)
	depthBand := depthDS.RasterBand(1)
	deptherr := depthBand.IO(gdal.Write, 0, 0, depthx, depthy, darr, depthx, depthy, 0, 0)
	if notnil(deptherr) {
		panic(deptherr)
	}
	bathyDS := memdrv.Create("bathy", depthx, depthy, 1, gdal.Int16, nil)
	bathyDS.SetGeoTransform(gt)
	bathyDS.SetProjection(proj)
	bathyBand := bathyDS.RasterBand(1)

	dmap := make(map[int16]bool)
	for _, v := range darr {
		if _, ok := dmap[v]; !ok {
			dmap[v] = true
		}
	}
	not11 := []string{}
	for k := range dmap {
		if k != 11 {
			not11 = append(not11, fmt.Sprintf("%d", k))
		}
	}
	log.Print("not11: ", not11)

	options := []string{fmt.Sprintf("MAXDIST=%d", r.maxdepth), fmt.Sprintf("NODATA=%d", r.maxdepth), fmt.Sprintf("VALUES=%s", strings.Join(not11, ","))}
	if Debug {
		log.Print("options: ", options)
	}

	// JMT: failure error gets thrown.  source says that happens when
	// pointers fail validation and when distunits is set incorrectly.
	cperr := depthBand.ComputeProximity(bathyBand, options, gdal.DummyProgress, nil)
	if cperr != nil {
		log.Panicf(cperr.Error())
	}

	resXsize := depthx - 2*r.maxdepth
	resYsize := depthy - 2*r.maxdepth
	results := make([]int16, resXsize*resYsize)
	bathyerr := bathyBand.IO(gdal.Read, r.maxdepth, r.maxdepth, resXsize, resYsize, results, resXsize, resYsize, 0, 0)
	if notnil(bathyerr) {
		panic(bathyerr)
	}
	return results
}

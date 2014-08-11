package carto

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path"

	"github.com/lukeroth/gdal"
)

var Debug = false

type Layer int

const (
	LandCover = iota
	Elevation
	Bathy
	Crust
)

// world.Region may get renamed to MCRegion
type Region struct {
	// variables
	name string

	ll FloatExtents

	// mostly-constant variables
	tilesize int
	scale    int
	vscale   int
	trim     int
	sealevel int
	maxdepth int

	albers map[string]IntExtents
	wgs84  map[string]FloatExtents

	vrts  map[string]string
	files map[string]string
}

var keys = []string{"elevation", "landcover"}

const (
	wgs84_proj  = "+proj=longlat +ellps=WGS84 +datum=WGS84 +no_defs"
	albers_proj = "+proj=aea +datum=NAD83 +lat_1=29.5 +lat_2=45.5 +lat_0=23 +lon_0=-96 +x_0=0 +y_0=0 +units=m"
)

const (
	tileheight = 256
	headroom   = 16
)

func MakeRegion(name string, ll FloatExtents) Region {
	// firm defaults
	scale := 6
	vscale := 6
	trim := 0
	tilesize := 256
	sealevel := 62
	maxdepth := 30
	vrts := map[string]string{}
	files := map[string]string{}
	albers := map[string]IntExtents{}
	wgs84 := map[string]FloatExtents{}
	for _, key := range keys {
		vrts[key] = ""
		files[key] = ""
		albers[key] = IntExtents{}
		wgs84[key] = FloatExtents{}
	}

	r := Region{name: name, ll: ll, tilesize: tilesize, scale: scale, vscale: vscale, trim: trim, sealevel: sealevel, maxdepth: maxdepth, vrts: vrts, files: files, albers: albers, wgs84: wgs84}
	r.generateExtents()
	return r
}

func (r Region) maybemaketiffs() {
	// hardcoded inputs:
	// - elevation IMG
	// - Block Island region (hardcoded extents)

	// outputs:
	// - an appropriately scaled array

	// to warp, we need:
	// - source WKT (check!)
	// - destination WKT (check!)
	// - resample algorithm (check!)
	// - input dataset (check!)
	// - output dataset

	// td, nerr := ioutil.TempDir("", "")
	// if nerr != nil {
	// 	panic(nerr)
	// }
	// defer os.RemoveAll(td)
	td := "."
	tf := path.Join(td, "test.tif")

	albers := albersExtents["elevation"]
	// wgs84 := wgs84Extents["elevation"]
	txrawSize := (albers[xMax] - albers[xMin])
	tyrawSize := (albers[yMax] - albers[yMin])
	txarrSize := txrawSize / xscale
	tyarrSize := tyrawSize / yscale

	elDS, err := gdal.Open(elFile, gdal.ReadOnly)
	if err != nil {
		panic(err)
	}
	defer elDS.Close()

	elProj := elDS.Projection()

	vDS, err := elDS.AutoCreateWarpedVRT(elProj, dstWKT, resampleAlg)
	if err != nil {
		panic(err)
	}
	defer vDS.Close()
	// log.Print(vDS.Projection())
	rXsize := vDS.RasterXSize()
	rYsize := vDS.RasterYSize()
	log.Printf("Dataset size: %d, %d", rXsize, rYsize)

	vGeoTransform := vDS.GeoTransform()
	log.Printf("Albers x range %d - %d, y range %d - %d", albers[xMin], albers[xMax], albers[yMin], albers[yMax])
	log.Printf("Albers size %d, %d", txrawSize, tyrawSize)
	log.Printf("xscale %d, yscale %d", xscale, yscale)
	log.Printf("Array size %d, %d", txrawSize/xscale, tyrawSize/yscale)
	log.Printf("Origin: %f, %f", vGeoTransform[0], vGeoTransform[3])
	log.Printf("Pixel size: %f, %f", vGeoTransform[1], vGeoTransform[5])
	log.Printf("Final: %f, %f", vGeoTransform[0]+vGeoTransform[1]*float64(rXsize), vGeoTransform[3]+vGeoTransform[5]*float64(rYsize))
	xOff := int((float64(albers[xMin]) - vGeoTransform[0]) / vGeoTransform[1])
	yOff := int((float64(albers[xMax]) - vGeoTransform[3]) / vGeoTransform[5])
	xSize := int((float64(txrawSize)) / vGeoTransform[1])
	ySize := int((float64(tyrawSize)) / vGeoTransform[5])
	log.Printf("Offset: %d, %d", xOff, yOff)
	log.Printf("Size: %d, %d", xSize, ySize)

	vBand := vDS.RasterBand(1)
	log.Print("Band type: ", vBand.RasterDataType().Name())
	vMax, maxerr := vBand.GetMaximum()
	if maxerr {
		log.Print("Max: ", vMax)
	}
	vMin, minerr := vBand.GetMinimum()
	if minerr {
		log.Print("Min: ", vMin)
	}
	vNDV, ok := vBand.NoDataValue()
	if ok {
		log.Print("Nodata value: ", vNDV)
	}
	log.Print("before first IO")
	// vBuffer := make([]float32, txarrSize*tyarrSize)
	vBuffer := make([]float32, rXsize*rYsize)
	// ioerr := vBand.IO(gdal.Read, xOff, yOff, xSize, ySize, vBuffer, txarrSize, tyarrSize, 0, 0)
	ioerr := vBand.IO(gdal.Read, 0, 0, rXsize, rYsize, vBuffer, rXsize, rYsize, 0, 0)
	if notnil(ioerr) {
		panic(ioerr)
	}
	log.Print("after first IO")

	for i, val := range vBuffer {
		if val == srcNodata {
			vBuffer[i] = 0
		}
	}

	for i, val := range vBuffer {
		if val != 0 {
			log.Printf("vBuffer[%d] = %f", i, val)
			break
		}
	}

	tDriver, err := gdal.GetDriverByName("GTiff")
	if err != nil {
		panic(err)
	}

	tDS := tDriver.Create(tf, txarrSize, tyarrSize, 1, gdal.Float32, nil)
	defer tDS.Close()

	tDS.SetGeoTransform(vGeoTransform)

	log.Print("before second IO")
	tBand := tDS.RasterBand(1)
	tBand.IO(gdal.Write, 0, 0, txarrSize, tyarrSize, vBuffer, txarrSize, tyarrSize, 0, 0)
	log.Print("after second IO")

}

func (r Region) buildMap() {
	elextents := r.albers["elevation"]

	path, err := exec.LookPath("gdalwarp")
	if err != nil {
		panic(err)
	}

	warpcmd := exec.Command(path, `-q`, `-multi`, `-t_srs`, albers_proj, `-tr`, fmt.Sprintf("%d", r.scale), fmt.Sprintf("%d", r.scale), `-te`, fmt.Sprintf("%d", elextents[xMin]), fmt.Sprintf("%d", elextents[yMin]), fmt.Sprintf("%d", elextents[xMax]), fmt.Sprintf("%d", elextents[yMax]), `-r`, `cubic`, `-srcnodata`, `"-340282346638529993179660072199368212480.000"`, `-dstnodata`, `0`, fmt.Sprintf(`%s`, r.vrts["elevation"]), fmt.Sprintf(`%s`, r.files["elevation"]))

	// remove elevation file if necessary
	stat, _ := os.Stat(r.files["elevation"])
	if stat != nil {
		rerr := os.Remove(r.files["elevation"])
		if rerr != nil {
			// JMT: I don't think I need to check this
			log.Printf(rerr.Error())
		}
	}

	// run the command
	_, nerr := warpcmd.Output()
	if nerr != nil {
		panic(nerr)
	}

	// open elds
	elDS, err := gdal.Open(elFile, gdal.ReadOnly)
	if err != nil {
		panic(err)
	}
	defer elDS.Close()
	rXsize := elDS.RasterXSize()
	rYsize := elDS.RasterYSize()
	log.Printf("Dataset size: %d, %d", rXsize, rYsize)

	// get transform
	elGT := elDS.GeoTransform()
	_ = elGT

	// get band
	elBand := elDS.RasterBand(1)
	log.Print("Band type: ", elBand.RasterDataType().Name())
	xBlock, yBlock := elBand.BlockSize()
	log.Print("Block size: ", xBlock, ", ", yBlock)

	// get array
	elBuffer := make([]float32, rXsize*rYsize)
	elrerr := elBand.IO(gdal.Read, 0, 0, rXsize, rYsize, elBuffer, rXsize, rYsize, 0, 0)
	if notnil(elrerr) {
		panic(elrerr)
	}
	// get sizes

	// get elmin and elmax
	elMin, minok := elBand.GetMinimum()
	elMax, maxok := elBand.GetMaximum()
	// if none, compute
	if !minok || !maxok {
		elMin, elMax = elBand.ComputeMinMax(0)
	}
	log.Print("Min = ", elMin)
	log.Print("Max = ", elMax)
	// close elband
	// close elds
	// (covered by defers)

	// check sealevel against elmin
	minsealevel := 2
	if elMin < 0 {
		minsealevel = minsealevel + int(-1.0*elMin/float64(r.scale))
	}
	maxsealevel := tileheight - headroom

	r.sealevel = setIntValue("sealevel", r.sealevel, maxsealevel, minsealevel)
	log.Print("sealevel: ", r.sealevel)

	// check maxdepth against sealevel
	minmaxdepth := 1
	maxmaxdepth := r.sealevel - 1
	r.maxdepth = setIntValue("maxdepth", r.maxdepth, maxmaxdepth, minmaxdepth)
	log.Print("maxdepth: ", r.maxdepth)

	// check trim against elmin
	mintrim := 0
	maxtrim := max(int(elMin), mintrim)
	r.trim = setIntValue("trim", r.trim, maxtrim, mintrim)
	log.Print("trim: ", r.trim)

	// vscale depends on sealevel, trim, elmax
	eltrimmed := float64(elMax - float64(r.trim))
	elroom := float64(tileheight - headroom - r.sealevel)
	minvscale := int(math.Ceil(eltrimmed / elroom))
	// NB: no real maximum vscale
	maxvscale := 99999
	r.vscale = setIntValue("vscale", r.vscale, maxvscale, minvscale)
	log.Print("vscale: ", r.vscale)

	// BUILD A DAMNED GEOTIFF
	// still have to generate elevation array, crust array, landcover array, depth array!
}

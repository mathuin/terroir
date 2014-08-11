package carto

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"path"

	"github.com/lukeroth/gdal"
)

var Debug = false

type Layer int

const (
	NoLayer = iota
	Landcover
	Elevation
	Bathy
	Crust
)

var keys = []string{"elevation", "landcover"}

const (
	wgs84_proj  = "+proj=longlat +ellps=WGS84 +datum=WGS84 +no_defs"
	albers_proj = "+proj=aea +datum=NAD83 +lat_1=29.5 +lat_2=45.5 +lat_0=23 +lon_0=-96 +x_0=0 +y_0=0 +units=m"
)

const (
	tileheight = 256
	headroom   = 16
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

	mapfile string
}

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
	// region files will end up being stored in a directory
	// this will be stored there too.
	mapfile := "/tmp/map.tif"

	r := Region{name: name, ll: ll, tilesize: tilesize, scale: scale, vscale: vscale, trim: trim, sealevel: sealevel, maxdepth: maxdepth, vrts: vrts, files: files, albers: albers, wgs84: wgs84, mapfile: mapfile}
	r.generateExtents()
	return r
}

// JMT: exists as a potential example to convert things to warp API
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
	if Debug {
		log.Print(vDS.Projection())
	}
	rXsize := vDS.RasterXSize()
	rYsize := vDS.RasterYSize()
	if Debug {
		log.Printf("Dataset size: %d, %d", rXsize, rYsize)
	}

	vGeoTransform := vDS.GeoTransform()
	if Debug {
		log.Printf("Albers x range %d - %d, y range %d - %d", albers[xMin], albers[xMax], albers[yMin], albers[yMax])
		log.Printf("Albers size %d, %d", txrawSize, tyrawSize)
		log.Printf("xscale %d, yscale %d", xscale, yscale)
		log.Printf("Array size %d, %d", txrawSize/xscale, tyrawSize/yscale)
		log.Printf("Origin: %f, %f", vGeoTransform[0], vGeoTransform[3])
		log.Printf("Pixel size: %f, %f", vGeoTransform[1], vGeoTransform[5])
		log.Printf("Final: %f, %f", vGeoTransform[0]+vGeoTransform[1]*float64(rXsize), vGeoTransform[3]+vGeoTransform[5]*float64(rYsize))
	}
	xOff := int((float64(albers[xMin]) - vGeoTransform[0]) / vGeoTransform[1])
	yOff := int((float64(albers[xMax]) - vGeoTransform[3]) / vGeoTransform[5])
	xSize := int((float64(txrawSize)) / vGeoTransform[1])
	ySize := int((float64(tyrawSize)) / vGeoTransform[5])
	if Debug {
		log.Printf("Offset: %d, %d", xOff, yOff)
		log.Printf("Size: %d, %d", xSize, ySize)
	}
	vBand := vDS.RasterBand(1)
	if Debug {
		log.Print("Band type: ", vBand.RasterDataType().Name())
	}
	vMax, maxerr := vBand.GetMaximum()
	if maxerr {
		if Debug {
			log.Print("Max: ", vMax)
		}
	}
	vMin, minerr := vBand.GetMinimum()
	if minerr {
		if Debug {
			log.Print("Min: ", vMin)
		}
	}
	vNDV, ok := vBand.NoDataValue()
	if ok {
		if Debug {
			log.Print("Nodata value: ", vNDV)
		}
	}
	if Debug {
		log.Print("before first IO")
	}
	// vBuffer := make([]float32, txarrSize*tyarrSize)
	vBuffer := make([]float32, rXsize*rYsize)
	// ioerr := vBand.IO(gdal.Read, xOff, yOff, xSize, ySize, vBuffer, txarrSize, tyarrSize, 0, 0)
	ioerr := vBand.IO(gdal.Read, 0, 0, rXsize, rYsize, vBuffer, rXsize, rYsize, 0, 0)
	if notnil(ioerr) {
		panic(ioerr)
	}
	if Debug {
		log.Print("after first IO")
	}

	for i, val := range vBuffer {
		if val == srcNodata {
			vBuffer[i] = 0
		}
	}

	for i, val := range vBuffer {
		if val != 0 {
			if Debug {
				log.Printf("vBuffer[%d] = %f", i, val)
			}
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

	if Debug {
		log.Print("before second IO")
	}
	tBand := tDS.RasterBand(1)
	tBand.IO(gdal.Write, 0, 0, txarrSize, tyarrSize, vBuffer, txarrSize, tyarrSize, 0, 0)
	if Debug {
		log.Print("after second IO")
	}

}

func (r Region) buildMap() {
	elExtents := r.albers["elevation"]

	path, err := exec.LookPath("gdalwarp")
	if err != nil {
		panic(err)
	}

	warpcmd := exec.Command(path, `-q`, `-multi`, `-t_srs`, albers_proj, `-tr`, fmt.Sprintf("%d", r.scale), fmt.Sprintf("%d", r.scale), `-te`, fmt.Sprintf("%d", elExtents[xMin]), fmt.Sprintf("%d", elExtents[yMin]), fmt.Sprintf("%d", elExtents[xMax]), fmt.Sprintf("%d", elExtents[yMax]), `-r`, `cubic`, `-srcnodata`, `"-340282346638529993179660072199368212480.000"`, `-dstnodata`, `0`, fmt.Sprintf(`%s`, r.vrts["elevation"]), fmt.Sprintf(`%s`, r.files["elevation"]))

	// remove elevation file if necessary
	remove(r.files["elevation"])

	// run the command
	_, nerr := warpcmd.Output()
	if nerr != nil {
		panic(nerr)
	}

	// open elds
	elDS, err := gdal.Open(r.files["elevation"], gdal.ReadOnly)
	if err != nil {
		panic(err)
	}
	defer elDS.Close()
	rXsize := elDS.RasterXSize()
	rYsize := elDS.RasterYSize()
	if Debug {
		log.Printf("Dataset size: %d, %d", rXsize, rYsize)
	}

	// get transform
	elGT := elDS.GeoTransform()
	if Debug {
		log.Printf("Origin: %f, %f", elGT[0], elGT[3])
		log.Printf("Pixel Size: %f, %f", elGT[1], elGT[5])
	}
	xOff := (float64(elExtents[xMin]) - elGT[0]) / elGT[1]
	yOff := (float64(elExtents[yMax]) - elGT[3]) / elGT[5]
	if Debug {
		log.Printf("Offset: %f, %f", xOff, yOff)
	}
	xSize := float64(elExtents[xMax]-elExtents[xMin]) / elGT[1]
	ySize := float64(elExtents[yMin]-elExtents[yMax]) / elGT[5]
	if Debug {
		log.Printf("Size: %f, %f", xSize, ySize)
	}

	// get band
	elBand := elDS.RasterBand(1)
	if Debug {
		log.Print("Band type: ", elBand.RasterDataType().Name())
	}

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
	if Debug {
		log.Print("Min = ", elMin)
		log.Print("Max = ", elMax)
	}
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
	if Debug {
		log.Print("sealevel: ", r.sealevel)
	}

	// check maxdepth against sealevel
	minmaxdepth := 1
	maxmaxdepth := r.sealevel - 1
	r.maxdepth = setIntValue("maxdepth", r.maxdepth, maxmaxdepth, minmaxdepth)
	if Debug {
		log.Print("maxdepth: ", r.maxdepth)
	}

	// check trim against elmin
	mintrim := 0
	maxtrim := max(int(elMin), mintrim)
	r.trim = setIntValue("trim", r.trim, maxtrim, mintrim)
	if Debug {
		log.Print("trim: ", r.trim)
	}

	// vscale depends on sealevel, trim, elmax
	eltrimmed := float64(elMax - float64(r.trim))
	elroom := float64(tileheight - headroom - r.sealevel)
	minvscale := int(math.Ceil(eltrimmed / elroom))
	// NB: no real maximum vscale
	maxvscale := 99999
	r.vscale = setIntValue("vscale", r.vscale, maxvscale, minvscale)
	if Debug {
		log.Print("vscale: ", r.vscale)
	}

	// build a four-layer GeoTIFF
	driver, derr := gdal.GetDriverByName("GTiff")
	if derr != nil {
		panic(err)
	}

	// remove it if it already exists
	remove(r.mapfile)

	mapDS := driver.Create(r.mapfile, rXsize, rYsize, 4, gdal.Int16, nil)
	defer mapDS.Close()

	mapDS.SetGeoTransform(elGT)

	mapSRS := gdal.CreateSpatialReference("")
	mapSRS.FromProj4(albers_proj)
	mapWKT, werr := mapSRS.ToWKT()
	if notnil(werr) {
		panic(werr)
	}
	mapDS.SetProjection(mapWKT)

	// transform the elevation erray
	newelBuffer := make([]int16, len(elBuffer))
	for i, v := range elBuffer {
		newelBuffer[i] = int16((int(v-float32(r.trim)) / r.vscale) + r.sealevel)
	}

	// write the elevation array to the mapDS raster
	elRaster := mapDS.RasterBand(Elevation)
	elRaster.IO(gdal.Write, 0, 0, rXsize, rYsize, newelBuffer, rXsize, rYsize, 0, 0)

	crustBuffer := make([]int16, len(elBuffer))
	_ = crustBuffer
	// ... screeching halt until I write IDTs in Go.
}

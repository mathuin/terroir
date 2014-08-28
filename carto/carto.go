package carto

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	"github.com/mathuin/gdal"
	"github.com/mathuin/terroir/idt"
)

var Debug = false

type Layer int

const (
	NoLayer = iota
	Landcover
	Elevation
	Bathy
	Crust
	NumLayers = iota - 1
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
	albers   map[string]IntExtents
	wgs84    map[string]FloatExtents
	vrts     map[string]string

	mapfile string
}

func MakeRegion(name string, ll FloatExtents, elname string, lcname string) Region {
	scale := 6
	vscale := 6
	trim := 0
	tilesize := 256
	sealevel := 62
	maxdepth := 30
	return MakeRegionFull(name, ll, elname, lcname, scale, vscale, trim, tilesize, sealevel, maxdepth)
}

// JMT: leading dot is bad
var datasetDir = "./datasets"
var mapsDir = "./maps"

func MakeRegionFull(name string, ll FloatExtents, elname string, lcname string, scale int, vscale int, trim int, tilesize int, sealevel int, maxdepth int) Region {
	vrts := map[string]string{}
	albers := map[string]IntExtents{}
	wgs84 := map[string]FloatExtents{}
	for _, key := range keys {
		albers[key] = IntExtents{}
		wgs84[key] = FloatExtents{}
	}
	vrts["elevation"] = path.Join(datasetDir, name, elname)
	vrts["landcover"] = path.Join(datasetDir, name, lcname)
	mapfile := path.Join(mapsDir, fmt.Sprintf("%s.tif", name))

	r := Region{name: name, ll: ll, tilesize: tilesize, scale: scale, vscale: vscale, trim: trim, sealevel: sealevel, maxdepth: maxdepth, vrts: vrts, albers: albers, wgs84: wgs84, mapfile: mapfile}
	r.generateExtents()
	return r
}

func (r Region) BuildMap() {
	td, nerr := ioutil.TempDir("", r.name)
	if nerr != nil {
		panic(nerr)
	}
	defer os.RemoveAll(td)
	elfile := path.Join(td, "elevation.tif")

	elExtents := r.albers["elevation"]

	path, err := exec.LookPath("gdalwarp")
	if err != nil {
		panic(err)
	}

	warpcmd := exec.Command(path, `-q`, `-multi`, `-t_srs`, albers_proj, `-tr`, fmt.Sprintf("%d", r.scale), fmt.Sprintf("%d", r.scale), `-te`, fmt.Sprintf("%d", elExtents[xMin]), fmt.Sprintf("%d", elExtents[yMin]), fmt.Sprintf("%d", elExtents[xMax]), fmt.Sprintf("%d", elExtents[yMax]), `-r`, `cubic`, `-srcnodata`, `"-340282346638529993179660072199368212480.000"`, `-dstnodata`, `0`, fmt.Sprintf(`%s`, r.vrts["elevation"]), fmt.Sprintf(`%s`, elfile))

	// run the command
	out, nerr := warpcmd.CombinedOutput()
	if notnil(nerr) {
		log.Print(string(out))
		panic(nerr)
	}
	_ = out

	// open elds
	elDS, err := gdal.Open(elfile, gdal.ReadOnly)
	if err != nil {
		panic(err)
	}
	defer elDS.Close()
	if Debug {
		datasetInfo(elDS, "Elevation")
	}
	rXsize := elDS.RasterXSize()
	rYsize := elDS.RasterYSize()

	// get transform
	elGT := elDS.GeoTransform()
	if Debug {
		regionInfo(elGT, elExtents)
	}

	// get band
	elBand := elDS.RasterBand(1)

	// get array
	bufferLen := rXsize * rYsize
	elBuffer := make([]float32, bufferLen)
	elrerr := elBand.IO(gdal.Read, 0, 0, rXsize, rYsize, elBuffer, rXsize, rYsize, 0, 0)
	if notnil(elrerr) {
		panic(elrerr)
	}

	// get elmin and elmax
	elMin, minok := elBand.GetMinimum()
	elMax, maxok := elBand.GetMaximum()
	// if none, compute
	if !minok || !maxok {
		elMin, elMax = elBand.ComputeMinMax(0)
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
	remove(fmt.Sprintf("%s.aux.xml", r.mapfile))

	mapDS := driver.Create(r.mapfile, rXsize, rYsize, NumLayers, gdal.Int16, nil)
	defer mapDS.Close()

	mapDS.SetGeoTransform(elGT)

	mapSRS := gdal.CreateSpatialReference("")
	mapSRS.FromProj4(albers_proj)
	mapWKT, werr := mapSRS.ToWKT()
	if notnil(werr) {
		panic(werr)
	}
	mapDS.SetProjection(mapWKT)

	// transform the elevation array
	elevarr := r.elev(elBuffer)
	elRaster := mapDS.RasterBand(Elevation)
	eioerr := elRaster.IO(gdal.Write, 0, 0, rXsize, rYsize, elevarr, rXsize, rYsize, 0, 0)
	if notnil(eioerr) {
		panic(eioerr)
	}

	// write the crust array to the raster
	crustarr := r.crust(rXsize, rYsize)
	crustRaster := mapDS.RasterBand(Crust)
	crusterr := crustRaster.IO(gdal.Write, 0, 0, rXsize, rYsize, crustarr, rXsize, rYsize, 0, 0)
	if notnil(crusterr) {
		panic(crusterr)
	}

	// landcover and depth follow
	lcExtents := r.albers["landcover"]

	lcDS, err := gdal.Open(r.vrts["landcover"], gdal.ReadOnly)
	if err != nil {
		panic(err)
	}
	defer lcDS.Close()
	if Debug {
		datasetInfo(lcDS, "Landcover")
	}

	// get transform
	lcGT := lcDS.GeoTransform()
	lcf := lcExtents.floats()
	lcxmin := int((lcf[xMin] - lcGT[0]) / lcGT[1])
	lcxmax := int((lcf[xMax] - lcGT[0]) / lcGT[1])
	lcymin := int((lcf[yMax] - lcGT[3]) / lcGT[5])
	lcymax := int((lcf[yMin] - lcGT[3]) / lcGT[5])
	lcxlen := lcxmax - lcxmin
	lcylen := lcymax - lcymin
	if Debug {
		regionInfo(lcGT, lcExtents)
	}

	// get array
	lcBand := lcDS.RasterBand(1)
	lcbufferLen := lcxlen * lcylen
	lcBuffer := make([]byte, lcbufferLen)
	lcrerr := lcBand.IO(gdal.Read, lcxmin, lcymin, lcxlen, lcylen, lcBuffer, lcxlen, lcylen, 0, 0)
	if notnil(lcrerr) {
		panic(lcrerr)
	}

	lcNodata, ok := lcBand.NoDataValue()
	if !ok {
		lcNodata = 0
	}
	// nodata is treated as water, which is 11
	newlcBuffer := make([]int, lcbufferLen)
	for i, v := range lcBuffer {
		// JMT: 0 also represents nodata in NLCD
		if v == byte(lcNodata) || v == byte(0) {
			lcBuffer[i] = 11
		}
		newlcBuffer[i] = int(lcBuffer[i])
	}

	lccoordslen := lcxlen * lcylen
	lcCoords := make([][2]float64, lccoordslen)
	for i := 0; i < lccoordslen; i++ {
		lcCoords[i] = [2]float64{lcGT[0] + lcGT[1]*float64(lcxmin+i%lcxlen), lcGT[3] + lcGT[5]*float64(lcymin+i/lcxlen)}
	}

	// depth coords
	depthxlen := (lcExtents[xMax] - lcExtents[xMin]) / r.scale
	depthylen := (lcExtents[yMax] - lcExtents[yMin]) / r.scale
	depthLen := depthylen * depthxlen
	depthCoords := make([][2]int, depthLen)
	for i := 0; i < depthLen; i++ {
		depthCoords[i] = [2]int{lcExtents[xMin] + r.scale*(i%depthxlen), lcExtents[yMax] - r.scale*(i/depthxlen)}
	}

	// IDT!
	lcIDT, lcerr := idt.NewIDT(lcCoords, newlcBuffer)
	if lcerr != nil {
		panic(lcerr)
	}

	deptharr, derr := lcIDT.Call(depthCoords, 31, true)
	if derr != nil {
		panic(derr)
	}

	bathyBuffer := r.bathy(deptharr, depthxlen, depthylen)

	lcarr := []int16{}
	bathyarr := []int16{}
	for i := 0; i < depthLen; i++ {
		if i%depthxlen >= r.maxdepth && i%depthxlen < depthxlen-r.maxdepth &&
			i/depthxlen >= r.maxdepth && i/depthxlen < depthylen-r.maxdepth {
			lcarr = append(lcarr, deptharr[i])
			bathyarr = append(bathyarr, bathyBuffer[i])
		}
	}

	bathyRaster := mapDS.RasterBand(Bathy)
	bathyerr := bathyRaster.IO(gdal.Write, 0, 0, rXsize, rYsize, bathyarr, rXsize, rYsize, 0, 0)
	if notnil(bathyerr) {
		panic(bathyerr)
	}

	lcRaster := mapDS.RasterBand(Landcover)
	lcrerr = lcRaster.IO(gdal.Write, 0, 0, rXsize, rYsize, lcarr, rXsize, rYsize, 0, 0)
	if notnil(lcrerr) {
		panic(lcrerr)
	}

	if Debug {
		datasetInfo(mapDS, "Output")
	}
}

func (r Region) elev(orig []float32) []int16 {
	elBuffer := make([]int16, len(orig))
	for i, v := range orig {
		elBuffer[i] = int16((int(v-float32(r.trim)) / r.vscale) + r.sealevel)
	}
	return elBuffer
}

func datasetInfo(ds gdal.Dataset, name string) {
	log.Printf("%s dataset", name)
	log.Printf("  Dataset size: %d, %d", ds.RasterXSize(), ds.RasterYSize())
	gt := ds.GeoTransform()
	log.Printf("  Origin: %f, %f", gt[0], gt[3])
	log.Printf("  Pixel size: %f, %f", gt[1], gt[5])
	histos := datasetHistograms(ds)
	for i, v := range histos {
		log.Printf("  Band %d: %s", i+1, v)
	}
}

func regionInfo(gt [6]float64, extents Extents) {
	ef := extents.floats()
	ei := extents.ints()
	log.Printf("  Start: xmin %d, xmax %d", ei[xMin], ei[xMax])
	log.Printf("         ymin %d, ymax %d", ei[yMin], ei[yMax])
	xOff := int((ef[xMin] - gt[0]) / gt[1])
	yOff := int((ef[yMax] - gt[3]) / gt[5])
	log.Printf("  Offset: %d, %d", xOff, yOff)
	xSize := int((ef[xMax] - ef[xMin]) / gt[1])
	ySize := int((ef[yMin] - ef[yMax]) / gt[5])
	log.Printf("  Size: %d, %d", xSize, ySize)
}

type RasterInfo struct {
	datatype string
	min      float64
	max      float64
}

func (ri RasterInfo) String() string {
	return fmt.Sprintf("%s (%f, %f)", ri.datatype, ri.min, ri.max)
}

func datasetMinMaxes(ds gdal.Dataset) []RasterInfo {
	bandCount := ds.RasterCount()
	retval := make([]RasterInfo, bandCount)
	for i := 0; i < bandCount; i++ {
		rbi := i + 1
		rb := ds.RasterBand(rbi)
		rbdt := rb.RasterDataType().Name()
		rbmin, minok := rb.GetMinimum()
		rbmax, maxok := rb.GetMaximum()
		if !minok || !maxok {
			rbmin, rbmax = rb.ComputeMinMax(0)
		}
		retval[i] = RasterInfo{datatype: rbdt, min: rbmin, max: rbmax}
	}
	return retval
}

type RasterHInfo struct {
	datatype string
	buckets  map[int]int
}

func (rhi RasterHInfo) String() string {
	retval := fmt.Sprintf("%s: ", rhi.datatype)
	bucketlist := make([]string, len(rhi.buckets))

	var keys []int
	for k, _ := range rhi.buckets {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for i, k := range keys {
		bucketlist[i] = fmt.Sprintf("%d: %d", k, rhi.buckets[k])
	}
	retval += strings.Join(bucketlist, ", ")
	return retval
}

func datasetHistograms(ds gdal.Dataset) []RasterHInfo {
	bandCount := ds.RasterCount()
	retval := make([]RasterHInfo, bandCount)
	for i := 0; i < bandCount; i++ {
		rbi := i + 1
		rb := ds.RasterBand(rbi)
		rbdt := rb.RasterDataType().Name()

		// read in full array
		dsx := ds.RasterXSize()
		dsy := ds.RasterYSize()
		dsprod := dsx * dsy
		rball := make([]int16, dsprod)
		rbrerr := rb.IO(gdal.Read, 0, 0, dsx, dsy, rball, dsx, dsy, 0, 0)
		if notnil(rbrerr) {
			panic(rbrerr)
		}
		rbh := make(map[int]int)
		for _, v := range rball {
			rbh[int(v)]++
		}
		retval[i] = RasterHInfo{datatype: rbdt, buckets: rbh}
	}
	return retval
}

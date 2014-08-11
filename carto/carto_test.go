package carto

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/lukeroth/gdal"
)

// unit tests

// more like integration tests

// VRTs don't yet work...
var elFile = "/home/jmt/git/mathuin/TopoMC/downloads/elevation/imgn42w072_13.img"
var lcFile = "/media/jmt/My Book/data/landcover/2011/nlcd_2011_landcover_2011_edition_2014_03_31.img"

// Albers equal area centered over US
var dstWKT = `PROJCS["USA_Contiguous_Albers_Equal_Area_Conic",GEOGCS["GCS_North_American_1983",DATUM["North_American_Datum_1983",SPHEROID["GRS_1980",6378137,298.257222101]],PRIMEM["Greenwich",0],UNIT["Degree",0.017453292519943295]],PROJECTION["Albers_Conic_Equal_Area"],PARAMETER["False_Easting",0],PARAMETER["False_Northing",0],PARAMETER["longitude_of_center",-96],PARAMETER["Standard_Parallel_1",29.5],PARAMETER["Standard_Parallel_2",45.5],PARAMETER["latitude_of_center",37.5],UNIT["Meter",1],AUTHORITY["EPSG","102003"]]`

var resampleAlg = gdal.ResampleAlg(2)

var xscale = 6
var yscale = 6

var srcNodata = float32(-340282346638529993179660072199368212480.0)

var dstNodata = 0

var llextents = FloatExtents{-71.533, -71.62, 41.238, 41.142}

var albersExtents = map[string]IntExtents{
	"elevation": {2015232, 2002944, 2291712, 2273280},
	"landcover": {2015412, 2002764, 2291892, 2273100},
}

var wgs84Extents = map[string]FloatExtents{
	"elevation": {-71.48245416519951, -71.68161580988412, 41.30762685177235, 41.12049902574309},
	"landcover": {-71.47980755995131, -71.68425418133639, 41.309590101757536, 41.11853504487958},
}

var buildMap_tests = []struct {
	ll     FloatExtents
	elvrt  string
	elfile string
}{
	{
		FloatExtents{-71.533, -71.62, 41.238, 41.142},
		"/media/jmt/My Book/data/elevation/13/elevation13.vrt",
		"elevation.tif",
	},
}

func Test_buildMap(t *testing.T) {
	td, nerr := ioutil.TempDir("", "")
	if nerr != nil {
		panic(nerr)
	}
	// defer os.RemoveAll(td)
	for _, tt := range buildMap_tests {
		realfile := path.Join(td, tt.elfile)
		r := MakeRegion("Pie", tt.ll)
		r.vrts["elevation"] = tt.elvrt
		r.files["elevation"] = realfile
		r.buildMap()
	}
}

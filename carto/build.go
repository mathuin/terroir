package carto

import (
	"fmt"
	"log"
	"runtime"

	"sync"

	"github.com/mathuin/gdal"
	"github.com/mathuin/terroir/world"
)

// JMT: convert to XZ, biome value, and list 0->n of points
type Column struct {
	xz     world.XZ
	biome  int           // was string
	blocks []world.Block // []string
}

func makeColumn(xz world.XZ, biome string, blocks []string) Column {
	bval, ok := world.Biome[biome]
	if !ok {
		log.Panicf("%s not found in world.Biome", biome)
	}
	bvals := make([]world.Block, len(blocks))
	for i, v := range blocks {
		bval, err := world.BlockNamed(v)
		if err != nil {
			panic(err)
		}
		bvals[i] = *bval
	}
	return Column{xz: xz, biome: bval, blocks: bvals}
}

func (r Region) genFeatures(in chan Feature) {
	ds, err := gdal.Open(r.mapfile, gdal.ReadOnly)
	if err != nil {
		panic(err)
	}
	if Debug {
		datasetInfo(ds, "genFeatures Input")
	}
	inx := ds.RasterXSize()
	iny := ds.RasterYSize()
	srs := gdal.CreateSpatialReference(ds.Projection())
	bufferLen := inx * iny

	lcarr := make([]int16, bufferLen)
	lcBand := ds.RasterBand(Landcover)
	lcrerr := lcBand.IO(gdal.Read, 0, 0, inx, iny, lcarr, inx, iny, 0, 0)
	if notnil(lcrerr) {
		panic(lcrerr)
	}

	// shapefile driver
	outdrv := gdal.OGRDriverByName("Memory")
	outDS, ok := outdrv.Create("out", nil)
	if !ok {
		panic(fmt.Errorf("OGR Driver Create Fail"))
	}
	outLayer := outDS.CreateLayer("polygons", srs, gdal.GT_Polygon, nil)

	// field definition
	outField := gdal.CreateFieldDefinition("lc", gdal.FT_Integer)
	outLayer.CreateField(outField, false)
	field := 0

	// options!
	options := []string{}

	// do it!
	err = lcBand.Polygonize(lcBand, outLayer, field, options, gdal.DummyProgress, nil)
	if notnil(err) {
		panic(err)
	}

	// iterate over features
	fc, ok := outLayer.FeatureCount(true)
	if !ok {
		panic(fmt.Errorf("outLayer.FeatureCount NOT OK"))
	}
	if Debug {
		log.Print("outLayer.FeatureCount(true): ", fc)
	}
	outLayer.ResetReading()
	for i := 0; i < fc; i++ {
		in <- Feature{outLayer.NextFeature()}
	}
	close(in)
}

func (r *Region) BuildWorld() (*world.World, error) {
	w := world.MakeWorld(r.name)
	w.SetRandomSeed(0)
	spawnpt := world.MakePoint(0, 0, 0)

	in := make(chan Feature)
	out := make(chan Column)

	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU()
	// if Debug {
	// 	log.Print("debug mode - only starting one worker")
	// 	numWorkers = 1
	// }

	// JMT: consider moving memo up here

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r.processFeatures(in, out, i)
		}(i)
	}
	go func() { wg.Wait(); close(out) }()
	go r.genFeatures(in)

	columncount := 0
	for column := range out {
		columncount++

		w.SetBiome(column.xz, byte(column.biome))

		pt := column.xz.Point(int32(0))
		for k, v := range column.blocks {
			pt.Y = int32(k)
			w.SetBlock(pt, v)
		}

		if pt.Y > spawnpt.Y {
			if Debug {
				log.Printf("new spawn: %s", pt)
			}
			spawnpt = pt
		}

		// JMT: naive lighting here
		w.SetSkyLight(pt, 15)
	}

	w.SetSpawn(spawnpt)

	return &w, nil
}

var Arrind2Debug = false

func arrind2(x int32, y int32, inx int32, iny int32, gti [6]int32) (int32, error) {
	if Arrind2Debug {
		log.Printf("x: %d, y: %d, inx: %d", x, y, inx)
		log.Printf("0: %d, 1: %d, 2: %d, 3: %d, 4: %d, 5: %d",
			gti[0], gti[1], gti[2], gti[3], gti[4], gti[5])
		log.Printf("x-0: %d, y-3: %d", x-gti[0], y-gti[3])
		log.Printf("x-0/1: %d, y-3/5: %d", (x-gti[0])/gti[1], (y-gti[3])/gti[5])
	}
	realx := (x - gti[0]) / gti[1]
	if realx < 0 {
		return 0, fmt.Errorf("realx %d < 0", realx)
	}
	if realx > inx {
		return 0, fmt.Errorf("realx %d >= inx %d", realx, inx)
	}
	realy := (y-gti[3])/gti[5] - 1
	if realy < 0 {
		return 0, fmt.Errorf("realx %d < 0", realx)
	}
	if realy > iny {
		return 0, fmt.Errorf("realy %d > iny %d", realy, iny)
	}
	return realx + realy*inx, nil
}

func (r *Region) processFeatures(in chan Feature, out chan Column, i int) {
	ds, err := gdal.Open(r.mapfile, gdal.ReadOnly)
	if err != nil {
		panic(err)
	}
	if Debug && i == 0 {
		datasetInfo(ds, "processFeatures")
	}
	inx := ds.RasterXSize()
	iny := ds.RasterYSize()
	bufferLen := inx * iny

	lcarr := make([]int16, bufferLen)
	lcBand := ds.RasterBand(Landcover)
	lcrerr := lcBand.IO(gdal.Read, 0, 0, inx, iny, lcarr, inx, iny, 0, 0)
	if notnil(lcrerr) {
		panic(lcrerr)
	}

	elevarr := make([]int16, bufferLen)
	elevBand := ds.RasterBand(Elevation)
	elevrerr := elevBand.IO(gdal.Read, 0, 0, inx, iny, elevarr, inx, iny, 0, 0)
	if notnil(elevrerr) {
		panic(elevrerr)
	}

	bathyarr := make([]int16, bufferLen)
	bathyBand := ds.RasterBand(Bathy)
	bathyrerr := bathyBand.IO(gdal.Read, 0, 0, inx, iny, bathyarr, inx, iny, 0, 0)
	if notnil(bathyrerr) {
		panic(bathyrerr)
	}

	crustarr := make([]int16, bufferLen)
	crustBand := ds.RasterBand(Crust)
	crustrerr := crustBand.IO(gdal.Read, 0, 0, inx, iny, crustarr, inx, iny, 0, 0)
	if notnil(crustrerr) {
		panic(crustrerr)
	}

	processed := 0

	// first pass at memoize
	// memo := make(map[string]Column)

	for f := range in {
		processed++

		head := fmt.Sprintf("%d: feature #%d", i, processed)

		// if Debug {
		// 	log.Printf("%s begins", head)
		// }
		pts := f.Points(ds, head)
		if len(pts) == 0 {
			log.Printf("%s: No points in geometry!", head)
			log.Print("SCRATCH ONE FEATURE")
			continue
		}

		lc := f.LCValue()
		switch lc {
		case 11:
			// "open water"
			for _, pt := range pts {
				elev := elevarr[pt.index]
				bathy := bathyarr[pt.index]
				crust := crustarr[pt.index]

				// key := fmt.Sprintf("%d|%d|%d|%d", lc, elev, bathy, crust)

				var biome string
				if int(bathy) >= r.maxdepth-1 {
					biome = "Deep Ocean"
				} else {
					biome = "Ocean"
				}

				// if col, ok := memo[key]; ok {
				// 	col.xz = pt.xz
				// 	out <- col
				// 	continue
				// }
				blocks := make([]string, elev)
				for y := int16(0); y < elev; y++ {
					if y == 0 {
						blocks[y] = "Bedrock"
					} else if y < (elev - bathy - crust) {
						blocks[y] = "Stone"
					} else if y < (elev - bathy) {
						blocks[y] = "Gravel"
					} else {
						blocks[y] = "Water"
					}
				}
				col := makeColumn(pt.xz, biome, blocks)
				// memo[key] = col
				out <- col
			}
		case 31:
			// "barren land"
			for _, pt := range pts {
				elev := elevarr[pt.index]
				// bathy := bathyarr[pt.index]
				crust := crustarr[pt.index]

				// key := fmt.Sprintf("%d|%d|%d|%d", lc, elev, bathy, crust)

				var biome string
				if elev > 92 {
					biome = "Desert Hills"
				} else {
					biome = "Desert"
				}

				// if col, ok := memo[key]; ok {
				// 	col.xz = pt.xz
				// 	out <- col
				// 	continue
				// }
				blocks := make([]string, elev)
				for y := int16(0); y < elev; y++ {
					if y == 0 {
						blocks[y] = "Bedrock"
					} else if y < (elev - crust - 1) {
						blocks[y] = "Stone"
					} else if y < elev-1 {
						blocks[y] = "Sandstone"
					} else {
						blocks[y] = "Sand"
					}
				}
				col := makeColumn(pt.xz, biome, blocks)
				// memo[key] = col
				out <- col
			}
		case 41:
			fallthrough
		case 42:
			fallthrough
		case 43:
			// "forest"
			for _, pt := range pts {
				elev := elevarr[pt.index]
				// bathy := bathyarr[pt.index]
				crust := crustarr[pt.index]

				// key := fmt.Sprintf("%d|%d|%d|%d", lc, elev, bathy, crust)

				var biome string
				if elev > 92 {
					biome = "Forest Hills"
				} else {
					biome = "Forest"
				}

				// if col, ok := memo[key]; ok {
				// 	col.xz = pt.xz
				// 	out <- col
				// 	continue
				// }
				blocks := make([]string, elev)
				for y := int16(0); y < elev; y++ {
					if y == 0 {
						blocks[y] = "Bedrock"
					} else if y < (elev - crust - 1) {
						blocks[y] = "Stone"
					} else if y < elev-1 {
						blocks[y] = "Dirt"
					} else {
						blocks[y] = "Grass Block"
					}
				}
				col := makeColumn(pt.xz, biome, blocks)
				// memo[key] = col
				out <- col
			}
		case 90:
			fallthrough
		case 95:
			// "swampland"
			for _, pt := range pts {
				elev := elevarr[pt.index]
				// bathy := bathyarr[pt.index]
				crust := crustarr[pt.index]

				// key := fmt.Sprintf("%d|%d|%d|%d", lc, elev, bathy, crust)

				var biome string
				if elev > 92 {
					biome = "Swampland M"
				} else {
					biome = "Swampland"
				}

				// if col, ok := memo[key]; ok {
				// 	col.xz = pt.xz
				// 	out <- col
				// 	continue
				// }
				blocks := make([]string, elev)
				for y := int16(0); y < elev; y++ {
					if y == 0 {
						blocks[y] = "Bedrock"
					} else if y < (elev - crust - 1) {
						blocks[y] = "Stone"
					} else if y < elev-1 {
						blocks[y] = "Dirt"
					} else {
						blocks[y] = "Grass Block"
					}
				}
				col := makeColumn(pt.xz, biome, blocks)
				// memo[key] = col
				out <- col
			}
		default:
			// anything else
			for _, pt := range pts {
				elev := elevarr[pt.index]
				// bathy := bathyarr[pt.index]
				crust := crustarr[pt.index]

				// key := fmt.Sprintf("%d|%d|%d|%d", lc, elev, bathy, crust)

				var biome string
				if elev > 152 {
					// Hills+ too
					biome = "Extreme Hills M"
				} else if elev > 122 {
					// Hills+ too
					biome = "Extreme Hills"
				} else if elev > 92 {
					biome = "Extreme Hills Edge"
				} else {
					// Sunflower Plains?
					biome = "Plains"
				}

				// if col, ok := memo[key]; ok {
				// 	col.xz = pt.xz
				// 	out <- col
				// 	continue
				// }
				blocks := make([]string, elev)
				for y := int16(0); y < elev; y++ {
					if y == 0 {
						blocks[y] = "Bedrock"
					} else if y < (elev - crust - 1) {
						blocks[y] = "Stone"
					} else if y < elev-1 {
						blocks[y] = "Dirt"
					} else {
						blocks[y] = "Grass Block"
					}
				}
				col := makeColumn(pt.xz, biome, blocks)
				// memo[key] = col
				out <- col
			}
		}
	}
}

type Feature struct {
	gdal.Feature
}

// returns true if the feature is valid
func (f Feature) isValid() bool {
	return f.Geometry().IsValid()
}

// returns the landcover value
func (f Feature) LCValue() int {
	// field 0 is the only field here
	return f.FieldAsInteger(0)
}

// generates list of points and sends them to a channel
func (f Feature) genPoints(in chan world.XZ, gti [6]int32) {
	g := f.Geometry()
	e := g.Envelope()
	eminx := int32(e.MinX())
	eminy := int32(e.MinY())
	emaxx := int32(e.MaxX())
	emaxy := int32(e.MaxY())

	for y := eminy; y < emaxy; y -= gti[5] {
		for x := eminx; x < emaxx; x += gti[1] {
			in <- world.XZ{X: x, Z: y}
		}
	}
	close(in)
}

type XZIndex struct {
	xz    world.XZ
	index int32
}

func makeXZIndex(xz world.XZ, index int32, gti [6]int32) XZIndex {
	return XZIndex{xz: world.XZ{X: xz.X / gti[1], Z: xz.Z / gti[5]}, index: index}
}

func (f Feature) Points(ds gdal.Dataset, head string) []XZIndex {
	inx := ds.RasterXSize()
	iny := ds.RasterYSize()
	gt := ds.GeoTransform()
	srs := gdal.CreateSpatialReference(ds.Projection())

	lcarr := make([]int16, inx*iny)
	lcBand := ds.RasterBand(Landcover)
	lcrerr := lcBand.IO(gdal.Read, 0, 0, inx, iny, lcarr, inx, iny, 0, 0)
	if notnil(lcrerr) {
		panic(lcrerr)
	}

	var gti [6]int32
	for i, v := range gt {
		gti[i] = int32(v)
	}

	in := make(chan world.XZ)
	out := make(chan XZIndex)

	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ppHead := fmt.Sprintf("%s %d", head, i)
			f.processPoints(in, out, inx, iny, gti, srs, lcarr, ppHead)
		}(i)
	}
	go func() { wg.Wait(); close(out) }()
	go f.genPoints(in, gti)

	pts := []XZIndex{}
	ptcount := 0
	for pt := range out {
		if Debug {
			if ptcount > 1 && ptcount%10000 == 0 {
				log.Printf("%s: %d columns", head, ptcount)
			}
		}
		ptcount++
		pts = append(pts, pt)
	}
	return pts
}

func (f Feature) processPoints(in chan world.XZ, out chan XZIndex, inx int, iny int, gti [6]int32, srs gdal.SpatialReference, lcarr []int16, ppHead string) {
	g := f.Geometry()

	lc := f.LCValue()

	for xz := range in {
		index, aerr := arrind2(xz.X, xz.Z, int32(inx), int32(iny), gti)
		// JMT: arrind returns nil if coordinates are invalid
		if aerr != nil {
			log.Printf("%s: aerr was not nil: %s", ppHead, aerr.Error())
			continue
		}
		if lcarr[index] != int16(lc) {
			continue
		}
		wkt := fmt.Sprintf("POINT (%f %f)", float64(xz.X)+0.5*float64(gti[1]), float64(xz.Z)-0.5*float64(gti[5]))
		pt, err := gdal.CreateFromWKT(wkt, srs)
		if notnil(err) {
			log.Printf("%s: gdal.CreateFromWKT() error: %s", ppHead, err.Error())
			continue
		}
		if g.Contains(pt) {
			out <- makeXZIndex(xz, index, gti)
		}
	}
}

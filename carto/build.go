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
type Out struct {
	blocks map[world.Point]world.Block
	spawn  world.Point
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

func (r *Region) buildWorld() (*world.World, error) {
	w := world.MakeWorld(r.name)
	w.SetRandomSeed(0)
	// JMT: need a sane storage location for files
	w.SetSaveDir(".")
	spawnpt := world.MakePoint(0, 0, 0)

	in := make(chan Feature)
	out := make(chan Out)

	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU() * runtime.NumCPU()
	if Debug {
		log.Print("debug mode - only starting one worker")
		numWorkers = 1
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r.processFeatures(in, out, i)
		}(i)
	}
	go func() { wg.Wait(); close(out) }()
	go r.genFeatures(in)

	blockcount := 0
	for block := range out {
		blockcount++
		if Debug {
			// log.Printf("set #%d of blocks received: %d blocks", blockcount, len(block.blocks))
		}
		for k, v := range block.blocks {
			w.SetBlock(k, &v)
		}

		if spawnpt.Y < block.spawn.Y || spawnpt.Z == 0 {
			if Debug {
				log.Printf("new spawn: %s", block.spawn)
			}
			spawnpt = block.spawn
		}
		if Debug {
			// log.Printf("set #%d of %d blocks processed", blockcount, len(block.blocks))
		}
	}

	w.SetSpawn(spawnpt)

	return &w, nil
}

var Arrind2Debug = false

func arrind2(x int, y int, inx int, iny int, gt [6]float64) (int, error) {
	if Arrind2Debug {
		log.Printf("x: %d, y: %d, inx: %d", x, y, inx)
		log.Printf("0: %d, 1: %d, 2: %d, 3: %d, 4: %d, 5: %d",
			int(gt[0]), int(gt[1]), int(gt[2]), int(gt[3]), int(gt[4]), int(gt[5]))
		log.Printf("x-0: %d, y-3: %d", x-int(gt[0]), y-int(gt[3]))
		log.Printf("x-0/1: %d, y-3/5: %d", (x-int(gt[0]))/int(gt[1]), (y-int(gt[3]))/int(gt[5]))
	}
	realx := (x - int(gt[0])) / int(gt[1])
	if realx < 0 {
		return 0, fmt.Errorf("realx %d < 0", realx)
	}
	if realx > inx {
		return 0, fmt.Errorf("realx %d >= inx %d", realx, inx)
	}
	realy := (y-int(gt[3]))/int(gt[5]) - 1
	if realy < 0 {
		return 0, fmt.Errorf("realx %d < 0", realx)
	}
	if realy > iny {
		return 0, fmt.Errorf("realy %d > iny %d", realy, iny)
	}
	return realx + realy*inx, nil
}

func (r *Region) processFeatures(in chan Feature, out chan Out, i int) {
	ds, err := gdal.Open(r.mapfile, gdal.ReadOnly)
	if err != nil {
		panic(err)
	}
	if Debug {
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

	for f := range in {
		// if !f.isValid() {
		// 	g := f.Geometry()
		// 	if !g.IsValid() {
		// 		newg := g.SimplifyPreservingTopology(4)
		// 		if !newg.IsValid() {
		// 			panic("WTF D00D")
		// 		}
		// 	}
		// 	// JMT: for invalid geometry, exit with no points
		// 	log.Printf("%d: Invalid geometry!", i)
		// 	continue
		// }
		processed++

		head := fmt.Sprintf("%d: feature #%d", i, processed)

		if Debug {
			log.Printf("%s begins", head)
		}
		pts := f.Points(ds, head, lcarr)
		if len(pts) == 0 {
			log.Printf("%s: No points in geometry!", head)
			log.Print("SCRATCH ONE FEATURE")
			continue
		}

		lc := f.LCValue()

		lame := map[int]string{
			11: "Water",
			21: "Stone",
			22: "Stone",
			23: "Stone",
			24: "Stone",
			31: "Sand",
			41: "Dirt",
			42: "Dirt",
			43: "Dirt",
			71: "Grass Block",
			90: "Water",
			95: "Water",
		}

		var putit *world.Block
		if val, ok := lame[lc]; ok {
			putit, err = world.BlockNamed(val)
			if err != nil {
				panic(err)
			}
		} else {
			log.Printf("%s: LC value %d not found!", head, lc)
			putit, err = world.BlockNamed("Dirt")
			if err != nil {
				panic(err)
			}
		}

		// blocks := make(map[world.Point]world.Block)
		spawn := world.Point{}

		bedrockBlock, _ := world.BlockNamed("Bedrock")

		totblks := 0
		for _, pt := range pts {
			blocks := make(map[world.Point]world.Block)

			elev := int32(elevarr[pt[2]])
			// bathy := bathyarr[pt[2]]
			// crust := crustarr[pt[2]]

			// build xz from pt
			xz := world.XZ{X: int32(pt[0]), Z: int32(pt[1])}
			blocks[xz.Point(0)] = *bedrockBlock
			for y := 1; y <= int(elev); y++ {
				blocks[xz.Point(int32(y))] = *putit
			}

			if elev > spawn.Y {
				spawn = xz.Point(elev)
			}
			totblks += len(blocks)
			out <- Out{blocks: blocks, spawn: spawn}
		}
		if Debug {
			log.Printf("%s ends, %d blocks, spawn: %v", head, totblks, spawn)
		}
		// out <- Out{blocks: blocks, spawn: spawn}
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

// returns a list of points from the feature
// x, y, index
func (f Feature) Points(ds gdal.Dataset, head string, lcarr []int16) [][3]int {
	inx := ds.RasterXSize()
	iny := ds.RasterYSize()
	gt := ds.GeoTransform()
	srs := gdal.CreateSpatialReference(ds.Projection())

	pts := [][3]int{}
	lc := f.LCValue()

	if Debug {
		log.Printf("%s: lc: %d", head, lc)
	}

	g := f.Geometry()
	e := g.Envelope()
	eminx := int(e.MinX())
	eminy := int(e.MinY())
	emaxx := int(e.MaxX())
	emaxy := int(e.MaxY())
	totcols := ((emaxy - eminy) / -1 * int(gt[5])) * ((emaxx - eminx) / int(gt[1]))
	// JMT: is there a way to determine whether envelopes exceed bounds?
	if Debug {
		log.Printf("%s: Envelope: (%d, %d) -> (%d, %d), %d total possible columns", head, eminx, eminy, emaxx, emaxy, totcols)
	}

	posscols := 0
	tenth := int(float64(totcols) / float64(10.0))
	for y := eminy; y < emaxy; y -= int(gt[5]) {
		for x := eminx; x < emaxx; x += int(gt[1]) {
			posscols++
			if totcols > 10000 && posscols%tenth == 0 {
				log.Printf("%s: %d of %d cols", head, posscols, totcols)
			}
			index, aerr := arrind2(x, y, inx, iny, gt)
			// JMT: arrind returns nil if coordinates are invalid
			if aerr != nil {
				log.Printf("%s: aerr was not nil: %s", head, aerr.Error())
				continue
			}
			if Debug {
				// log.Printf("%d: x: %d, y: %d, index: %d", i, x, y, index)
			}
			if lcarr[index] != int16(lc) {
				continue
			}
			wkt := fmt.Sprintf("POINT (%f %f)", float64(x)+0.5*gt[1], float64(y)-0.5*gt[5])
			pt, err := gdal.CreateFromWKT(wkt, srs)
			if notnil(err) {
				log.Printf("%s: gdal.CreateFromWKT() error: %s", head, err.Error())
				continue
			}
			if g.Contains(pt) {
				pts = append(pts, [3]int{x, y, index})
			}
		}
	}

	if Debug {
		log.Printf("%s: %d pts (%d possible cols) in feature", head, len(pts), posscols)
	}

	return pts
}

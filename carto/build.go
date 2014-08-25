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
	biome  string
	blocks []string
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
	out := make(chan Column)

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

	columncount := 0
	for column := range out {
		columncount++

		b, ok := world.Biome[column.biome]
		if !ok {
			err := fmt.Errorf("biome %s not in world.Biome", column.biome)
			panic(err)
		}
		w.SetBiome(column.xz, byte(b))

		for k, v := range column.blocks {
			pt := world.Point{X: column.xz.X, Y: int32(k), Z: column.xz.Z}
			b, err := world.BlockNamed(v)
			if err != nil {
				panic(err)
			}
			w.SetBlock(pt, b)
		}

		topBlock := world.Point{X: column.xz.X, Y: int32(len(column.blocks)), Z: column.xz.Z}

		if topBlock.Y > spawnpt.Y {
			if Debug {
				log.Printf("new spawn: %s", topBlock)
			}
			spawnpt = topBlock
		}

		// JMT: naive lighting here
		w.SetSkyLight(topBlock, 15)
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

func (r *Region) processFeatures(in chan Feature, out chan Column, i int) {
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

		switch lc {
		case 11:
			// "open water"
			for _, pt := range pts {
				elev := elevarr[pt[2]]
				bathy := bathyarr[pt[2]]
				crust := crustarr[pt[2]]

				col := Column{}
				col.xz = world.XZ{X: int32(pt[0]), Z: int32(pt[1])}

				if int(bathy) <= r.maxdepth-1 {
					col.biome = "Deep Ocean"
				} else {
					col.biome = "Ocean"
				}

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
				col.blocks = blocks
				out <- col
			}
		default:
			// anything else
			for _, pt := range pts {
				elev := elevarr[pt[2]]
				crust := crustarr[pt[2]]

				col := Column{}
				col.xz = world.XZ{X: int32(pt[0]), Z: int32(pt[1])}
				col.biome = "Plains"

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
				col.blocks = blocks
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
				if Debug {
					log.Printf("%s: %d of %d cols", head, posscols, totcols)
				}
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

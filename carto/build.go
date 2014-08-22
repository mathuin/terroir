package carto

import (
	"fmt"
	"log"
	"runtime"

	"sort"
	"strings"
	"sync"

	"github.com/mathuin/gdal"
	"github.com/mathuin/terroir/world"
)

type Out struct {
	blocks map[world.Point]world.Block
	spawn  world.Point
}

func (r Region) genFeatures(in chan gdal.Feature) {
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
	options := []string{""}

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
		in <- outLayer.NextFeature()
	}
	close(in)
}

func (r *Region) buildWorld() (*world.World, error) {
	w := world.MakeWorld(r.name)
	w.SetRandomSeed(0)
	// JMT: need a sane storage location for files
	w.SetSaveDir(".")
	spawnpt := world.MakePoint(0, 0, 0)

	in := make(chan gdal.Feature)
	out := make(chan Out)

	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU()
	if Debug {
		log.Print("debug mode - only starting one worker")
		numWorkers = 1
	}

	for i := 0; i < numWorkers; i++ {
		go func(i int) {
			wg.Add(1)
			r.processFeatures(in, out, i)
			wg.Done()
		}(i)
	}
	go func() { wg.Wait(); close(out) }()
	go r.genFeatures(in)

	flag := false
	for block := range out {
		if !flag {
			flag = true
			logout := fmt.Sprintf("First block: ")
			blocklist := make([]string, len(block.blocks))
			var keys world.Points
			for k, _ := range block.blocks {
				keys = append(keys, k)
			}
			sort.Sort(keys)
			for i, k := range keys {
				blocklist[i] = fmt.Sprintf("%s: %s", k, block.blocks[k])
			}
			logout += strings.Join(blocklist, ", ")
			// log.Print(logout)
		}
		if Debug {
			log.Printf("new set of blocks received: %d blocks", len(block.blocks))
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

func (r *Region) processFeatures(in chan gdal.Feature, out chan Out, i int) {
	ds, err := gdal.Open(r.mapfile, gdal.ReadOnly)
	if err != nil {
		panic(err)
	}
	if Debug {
		datasetInfo(ds, "processFeatures")
	}
	inx := ds.RasterXSize()
	iny := ds.RasterYSize()
	gt := ds.GeoTransform()
	srs := gdal.CreateSpatialReference(ds.Projection())
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

	for f := range in {
		pts := [][3]int{}
		// field 0 is the only field here
		lc := f.FieldAsInteger(0)

		if Debug {
			log.Printf("%d: lc: %d", i, lc)
		}

		g := f.Geometry()
		if !g.IsValid() {
			// JMT: for invalid geometry, exit with no points
			log.Printf("%d: Invalid geometry!", i)
			continue
		}
		e := g.Envelope()
		eminx := int(e.MinX())
		eminy := int(e.MinY())
		emaxx := int(e.MaxX())
		emaxy := int(e.MaxY())
		// JMT: is there a way to determine whether envelopes exceed bounds?
		if Debug {
			log.Printf("%d: Envelope: (%d, %d) -> (%d, %d)", i, eminx, eminy, emaxx, emaxy)
		}

		// JMT: for now, just iterate through all the points in this particular polygon
		posspts := 0
		for y := eminy; y < emaxy; y -= int(gt[5]) {
			for x := eminx; x < emaxx; x += int(gt[1]) {
				posspts++
				index, aerr := arrind2(x, y, inx, iny, gt)
				// JMT: arrind returns nil if coordinates are invalid
				if aerr != nil {
					log.Printf("%d: aerr was not nil: %s", i, aerr.Error())
					continue
				}
				inpt := [3]int{x, y, index}
				if Debug {
					// log.Printf("%d: x: %d, y: %d, index: %d (%d)", i, x, y, index, lcarr[index])
				}
				if lcarr[index] != int16(lc) {
					continue
				}
				wkt := fmt.Sprintf("POINT (%f %f)", float64(x)+0.5*gt[1], float64(y)-0.5*gt[5])
				pt, err := gdal.CreateFromWKT(wkt, srs)
				if err != nil && err.Error() != "No Error" {
					log.Printf("%d: gdal.CreateFromWKT() error: %s", i, err.Error())
					continue
				}
				if g.Contains(pt) {
					pts = append(pts, inpt)
				}
			}
		}

		if Debug {
			log.Printf("%d: %d pts (%d possible) in feature", i, len(pts), posspts)
		}

		cout := new(Out)
		cout.blocks = make(map[world.Point]world.Block)
		cout.spawn = world.Point{}
		bedrockBlock, _ := world.BlockNamed("Bedrock")
		// waterBlock, _ := world.BlockNamed("Water")
		dirtBlock, _ := world.BlockNamed("Dirt")
		stoneBlock, _ := world.BlockNamed("Stone")
		for _, pt := range pts {
			var putit *world.Block
			switch lc {
			case 11:
				putit = stoneBlock
			default:
				putit = dirtBlock
			}
			// elev := int32(elevarr[pt[2]])
			elev := int32(62)
			// bathy := bathyarr[pt[2]]
			// crust := crustarr[pt[2]]

			// build xz from pt
			xz := world.XZ{X: int32(pt[0]), Z: int32(pt[1])}
			cout.blocks[xz.Point(0)] = *bedrockBlock
			for y := 1; y <= int(elev); y++ {
				cout.blocks[xz.Point(int32(y))] = *putit
			}

			if elev > cout.spawn.Y {
				cout.spawn = xz.Point(elev)
			}
		}
		if Debug {
			log.Printf("%d: %d pts in blocks", i, len(cout.blocks))
			log.Printf("%d: spawn: %v", i, cout.spawn)
		}
		out <- *cout
	}
}

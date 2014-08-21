package carto

import (
	"fmt"
	"log"

	"github.com/mathuin/gdal"
	"github.com/mathuin/terroir/world"
)

func (r *Region) newbiome(inx int, iny int, gt [6]float64, lcarr []int16, elevarr []int16, bathyarr []int16) (biomearr []int16, berr error) {
	product := "NLCD 2011"
	bufferLen := len(lcarr)
	biomearr = make([]int16, bufferLen)

	// mem driver!
	memdrv, err1 := gdal.GetDriverByName("MEM")
	if err1 != nil {
		return biomearr, err1
	}

	srcDS := memdrv.Create("src", inx, iny, 1, gdal.Int16, nil)
	// top left x, west-east resolution, 0,
	// top left y, 0, north-south resolution (negative)
	// JMT: check this with production
	srcDS.SetGeoTransform(gt)

	// create a source band from that array
	srcBand := srcDS.RasterBand(1)
	err := srcBand.IO(gdal.Write, 0, 0, inx, iny, lcarr, inx, iny, 0, 0)
	if err != nil && err.Error() != "No Error" {
		return biomearr, err
	}

	// shapefile driver
	outdrv := gdal.OGRDriverByName("Memory")
	outDS, ok := outdrv.Create("out", nil)
	if !ok {
		return biomearr, fmt.Errorf("OGR Driver Create Fail")
	}

	// projection sigh
	outSRS := gdal.CreateSpatialReference("")
	// outSRS.FromProj4(albers_proj)
	outLayer := outDS.CreateLayer("polygons", outSRS, gdal.GT_Polygon, nil)

	// field definition
	outField := gdal.CreateFieldDefinition("lc", gdal.FT_Integer)
	outLayer.CreateField(outField, false)
	field := 0

	// options!
	options := []string{""}

	// do it!
	err = srcBand.Polygonize(srcBand, outLayer, field, options, gdal.DummyProgress, nil)
	if notnil(err) {
		return biomearr, err
	}

	// iterate over features
	fc, ok := outLayer.FeatureCount(true)
	if !ok {
		return biomearr, fmt.Errorf("outLayer.FeatureCount NOT OK")
	}
	if Debug {
		log.Print("outLayer.FeatureCount(true): ", fc)
	}
	outLayer.ResetReading()
	for i := 0; i < fc; i++ {
		if Debug {
			log.Printf("Feature %d", i)
		}
		f := outLayer.NextFeature()
		lc, pts, err := traverse_feature(f, field, inx, iny, gt)
		if err != nil {
			panic(err)
		}
		// both class and type information available
		lcct, err := TerrainType(lc).Data(product)
		lcc, lct := lcct[0], lcct[1]
		switch lcc {
		case "Water":
			switch lct {
			case "Open Water":
				// areas of open water, generally with less than 25%
				// cover of vegetation or soil

				// Biomes to consider:
				// - "River": if we ever get hydro data...

				// repeated bit of code here
				// pts: list of points
				// bathyarr, elevarr: need pointer
				for _, pt := range pts {
					var biome string
					// unique to us here
					bathy := bathyarr[pt]
					if bathy > int16(r.maxdepth-1) {
						biome = "Deep Ocean"
					} else {
						biome = "Ocean"
					}
					//not unique to us here
					val, ok := world.Biome[biome]
					if !ok {
						log.Printf("%s is not a valid biome for %d!", biome, lc)
						val = -1
					}
					biomearr[pt] = int16(val)
				}
			}
		case "Barren":
			// repeated bit of code here
			// pts: list of points
			// bathyarr, elevarr: need pointer
			for _, pt := range pts {
				var biome string
				// unique to us here
				elev := elevarr[pt]
				if elev > 92 {
					biome = "Desert Hills"
				} else {
					// or Desert M
					biome = "Desert"
				}
				//not unique to us here
				val, ok := world.Biome[biome]
				if !ok {
					log.Printf("%s is not a valid biome for %d!", biome, lc)
					val = -1
				}
				biomearr[pt] = int16(val)
			}
		case "Forest":
			// repeated bit of code here
			// pts: list of points
			// bathyarr, elevarr: need pointer
			for _, pt := range pts {
				var biome string
				// unique to us here
				elev := elevarr[pt]
				if elev > 92 {
					biome = "Forest Hills"
				} else {
					biome = "Forest"
				}
				//not unique to us here
				val, ok := world.Biome[biome]
				if !ok {
					log.Printf("%s is not a valid biome for %d!", biome, lc)
					val = -1
				}
				biomearr[pt] = int16(val)
			}
		default:
			// plains woo
			for _, pt := range pts {
				var biome string
				elev := elevarr[pt]
				if elev > 152 {
					// also "Extreme Hills+ M"
					biome = "Extreme Hills M"
				} else if elev > 122 {
					// Also "Extreme Hills+"
					biome = "Extreme Hills"
				} else if elev > 92 {
					biome = "Extreme Hills Edge"
				} else {
					// Rarely "Sunflower Plains"
					biome = "Plains"
				}
				val, ok := world.Biome[biome]
				if !ok {
					log.Printf("%s is not a valid biome for %d!", biome, lc)
					val = -1
				}
				biomearr[pt] = int16(val)
			}
		}
	}
	return biomearr, nil
}

func arrind(pt [2]int, inx int, iny int, gt [6]float64) (int, error) {
	if Debug {
		// log.Printf("x: %d, y: %d, inx: %d", pt[0], pt[1], inx)
		// log.Printf("0: %d, 1: %d, 2: %d, 3: %d, 4: %d, 5: %d",
		// int(gt[0]), int(gt[1]), int(gt[2]), int(gt[3]), int(gt[4]), int(gt[5]))
		// 		log.Printf("x-0: %d, y-3: %d", pt[0]-int(gt[0]), pt[1]-int(gt[3]))
		// log.Printf("x-0/1: %d, y-3/5: %d", (pt[0]-int(gt[0]))/int(gt[1]), (pt[1]-int(gt[3]))/int(gt[5]))
	}
	realx := (pt[0] - int(gt[0])) / int(gt[1])
	if realx < 0 {
		return 0, fmt.Errorf("realx %d < 0", realx)
	}
	if realx >= inx {
		return 0, fmt.Errorf("realx %d >= inx %d", realx, inx)
	}
	realy := (pt[1] - int(gt[3])) / int(gt[5])
	if realy <= 0 {
		return 0, fmt.Errorf("realx %d < 0", realx)
	}
	if realy >= iny {
		return 0, fmt.Errorf("realy %d >= iny %d", realy, iny)
	}
	return realx + realy*inx, nil
}

func traverse_feature(f gdal.Feature, field int, inx int, iny int, gt [6]float64) (lc int, pts []int, err error) {
	lc = f.FieldAsInteger(field)
	if Debug {
		log.Printf("lc: %d", lc)
	}
	g := f.Geometry()
	if !g.IsValid() {
		// JMT: for invalid geometry, exit with no points
		return lc, pts, err
	}
	e := g.Envelope()
	eminx := int(e.MinX())
	eminy := int(e.MinY())
	emaxx := int(e.MaxX())
	emaxy := int(e.MaxY())
	// JMT: is there a way to determine whether envelopes exceed bounds?
	if Debug {
		log.Printf("  Envelope: (%d, %d) -> (%d, %d)", eminx, eminy, emaxx, emaxy)
	}
	outSRS := gdal.CreateSpatialReference("")
	for y := eminy; y < emaxy; y -= int(gt[5]) {
		for x := eminx; x < emaxx; x += int(gt[1]) {
			inpt := [2]int{x, y}
			index, aerr := arrind(inpt, inx, iny, gt)
			// JMT: arrind returns nil if coordinates are invalid
			if aerr != nil {
				continue
			}
			if Debug {
				// log.Printf("x: %d, y: %d, index: %d", x, y, index)
			}
			wkt := fmt.Sprintf("POINT (%f %f)", float64(x)+0.5*gt[1], float64(y)+0.5*gt[5])
			pt, err := gdal.CreateFromWKT(wkt, outSRS)
			if err != nil && err.Error() != "No Error" {
				return lc, pts, err
			}
			if g.Contains(pt) {
				pts = append(pts, index)
				if Debug {
					// log.Print(pts)
				}
			}
		}
	}
	return lc, pts, nil
}

// NLCD 2011 also has canopy and impervious surfaces
// might be "fun" to add those eventually

func (r *Region) biome(lcarr []int16, elevarr []int16, bathyarr []int16) (biomearr []int16) {
	product := "NLCD 2011"
	bufferLen := len(lcarr)
	biomearr = make([]int16, bufferLen)
	for i := 0; i < bufferLen; i++ {
		lc := TerrainType(lcarr[i])
		elev := elevarr[i]
		bathy := bathyarr[i]
		if Debug {
			// log.Printf("lc: %d, elev: %d, bathy: %d", lc, elev, bathy)
		}
		var biome string
		lcct, err := lc.Data(product)
		if err != nil {
			panic(err)
		}
		lcc, lct := lcct[0], lcct[1]
		switch lcc {
		case "Water":
			switch lct {
			case "Open Water":

				// areas of open water, generally with less than 25%
				// cover of vegetation or soil

				// Biomes to consider:
				// - "River": if we ever get hydro data...

				if bathy > int16(r.maxdepth-1) {
					biome = "Deep Ocean"
				} else {
					biome = "Ocean"
				}
			case "Perennial Ice/Snow":

				// areas characterized by a perennial cover of ice
				// and/or snow, generally greater than 25% of total
				// cover

				// Biomes to consider:
				// - "Frozen River": any river in "ice plains"
				// - "Ice Plains": land covered with snow
				// - "Ice Mountains": like ice plains but higher

				biome = "Frozen Ocean"
			}
		case "Developed":
			switch lct {
			case "Open Space":

				// areas with a mixture of some constructed materials,
				// but mostly vegetation in the form of lawn grasses.
				// Impervious surfaces account for less than 20% of
				// total cover.  These areas most commonly include
				// large-lot single-family housing units, parks, golf
				// courses, and vegetation planted in developed
				// settings for recreation, erosion control, or
				// aesthetic purposes.

			case "Low Intensity":

				// areas with a mixture of constructed materials and
				// vegetation.  Impervious surfaces account for 20% to
				// 49% of total cover.  These areas most commonly
				// include single-family housing units.

			case "Medium Intensity":

				// areas with a mixture of constructed materials and
				// vegetation.  Impervious surfaces account for 50% to
				// 79% of total cover.  These areas most commonly
				// include single-family housing units.

			case "High Intensity":
				// highly developed areas where people reside or work
				// in high numbers.  Examples include apartment
				// complexes, row houses and
				// commercial/industrial. Impervious surfaces account
				// for 80% to 100% of the total cover.

			}

			if elev > 152 {
				// also "Extreme Hills+ M"
				biome = "Extreme Hills M"
			} else if elev > 122 {
				// Also "Extreme Hills+"
				biome = "Extreme Hills"
			} else if elev > 92 {
				biome = "Extreme Hills Edge"
			} else {
				// Rarely "Sunflower Plains"
				biome = "Plains"
			}
		case "Barren":

			// areas of bedrock, desert pavement, scarps, talus,
			// slides, volcanic material, glacial debris, sand dunes,
			// strip mines, gravel pits and other accumulations of
			// earthen material.  Generally, vegetation accounts for
			// less than 15% of total cover

			// Biomes to consider:
			// - if bordering ocean...
			//   "Cold Beach": ice plains for neighbors
			//   "Stone Beach": beach with stones.
			//   "Beach": beach _without_ stones.
			// - if high up?
			//   "Mesa" or "Mesa (Bryce)": maybe!
			//   "Mesa Plateau": flat top
			//   "Mesa Plateau F": forest on top

			if elev > 92 {
				biome = "Desert Hills"
			} else {
				// or Desert M
				biome = "Desert"
			}
		case "Forest":
			// Biomes to consider:
			// - "Forest": lots of trees, occasional hills, tall grass.
			//           oak and birch
			// - "Forest Hills" higher up
			// Other possibilities:
			// "Flower Forest": Hah
			// "Birch Forest"
			// "Birch Forest M": taller trees
			// "Birch Forest Hills"
			// "Birch Forest Hills M": large mountains and tall trees
			// "Roofed Forest": 100% coverage "dark oak"
			// "Roofed Forest M": mountains!

			// all forests have this in common:

			// areas dominated by trees generally greater than 5
			// meters tall, and greater than 20% of total vegetation
			// cover.

			// Trees and their terrains:
			// Oak: any forest
			// Spruce: evergreen forest
			// Birch: deciduous forest
			// Jungle: jungle (unknown)
			// Acacia: savannah (unknown)
			// Dark Oak: roofed forest (unknown)

			switch lct {
			case "Deciduous Forest":

				// More than 75% of the tree species shed foliage
				// simultaneously in response to seasonal change.

			case "Evergreen Forest":

				// More than 75% of the tree species maintain their
				// leaves all year.  Canopy is never without green
				// foliage.

			case "Mixed Forest":

				// Neither deciduous nor evergreen
				// species are greater than 75% of total tree cover.

			}
			if elev > 92 {
				biome = "Forest Hills"
			} else {
				biome = "Forest"
			}
		case "Scrubland":
			// NB: Only in Alaska?!
			// "Cold Taiga": snowy Taiga with fern and large fern
			// "Cold Taiga Hills": same with hills
			// "Cold Taiga M": same with mountains
			// "Taiga": no snow, spruce, ferns,et c.
			// "Taiga Hills": hills
			// "Taiga M": big hills
			// "Mega Taiga": rare version of taiga
			// "Mega Taiga Hills": with hills
			// "Mega Spruce Taiga": ...
			// maybe:
			// "Jungle"
			// "Jungle Edge": boundary with non-jungle territory
			// "Jungle Hills": ...
			// "Jungle M": mountainous
			// "JungleEdge M": boundary with non-jungle territory
			// possibly:
			// "Savannah": flat dry hot grass
			// "Savannah Plateau": flat high dry hot grass
			// "Savannah M": mountains to height limit?!
			switch lct {
			case "Dwarf Scrub":

				// Alaska only areas dominated by shrubs less than 20
				// centimeters tall with shrub canopy typically
				// greater than 20% of total vegetation.  This type is
				// often co-associated with grasses, sedges, herbs,
				// and non-vascular vegetation

			case "Shrub":

				// areas dominated by shrubs; less than 5 meters tall
				// with shrub canopy typically greater than 20% of
				// total vegetation.  This class includes true shrubs,
				// young trees in an early successional stage or trees
				// stunted from environmental conditions.

			}
		case "Herbaceous":
			switch lct {
			case "Grassland":

				// areas dominated by gramanoid or herbaceous
				// vegetation, generally greater than 80% of total
				// vegetation. These areas are not subject to
				// intensive management such as tilling, but can be
				// utilized for grazing.

				// 62 is sea level, +30 = hills, +60 = extreme
				if elev > 152 {
					// also "Extreme Hills+ M"
					biome = "Extreme Hills M"
				} else if elev > 122 {
					// Also "Extreme Hills+"
					biome = "Extreme Hills"
				} else if elev > 92 {
					biome = "Extreme Hills Edge"
				} else {
					// Rarely "Sunflower Plains"
					biome = "Plains"
				}
			case "Sedge":

				// Alaska only areas dominated by sedges and forbs,
				// generally greater than 80% of total
				// vegetation. This type can occur with significant
				// other grasses or other grass like plants, and
				// includes sedge tundra, and sedge tussock tundra.

			case "Lichens":

				// Alaska only areas dominated by fruticose or foliose
				// lichens generally greater than 80% of total
				// vegetation.

			case "Moss":

				// Alaska only areas dominated by mosses, generally
				// greater than 80% of total vegetation.

			}
		case "Cultivated":
			switch lct {
			case "Pasture":

				// areas of grasses, legumes, or grass-legume mixtures
				// planted for livestock grazing or the production of
				// seed or hay crops, typically on a perennial
				// cycle. Pasture/hay vegetation accounts for greater
				// than 20% of total vegetation.

			case "Crops":

				// areas used for the production of annual crops, such
				// as corn, soybeans, vegetables, tobacco, and cotton,
				// and also perennial woody crops such as orchards and
				// vineyards. Crop vegetation accounts for greater
				// than 20% of total vegetation. This class also
				// includes all land being actively tilled.

			}
		case "Wetlands":
			// "Swampland": clay sand dirt in pools, trees have vines,
			//              lots of mushrooms and sugarcanes
			// "Swampland M": slightly hillier, otherwise same
			switch lct {
			case "Woody Wetlands":

				// areas where forest or shrubland vegetation accounts
				// for greater than 20% of vegetative cover and the
				// soil or substrate is periodically saturated with or
				// covered with water.

			case "Emergent Herbaceous Wetlands":

				// Areas where perennial herbaceous vegetation
				// accounts for greater than 80% of vegetative cover
				// and the soil or substrate is periodically saturated
				// with or covered with water.

			}
			if elev > 92 {
				biome = "Swampland M"
			} else {
				biome = "Swampland"
			}
		default:
			log.Printf("Landcover type %d does not have a valid class!", lc)
			if elev > 152 {
				biome = "Extreme Hills M"
			} else if elev > 122 {
				biome = "Extreme Hills"
			} else if elev > 92 {
				biome = "Extreme Hills Edge"
			} else {
				biome = "Plains"
			}
		}
		val, ok := world.Biome[biome]
		if !ok {
			log.Printf("%s is not a valid biome for %d!", biome, lc)
			val = -1
		}
		biomearr[i] = int16(val)
	}
	return biomearr
}

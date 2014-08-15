package carto

import (
	"log"

	"github.com/mathuin/terroir/world"
)

// NLCD 2011 also has canopy and impervious surfaces
// might be "fun" to add those eventually

// tests should be written to check for every assignable biome.

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

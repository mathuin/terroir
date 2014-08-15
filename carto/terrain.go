package carto

import "fmt"

type TerrainType int

func (tt TerrainType) Data(which string) ([2]string, error) {
	var bad [2]string
	product, ok := terrainData[which]
	if !ok {
		return bad, fmt.Errorf("Terrain data for product %s not found!", which)
	}
	val, ok := product[tt]
	if !ok {
		return bad, fmt.Errorf("Terrain data for value %d not found in product %s!", tt, which)
	}
	return val, nil
}

var terrainData = map[string]map[TerrainType][2]string{
	"NLCD 2011": {
		11: [2]string{"Water", "Open Water"},
		12: [2]string{"Water", "Perennial Ice/Snow"},
		21: [2]string{"Developed", "Open Space"},
		22: [2]string{"Developed", "Low Intensity"},
		23: [2]string{"Developed", "Medium Intensity"},
		24: [2]string{"Developed", "High Intensity"},
		31: [2]string{"Barren", "Barren Land"},
		41: [2]string{"Forest", "Deciduous Forest"},
		42: [2]string{"Forest", "Evergreen Forest"},
		43: [2]string{"Forest", "Mixed Forest"},
		51: [2]string{"Scrubland", "Dwarf Scrub"},
		52: [2]string{"Scrubland", "Shrub"},
		71: [2]string{"Herbaceous", "Grassland"},
		72: [2]string{"Herbaceous", "Sedge"},
		73: [2]string{"Herbaceous", "Lichens"},
		74: [2]string{"Herbaceous", "Moss"},
		81: [2]string{"Cultivated", "Pasture"},
		82: [2]string{"Cultivated", "Crops"},
		90: [2]string{"Wetlands", "Woody Wetlands"},
		95: [2]string{"Wetlands", "Emergent Herbaceous Wetlands"},
	},
}

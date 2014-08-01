// chunk information

package world

type Section struct {
	Blocks   [4096]byte
	DataAdd  [2048]byte
	BlockSky [2048]byte
}

type Chunk struct {
	xPos      int32
	zPos      int32
	Biomes    [256]byte
	HeightMap [256]int32
	Sections  [16]Section
}

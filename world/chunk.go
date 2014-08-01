// chunk information

package world

import (
	"math"

	"github.com/mathuin/terroir/nbt"
)

type Section struct {
	blocks   [4096]byte
	addData  [4096]byte
	blockSky [4096]byte
}

func Half(arrin [4096]byte, top bool) (arrout [2048]byte) {
	for i, valin := range arrin {
		var val byte
		if top {
			val = valin >> 4
		} else {
			val = ((valin << 4) >> 4)
		}
		if math.Mod(float64(i), 2) == 1 {
			val = val << 4
		}
		halfi := i / 2
		arrout[halfi] = arrout[halfi] + val
	}
	return arrout
}

func (s Section) add() [2048]byte {
	return Half(s.addData, true)
}

func (s Section) data() [2048]byte {
	return Half(s.addData, false)
}

func (s Section) blockLight() [2048]byte {
	return Half(s.blockSky, true)
}

func (s Section) skyLight() [2048]byte {
	return Half(s.blockSky, false)
}

func (s Section) write(y int) []nbt.Tag {
	tarr := make([]nbt.Tag, 0)

	yTag := nbt.MakeTag(nbt.TAG_Byte, "Y")
	yTag.SetPayload(byte(y))

	blocksTag := nbt.MakeTag(nbt.TAG_Byte_Array, "Blocks")
	blocksTag.SetPayload(s.blocks)

	addTag := nbt.MakeTag(nbt.TAG_Byte_Array, "Add")
	addTag.SetPayload(s.add())
	dataTag := nbt.MakeTag(nbt.TAG_Byte_Array, "Data")
	dataTag.SetPayload(s.data())

	blockLightTag := nbt.MakeTag(nbt.TAG_Byte_Array, "BlockLight")
	blockLightTag.SetPayload(s.blockLight())
	skyLightTag := nbt.MakeTag(nbt.TAG_Byte_Array, "SkyLight")
	skyLightTag.SetPayload(s.skyLight())

	return tarr
}

// need to write something for reading sections!

// need to write func Double([2048]byte, [2048]byte) [4096]byte

// Chunk coords here are "world-relative"
type Chunk struct {
	xPos      int32
	zPos      int32
	biomes    [256]byte
	heightMap [256]int32
	sections  [16]Section
}

func (c Chunk) write() nbt.Tag {
	xPosTag := nbt.MakeTag(nbt.TAG_Int, "xPos")
	xPosTag.SetPayload(c.xPos)
	zPosTag := nbt.MakeTag(nbt.TAG_Int, "zPos")
	zPosTag.SetPayload(c.zPos)
	lastUpdateTag := nbt.MakeTag(nbt.TAG_Long, "LastUpdate")
	lastUpdateTag.SetPayload(int64(0))
	lightPopulatedTag := nbt.MakeTag(nbt.TAG_Byte, "LightPopulated")
	lightPopulatedTag.SetPayload(byte(0))
	terrainPopulatedTag := nbt.MakeTag(nbt.TAG_Byte, "TerrainPopulated")
	terrainPopulatedTag.SetPayload(byte(1))
	vTag := nbt.MakeTag(nbt.TAG_Byte, "V")
	vTag.SetPayload(byte(1))
	inhabitedTimeTag := nbt.MakeTag(nbt.TAG_Long, "InhabitedTime")
	inhabitedTimeTag.SetPayload(int64(0))
	biomesTag := nbt.MakeTag(nbt.TAG_Byte_Array, "Biomes")
	biomesTag.SetPayload(c.biomes)
	heightMapTag := nbt.MakeTag(nbt.TAG_Int_Array, "HeightMap")
	heightMapTag.SetPayload(c.heightMap)

	sectionsPayload := make([]interface{}, 16)
	for i, s := range c.sections {
		sectionsPayload[i] = s.write(i)
	}
	sectionsTag := nbt.MakeTag(nbt.TAG_List, "Sections")
	sectionsTag.SetPayload(sectionsPayload)

	// leaving the rest out for now...

	levelTag := nbt.MakeTag(nbt.TAG_Compound, "Level")
	levelTagPayload := []nbt.Tag{xPosTag, zPosTag, lastUpdateTag, lightPopulatedTag, terrainPopulatedTag, vTag, inhabitedTimeTag, biomesTag, heightMapTag, sectionsTag}
	levelTag.SetPayload(levelTagPayload)

	// not sure if this needs further wrapping or what...
	return levelTag
}

// need to write something for reading chunks

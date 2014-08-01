// chunk information

package world

import (
	"log"

	"github.com/mathuin/terroir/nbt"
)

type FullByte [4096]byte
type HalfByte [2048]byte

type Section struct {
	blocks   FullByte
	addData  FullByte
	blockSky FullByte
}

func split(in byte) (byte, byte) {
	return in >> 4, ((in << 4) >> 4)
}

func unsplit(top byte, bot byte) byte {
	return top*16 + bot
}

func toHalf(inlow byte, inhigh byte) (outtop byte, outbot byte) {
	inlowtop, inlowbot := split(inlow)
	inhightop, inhighbot := split(inhigh)

	outtop = unsplit(inhightop, inlowtop)
	outbot = unsplit(inhighbot, inlowbot)
	return
}

func toDouble(intop byte, inbot byte) (outlow byte, outhigh byte) {
	intoptop, intopbot := split(intop)
	inbottop, inbotbot := split(inbot)

	outlow = unsplit(intopbot, inbotbot)
	outhigh = unsplit(intoptop, inbottop)
	return
}

func Half(arrin FullByte, top bool) (arrout HalfByte) {
	for i := range arrout {
		outtop, outbot := toHalf(arrin[i/2], arrin[i/2+1])
		if top {
			arrout[i] = outtop
		} else {
			arrout[i] = outbot
		}
	}
	return
}

func Double(top HalfByte, bot HalfByte) (full FullByte) {
	for i := range top {
		di := i * 2
		full[di], full[di+1] = toDouble(top[i], bot[i])
	}
	return
}

func (s Section) add() HalfByte {
	return Half(s.addData, true)
}

func (s Section) data() HalfByte {
	return Half(s.addData, false)
}

func (s Section) blockLight() HalfByte {
	return Half(s.blockSky, true)
}

func (s Section) skyLight() HalfByte {
	return Half(s.blockSky, false)
}

func (s Section) write(y int) nbt.Tag {
	sElems := []nbt.CompoundElem{
		{"Y", nbt.TAG_Byte, byte(y)},
		{"Blocks", nbt.TAG_Byte_Array, s.blocks},
		{"Add", nbt.TAG_Byte_Array, s.add()},
		{"Data", nbt.TAG_Byte_Array, s.data()},
		{"BlockLight", nbt.TAG_Byte_Array, s.blockLight()},
		{"SkyLight", nbt.TAG_Byte_Array, s.skyLight()},
	}

	sTag := nbt.MakeCompound("", sElems)

	return sTag
}

func (s *Section) read(t nbt.Tag) {
	requiredTags := map[string]bool{
		"Y":          false,
		"Blocks":     false,
		"Add":        false,
		"Data":       false,
		"BlockLight": false,
		"SkyLight":   false,
	}
	var addTemp, dataTemp, blockTemp, skyTemp HalfByte

	for _, tval := range t.Payload.([]nbt.Tag) {
		if _, ok := requiredTags[tval.Name]; ok {
			requiredTags[tval.Name] = true
			switch tval.Name {
			case "Y":
				// Y tags are checked on the chunk level.
			case "Blocks":
				s.blocks = tval.Payload.(FullByte)
			case "Add":
				addTemp = tval.Payload.(HalfByte)
			case "Data":
				dataTemp = tval.Payload.(HalfByte)
			case "BlockLight":
				blockTemp = tval.Payload.(HalfByte)
			case "SkyLight":
				skyTemp = tval.Payload.(HalfByte)
			}
		} else {
			log.Fatalf("tag name %s not required for section", tval.Name)
		}
	}

	for rtkey, rtval := range requiredTags {
		if rtval == false {
			log.Fatalf("tag name %s required for section but not found", rtkey)
		}
	}

	s.addData = Double(addTemp, dataTemp)
	s.blockSky = Double(blockTemp, skyTemp)
}

// Chunk coords here are "world-relative"
type Chunk struct {
	xPos      int32
	zPos      int32
	biomes    [256]byte
	heightMap [256]int32
	sections  [16]Section
}

func (c Chunk) write() nbt.Tag {

	sectionsPayload := make([]interface{}, 16)
	for i, s := range c.sections {
		sectionsPayload[i] = s.write(i)
	}

	var levelElems = []nbt.CompoundElem{
		{"xPos", nbt.TAG_Int, c.xPos},
		{"zPos", nbt.TAG_Int, c.zPos},
		{"LastUpdate", nbt.TAG_Long, int64(0)},
		{"LightPopulated", nbt.TAG_Byte, byte(0)},
		{"TerrainPopulated", nbt.TAG_Byte, byte(1)},
		{"V", nbt.TAG_Byte, byte(1)},
		{"InhabitedTime", nbt.TAG_Long, int64(0)},
		{"Biomes", nbt.TAG_Byte_Array, c.biomes},
		{"HeightMap", nbt.TAG_Int_Array, c.heightMap},
		{"Sections", nbt.TAG_List, sectionsPayload},
	}

	levelTag := nbt.MakeCompound("Level", levelElems)

	// not sure if this needs further wrapping or what...
	return levelTag
}

func (c *Chunk) read(tarr []nbt.Tag) {
	requiredTags := map[string]bool{
		"xPos":      false,
		"zPos":      false,
		"Biomes":    false,
		"HeightMap": false,
		"Sections":  false,
	}
	for _, tval := range tarr {
		if _, ok := requiredTags[tval.Name]; ok {
			requiredTags[tval.Name] = true
			switch tval.Name {
			case "xPos":
				c.xPos = tval.Payload.(int32)
			case "zPos":
				c.zPos = tval.Payload.(int32)
			case "Biomes":
				c.biomes = tval.Payload.([256]byte)
			case "HeightMap":
				c.heightMap = tval.Payload.([256]int32)
			case "Sections":
				var yVals map[int]bool
				for _, stag := range tval.Payload.([16]nbt.Tag) {
					var yValFound bool
					var yVal int
					for _, subtags := range stag.Payload.([]nbt.Tag) {
						if subtags.Name == "Y" {
							yValFound = true
							yVal = int(subtags.Payload.(byte))
						}
					}
					if !yValFound {
						panic("no yVal found")
					}
					if _, ok := yVals[yVal]; ok {
						panic("yVal already found")
					}
					yVals[yVal] = true
					c.sections[yVal].read(stag)
				}
			}
		} else {
			log.Fatalf("tag name %s not required for chunk", tval.Name)
		}
	}

	for rtkey, rtval := range requiredTags {
		if rtval == false {
			log.Fatalf("tag name %s required for chunk but not found", rtkey)
		}
	}

}

// chunk information

package world

import (
	"log"

	"github.com/mathuin/terroir/nbt"
)

type FullByte [4096]byte
type HalfByte [2048]byte

type Section struct {
	blocks   []byte
	addData  []byte
	blockSky []byte
}

func NewSection() *Section {
	if Debug {
		log.Printf("NEW SECTION")
	}
	return &Section{blocks: make([]byte, 4096), addData: make([]byte, 4096), blockSky: make([]byte, 4096)}
}

func MakeSection() Section {
	if Debug {
		log.Printf("MAKE SECTION")
	}
	return Section{blocks: make([]byte, 4096), addData: make([]byte, 4096), blockSky: make([]byte, 4096)}
}

func (s Section) write(y int) nbt.Tag {
	add, data := Halve(s.addData)
	blockLight, skyLight := Halve(s.blockSky)

	sElems := []nbt.CompoundElem{
		{"Y", nbt.TAG_Byte, byte(y)},
		{"Blocks", nbt.TAG_Byte_Array, s.blocks},
		{"Add", nbt.TAG_Byte_Array, add},
		{"Data", nbt.TAG_Byte_Array, data},
		{"BlockLight", nbt.TAG_Byte_Array, blockLight},
		{"SkyLight", nbt.TAG_Byte_Array, skyLight},
	}

	sTag := nbt.MakeCompound("", sElems)

	return sTag
}

func (s *Section) read(tarr []nbt.Tag) {
	addTemp := make([]byte, 2048)
	dataTemp := make([]byte, 2048)
	blockTemp := make([]byte, 2048)
	skyTemp := make([]byte, 2048)

	for _, tval := range tarr {
		switch tval.Name {
		case "Y":
		// Y tags are checked on the chunk level.
		case "Blocks":
			s.blocks = tval.Payload.([]byte)
		case "Add":
			addTemp = tval.Payload.([]byte)
		case "Data":
			dataTemp = tval.Payload.([]byte)
		case "BlockLight":
			blockTemp = tval.Payload.([]byte)
		case "SkyLight":
			skyTemp = tval.Payload.([]byte)
		default:
			log.Fatalf("tag name %s not required for section", tval.Name)
		}
	}

	s.addData = Double(addTemp, dataTemp)
	s.blockSky = Double(blockTemp, skyTemp)
}

// JMT: learn more about coords to write better tests
// coords are 0-31 for each, I think
// since region files hold 32x32 chunks
type Chunk struct {
	xPos      int32
	zPos      int32
	biomes    []byte
	heightMap []int32
	sections  []Section
}

func NewChunk(xPos int32, zPos int32) *Chunk {
	if Debug {
		log.Printf("NEW CHUNK: xPos %d, zPos %d", xPos, zPos)
	}
	sections := make([]Section, 0)
	for i := 0; i < 16; i++ {
		sections = append(sections, MakeSection())
	}
	return &Chunk{xPos: xPos, zPos: zPos, biomes: make([]byte, 256), heightMap: make([]int32, 256), sections: sections}
}

func MakeChunk(xPos int32, zPos int32) Chunk {
	if Debug {
		log.Printf("MAKE CHUNK: xPos %d, zPos %d", xPos, zPos)
	}
	sections := make([]Section, 0)
	for i := 0; i < 16; i++ {
		sections = append(sections, MakeSection())
	}
	return Chunk{xPos: xPos, zPos: zPos, biomes: make([]byte, 256), heightMap: make([]int32, 256), sections: sections}
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
		{"Biomes", nbt.TAG_Byte_Array, []byte(c.biomes)},
		{"HeightMap", nbt.TAG_Int_Array, []int32(c.heightMap)},
		{"Sections", nbt.TAG_List, sectionsPayload},
	}

	levelTag := nbt.MakeCompound("Level", levelElems)

	// not sure if this needs further wrapping or what...
	return levelTag
}

func (c *Chunk) Read(t nbt.Tag) {
	requiredTags := map[string]bool{
		"xPos":      false,
		"zPos":      false,
		"Biomes":    false,
		"HeightMap": false,
		"Sections":  false,
	}

	if t.Type != nbt.TAG_Compound {
		log.Panic("chunk read tag type not TAG_Compound!")
	}

	levelTag := t.Payload.([]nbt.Tag)[0]

	if Debug {
		log.Printf("Chunk Read 2: tag type %s name %s", nbt.Names[levelTag.Type], levelTag.Name)
	}

	for _, tval := range levelTag.Payload.([]nbt.Tag) {
		if Debug {
			log.Printf("tag type %s name %s found", nbt.Names[tval.Type], tval.Name)
		}
		if _, ok := requiredTags[tval.Name]; ok {
			if Debug {
				log.Printf(" -- tag is required")
			}
			requiredTags[tval.Name] = true
			switch tval.Name {
			case "xPos":
				c.xPos = tval.Payload.(int32)
			case "zPos":
				c.zPos = tval.Payload.(int32)
			case "Biomes":
				c.biomes = tval.Payload.([]byte)
			case "HeightMap":
				c.heightMap = tval.Payload.([]int32)
			case "Sections":
				yVals := make(map[int]bool, 0)
				for _, s := range tval.Payload.([][]nbt.Tag) {
					var yValFound bool
					var yVal int
					for _, subtag := range s {
						if subtag.Name == "Y" {
							yValFound = true
							yVal = int(subtag.Payload.(byte))
						}
					}
					if !yValFound {
						panic("no yVal found")
					}
					if _, ok := yVals[yVal]; ok {
						panic("yVal already found")
					}
					yVals[yVal] = true
					c.sections[yVal].read(s)
				}
			}
			//		} else {
			//			log.Fatalf("tag name %s not required for chunk", tval.Name)
		}
	}

	for rtkey, rtval := range requiredTags {
		if rtval == false {
			log.Fatalf("tag name %s required for chunk but not found", rtkey)
		}
	}

}

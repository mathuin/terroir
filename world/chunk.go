package world

import (
	"fmt"
	"log"

	"github.com/mathuin/terroir/nbt"
)

type Chunk struct {
	xPos      int32
	zPos      int32
	biomes    []byte
	heightMap []int32
	Sections  map[int]Section
	// living things
	entities     []Entity
	tileEntities []TileEntity
	tileTicks    []TileTick
}

func MakeChunk(xPos int32, zPos int32) Chunk {
	if Debug {
		log.Printf("MAKE CHUNK: xPos %d, zPos %d", xPos, zPos)
	}
	biomes := make([]byte, 256)
	heightMap := make([]int32, 256)
	Sections := map[int]Section{}
	entities := []Entity{}
	tileEntities := []TileEntity{}
	tileTicks := []TileTick{}
	return Chunk{xPos: xPos, zPos: zPos, biomes: biomes, heightMap: heightMap, Sections: Sections, entities: entities, tileEntities: tileEntities, tileTicks: tileTicks}
}

func (c Chunk) Name() string {
	return fmt.Sprintf("%d, %d", c.xPos, c.zPos)
}

func (c Chunk) write() nbt.Tag {

	sectionsPayload := [][]nbt.Tag{}
	for i, s := range c.Sections {
		st := s.write(i)
		sectionsPayload = append(sectionsPayload, st)
	}

	entitiesPayload := [][]nbt.Tag{}
	for _, e := range c.entities {
		entitiesPayload = append(entitiesPayload, e.write())
	}

	tileEntitiesPayload := [][]nbt.Tag{}
	for _, te := range c.tileEntities {
		tileEntitiesPayload = append(tileEntitiesPayload, te.write())
	}

	tileTicksPayload := [][]nbt.Tag{}
	for _, tt := range c.tileTicks {
		ttt := tt.write()
		tileTicksPayload = append(tileTicksPayload, ttt.Payload.([]nbt.Tag))
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
		{"Entities", nbt.TAG_List, entitiesPayload},
		{"TileEntities", nbt.TAG_List, tileEntitiesPayload},
	}

	if len(c.tileTicks) > 0 {
		tickElem := nbt.CompoundElem{"TileTicks", nbt.TAG_List, tileTicksPayload}
		levelElems = append(levelElems, tickElem)
	}

	levelTag := nbt.MakeCompound("Level", levelElems)

	topTag := nbt.MakeTag(nbt.TAG_Compound, "")
	if err := topTag.SetPayload([]nbt.Tag{levelTag}); err != nil {
		panic(err)
	}

	return topTag
}

func (c *Chunk) Read(t nbt.Tag) {
	requiredTags := map[string]bool{
		"xPos":         false,
		"zPos":         false,
		"Biomes":       false,
		"HeightMap":    false,
		"Sections":     false,
		"Entities":     false,
		"TileEntities": false,
	}

	if t.Type != nbt.TAG_Compound {
		log.Panic("top tag type not TAG_Compound!")
	}

	if t.Name != "" {
		log.Panic("top tag not unnamed")
	}

	levelTag := t.Payload.([]nbt.Tag)[0]

	if levelTag.Type != nbt.TAG_Compound {
		log.Panic("level tag type not TAG_Compound!")
	}

	if levelTag.Name != "Level" {
		log.Panic("level tag not named Level")
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
				xPos := tval.Payload.(int32)
				if xPos != c.xPos {
					log.Fatalf("xPos %d does not match c.xPos %d", xPos, c.xPos)
				}
			case "zPos":
				zPos := tval.Payload.(int32)
				if zPos != c.zPos {
					log.Fatalf("zPos %d does not match c.zPos %d", zPos, c.zPos)
				}
			case "Biomes":
				c.biomes = tval.Payload.([]byte)
			case "HeightMap":
				c.heightMap = tval.Payload.([]int32)
			case "Sections":
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
					if _, ok := c.Sections[yVal]; ok {
						panic("yVal already found")
					}
					c.Sections[yVal] = ReadSection(s)
				}
			case "Entities":
				if tval.Payload != nil {
					es := make([]Entity, 0)
					for _, e := range tval.Payload.([][]nbt.Tag) {
						es = append(es, ReadEntity(e))
					}
					c.entities = es
				}
			case "TileEntities":
				if tval.Payload != nil {
					tes := make([]TileEntity, 0)
					for _, te := range tval.Payload.([][]nbt.Tag) {
						tes = append(tes, ReadTileEntity(te))
					}
					c.tileEntities = tes
				}
			}
		} else {
			// optional tags
			switch tval.Name {
			case "TileTicks":
				tts := make([]TileTick, 0)
				for _, tt := range tval.Payload.([][]nbt.Tag) {
					tts = append(tts, ReadTileTick(tt))
				}
				c.tileTicks = tts
			}
		}
	}

	for rtkey, rtval := range requiredTags {
		if rtval == false {
			log.Fatalf("tag name %s required for chunk but not found", rtkey)
		}
	}

}

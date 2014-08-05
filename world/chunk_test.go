package world

import (
	"log"
	"testing"

	"github.com/mathuin/terroir/nbt"
)

func Test_sectionWrite(t *testing.T) {
	s := NewSection()

	// populate it
	for i := range s.blocks {
		s.blocks[i] = 0x31
		s.addData[i] = 0x12
		s.blockSky[i] = 0x34
	}

	sTagPayload := s.write(1)

	requiredTags := map[string]bool{
		"Y":          false,
		"Blocks":     false,
		"Add":        false,
		"Data":       false,
		"BlockLight": false,
		"SkyLight":   false,
	}

	for _, tag := range sTagPayload {
		if _, ok := requiredTags[tag.Name]; ok {
			requiredTags[tag.Name] = true
			switch tag.Name {
			case "Y":
				if tag.Payload.(byte) != 0x01 {
					t.Errorf("Y value does not match")
				}
			case "Blocks":
				if tag.Payload.([]byte)[0] != 0x31 {
					t.Errorf("Block value does not match")
				}
			case "Add":
				if tag.Payload.([]byte)[0] != 0x11 {
					t.Errorf("Add value does not match")
				}
			case "Data":
				if tag.Payload.([]byte)[0] != 0x22 {
					t.Errorf("Data value does not match")
				}
			case "BlockLight":
				if tag.Payload.([]byte)[0] != 0x33 {
					t.Errorf("BlockLight value does not match")
				}
			case "SkyLight":
				if tag.Payload.([]byte)[0] != 0x44 {
					t.Errorf("SkyLight value does not match")
				}
			}
		} else {
			t.Errorf("tag name %s not required for section", tag.Name)
		}
	}

	for rtkey, rtval := range requiredTags {
		if rtval == false {
			t.Errorf("tag name %s required for section but not found", rtkey)
		}
	}
}

func Test_chunkWrite(t *testing.T) {
	var cxPos, czPos int32
	var cbiomes []byte
	var cheightMap []int32

	cxPos = 1
	czPos = 2
	cbiomes = make([]byte, 256)
	cheightMap = make([]int32, 256)

	for i := range cbiomes {
		cbiomes[i] = byte(255 - i)
		cheightMap[i] = int32(255 - i)
	}

	c := NewChunk(cxPos, czPos)
	for i := range cbiomes {
		c.biomes[i] = cbiomes[i]
		c.heightMap[i] = cheightMap[i]
	}

	for i := 0; i < 16; i++ {
		s := MakeSection()
		for j := range s.blocks {
			s.blocks[j] = 0x30 & byte(i)
			s.addData[j] = 0x12
			s.blockSky[j] = 0x34
		}
		c.sections[i] = s
	}

	cTag := c.write()

	if cTag.Type != nbt.TAG_Compound {
		t.Errorf("chunk tab not of type TAG_Compound")
	}

	if cTag.Name != "" {
		t.Errorf("chunk tag not unnamed")
	}

	lTag := cTag.Payload.([]nbt.Tag)[0]

	if lTag.Type != nbt.TAG_Compound {
		t.Errorf("level tag not of type TAG_Compound")
	}

	if lTag.Name != "Level" {
		t.Errorf("level tag not named Level")
	}

	requiredTags := map[string]bool{
		"xPos":             false,
		"zPos":             false,
		"LastUpdate":       false,
		"LightPopulated":   false,
		"TerrainPopulated": false,
		"V":                false,
		"InhabitedTime":    false,
		"Biomes":           false,
		"HeightMap":        false,
		"Sections":         false,
		"Entities":         false,
		"TileEntities":     false,
	}

	for _, tag := range lTag.Payload.([]nbt.Tag) {
		if _, ok := requiredTags[tag.Name]; ok {
			requiredTags[tag.Name] = true
			switch tag.Name {
			case "xPos":
				if tag.Payload.(int32) != cxPos {
					t.Errorf("xPos value does not match")
				}
			case "zPos":
				if tag.Payload.(int32) != czPos {
					t.Errorf("zPos value does not match")
				}
			case "Biomes":
				if tag.Payload.([]byte)[0] != cbiomes[0] {
					t.Errorf("biomes value does not match")
				}
			case "HeightMap":
				if tag.Payload.([]int32)[0] != cheightMap[0] {
					t.Errorf("heightMap value does not match")
				}
				// sections are checked elsewhere
			}
		} else {
			log.Fatalf("tag name %s not required for chunk", tag.Name)
		}
	}
}

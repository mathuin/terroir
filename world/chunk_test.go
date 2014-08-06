package world

import (
	"testing"

	"github.com/mathuin/terroir/nbt"
)

func Test_chunkWriteRead(t *testing.T) {
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

	c := MakeChunk(cxPos, czPos)
	for i := range cbiomes {
		c.biomes[i] = cbiomes[i]
		c.heightMap[i] = cheightMap[i]
	}

	for i := 0; i < 16; i++ {
		s := MakeSection()
		for j := range s.Blocks {
			s.Blocks[j] = 0x30 & byte(i)
		}
		for j := range s.Add {
			s.Add[j] = 0x12
			s.Data[j] = 0x34
			s.BlockLight[j] = 0x56
			s.SkyLight[j] = 0x78
		}
		c.Sections[i] = s
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
			t.Errorf("tag name %s not required for chunk", tag.Name)
		}
	}

	// Now read it back!
	newc := MakeChunk(cxPos, czPos)
	if err := newc.Read(cTag); err != nil {
		t.Fail()
	}

	// if this succeeds without throwing an error, it was able to read
	// the level tags xPos and zPos, so that's good enough for me.
}

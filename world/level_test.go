package world

import (
	"testing"

	"github.com/mathuin/terroir/nbt"
)

func Test_newLevel(t *testing.T) {
	// initial parameters
	pname := "ExampleLevel"
	px := int32(-208)
	py := int32(64)
	pz := int32(220)
	prseed := int64(2603059821051629081)

	w := MakeWorld(pname)
	w.SetRandomSeed(prseed)
	w.SetSpawn(px, py, pz)

	requiredTags := map[string]bool{
		"LevelName":  false,
		"SpawnX":     false,
		"SpawnY":     false,
		"SpawnZ":     false,
		"RandomSeed": false,
		"version":    false,
		"GameRules":  false,
	}

	gameRulesRequiredTags := map[string]bool{
		"commandBlockOutput": false,
		"doTileDrops":        false,
		"keepInventory":      false,
	}

	testNewLevel, err := w.level()
	if err != nil {
		t.Fail()
	}
	topPayload := testNewLevel.Payload.([]nbt.Tag)
	if len(topPayload) != 1 {
		t.Errorf("topPayload does not contain only one tag\n")
	}
	dataTag := topPayload[0]
	if dataTag.Name != "Data" {
		t.Errorf("Data is not name of top inner tag!\n")
	}
	for _, tag := range dataTag.Payload.([]nbt.Tag) {
		switch tag.Name {
		case "LevelName":
			if tag.Payload.(string) != pname {
				t.Errorf("Name does not match\n")
			}
			requiredTags[tag.Name] = true
		case "SpawnX":
			if tag.Payload.(int32) != px {
				t.Errorf("SpawnX does not match\n")
			}
			requiredTags[tag.Name] = true
		case "SpawnY":
			if tag.Payload.(int32) != py {
				t.Errorf("SpawnY does not match\n")
			}
			requiredTags[tag.Name] = true
		case "SpawnZ":
			if tag.Payload.(int32) != pz {
				t.Errorf("SpawnZ does not match\n")
			}
			requiredTags[tag.Name] = true
		case "RandomSeed":
			if tag.Payload.(int64) != prseed {
				t.Errorf("SpawnX does not match\n")
			}
			requiredTags[tag.Name] = true
		case "version":
			if tag.Payload.(int32) != int32(19133) {
				t.Errorf("version does not match\n")
			}
			requiredTags[tag.Name] = true
		case "commandBlockOutput":
			t.Errorf("commandBlockOutput does not belong in top level\n")
		case "GameRules":
			for _, ruletag := range tag.Payload.([]nbt.Tag) {
				switch ruletag.Name {
				case "commandBlockOutput":
					if ruletag.Payload.(string) != "true" {
						t.Errorf("commandBlockOutput does not match\n")
					}
					gameRulesRequiredTags[ruletag.Name] = true
				case "doTileDrops":
					if ruletag.Payload.(string) != "true" {
						t.Errorf("doTileDrops does not match\n")
					}
					gameRulesRequiredTags[ruletag.Name] = true
				case "keepInventory":
					if ruletag.Payload.(string) != "false" {
						t.Errorf("keepInventory does not match\n")
					}
					gameRulesRequiredTags[ruletag.Name] = true
				}
			}
			requiredTags[tag.Name] = true
		}
	}
	for rtkey, rtval := range requiredTags {
		if rtval == false {
			t.Errorf("tag name %s required for section but not found", rtkey)
		}
	}
	for rtkey, rtval := range gameRulesRequiredTags {
		if rtval == false {
			t.Errorf("tag name %s required for game rules but not found", rtkey)
		}
	}
}

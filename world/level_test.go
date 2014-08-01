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

	// JMT: make this a map like the chunk read checker!
	levelNameCheck := false
	spawnXCheck := false
	gameRulesCheck := false

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
			levelNameCheck = true
		case "SpawnX":
			if tag.Payload.(int32) != px {
				t.Errorf("SpawnX does not match\n")
			}
			spawnXCheck = true
		case "SpawnY":
			if tag.Payload.(int32) != py {
				t.Errorf("SpawnY does not match\n")
			}
		case "SpawnZ":
			if tag.Payload.(int32) != pz {
				t.Errorf("SpawnZ does not match\n")
			}
		case "RandomSeed":
			if tag.Payload.(int64) != prseed {
				t.Errorf("SpawnX does not match\n")
			}
		case "version":
			if tag.Payload.(int32) != int32(19133) {
				t.Errorf("version does not match\n")
			}
		case "commandBlockOutput":
			t.Errorf("commandBlockOutput does not belong in top level\n")
		case "GameRules":
			for _, ruletag := range tag.Payload.([]nbt.Tag) {
				switch ruletag.Name {
				case "commandBlockOutput":
					if ruletag.Payload.(string) != "true" {
						t.Errorf("commandBlockOutput does not match\n")
					}
				case "doTileDrops":
					if ruletag.Payload.(string) != "true" {
						t.Errorf("doTileDrops does not match\n")
					}
				case "keepInventory":
					if ruletag.Payload.(string) != "false" {
						t.Errorf("keepInventory does not match\n")
					}
				}
			}
			gameRulesCheck = true
		}
	}
	if !levelNameCheck {
		t.Errorf("missing LevelName tag")
	}
	if !spawnXCheck {
		t.Errorf("missing SpawnX tag")
	}
	if !gameRulesCheck {
		t.Errorf("missing GameRules tag")
	}
}

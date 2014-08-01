// Stuff related to level.dat

package world

import (
	"fmt"

	"github.com/mathuin/terroir/nbt"
)

func (w World) level() (t nbt.Tag, err error) {
	if w.spawnSet == false {
		return t, fmt.Errorf("Spawn must be set before creating level")
	}
	var gamerules = []struct {
		key   string
		value string
	}{
		{"commandBlockOutput", "true"},
		{"doDaylightCycle", "true"},
		{"doFireTick", "true"},
		{"doMobLoot", "true"},
		{"doMobSpawning", "true"},
		{"doTileDrops", "true"},
		{"keepInventory", "false"},
		{"mobGriefing", "true"},
		{"naturalRegeneration", "true"},
	}

	gameRulesPayload := make([]nbt.Tag, 0)

	for _, rule := range gamerules {
		newTag := nbt.MakeTag(nbt.TAG_String, rule.key)
		newTag.SetPayload(rule.value)
		gameRulesPayload = append(gameRulesPayload, newTag)
	}

	dataPayload := make([]nbt.Tag, 0)

	var tags = []struct {
		key   string
		tag   byte
		value interface{}
	}{
		{"allowCommands", nbt.TAG_Byte, byte(0)},
		{"generatorName", nbt.TAG_String, "default"},
		{"generatorOptions", nbt.TAG_String, ""},
		{"generatorVersion", nbt.TAG_Int, int32(1)},
		{"hardcore", nbt.TAG_Byte, byte(0)},
		{"initialized", nbt.TAG_Byte, byte(1)},
		{"raining", nbt.TAG_Byte, byte(0)},
		{"rainTime", nbt.TAG_Int, int32(0)},
		{"thundering", nbt.TAG_Byte, byte(0)},
		{"thunderTime", nbt.TAG_Int, int32(0)},
		{"version", nbt.TAG_Int, int32(19133)},
		{"DayTime", nbt.TAG_Long, int64(0)},
		{"GameRules", nbt.TAG_Compound, gameRulesPayload},
		{"GameType", nbt.TAG_Int, int32(0)},
		{"LastPlayed", nbt.TAG_Long, int64(0)},
		{"LevelName", nbt.TAG_String, w.Name},
		{"MapFeatures", nbt.TAG_Byte, byte(1)},
		{"RandomSeed", nbt.TAG_Long, w.randomSeed},
		{"SizeOnDisk", nbt.TAG_Long, int64(0)},
		{"SpawnX", nbt.TAG_Int, w.spawnX},
		{"SpawnY", nbt.TAG_Int, w.spawnY},
		{"SpawnZ", nbt.TAG_Int, w.spawnZ},
		{"Time", nbt.TAG_Long, int64(0)},
	}
	for _, tag := range tags {
		newTag := nbt.MakeTag(tag.tag, tag.key)
		newTag.SetPayload(tag.value)
		dataPayload = append(dataPayload, newTag)
	}

	dataTag := nbt.MakeTag(nbt.TAG_Compound, "Data")
	dataTag.SetPayload(dataPayload)

	t = nbt.MakeTag(nbt.TAG_Compound, "")
	t.SetPayload([]nbt.Tag{dataTag})
	return
}

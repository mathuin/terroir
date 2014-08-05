// Stuff related to level.dat

package world

import (
	"fmt"
	"log"
	"path"

	"github.com/mathuin/terroir/nbt"
)

func (w World) level() (t nbt.Tag, err error) {
	if w.spawnSet == false {
		return t, fmt.Errorf("Spawn must be set before creating level")
	}
	gameRulesElems := []nbt.CompoundElem{
		{"commandBlockOutput", nbt.TAG_String, "true"},
		{"doDaylightCycle", nbt.TAG_String, "true"},
		{"doFireTick", nbt.TAG_String, "true"},
		{"doMobLoot", nbt.TAG_String, "true"},
		{"doMobSpawning", nbt.TAG_String, "true"},
		{"doTileDrops", nbt.TAG_String, "true"},
		{"keepInventory", nbt.TAG_String, "false"},
		{"mobGriefing", nbt.TAG_String, "true"},
		{"naturalRegeneration", nbt.TAG_String, "true"},
	}

	gameRulesPayload := nbt.MakeCompoundPayload(gameRulesElems)

	dataElems := []nbt.CompoundElem{
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
		{"RandomSeed", nbt.TAG_Long, w.RandomSeed},
		{"SizeOnDisk", nbt.TAG_Long, int64(0)},
		{"SpawnX", nbt.TAG_Int, w.Spawn.X},
		{"SpawnY", nbt.TAG_Int, w.Spawn.Y},
		{"SpawnZ", nbt.TAG_Int, w.Spawn.Z},
		{"Time", nbt.TAG_Long, int64(0)},
	}
	dataTag := nbt.MakeCompound("Data", dataElems)

	t = nbt.MakeTag(nbt.TAG_Compound, "")
	t.SetPayload([]nbt.Tag{dataTag})
	return
}

func (w World) writelevel(dir string) error {
	if path.Base(dir) != w.Name {
		return fmt.Errorf("directory does not contain world name")
	}
	levelTag, err := w.level()
	if err != nil {
		return err
	}
	levelFile := path.Join(dir, "level.dat")
	if Debug {
		log.Printf("Writing level file %s", levelFile)
	}
	return nbt.WriteCompressedFile(levelFile, levelTag)
}

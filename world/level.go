// Stuff related to level.dat

package world

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mathuin/terroir/nbt"
)

func (w World) level() (t nbt.Tag, err error) {
	if w.spawnSet == false {
		return t, fmt.Errorf("Spawn must be set before creating level")
	}
	gameRulesElems := []nbt.CompoundElem{
		{Key: "commandBlockOutput", Tag: nbt.TAG_String, Value: "true"},
		{Key: "doDaylightCycle", Tag: nbt.TAG_String, Value: "true"},
		{Key: "doFireTick", Tag: nbt.TAG_String, Value: "true"},
		{Key: "doMobLoot", Tag: nbt.TAG_String, Value: "true"},
		{Key: "doMobSpawning", Tag: nbt.TAG_String, Value: "true"},
		{Key: "doTileDrops", Tag: nbt.TAG_String, Value: "true"},
		{Key: "keepInventory", Tag: nbt.TAG_String, Value: "false"},
		{Key: "mobGriefing", Tag: nbt.TAG_String, Value: "true"},
		{Key: "naturalRegeneration", Tag: nbt.TAG_String, Value: "true"},
	}

	gameRulesPayload := nbt.MakeCompoundPayload(gameRulesElems)

	dataElems := []nbt.CompoundElem{
		{Key: "allowCommands", Tag: nbt.TAG_Byte, Value: byte(0)},
		{Key: "generatorName", Tag: nbt.TAG_String, Value: "default"},
		{Key: "generatorOptions", Tag: nbt.TAG_String, Value: ""},
		{Key: "generatorVersion", Tag: nbt.TAG_Int, Value: int32(1)},
		{Key: "hardcore", Tag: nbt.TAG_Byte, Value: byte(0)},
		{Key: "initialized", Tag: nbt.TAG_Byte, Value: byte(1)},
		{Key: "raining", Tag: nbt.TAG_Byte, Value: byte(0)},
		{Key: "rainTime", Tag: nbt.TAG_Int, Value: int32(0)},
		{Key: "thundering", Tag: nbt.TAG_Byte, Value: byte(0)},
		{Key: "thunderTime", Tag: nbt.TAG_Int, Value: int32(0)},
		{Key: "version", Tag: nbt.TAG_Int, Value: int32(19133)},
		{Key: "DayTime", Tag: nbt.TAG_Long, Value: int64(0)},
		{Key: "GameRules", Tag: nbt.TAG_Compound, Value: gameRulesPayload},
		{Key: "GameType", Tag: nbt.TAG_Int, Value: int32(0)},
		{Key: "LastPlayed", Tag: nbt.TAG_Long, Value: int64(0)},
		{Key: "LevelName", Tag: nbt.TAG_String, Value: w.Name},
		{Key: "MapFeatures", Tag: nbt.TAG_Byte, Value: byte(1)},
		{Key: "RandomSeed", Tag: nbt.TAG_Long, Value: w.RandomSeed},
		{Key: "SizeOnDisk", Tag: nbt.TAG_Long, Value: int64(0)},
		{Key: "SpawnX", Tag: nbt.TAG_Int, Value: w.Spawn.X},
		{Key: "SpawnY", Tag: nbt.TAG_Int, Value: w.Spawn.Y},
		{Key: "SpawnZ", Tag: nbt.TAG_Int, Value: w.Spawn.Z},
		{Key: "Time", Tag: nbt.TAG_Long, Value: int64(0)},
	}
	dataTag := nbt.MakeCompound("Data", dataElems)

	t = nbt.MakeTag(nbt.TAG_Compound, "")
	t.SetPayload([]nbt.Tag{dataTag})
	return
}

func (w World) writeLevel() error {
	worldDir := path.Join(w.SaveDir, w.Name)
	if err := os.MkdirAll(worldDir, 0775); err != nil {
		return err
	}
	levelTag, err := w.level()
	if err != nil {
		return err
	}
	levelFile := path.Join(worldDir, "level.dat")
	if Debug {
		log.Printf("Writing level file %s", levelFile)
	}
	return nbt.WriteCompressedFile(levelFile, levelTag)
}

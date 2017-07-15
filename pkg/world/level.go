// Stuff related to level.dat

package world

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mathuin/terroir/pkg/nbt"
)

func (w World) level() (t nbt.Tag, err error) {

	dataversion := 1139

	if w.spawnSet == false {
		return t, fmt.Errorf("Spawn must be set before creating level")
	}

	gameRulesElems := []nbt.CompoundElem{
		{Key: "announceAdvancements", Tag: nbt.TAG_String, Value: "true"},
		{Key: "commandBlockOutput", Tag: nbt.TAG_String, Value: "true"},
		{Key: "disableElytraMovementCheck", Tag: nbt.TAG_String, Value: "false"},
		{Key: "doDaylightCycle", Tag: nbt.TAG_String, Value: "true"},
		{Key: "doEntityDrops", Tag: nbt.TAG_String, Value: "true"},
		{Key: "doFireTick", Tag: nbt.TAG_String, Value: "true"},
		{Key: "doLimitedCrafting", Tag: nbt.TAG_String, Value: "false"},
		{Key: "doMobLoot", Tag: nbt.TAG_String, Value: "true"},
		{Key: "doMobSpawning", Tag: nbt.TAG_String, Value: "true"},
		{Key: "doTileDrops", Tag: nbt.TAG_String, Value: "true"},
		{Key: "doWeatherCycle", Tag: nbt.TAG_String, Value: "true"},
		{Key: "gameLoopFunction", Tag: nbt.TAG_String, Value: "-"},
		{Key: "keepInventory", Tag: nbt.TAG_String, Value: "false"},
		{Key: "logAdminCommands", Tag: nbt.TAG_String, Value: "true"},
		{Key: "maxCommandChainLength", Tag: nbt.TAG_String, Value: "65536"},
		{Key: "maxEntityCramming", Tag: nbt.TAG_String, Value: "24"},
		{Key: "mobGriefing", Tag: nbt.TAG_String, Value: "true"},
		{Key: "naturalRegeneration", Tag: nbt.TAG_String, Value: "true"},
		{Key: "randomTickSpeed", Tag: nbt.TAG_String, Value: "3"},
		{Key: "reducedDebugInfo", Tag: nbt.TAG_String, Value: "false"},
		{Key: "sendCommandFeedback", Tag: nbt.TAG_String, Value: "true"},
		{Key: "showDeathMessages", Tag: nbt.TAG_String, Value: "true"},
		{Key: "spawnRadius", Tag: nbt.TAG_String, Value: "10"},
		{Key: "spectatorsGenerateChunks", Tag: nbt.TAG_String, Value: "true"},
	}

	gameRulesPayload := nbt.MakeCompoundPayload(gameRulesElems)

	versionElems := []nbt.CompoundElem{
		{Key: "Id", Tag: nbt.TAG_Int, Value: int32(dataversion)},
		{Key: "Name", Tag: nbt.TAG_String, Value: "1.12"},
		{Key: "Snapshot", Tag: nbt.TAG_Byte, Value: byte(0)},
	}

	versionPayload := nbt.MakeCompoundPayload(versionElems)

	dataElems := []nbt.CompoundElem{
		{Key: "allowCommands", Tag: nbt.TAG_Byte, Value: byte(0)},
		{Key: "clearWeatherTime", Tag: nbt.TAG_Int, Value: int32(0)},
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
		{Key: "BorderCenterX", Tag: nbt.TAG_Double, Value: float64(0.0)},
		{Key: "BorderCenterZ", Tag: nbt.TAG_Double, Value: float64(0.0)},
		{Key: "BorderDamagePerBlock", Tag: nbt.TAG_Double, Value: float64(0.2)},
		{Key: "BorderSafeZone", Tag: nbt.TAG_Double, Value: float64(5.0)},
		{Key: "BorderSize", Tag: nbt.TAG_Double, Value: float64(60000000.0)},
		{Key: "BorderSizeLerpTarget", Tag: nbt.TAG_Double, Value: float64(60000000.0)},
		{Key: "BorderSizeLerpTime", Tag: nbt.TAG_Long, Value: int64(0)},
		{Key: "BorderWarningBlocks", Tag: nbt.TAG_Double, Value: float64(5.0)},
		{Key: "BorderWarningTime", Tag: nbt.TAG_Double, Value: float64(15.0)},
		{Key: "DataVersion", Tag: nbt.TAG_Int, Value: int32(1139)},
		{Key: "DayTime", Tag: nbt.TAG_Long, Value: int64(0)},
		{Key: "Difficulty", Tag: nbt.TAG_Byte, Value: byte(1)},
		{Key: "DifficultyLocked", Tag: nbt.TAG_Byte, Value: byte(0)},
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
		{Key: "Version", Tag: nbt.TAG_Compound, Value: versionPayload},
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

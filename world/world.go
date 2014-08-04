// Minecraft world package.

package world

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/mathuin/terroir/nbt"
)

var Debug = false

type World struct {
	Name       string
	spawnX     int32
	spawnY     int32
	spawnZ     int32
	spawnSet   bool
	randomSeed int64
}

func NewWorld(Name string) *World {
	if Debug {
		log.Printf("NEW WORLD: %s\n", Name)
	}
	return &World{Name: Name}
}

func MakeWorld(Name string) World {
	if Debug {
		log.Printf("MAKE WORLD: %s\n", Name)
	}
	return World{Name: Name}
}

func (w World) String() string {
	return fmt.Sprintf("World{Name: %s, Spawn: (%d, %d, %d), RandomSeed: %d}", w.Name, w.spawnX, w.spawnY, w.spawnZ, w.randomSeed)
}

func (w *World) SetRandomSeed(seed int64) {
	if Debug {
		log.Printf("SET SEED: %s: %d\n", w.Name, seed)
	}
	w.randomSeed = seed
}

func (w *World) SetSpawn(x int32, y int32, z int32) {
	if Debug {
		if w.spawnSet {
			log.Printf("CHANGE SPAWN: %s: from (%d, %d, %d) to (%d, %d, %d)\n", w.Name, w.spawnX, w.spawnY, w.spawnZ, x, y, z)
		} else {
			log.Printf("SET SPAWN: %s: (%d, %d, %d)\n", w.Name, x, y, z)
		}
	}
	w.spawnX = x
	w.spawnY = y
	w.spawnZ = z
	w.spawnSet = true
}

func (w World) Write(dir string) error {
	// make sure the directory exists and is writeable
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(dir, 0775)
		} else {
			return err
		}
	}
	worldDir := path.Join(dir, w.Name)
	if _, err := os.Stat(worldDir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(worldDir, 0775)
		} else {
			return err
		}
	}

	// write level
	if err := w.writelevel(worldDir); err != nil {
		return err
	}

	return nil
}

func ReadWorld(dir string, name string) (*World, error) {
	// does dir exist?
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("save dir does not exist")
		} else {
			return nil, err
		}
	}
	// does dir+name exist?
	worldDir := path.Join(dir, name)
	if _, err := os.Stat(worldDir); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("world dir does not exist")
		} else {
			return nil, err
		}
	}
	// does dir+name+region exist?
	regionDir := path.Join(worldDir, "region")
	if _, err := os.Stat(regionDir); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("region dir does not exist")
		} else {
			return nil, err
		}
	}
	// read level file
	levelFile := path.Join(worldDir, "level.dat")
	levelTag, err := nbt.ReadCompressedFile(levelFile)
	if err != nil {
		return nil, err
	}
	_ = levelTag
	// make a new world
	w := NewWorld(name)
	// for file in region dir
	rd, err := ioutil.ReadDir(regionDir)
	if err != nil {
		return nil, err
	}
	for _, fi := range rd {
		if strings.HasSuffix(fi.Name(), "mca") {
			// read file and import chunks
		}
	}

	return w, nil
}

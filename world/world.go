// Minecraft world package.

package world

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/mathuin/terroir/nbt"
)

var Debug = false
var NowDebug = false

type World struct {
	Name       string
	spawn      Point
	spawnSet   bool
	randomSeed int64
	chunkMap   map[string]Chunk
}

func NewWorld(Name string) *World {
	if Debug {
		log.Printf("NEW WORLD: %s", Name)
	}
	return &World{Name: Name}
}

func MakeWorld(Name string) World {
	if Debug {
		log.Printf("MAKE WORLD: %s", Name)
	}
	return World{Name: Name}
}

func (w World) String() string {
	return fmt.Sprintf("World{Name: %s, Spawn: %v, RandomSeed: %d}", w.Name, w.spawn, w.randomSeed)
}

func (w *World) SetRandomSeed(seed int64) {
	if Debug {
		log.Printf("SET SEED: %s: %d", w.Name, seed)
	}
	w.randomSeed = seed
}

func (w *World) SetSpawn(p Point) {
	if Debug {
		if w.spawnSet {
			log.Printf("CHANGE SPAWN: %s: from (%d, %d, %d) to (%d, %d, %d)", w.Name, w.spawn.X, w.spawn.Y, w.spawn.Z, p.X, p.Y, p.Z)
		} else {
			log.Printf("SET SPAWN: %s: (%d, %d, %d)", w.Name, p.X, p.Y, p.Z)
		}
	}
	w.spawn = p
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
	regionRE, err := regexp.Compile("r\\.(-?\\d*)\\.(-?\\d*)\\.mca")
	rd, err := ioutil.ReadDir(regionDir)
	if err != nil {
		return nil, err
	}
	for _, fi := range rd {
		matches := regionRE.FindAllStringSubmatch(fi.Name(), -1)
		if matches == nil {
			continue
		}
		match := matches[0]
		mfn := match[0]
		outx, xerr := strconv.ParseInt(match[1], 10, 32)
		if xerr != nil {
			panic(xerr)
		}
		outz, zerr := strconv.ParseInt(match[2], 10, 32)
		if zerr != nil {
			panic(zerr)
		}
		rname := path.Join(regionDir, mfn)
		log.Printf("regionfile name %s", rname)
		r, rerr := os.Open(rname)
		if rerr != nil {
			panic(rerr)
		}
		defer r.Close()
		wooregion := ReadRegion(r, int32(outx), int32(outz))
		_ = wooregion
	}

	return w, nil
}

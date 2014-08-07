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

type XZ struct {
	X int32
	Z int32
}

func (xz XZ) String() string {
	return fmt.Sprintf("(%d, %d)", xz.X, xz.Z)
}

type World struct {
	Name       string
	Spawn      Point
	spawnSet   bool
	RandomSeed int64
	ChunkMap   map[XZ]Chunk
	RegionMap  map[XZ][]XZ
}

func MakeWorld(Name string) World {
	if Debug {
		log.Printf("MAKE WORLD: %s", Name)
	}
	ChunkMap := map[XZ]Chunk{}
	RegionMap := map[XZ][]XZ{}
	return World{Name: Name, ChunkMap: ChunkMap, RegionMap: RegionMap}
}

func (w World) String() string {
	return fmt.Sprintf("World{Name: %s, Spawn: %v, RandomSeed: %d}", w.Name, w.Spawn, w.RandomSeed)
}

func (w *World) SetRandomSeed(seed int64) {
	if Debug {
		log.Printf("SET SEED: %s: %d", w.Name, seed)
	}
	w.RandomSeed = seed
}

func (w *World) SetSpawn(p Point) {
	if Debug {
		if w.spawnSet {
			log.Printf("CHANGE SPAWN: %s: from (%d, %d, %d) to (%d, %d, %d)", w.Name, w.Spawn.X, w.Spawn.Y, w.Spawn.Z, p.X, p.Y, p.Z)
		} else {
			log.Printf("SET SPAWN: %s: (%d, %d, %d)", w.Name, p.X, p.Y, p.Z)
		}
	}
	w.Spawn = p
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
	regionDir := path.Join(worldDir, "region")
	if _, err := os.Stat(regionDir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(regionDir, 0775)
		} else {
			return err
		}
	}

	// write level
	if err := w.writelevel(worldDir); err != nil {
		return err
	}

	for key := range w.RegionMap {
		if err := w.WriteRegion(regionDir, key); err != nil {
			return err
		}
	}

	return nil
}

func ReadWorld(dir string, name string) (*World, error) {

	var spawn Point
	var rSeed int64

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
	if Debug {
		log.Printf("Reading level file %s", levelFile)
	}
	levelTag, err := nbt.ReadCompressedFile(levelFile)
	if err != nil {
		return nil, err
	}

	// sanity check level
	requiredTags := map[string]bool{
		"LevelName":  false,
		"SpawnX":     false,
		"SpawnY":     false,
		"SpawnZ":     false,
		"RandomSeed": false,
		"version":    false,
	}

	topPayload := levelTag.Payload.([]nbt.Tag)
	if len(topPayload) != 1 {
		return nil, fmt.Errorf("levelTag does not contain only one tag\n")
	}
	dataTag := topPayload[0]
	if dataTag.Name != "Data" {
		return nil, fmt.Errorf("Data is not name of top inner tag!\n")
	}
	for _, tag := range dataTag.Payload.([]nbt.Tag) {
		switch tag.Name {
		case "LevelName":
			if tag.Payload.(string) != name {
				return nil, fmt.Errorf("Name does not match\n")
			}
			requiredTags[tag.Name] = true
		case "SpawnX":
			spawn.X = tag.Payload.(int32)
			requiredTags[tag.Name] = true
		case "SpawnY":
			spawn.Y = tag.Payload.(int32)
			requiredTags[tag.Name] = true
		case "SpawnZ":
			spawn.Z = tag.Payload.(int32)
			requiredTags[tag.Name] = true
		case "RandomSeed":
			rSeed = tag.Payload.(int64)
			requiredTags[tag.Name] = true
		case "version":
			if tag.Payload.(int32) != int32(19133) {
				return nil, fmt.Errorf("version does not match\n")
			}
			requiredTags[tag.Name] = true
		}
	}
	for rtkey, rtval := range requiredTags {
		if rtval == false {
			return nil, fmt.Errorf("tag name %s required for section but not found", rtkey)
		}
	}

	// make a new world
	w := MakeWorld(name)
	w.SetRandomSeed(rSeed)
	w.SetSpawn(spawn)

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
		if Debug {
			log.Printf("Reading region file %s", rname)
		}
		r, rerr := os.Open(rname)
		if rerr != nil {
			panic(rerr)
		}
		defer r.Close()
		n, rerr := w.ReadRegion(r, int32(outx), int32(outz))
		if rerr != nil {
			panic(rerr)
		}
		if Debug {
			log.Printf("... read %d chunks", n)
		}
	}

	return &w, nil
}

func (w *World) Block(pt Point) (int, error) {
	s, err := w.Section(pt)
	if err != nil {
		return 0, err
	}
	base := int(s.Blocks[pt.Index()])
	add := int(Nibble(s.Add, pt.Index()))
	retval := base + add*256
	return retval, nil
}

func (w *World) SetBlock(pt Point, b int) error {
	base := byte(b % 256)
	add := byte(b / 256)
	s, err := w.Section(pt)
	if err != nil {
		return err
	}
	i := pt.Index()
	s.Blocks[i] = byte(base)
	WriteNibble(s.Add, i, add)
	return nil
}

// TODO: generalize this
func (w *World) Data(pt Point) (byte, error) {
	s, err := w.Section(pt)
	if err != nil {
		return 0, err
	}
	return Nibble(s.Data, pt.Index()), nil
}

func (w *World) SetData(pt Point, b byte) error {
	s, err := w.Section(pt)
	if err != nil {
		return err
	}
	WriteNibble(s.Data, pt.Index(), b)
	return nil
}

func (w *World) BlockLight(pt Point) (byte, error) {
	s, err := w.Section(pt)
	if err != nil {
		return 0, err
	}
	return Nibble(s.BlockLight, pt.Index()), nil
}

func (w *World) SetBlockLight(pt Point, b byte) error {
	s, err := w.Section(pt)
	if err != nil {
		return err
	}
	WriteNibble(s.BlockLight, pt.Index(), b)
	return nil
}

func (w *World) SkyLight(pt Point) (byte, error) {
	s, err := w.Section(pt)
	if err != nil {
		return 0, err
	}
	return Nibble(s.SkyLight, pt.Index()), nil
}

func (w *World) SetSkyLight(pt Point, b byte) error {
	s, err := w.Section(pt)
	if err != nil {
		return err
	}
	WriteNibble(s.SkyLight, pt.Index(), b)
	return nil
}

func (w World) Section(pt Point) (*Section, error) {
	cXZ := pt.ChunkXZ()
	yf := int(floor(pt.Y, 16))
	if c, ok := w.ChunkMap[cXZ]; ok {
		if s, ok := c.Sections[yf]; ok {
			return &s, nil
		}
		return nil, fmt.Errorf("section %d of chunk %v not present", yf, cXZ)
	}
	return nil, fmt.Errorf("chunk %v not present", cXZ)
}

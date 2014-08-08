// Minecraft world package.

package world

import (
	"encoding/binary"
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
	SaveDir    string
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

func (w *World) SetSaveDir(dir string) error {
	if Debug {
		log.Printf("SET SAVE DIR: %s: %s", w.Name, dir)
	}

	if err := os.MkdirAll(dir, 0775); err != nil {
		return err
	}
	w.SaveDir = dir
	return nil
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

func (w World) Write() error {
	if Debug {
		log.Printf("WRITE WORLD: %s", w.Name)
	}

	if w.SaveDir == "" {
		return fmt.Errorf("world savedir not set")
	}

	// write level
	if err := w.writeLevel(); err != nil {
		return err
	}

	// write regions
	if err := w.writeRegions(); err != nil {
		return err
	}
	return nil
}

func ReadWorld(dir string, name string, loadAllChunks bool) (*World, error) {

	var spawn Point
	var rSeed int64

	// read level file
	worldDir := path.Join(dir, name)
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
	w.SetSaveDir(dir)
	w.SetRandomSeed(rSeed)
	w.SetSpawn(spawn)

	if loadAllChunks {
		if err := w.loadAllChunksFromAllRegions(); err != nil {
			return nil, err
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
	c, ok := w.ChunkMap[cXZ]
	if !ok {
		cp, lerr := w.loadChunk(cXZ)
		if lerr != nil {
			var emptytag nbt.Tag
			mcp, merr := w.MakeChunk(cXZ, emptytag)
			if merr != nil {
				return nil, merr
			}
			cp = mcp
		}
		c = *cp

	}
	s, ok := c.Sections[yf]
	if !ok {
		c.Sections[yf] = MakeSection()
		s = c.Sections[yf]
	}
	return &s, nil
}

// arguments are in chunk coordinates
func (w *World) loadChunk(cXZ XZ) (*Chunk, error) {
	if Debug {
		log.Printf("LOAD CHUNK: %s: %v", w.Name, cXZ)
	}
	rXZ := XZ{X: floor(cXZ.X, 32), Z: floor(cXZ.Z, 32)}
	rname := w.regionFilename(rXZ)
	r, err := os.Open(rname)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	cindex := (cXZ.Z%32)*32 + (cXZ.X % 32)

	_, lerr := r.Seek(int64(cindex*4), os.SEEK_SET)
	if lerr != nil {
		return nil, lerr
	}

	var location int32

	err = binary.Read(r, binary.BigEndian, &location)
	if err != nil {
		return nil, err
	}

	if location == 0 {
		return nil, fmt.Errorf("location is zero: chunk not in region file")
	}

	return w.loadChunkFromRegion(r, location, cXZ)
}

func (w *World) MakeChunk(xz XZ, tag nbt.Tag) (*Chunk, error) {
	c := MakeChunk(xz.X, xz.Z)
	var emptytag nbt.Tag
	if tag != emptytag {
		c.Read(tag)
		if c.xPos != xz.X || c.zPos != xz.Z {
			return nil, fmt.Errorf("tag position (%d, %d) did not match XZ %v", c.xPos, c.zPos, xz)
		}
	}
	if err := w.addChunkToMaps(c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (w *World) addChunkToMaps(c Chunk) error {
	cXZ := XZ{X: c.xPos, Z: c.zPos}
	if _, ok := w.ChunkMap[cXZ]; ok {
		return fmt.Errorf("chunk %v already exists", cXZ)
	}
	w.ChunkMap[cXZ] = c
	rXZ := XZ{X: floor(cXZ.X, 32), Z: floor(cXZ.Z, 32)}
	for _, cptr := range w.RegionMap[rXZ] {
		if cptr == cXZ {
			return fmt.Errorf("chunk %v already in region map", cXZ)
		}
	}
	w.RegionMap[rXZ] = append(w.RegionMap[rXZ], cXZ)
	return nil
}

func (w *World) loadAllChunksFromAllRegions() error {
	rXZList := make([]XZ, 0)

	regionRE, err := regexp.Compile("r\\.(-?\\d*)\\.(-?\\d*)\\.mca")
	regionDir := path.Join(w.SaveDir, w.Name, "region")
	rd, err := ioutil.ReadDir(regionDir)
	if err != nil {
		return err
	}

	for _, fi := range rd {
		matches := regionRE.FindAllStringSubmatch(fi.Name(), -1)
		if matches == nil {
			continue
		}
		match := matches[0]
		outx, xerr := strconv.ParseInt(match[1], 10, 32)
		if xerr != nil {
			return xerr
		}
		outz, zerr := strconv.ParseInt(match[2], 10, 32)
		if zerr != nil {
			return zerr
		}

		rXZ := XZ{X: int32(outx), Z: int32(outz)}
		rXZList = append(rXZList, rXZ)
	}

	for _, rXZ := range rXZList {
		// TODO: parallelize this .. or skip it entirely!
		_, rerr := w.loadAllChunksFromRegion(rXZ)
		if rerr != nil {
			return rerr
		}
	}
	return nil
}

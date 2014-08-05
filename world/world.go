// Minecraft world package.

package world

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/mathuin/terroir/nbt"
)

var Debug = false

type XZ struct {
	X int32
	Z int32
}

type World struct {
	Name       string
	Spawn      Point
	spawnSet   bool
	RandomSeed int64
	ChunkMap   map[XZ]Chunk
	RegionMap  map[XZ][]XZ
}

func NewWorld(Name string) *World {
	if Debug {
		log.Printf("NEW WORLD: %s", Name)
	}
	ChunkMap := map[XZ]Chunk{}
	RegionMap := map[XZ][]XZ{}
	return &World{Name: Name, ChunkMap: ChunkMap, RegionMap: RegionMap}
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
	w := NewWorld(name)
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

	return w, nil
}

// this can go back to region.go in a post-region world
func (w *World) ReadRegion(r io.ReadSeeker, xCoord int32, zCoord int32) (int, error) {
	numchunks := 0

	// build the data structures
	locations := make([]int32, 1024)
	timestamps := make([]int32, 1024)

	// populate them
	err := binary.Read(r, binary.BigEndian, locations)
	if err != nil {
		return numchunks, err
	}
	err = binary.Read(r, binary.BigEndian, timestamps)
	if err != nil {
		return numchunks, err
	}

	for i := 0; i < 1024; i++ {
		// coordinates
		x := int(xCoord)*32 + i%32
		z := int(zCoord)*32 + i/32

		offcount := locations[i]
		offsetval := offcount / 256
		countval := offcount % 256
		timestamp := timestamps[i]
		if timestamp > 0 || offsetval > 0 || countval > 0 {
			if Debug {
				log.Printf("[%d, %d]", x, z)
				log.Printf("  offset %d sectors (%d bytes)", offsetval, offsetval*4096)
				log.Printf("  count %d sectors (%d bytes)", countval, countval*4096)
				log.Printf("  timestamp %d", timestamp)
			}
			pos, perr := r.Seek(int64(offsetval*4096), os.SEEK_SET)
			if perr != nil {
				panic(perr)
			}
			if Debug {
				log.Printf("Current seek position (read) %d", pos)
			}
			var chunklen int32
			err = binary.Read(r, binary.BigEndian, &chunklen)
			if err != nil {
				return numchunks, err
			}
			if Debug {
				log.Printf("Actual read: %d bytes (%d bytes padding)", chunklen, (countval*4096 - chunklen))
			}
			flag := make([]uint8, 1)
			_, err = io.ReadFull(r, flag)
			if err != nil {
				return numchunks, err
			}
			zchr := make([]byte, chunklen)
			var zr, unzr io.Reader
			zr = bytes.NewBuffer(zchr)
			ret, err := io.ReadFull(r, zchr)
			if err != nil {
				return numchunks, err
			}
			if Debug {
				log.Printf("%d compressed bytes read", ret)
			}
			if Debug {
				log.Printf("Compression:")
			}
			switch flag[0] {
			case 0:
				if Debug {
					log.Printf("  none?")
				}
				unzr = zr
			case 1:
				if Debug {
					log.Printf("  gzip")
				}
				unzr, err = gzip.NewReader(zr)
				if err != nil {
					return numchunks, err
				}
			case 2:
				if Debug {
					log.Printf("  zlib")
				}
				unzr, err = zlib.NewReader(zr)
				if err != nil {
					return numchunks, err
				}
			}
			zstr, err := ioutil.ReadAll(unzr)
			if err != nil {
				return numchunks, err
			}
			if Debug {
				log.Printf("uncompressed len %d", len(zstr))
			}
			writeNotParse := false
			var tag nbt.Tag
			tmpchunk := MakeChunk(int32(x), int32(z))
			if writeNotParse {
				writeFileName := fmt.Sprintf("chunk.%d.%d.dat", x, z)
				err = ioutil.WriteFile(writeFileName, zstr, 0755)
				if err != nil {
					return numchunks, err
				}
				if Debug {
					log.Println(writeFileName)
				}
			} else {
				zb := bytes.NewBuffer(zstr)
				tag, err = nbt.ReadTag(zb)
				if err != nil {
					return numchunks, err
				}
				tmpchunk.Read(tag)
			}
			cXZ := XZ{X: tmpchunk.xPos, Z: tmpchunk.zPos}
			w.ChunkMap[cXZ] = tmpchunk
			rXZ := XZ{X: int32(math.Floor(float64(cXZ.X) / 32.0)), Z: int32(math.Floor(float64(cXZ.Z) / 32.0))}
			w.RegionMap[rXZ] = append(w.RegionMap[rXZ], cXZ)
			numchunks = numchunks + 1
		}
	}
	return numchunks, nil
}

func (w *World) Block(pt Point) int {
	s := w.Section(pt)
	base := int(s.Blocks[pt.Index()])
	add := int(Nibble(s.Add, pt.Index()))
	return base + add*256
}

func (w *World) SetBlock(pt Point, b int) {
	base := byte(b % 256)
	add := byte(b / 256)
	s := w.Section(pt)
	i := pt.Index()
	s.Blocks[i] = byte(base)
	WriteNibble(s.Add, i, add)
}

func (w *World) Data(pt Point) byte {
	return Nibble(w.Section(pt).Data, pt.Index())
}

func (w *World) SetData(pt Point, b byte) {
	WriteNibble(w.Section(pt).Data, pt.Index(), b)
}

func (w *World) BlockLight(pt Point) byte {
	return Nibble(w.Section(pt).BlockLight, pt.Index())
}

func (w *World) SetBlockLight(pt Point, b byte) {
	WriteNibble(w.Section(pt).BlockLight, pt.Index(), b)
}

func (w *World) SkyLight(pt Point) byte {
	return Nibble(w.Section(pt).SkyLight, pt.Index())
}

func (w *World) SetSkyLight(pt Point, b byte) {
	WriteNibble(w.Section(pt).SkyLight, pt.Index(), b)
}

func (w *World) ReplaceBlock(from byte, to byte) int {
	count := 0
	for _, c := range w.ChunkMap {
		for _, s := range c.Sections {
			for i := range s.Blocks {
				if s.Blocks[i] == from {
					count = count + 1
					s.Blocks[i] = to
				}
			}
		}
	}
	return count
}

func (w *World) WriteRegion(dir string, key XZ) error {
	cb := new(bytes.Buffer)
	locations := make([]int32, 1024)
	timestamps := make([]int32, 1024)
	zlibcomp := make([]byte, 1)
	zlibcomp[0] = byte(2)
	offset := int32(2)

	numchunks := 0
	for _, v := range w.RegionMap[key] {
		if Debug {
			log.Printf("Writing %d, %d...", v.X, v.Z)
		}
		c := w.ChunkMap[v]
		cx := c.xPos % 32
		if cx < 0 {
			cx = cx + 32
		}
		cz := c.zPos % 32
		if cz < 0 {
			cz = cz + 32
		}
		arroff := cz*32 + cx
		if Debug {
			log.Printf("arroff: (%d, %d) -> %d * 32 + %d = %d", c.xPos, c.zPos, cz, cx, arroff)
		}

		// write chunk to compressed buffer
		var zb bytes.Buffer
		zw := zlib.NewWriter(&zb)
		ct := c.write()
		if err := ct.Write(zw); err != nil {
			return err
		}
		zw.Close()

		// - calculate lengths
		// (the extra byte is the compression byte)
		ccl := int32(zb.Len() + 1)
		count := int32(math.Ceil(float64(ccl) / 4096.0))
		pad := int32(4096*count) - ccl - 4
		whole := int(ccl + pad + 4)

		if Debug {
			log.Printf("Length of compressed chunk: %d", ccl)
			log.Printf("Count of sectors: %d", count)
			log.Printf("Padding: %d", pad)
			log.Printf("Whole amount written: %d", whole)
		}

		if pad > 4096 {
			return fmt.Errorf("pad %d > 4096", pad)
		}

		if (whole % 4096) != 0 {
			return fmt.Errorf("%d not even multiple of 4096", whole)
		}

		posb := cb.Len()
		if Debug {
			log.Printf("Chunk %d seek position before %d", numchunks, posb)
			log.Printf("(%d after two tables added)", 8192+posb)
		}
		posshould := (offset - 2) * 4096
		if posb != int(posshould) {
			return fmt.Errorf("posb for chunk %d is %d but should be %d!", numchunks, posb, posshould)
		}

		// - write chunk header and compressed chunk data to chunk writer
		err := binary.Write(cb, binary.BigEndian, ccl)
		if err != nil {
			return err
		}
		_, err = cb.Write(zlibcomp)
		if err != nil {
			return err
		}
		_, err = zb.WriteTo(cb)
		if err != nil {
			return err
		}

		// - write necessary padding of zeroes to chunks writer
		padb := make([]byte, pad)
		_, err = cb.Write(padb)
		if err != nil {
			return err
		}

		posa := cb.Len()
		if Debug {
			log.Printf("Chunk %d seek position after %d", numchunks, posa)
			log.Printf("(%d after two tables added)", 8192+posa)

		}
		if int32(posa-posb) != count*4096 {
			return fmt.Errorf("posa for %d, %d is %d -- not even multiple of 4096!", cx, cz, posa)
		}

		// - write current time to timestamp array
		timestamps[arroff] = int32(time.Now().Unix())
		if Debug {
			log.Printf("Timestamps are %d", timestamps[arroff])
		}

		// - write offset and count to locations array
		locations[arroff] = offset*256 + int32(count)
		if Debug {
			log.Printf("Locations are %d (%d * 256 + %d)", locations[arroff], offset, count)
		}
		offset = offset + count
		numchunks = numchunks + 1
	}

	// open actual region file for writing
	rfn := fmt.Sprintf("r.%d.%d.mca", key.X, key.Z)
	rname := path.Join(dir, rfn)
	if Debug {
		log.Printf("Writing region file %s...", rname)
	}
	iow, err := os.Create(rname)
	if err != nil {
		return err
	}
	defer iow.Close()

	// write locations array to real io.writer
	err = binary.Write(iow, binary.BigEndian, locations)
	if err != nil {
		return err
	}

	// write timestamps array to real io.writer
	err = binary.Write(iow, binary.BigEndian, timestamps)
	if err != nil {
		return err
	}
	// write chunks writer to real io.writer
	_, err = cb.WriteTo(iow)
	if err != nil {
		return err
	}
	if Debug {
		log.Printf("... wrote %d chunks", numchunks)
	}
	return nil
}

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
	"time"

	"github.com/mathuin/terroir/nbt"
)

type Chunk struct {
	xPos      int32
	zPos      int32
	biomes    []byte
	heightMap []int32
	Sections  map[int]Section
	// living things
	entities     []Entity
	tileEntities []TileEntity
	tileTicks    []TileTick
}

func MakeChunk(xPos int32, zPos int32) Chunk {
	if Debug {
		log.Printf("MAKE CHUNK: xPos %d, zPos %d", xPos, zPos)
	}
	biomes := make([]byte, 256)
	heightMap := make([]int32, 256)
	Sections := map[int]Section{}
	entities := []Entity{}
	tileEntities := []TileEntity{}
	tileTicks := []TileTick{}
	return Chunk{xPos: xPos, zPos: zPos, biomes: biomes, heightMap: heightMap, Sections: Sections, entities: entities, tileEntities: tileEntities, tileTicks: tileTicks}
}

func (c Chunk) Name() string {
	return fmt.Sprintf("%d, %d", c.xPos, c.zPos)
}

func (c Chunk) write() nbt.Tag {

	sectionsPayload := [][]nbt.Tag{}
	for i, s := range c.Sections {
		st := s.write(i)
		sectionsPayload = append(sectionsPayload, st)
	}

	entitiesPayload := [][]nbt.Tag{}
	for _, e := range c.entities {
		entitiesPayload = append(entitiesPayload, e.write())
	}

	tileEntitiesPayload := [][]nbt.Tag{}
	for _, te := range c.tileEntities {
		tileEntitiesPayload = append(tileEntitiesPayload, te.write())
	}

	tileTicksPayload := [][]nbt.Tag{}
	for _, tt := range c.tileTicks {
		ttt := tt.write()
		tileTicksPayload = append(tileTicksPayload, ttt.Payload.([]nbt.Tag))
	}

	var levelElems = []nbt.CompoundElem{
		{"xPos", nbt.TAG_Int, c.xPos},
		{"zPos", nbt.TAG_Int, c.zPos},
		{"LastUpdate", nbt.TAG_Long, int64(0)},
		{"LightPopulated", nbt.TAG_Byte, byte(0)},
		{"TerrainPopulated", nbt.TAG_Byte, byte(1)},
		{"V", nbt.TAG_Byte, byte(1)},
		{"InhabitedTime", nbt.TAG_Long, int64(0)},
		{"Biomes", nbt.TAG_Byte_Array, []byte(c.biomes)},
		{"HeightMap", nbt.TAG_Int_Array, []int32(c.heightMap)},
		{"Sections", nbt.TAG_List, sectionsPayload},
		{"Entities", nbt.TAG_List, entitiesPayload},
		{"TileEntities", nbt.TAG_List, tileEntitiesPayload},
	}

	if len(c.tileTicks) > 0 {
		tickElem := nbt.CompoundElem{"TileTicks", nbt.TAG_List, tileTicksPayload}
		levelElems = append(levelElems, tickElem)
	}

	levelTag := nbt.MakeCompound("Level", levelElems)

	topTag := nbt.MakeTag(nbt.TAG_Compound, "")
	if err := topTag.SetPayload([]nbt.Tag{levelTag}); err != nil {
		panic(err)
	}

	return topTag
}

func (c *Chunk) Read(t nbt.Tag) error {
	requiredTags := map[string]bool{
		"xPos":         false,
		"zPos":         false,
		"Biomes":       false,
		"HeightMap":    false,
		"Sections":     false,
		"Entities":     false,
		"TileEntities": false,
	}

	if t.Type != nbt.TAG_Compound {
		return fmt.Errorf("top tag type not TAG_Compound!")
	}

	if t.Name != "" {
		return fmt.Errorf("top tag not unnamed")
	}

	levelTag := t.Payload.([]nbt.Tag)[0]

	if levelTag.Type != nbt.TAG_Compound {
		return fmt.Errorf("level tag type not TAG_Compound!")
	}

	if levelTag.Name != "Level" {
		return fmt.Errorf("level tag not named Level")
	}

	for _, tval := range levelTag.Payload.([]nbt.Tag) {
		if Debug {
			log.Printf("tag type %s name %s found", nbt.Names[tval.Type], tval.Name)
		}
		if _, ok := requiredTags[tval.Name]; ok {
			if Debug {
				log.Printf(" -- tag is required")
			}
			requiredTags[tval.Name] = true
			switch tval.Name {
			case "xPos":
				xPos := tval.Payload.(int32)
				if xPos != c.xPos {
					return fmt.Errorf("xPos %d does not match c.xPos %d", xPos, c.xPos)
				}
			case "zPos":
				zPos := tval.Payload.(int32)
				if zPos != c.zPos {
					return fmt.Errorf("zPos %d does not match c.zPos %d", zPos, c.zPos)
				}
			case "Biomes":
				c.biomes = tval.Payload.([]byte)
			case "HeightMap":
				c.heightMap = tval.Payload.([]int32)
			case "Sections":
				for _, s := range tval.Payload.([][]nbt.Tag) {
					var yValFound bool
					var yVal int
					for _, subtag := range s {
						if subtag.Name == "Y" {
							yValFound = true
							yVal = int(subtag.Payload.(byte))
						}
					}
					if !yValFound {
						return fmt.Errorf("no yVal found")
					}
					if _, ok := c.Sections[yVal]; ok {
						return fmt.Errorf("yVal already found")
					}
					news, err := ReadSection(s)
					if err != nil {
						return err
					}
					c.Sections[yVal] = *news
				}
			case "Entities":
				if tval.Payload != nil {
					es := make([]Entity, 0)
					for _, e := range tval.Payload.([][]nbt.Tag) {
						es = append(es, ReadEntity(e))
					}
					c.entities = es
				}
			case "TileEntities":
				if tval.Payload != nil {
					tes := make([]TileEntity, 0)
					for _, te := range tval.Payload.([][]nbt.Tag) {
						tes = append(tes, ReadTileEntity(te))
					}
					c.tileEntities = tes
				}
			}
		} else {
			// optional tags
			switch tval.Name {
			case "TileTicks":
				tts := make([]TileTick, 0)
				for _, tt := range tval.Payload.([][]nbt.Tag) {
					tts = append(tts, ReadTileTick(tt))
				}
				c.tileTicks = tts
			}
		}
	}

	for rtkey, rtval := range requiredTags {
		if rtval == false {
			return fmt.Errorf("tag name %s required for chunk but not found", rtkey)
		}
	}

	return nil
}

type CTROut struct {
	arroff int32
	count  int32
	arrout []byte
	err    error
}

func WriteChunkToRegion(in chan Chunk, out chan CTROut, i int) {
	for c := range in {
		cout := new(CTROut)
		cb := new(bytes.Buffer)
		// JMT: not too hard to support gzip now...
		comptype := byte(2)

		cx := c.xPos % 32
		if cx < 0 {
			cx = cx + 32
		}
		cz := c.zPos % 32
		if cz < 0 {
			cz = cz + 32
		}
		cout.arroff = cz*32 + cx
		if Debug {
			log.Printf("arroff: (%d, %d) -> %d * 32 + %d = %d", i, c.xPos, c.zPos, cz, cx, cout.arroff)
		}

		// write chunk to compressed buffer
		ct := c.write()
		var zb bytes.Buffer
		zw := zlib.NewWriter(&zb)
		start := time.Now().UnixNano()
		cout.err = ct.Write(zw)
		end := time.Now().UnixNano()
		if Debug {
			log.Printf("ct.Write(zw) took %d nanoseconds", end-start)
		}
		if cout.err != nil {
			out <- *cout
			continue
		}
		zw.Close()

		// - calculate lengths
		// (the extra byte is the compression byte)
		ccl := int32(zb.Len() + 1)
		cout.count = int32(math.Ceil(float64(ccl) / 4096.0))
		pad := int32(4096*cout.count) - ccl - 4
		whole := int(ccl + pad + 4)

		if Debug {
			log.Printf("Length of compressed chunk: %d", ccl)
			log.Printf("Count of sectors: %d", cout.count)
			log.Printf("Padding: %d", pad)
			log.Printf("Whole amount written: %d", whole)
		}

		if pad > 4096 {
			cout.err = fmt.Errorf("pad %d > 4096", pad)
			out <- *cout
			continue
		}

		if (whole % 4096) != 0 {
			cout.err = fmt.Errorf("%d not even multiple of 4096", whole)
			out <- *cout
			continue
		}

		// - write chunk header and compressed chunk data to chunk writer
		cout.err = binary.Write(cb, binary.BigEndian, ccl)
		if cout.err != nil {
			out <- *cout
			continue
		}

		cout.err = cb.WriteByte(comptype)
		if cout.err != nil {
			out <- *cout
			continue
		}

		_, cout.err = zb.WriteTo(cb)
		if cout.err != nil {
			out <- *cout
			continue
		}

		// - write necessary padding of zeroes to chunks writer
		padb := make([]byte, pad)
		_, cout.err = cb.Write(padb)
		if cout.err != nil {
			out <- *cout
			continue
		}

		if cb.Len() != whole {
			cout.err = fmt.Errorf("cb.Len() %d does not match whole %d", cb.Len(), whole)
			out <- *cout
			continue
		}
		cout.arrout = cb.Bytes()

		out <- *cout
	}
}

func (w *World) loadChunkFromRegion(r io.ReadSeeker, location int32, cXZ XZ) (*Chunk, error) {
	offset := location / 256
	count := location % 256

	if Debug {
		log.Printf("%v: ", cXZ)
		log.Printf("  offset %d sectors (%d bytes)", offset, offset*4096)
		log.Printf("  count %d sectors (%d bytes)", count, count*4096)
	}

	_, perr := r.Seek(int64(offset*4096), os.SEEK_SET)
	if perr != nil {
		panic(perr)
	}

	var chunklen int32
	err := binary.Read(r, binary.BigEndian, &chunklen)
	if err != nil {
		return nil, err
	}

	if Debug {
		log.Printf("Actual read: %d bytes (%d bytes padding)", chunklen, (int32(count*4096) - chunklen))
	}

	flag := make([]uint8, 1)
	_, err = io.ReadFull(r, flag)
	if err != nil {
		return nil, err
	}
	zchr := make([]byte, chunklen)
	var zr, unzr io.Reader
	zr = bytes.NewBuffer(zchr)
	ret, err := io.ReadFull(r, zchr)
	if err != nil {
		return nil, err
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
			return nil, err
		}
	case 2:
		if Debug {
			log.Printf("  zlib")
		}
		unzr, err = zlib.NewReader(zr)
		if err != nil {
			return nil, err
		}
	}
	zstr, err := ioutil.ReadAll(unzr)
	if err != nil {
		return nil, err
	}
	if Debug {
		log.Printf("uncompressed len %d", len(zstr))
	}
	var tag nbt.Tag
	zb := bytes.NewBuffer(zstr)
	tag, err = nbt.ReadTag(zb)
	if err != nil {
		return nil, err
	}
	var c *Chunk
	c, err = w.MakeChunk(cXZ, tag)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (w *World) loadAllChunksFromRegion(rXZ XZ) (int, error) {
	if Debug {
		log.Printf("LOAD ALL CHUNKS: %s: %v", w.Name, rXZ)
	}
	rname := w.regionFilename(rXZ)
	if Debug {
		log.Printf("Reading region file %s", rname)
	}
	r, err := os.Open(rname)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	locations := make([]int32, 1024)

	err = binary.Read(r, binary.BigEndian, &locations)

	numchunks := 0
	for i, location := range locations {
		// eventually parallelize this
		// mutexes around chunkmap and regionmap of course
		if location != 0 {
			cXZ := XZ{X: rXZ.X*32 + int32(i%32), Z: rXZ.Z*32 + int32(i/32)}
			_, err := w.loadChunkFromRegion(r, location, cXZ)
			if err != nil {
				return 0, err
			}
			numchunks = numchunks + 1
		}
	}
	if Debug {
		log.Printf("... read %d chunks", numchunks)
	}
	return numchunks, nil
}

// region files and such

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

type Region struct {
	xCoord int32
	zCoord int32
	chunks map[string]Chunk
}

func NewRegion(xCoord int32, zCoord int32) *Region {
	if Debug {
		log.Printf("NEW REGION: xCoord %d, zCoord %d", xCoord, zCoord)
	}
	return &Region{xCoord: xCoord, zCoord: zCoord, chunks: make(map[string]Chunk, 0)}
}

func MakeRegion(xCoord int32, zCoord int32) Region {
	if Debug {
		log.Printf("MAKE REGION: xCoord %d, zCoord %d", xCoord, zCoord)
	}
	return Region{xCoord: xCoord, zCoord: zCoord, chunks: make(map[string]Chunk, 0)}
}

func (r *Region) Write(w io.Writer) {
	cb := new(bytes.Buffer)
	locations := make([]int32, 1024)
	timestamps := make([]int32, 1024)
	zlibcomp := make([]byte, 1)
	zlibcomp[0] = byte(2)
	offset := int32(2)

	// for each chunk in chunks list
	numchunks := 0
	for k, c := range r.chunks {
		if Debug {
			log.Printf("Writing %s...", k)
		}
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
			panic(err)
		}
		zw.Close()

		// - calculate lengths
		// (the extra byte is the compression byte)
		ccl := int32(zb.Len() + 1)
		count := int32(math.Ceil(float64(ccl) / 4096.0))
		pad := int32(4096*count) - ccl - 4

		if pad > 4096 {
			log.Printf("ccl %d count %d pad %d", ccl, count, pad)
			log.Panic("Too much padding somehow!")
		}
		if (int(ccl+pad+4) % 4096) != 0 {
			log.Printf("ccl %d count %d pad %d", ccl, count, pad)
			log.Printf("sum %d should be %d", (ccl + pad + 4), count*4096)
			log.Panic("Not an even page size!")
		}

		if Debug {
			log.Printf("Length of compressed chunk: %d", ccl)
			log.Printf("Count of sectors: %d", count)
			log.Printf("Padding: %d", pad)
		}

		posb := cb.Len()
		if Debug {
			log.Printf("Chunk %d seek position before %d", numchunks, posb)
			log.Printf("(%d after two tables added)", 8192+posb)
		}
		posshould := (offset - 2) * 4096
		if posb != int(posshould) {
			log.Panicf("posb for chunk %d is %d but should be %d!", numchunks, posb, posshould)
		}

		// - write chunk header and compressed chunk data to chunk writer
		err := binary.Write(cb, binary.BigEndian, ccl)
		if err != nil {
			panic(err)
		}
		_, err = cb.Write(zlibcomp)
		if err != nil {
			panic(err)
		}
		_, err = zb.WriteTo(cb)
		if err != nil {
			panic(err)
		}

		// - write necessary padding of zeroes to chunks writer
		padb := make([]byte, pad)
		_, err = cb.Write(padb)
		if err != nil {
			panic(err)
		}

		posa := cb.Len()
		if Debug {
			log.Printf("Chunk %d seek position after %d", numchunks, posa)
			log.Printf("(%d after two tables added)", 8192+posa)

		}
		if int32(posa-posb) != count*4096 {
			log.Panicf("posa for %d, %d is %d -- not even multiple of 4096!", cx, cz, posa)
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
	// write locations array to real io.writer
	err := binary.Write(w, binary.BigEndian, locations)
	if err != nil {
		panic(err)
	}

	// write timestamps array to real io.writer
	err = binary.Write(w, binary.BigEndian, timestamps)
	if err != nil {
		panic(err)
	}
	// write chunks writer to real io.writer
	_, err = cb.WriteTo(w)
	if err != nil {
		panic(err)
	}
}

func ReadRegion(r io.ReadSeeker, xCoord int32, zCoord int32) *Region {
	region := NewRegion(xCoord, zCoord)

	// build the data structures
	locations := make([]int32, 1024)
	timestamps := make([]int32, 1024)

	// populate them
	err := binary.Read(r, binary.BigEndian, locations)
	if err != nil {
		panic(err)
	}
	err = binary.Read(r, binary.BigEndian, timestamps)
	if err != nil {
		panic(err)
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
				panic(err)
			}
			if Debug {
				log.Printf("Actual read: %d bytes (%d bytes padding)", chunklen, (countval*4096 - chunklen))
			}
			flag := make([]uint8, 1)
			_, err = io.ReadFull(r, flag)
			if err != nil {
				panic(err)
			}
			zchr := make([]byte, chunklen)
			var zr, unzr io.Reader
			zr = bytes.NewBuffer(zchr)
			ret, err := io.ReadFull(r, zchr)
			if err != nil {
				panic(err)
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
					panic(err)
				}
			case 2:
				if Debug {
					log.Printf("  zlib")
				}
				unzr, err = zlib.NewReader(zr)
				if err != nil {
					panic(err)
				}
			}
			zstr, err := ioutil.ReadAll(unzr)
			if err != nil {
				panic(err)
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
					panic(err)
				}
				if Debug {
					log.Println(writeFileName)
				}
			} else {
				zb := bytes.NewBuffer(zstr)
				tag, err = nbt.ReadTag(zb)
				if err != nil {
					panic(err)
				}
				tmpchunk.Read(tag)
			}
			region.chunks[tmpchunk.Name()] = tmpchunk
		}
	}
	return region
}

func (r *Region) ReplaceBlock(from byte, to byte) int {
	count := 0
	for _, c := range r.chunks {
		for _, s := range c.sections {
			for i := range s.blocks {
				if s.blocks[i] == from {
					count = count + 1
					s.blocks[i] = to
				}
			}
		}
	}
	return count
}

func (r Region) Compare(newr Region) error {
	// compare xCoord and yCoord
	if r.xCoord != newr.xCoord || r.zCoord != newr.zCoord {
		return fmt.Errorf("region coordinates do not match!")
	}

	// check chunks
	for i, nc := range newr.chunks {
		if c, ok := r.chunks[i]; ok {
			newrtag := nc.write()
			newrtop := newrtag.Payload.([]nbt.Tag)[0].Payload.([]nbt.Tag)[0]
			rtag := c.write()
			rtop := rtag.Payload.([]nbt.Tag)[0].Payload.([]nbt.Tag)[0]
			if newrtop != rtop {
				return fmt.Errorf("chunks do not match!")
			}
		} else {
			return fmt.Errorf("r[%s] does not exist", i)
		}
	}
	// other way round
	for i := range r.chunks {
		if _, ok := newr.chunks[i]; !ok {
			return fmt.Errorf("newr[%s] does not exist", i)
		}
	}
	return nil
}

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
	"os"
	"path"
	"time"

	"github.com/mathuin/terroir/nbt"
)

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
			var tag nbt.Tag
			tmpchunk := MakeChunk(int32(x), int32(z))
			zb := bytes.NewBuffer(zstr)
			tag, err = nbt.ReadTag(zb)
			if err != nil {
				return numchunks, err
			}
			tmpchunk.Read(tag)
			cXZ := XZ{X: tmpchunk.xPos, Z: tmpchunk.zPos}
			w.ChunkMap[cXZ] = tmpchunk
			rXZ := XZ{X: floor(cXZ.X, 32), Z: floor(cXZ.Z, 32)}
			w.RegionMap[rXZ] = append(w.RegionMap[rXZ], cXZ)
			numchunks = numchunks + 1
		}
	}
	return numchunks, nil
}

func (w *World) WriteRegion(dir string, key XZ) error {
	cb := new(bytes.Buffer)
	locations := make([]int32, 1024)
	timestamps := make([]int32, 1024)
	offset := int32(2)

	numchunks := 0
	for _, v := range w.RegionMap[key] {
		if Debug {
			log.Printf("Writing %d, %d...", v.X, v.Z)
		}
		c := w.ChunkMap[v]
		arroff, count, arrout, err := c.WriteChunkToRegion()
		if err != nil {
			return err
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

		// - write bytes to master chunk buffer
		_, err = cb.Write(arrout)
		if err != nil {
			return err
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

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

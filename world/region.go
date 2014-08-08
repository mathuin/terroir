// region files and such

package world

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

func (w World) genChunks(key XZ, in chan Chunk) {
	for _, v := range w.RegionMap[key] {
		in <- w.ChunkMap[v]
	}
	close(in)
}

func (w *World) writeRegion(dir string, key XZ) error {
	chunks := 1024
	cb := new(bytes.Buffer)
	locations := make([]int32, chunks)
	timestamps := make([]int32, chunks)
	offset := int32(2)
	numchunks := 0

	in := make(chan Chunk)
	out := make(chan CTROut)

	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		go func(i int) {
			wg.Add(1)
			WriteChunkToRegion(in, out, i)
			wg.Done()
		}(i)
	}
	go func() { wg.Wait(); close(out) }()
	go w.genChunks(key, in)

	for cout := range out {
		if cout.err != nil {
			return cout.err
		}

		// - write current time to timestamp array
		timestamps[cout.arroff] = int32(time.Now().Unix())
		if Debug {
			log.Printf("Timestamps are %d", timestamps[cout.arroff])
		}

		// - write offset and count to locations array
		locations[cout.arroff] = offset*256 + int32(cout.count)
		if Debug {
			log.Printf("Locations are %d (%d * 256 + %d)", locations[cout.arroff], offset, cout.count)
		}

		// - write bytes to master chunk buffer
		_, err := cb.Write(cout.arrout)
		if err != nil {
			return err
		}

		offset = offset + cout.count
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

func (w *World) writeRegions() error {
	regionDir := path.Join(w.SaveDir, w.Name, "region")
	if err := os.MkdirAll(regionDir, 0775); err != nil {
		return err
	}

	for key := range w.RegionMap {
		if err := w.writeRegion(regionDir, key); err != nil {
			return err
		}
	}
	return nil
}

func (w World) regionFilename(rXZ XZ) string {
	return path.Join(w.SaveDir, w.Name, "region", fmt.Sprintf("r.%d.%d.mca", rXZ.Z, rXZ.Z))
}

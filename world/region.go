// region files and such

package world

import (
	"io"
	"log"
)

type Region struct {
	xCoord int32
	zCoord int32
	chunks []Chunk
}

func NewRegion(xCoord int32, zCoord int32) *Region {
	if Debug {
		log.Printf("NEW REGION: xCoord %d, zCoord %d", xCoord, zCoord)
	}
	sections := make([]Section, 0)
	for i := 0; i < 16; i++ {
		sections = append(sections, MakeSection())
	}
	return &Region{xCoord: xCoord, zCoord: zCoord, chunks: make([]Chunk, 0)}
}

func MakeRegion(xCoord int32, zCoord int32) Region {
	if Debug {
		log.Printf("MAKE REGION: xCoord %d, zCoord %d", xCoord, zCoord)
	}
	return Region{xCoord: xCoord, zCoord: zCoord, chunks: make([]Chunk, 0)}
}

func (r *Region) write(w io.Writer) {
	// create io.writer for chunks
	// set offset to 2
	// for each chunk in chunks list
	// - write current time to timestamp array (0, 0 to 31, 31 (z before x))
	// - write tag to new compressed dealiebobber
	// - calculate lengths
	// - write offset and count to locations array (0, 0 to 31, 31 (z before x))
	// - write compressed dealiebobber to chunks writer
	// - write necessary padding of zeroes to chunks writer

	// write locations array to real io.writer
	// write timestamps array to real io.writer
	// write chunks writer to real io.writer
}

func readRegion(r io.Reader) *Region {
	// get coords from where, file?
	// read locations array
	// read timestamps array (discard?)
	// traverse locations and make list of offsets and counts
	// for each chunk in locations array
	// read important data from region (skip padding?)
	// unwrap data
	// populate chunk list
	return nil
}

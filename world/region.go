// region files and such

package world

type ChunkLoc struct {
	Offset [3]uint8
	Count  uint8
}

type ChunkTime struct {
	Stamp int32
}

type Region struct {
	xCoord     int32
	zCoord     int32
	Locations  [1024]ChunkLoc
	Timestamps [1024]ChunkTime
	Chunks     []Chunk
}

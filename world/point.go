package world

import (
	"fmt"
	"log"
)

type Point struct {
	X int32
	Y int32
	Z int32
}

func MakePoint(X int32, Y int32, Z int32) Point {
	if Debug {
		log.Printf("MAKE POINT: %d, %d, %d", X, Y, Z)
	}
	return Point{X: X, Y: Y, Z: Z}
}

func (p Point) String() string {
	return fmt.Sprintf("Point{X: %d, Y: %d, Z: %d}", p.X, p.Y, p.Z)
}

func (pt Point) Index() int {
	return int(pt.X%16 + pt.Z%16*16 + pt.Y%16*16*16)
}

func (p Point) ChunkXZ() XZ {
	return XZ{X: floor(p.X, 16), Z: floor(p.Z, 16)}
}

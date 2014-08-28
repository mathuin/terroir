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
	remX := pt.X % 16
	if remX < 0 {
		remX += 16
	}
	remY := pt.Y % 16
	if remY < 0 {
		remY += 16
	}
	remZ := pt.Z % 16
	if remZ < 0 {
		remZ += 16
	}
	return int(remX%16 + remZ%16*16 + remY%16*16*16)
}

func (p Point) ChunkXZ() XZ {
	return XZ{X: floor(p.X, 16), Z: floor(p.Z, 16)}
}

type Points []Point

func (slice Points) Len() int {
	return len(slice)
}

func (slice Points) Less(i, j int) bool {
	if slice[i].X != slice[j].X {
		return slice[i].X < slice[j].X
	}
	if slice[i].Y != slice[j].Y {
		return slice[i].Y < slice[j].Y
	}
	return slice[i].Z < slice[j].Z
}

func (slice Points) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

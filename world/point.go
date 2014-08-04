package world

import (
	"fmt"
	"log"
	"math"
)

type Point struct {
	X int32
	Y int32
	Z int32
}

func NewPoint(X int32, Y int32, Z int32) *Point {
	if Debug {
		log.Printf("NEW POINT: %d, %d, %d", X, Y, Z)
	}
	return &Point{X: X, Y: Y, Z: Z}
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

// currently returns chunk *coordinates*
func (p Point) WhichChunk() (cx int32, cz int32) {
	cx = int32(math.Floor(float64(p.X) / 16.0))
	cz = int32(math.Floor(float64(p.Z) / 16.0))
	return
}

type Location struct {
	X float64
	Y float64
	Z float64
}

func NewLocation(X float64, Y float64, Z float64) *Location {
	if Debug {
		log.Printf("NEW LOCATION: %d, %d, %d", X, Y, Z)
	}
	return &Location{X: X, Y: Y, Z: Z}
}

func MakeLocation(X float64, Y float64, Z float64) Location {
	if Debug {
		log.Printf("MAKE LOCATION: %d, %d, %d", X, Y, Z)
	}
	return Location{X: X, Y: Y, Z: Z}
}

func (l Location) String() string {
	return fmt.Sprintf("Location{X: %d, Y: %d, Z: %d}", l.X, l.Y, l.Z)
}

func (l Location) ToPoint() Point {
	return MakePoint(int32(math.Floor(l.X)), int32(math.Floor(l.Y)), int32(math.Floor(l.Z)))
}

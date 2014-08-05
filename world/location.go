package world

import (
	"fmt"
	"log"
	"math"
)

type Location struct {
	X float64
	Y float64
	Z float64
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

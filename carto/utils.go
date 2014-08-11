package carto

import (
	"log"
	"os"
)

type CartoError error

func notnil(err error) bool {
	return (err != nil && err.Error() != "No Error")
}

type Float32Arr []float32

func (arr Float32Arr) min() (m float32) {
	m = arr[0]
	for _, v := range arr {
		if v < m {
			m = v
		}
	}
	return
}

func (arr Float32Arr) max() (m float32) {
	m = arr[0]
	for _, v := range arr {
		if v > m {
			m = v
		}
	}
	return
}

type Float64Arr []float64

func (arr Float64Arr) min() (m float64) {
	m = arr[0]
	for _, v := range arr {
		if v < m {
			m = v
		}
	}
	return
}

func (arr Float64Arr) max() (m float64) {
	m = arr[0]
	for _, v := range arr {
		if v > m {
			m = v
		}
	}
	return
}

func setIntValue(name string, old int, mymax int, mymin int) int {
	if old > mymax || old < mymin {
		log.Printf("warning: %s %d outside %d-%d range", name, old, mymin, mymax)
	}

	return int(min(max(old, mymin), mymax))
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func remove(name string) {
	mstat, _ := os.Stat(name)
	if mstat != nil {
		rerr := os.Remove(name)
		if rerr != nil {
			// JMT: I don't think I need to check this
			log.Printf(rerr.Error())
		}
	}
}
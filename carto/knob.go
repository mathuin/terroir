package carto

import "fmt"

type Knob interface {
	setValue()
}

type IntKnob struct {
	name  string
	value int
}

func (ik *IntKnob) setValue(mymax int, mymin int) (msg string) {
	if ik.value > mymax || ik.value < mymin {
		msg = fmt.Sprintf("warning: %s %d outside %d-%d range", ik.name, ik.value, mymin, mymax)
	}

	ik.value = min(max(ik.value, mymin), mymax)

	return
}

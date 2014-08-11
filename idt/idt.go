package idt

import (
	"fmt"

	"code.google.com/p/biogo.store/kdtree"
)

type IDT struct {
	coordmap map[string]int
	values   []int
	tree     kdtree.Tree
}

func NewIDT(coords [][2]float64, values []int) (*IDT, error) {
	if len(coords) != len(values) {
		return nil, fmt.Errorf("coords and values must be of the same length")
	}
	kdpts := kdtree.Points{}
	coordmap := make(map[string]int)
	for i, v := range coords {
		kdpt := kdtree.Point{v[0], v[1]}
		kdpts = append(kdpts, kdpt)
		coordmap[fmt.Sprint(kdpt)] = i
	}
	tree := kdtree.New(kdpts, false)
	return &IDT{coordmap: coordmap, values: values, tree: *tree}, nil
}

func (idt IDT) Call(base [][2]int, nnear int, majority bool) (outarr []int, err error) {
	basepts := kdtree.Points{}
	for _, v := range base {
		basepts = append(basepts, kdtree.Point{float64(v[0]), float64(v[1])})
	}

	outarr = make([]int, len(basepts))

	for i, q := range basepts {
		nk := kdtree.NewNKeeper(nnear)
		idt.tree.NearestSet(nk, q)
		distance := []float64{}
		index := []int{}
		for _, val := range nk.Heap {
			if val.Comparable != nil {
				distance = append(distance, val.Dist)
				index = append(index, idt.coordmap[fmt.Sprint(val.Comparable)])
			}
		}
		var wz int
		if nnear == 1 || distance[0] < 1e-10 {
			wz = idt.values[index[0]]
		} else {
			w := make([]float64, len(distance))
			var sumw float64
			for i, d := range distance {
				w[i] = 1 / d
				sumw = sumw + w[i]
			}
			for i := range distance {
				w[i] = w[i] / sumw
			}
			values := make([]int, len(index))
			for i, v := range index {
				values[i] = idt.values[v]
			}
			if majority {
				majordict := make(map[int]float64)
				for i, v := range values {
					majordict[v] += w[i]
				}
				var max float64
				for k, v := range majordict {
					if v > max {
						wz = k
					}
				}
			} else {
				wz, err = dot(w, values)
				if err != nil {
					panic(err)
				}
			}
		}
		outarr[i] = wz
	}
	return outarr, nil
}

func dot(x []float64, y []int) (int, error) {
	var retval int
	if len(x) != len(y) {
		return 0, fmt.Errorf("lengths do not match")
	}
	for i := range x {
		retval += int(x[i] * float64(y[i]))
	}
	return retval, nil
}

package idt

import (
	"fmt"
	"log"
	"runtime"
	"sync"

	"code.google.com/p/biogo.store/kdtree"
)

var Debug = false

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

type OrdPt struct {
	index int
	pt    kdtree.Point
}

type OrdVal struct {
	index int
	val   int16
}

func (idt IDT) genOrdPts(base [][2]int, in chan OrdPt) {
	for i, v := range base {
		in <- OrdPt{index: i, pt: kdtree.Point{float64(v[0]), float64(v[1])}}
	}
	close(in)
}

func (idt IDT) Call(base [][2]int, nnear int, majority bool) (outarr []int16, err error) {
	outarr = make([]int16, len(base))

	in := make(chan OrdPt)
	out := make(chan OrdVal)

	var wg sync.WaitGroup

	nCPU := runtime.NumCPU()
	numWorkers := nCPU * nCPU
	if Debug {
		log.Printf("debug mode -- only one worker!")
		numWorkers = 1
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			if Debug {
				log.Printf("Starting reducer #%d!", i)
			}
			defer wg.Done()
			idt.Reduce(in, out, nnear, majority, i)
		}(i)
	}
	go func() { wg.Wait(); close(out) }()
	go idt.genOrdPts(base, in)

	for ordval := range out {
		outarr[ordval.index] = ordval.val
	}

	return outarr, nil
}

func (idt IDT) Reduce(in chan OrdPt, out chan OrdVal, nnear int, majority bool, i int) {
	for ordpt := range in {
		ordval := new(OrdVal)
		ordval.index = ordpt.index
		q := ordpt.pt
		if Debug {
			log.Printf("pt: %+#v", q)
		}

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
		if Debug {
			log.Print("distance: ", distance)
			log.Print("index: ", index)
		}
		var wz int
		if nnear == 1 || distance[0] < 1e-10 {
			wz = idt.values[index[0]]
		} else {
			w := make([]float64, len(distance))
			var sumw float64
			for i, d := range distance {
				w[i] = 1.0 / d
				sumw += w[i]
			}
			if Debug {
				log.Print("w: ", w)
				log.Print("sumw: ", sumw)
			}
			for i := range distance {
				w[i] /= sumw
			}
			if Debug {
				log.Print("new w: ", w)
			}
			values := make([]int, len(index))
			for i, v := range index {
				values[i] = idt.values[v]
			}
			if Debug {
				log.Print("values: ", values)
			}
			if majority {
				majordict := make(map[int]float64)
				for i, v := range values {
					majordict[v] += w[i]
				}
				if Debug {
					log.Print("majordict: ")
					for k, v := range majordict {
						log.Printf(" %d: %f", k, v)
					}
				}
				major := float64(0)
				for k, v := range majordict {
					if v > major {
						if Debug {
							log.Printf("new max value %f with ind %d", v, k)
						}
						major = v
						wz = k
					}
				}
			} else {
				var err error
				wz, err = dot(w, values)
				if err != nil {
					panic(err)
				}
			}
		}
		if Debug {
			log.Print("wz: ", wz)
		}
		ordval.val = int16(wz)
		out <- *ordval
	}
}

func dot(x []float64, y []int) (int, error) {
	var retval float64
	if len(x) != len(y) {
		return 0, fmt.Errorf("lengths do not match")
	}
	for i := range x {
		retval += x[i] * float64(y[i])
	}
	return int(retval + 0.5), nil
}

package world

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func Test_regionWriteRead(t *testing.T) {
	saveDir := "."
	worldName := "TerroirTest"
	newWorldName := "TerroirTwo"
	rx := 0
	rz := 0
	expected_count := 1024

	w := MakeWorld(worldName)

	// open the file
	rpath := path.Join(saveDir, worldName, "region")
	rname := path.Join(rpath, fmt.Sprintf("r.%d.%d.mca", rx, rz))
	r, rerr := os.Open(rname)
	if rerr != nil {
		panic(rerr)
	}
	defer r.Close()

	// readregion here
	n, rerr := w.ReadRegion(r, int32(rx), int32(rz))
	if rerr != nil {
		panic(rerr)
	}

	if n != expected_count {
		t.Errorf("expected %d sectors, got %d", expected_count, n)
	}

	nw := MakeWorld(newWorldName)
	nw.SetSpawn(w.Spawn)
	nw.SetRandomSeed(w.RandomSeed)

	for k, v := range w.ChunkMap {
		nw.ChunkMap[k] = v
	}

	for k, v := range w.RegionMap {
		nw.RegionMap[k] = v
	}

	// set two points to obsidian
	obsidian := 49
	pt := MakePoint(0, 60, 0)
	pt2 := MakePoint(1, 60, 0)
	nw.SetBlock(pt, obsidian)
	nw.SetBlock(pt2, obsidian)

	newSaveDir, nerr := ioutil.TempDir("", "")
	if nerr != nil {
		panic(nerr)
	}
	defer os.RemoveAll(newSaveDir)
	newWorldDir := path.Join(newSaveDir, w.Name)
	if _, err := os.Stat(newWorldDir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(newWorldDir, 0775)
		} else {
			t.Fail()
		}
	}
	newRegionDir := path.Join(newWorldDir, "region")
	if _, err := os.Stat(newRegionDir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(newRegionDir, 0775)
		} else {
			t.Fail()
		}
	}

	// writeregion here (actually write world!!)
	nwXZ := XZ{X: int32(rx), Z: int32(rz)}
	if err := nw.WriteRegion(newRegionDir, nwXZ); err != nil {
		t.Fail()
	}

	// read it back
	sw := MakeWorld(newWorldName)

	// readregion here
	swrname := path.Join(newRegionDir, fmt.Sprintf("r.%d.%d.mca", rx, rz))
	swr, swrerr := os.Open(swrname)
	if swrerr != nil {
		panic(swrerr)
	}
	defer swr.Close()

	swn, rerr := sw.ReadRegion(swr, int32(rx), int32(rz))
	if rerr != nil {
		panic(rerr)
	}
	if swn != expected_count {
		t.Errorf("expected %d sectors, got %d", expected_count, swn)
	}

	// check value of some particular block
	nbval := sw.Block(pt)
	nb2val := sw.Block(pt2)

	if nbval != obsidian {
		t.Errorf("nbval %v is not equal to obsidian %v", nbval, obsidian)
	}
	if nb2val != obsidian {
		t.Errorf("nb2val %v is not equal to obsidian %v", nb2val, obsidian)
	}
}

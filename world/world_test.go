// Tests for Minecraft world package.

package world

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

func Test_worldWriteRead(t *testing.T) {
	saveDir := "."
	worldName := "TerroirTest"
	newWorldName := "TerroirTwo"
	rx := 0
	rz := 0
	expected_count := 1024

	td, nerr := ioutil.TempDir("", "")
	if nerr != nil {
		panic(nerr)
	}
	defer os.RemoveAll(td)

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

	// write a new copy of the world
	nw := MakeWorld(newWorldName)
	nw.SetSpawn(w.Spawn)
	nw.SetRandomSeed(w.RandomSeed)

	for k, v := range w.ChunkMap {
		nw.ChunkMap[k] = v
	}

	for k, v := range w.RegionMap {
		nw.RegionMap[k] = v
	}

	nw.SetSaveDir(td)
	nw.Write()

	sw, swerr := ReadWorld(td, newWorldName)
	if swerr != nil {
		t.Fail()
	}
	_ = sw
}

func Test_FullReadWriteRead(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	saveDir := "."
	worldName := "TerroirTest"
	newWorldName := "TerroirTwo"

	w, err := ReadWorld(saveDir, worldName)
	if err != nil {
		t.Fail()
	}

	nw := MakeWorld(newWorldName)

	// copy all but the name
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

	nw.SetSaveDir(newSaveDir)
	if err := nw.Write(); err != nil {
		t.Fail()
	}

	// read it back
	sw, serr := ReadWorld(newSaveDir, newWorldName)
	if serr != nil {
		t.Fail()
	}

	// check value of some particular block
	nbval, err := sw.Block(pt)
	if err != nil {
		log.Panic(err)
		t.Fail()
	}
	nb2val, err := sw.Block(pt2)
	if err != nil {
		log.Panic(err)
		t.Fail()
	}

	if nbval != obsidian {
		t.Errorf("nbval %v is not equal to obsidian %v", nbval, obsidian)
	}
	if nb2val != obsidian {
		t.Errorf("nb2val %v is not equal to obsidian %v", nb2val, obsidian)
	}
}

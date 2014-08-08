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
	rXZ := XZ{X: int32(rx), Z: int32(rz)}
	n, rerr := w.loadAllChunksFromRegion(rXZ)
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

	sw, swerr := ReadWorld(td, newWorldName, false)
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

	w, err := ReadWorld(saveDir, worldName, true)
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
	obsidian, err := BlockNamed("Obsidian")
	if err != nil {
		t.Fail()
	}
	pt := MakePoint(0, 60, 0)
	pt2 := MakePoint(1, 60, 0)
	nw.SetBlock(pt, obsidian)
	nw.SetBlock(pt2, obsidian)

	newSaveDir, nerr := ioutil.TempDir("", "")
	if nerr != nil {
		t.Fail()
	}
	defer os.RemoveAll(newSaveDir)

	nw.SetSaveDir(newSaveDir)
	if err := nw.Write(); err != nil {
		t.Fail()
	}

	// read it back
	sw, serr := ReadWorld(newSaveDir, newWorldName, true)
	if serr != nil {
		t.Fail()
	}

	// check value of some particular block
	nbval, err := sw.Block(pt)
	if err != nil {
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

func Test_tinyWorld(t *testing.T) {
	worldName := "PointTest"
	td, nerr := ioutil.TempDir("", "")
	if nerr != nil {
		t.Fail()
	}
	defer os.RemoveAll(td)

	// points and values for the world
	var points = []struct {
		p Point
		s string
	}{
		{Point{X: 7, Y: 85, Z: 7}, "Obsidian"},
		{Point{X: 8, Y: 85, Z: 7}, "Obsidian"},
		{Point{X: 7, Y: 86, Z: 7}, "Air"},
		{Point{X: 8, Y: 86, Z: 7}, "Air"},
		{Point{X: 7, Y: 87, Z: 7}, "Air"},
		{Point{X: 8, Y: 87, Z: 7}, "Air"},
	}

	w := MakeWorld(worldName)
	w.SetSaveDir(td)
	spawnPoint := points[0].p
	w.SetSpawn(spawnPoint)
	w.SetRandomSeed(0)

	// set the points
	for _, pv := range points {
		b, err := BlockNamed(pv.s)
		if err != nil {
			t.Fail()
		}
		w.SetBlock(pv.p, b)
	}

	// now write the level
	err := w.Write()
	if err != nil {
		t.Fail()
	}

	// now read the level again
	nw, err := ReadWorld(td, worldName, false)
	if err != nil {
		t.Fail()
	}

	// are those blocks still set?
	for i, pv := range points {
		b, err := nw.Block(pv.p)
		if err != nil {
			t.Fail()
		}
		bn, err := b.BlockName()
		if err != nil {
			t.Fail()
		}
		if bn != pv.s {
			t.Errorf("Point #%d: (%v) expected %s, got %s", i, pv.p, pv.s, bn)
		}
	}
}

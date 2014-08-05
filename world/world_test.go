// Tests for Minecraft world package.

package world

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/mathuin/terroir/nbt"
)

func Test_WriteWorld(t *testing.T) {
	pname := "ExampleLevel"
	px := int32(-208)
	py := int32(64)
	pz := int32(220)
	ppt := MakePoint(px, py, pz)
	prseed := int64(2603059821051629081)

	w := MakeWorld(pname)
	w.SetRandomSeed(prseed)
	w.SetSpawn(ppt)

	td, err := ioutil.TempDir("", "")
	defer os.RemoveAll(td)
	w.Write(td)

	// does directory exist?
	worldDir := path.Join(td, pname)
	checkDir, err := os.Stat(worldDir)
	if err != nil {
		t.Fail()
	}
	if !checkDir.IsDir() {
		t.Fail()
	}

	// does it contain a level.dat file?
	levelFile := path.Join(worldDir, "level.dat")
	checkFile, err := os.Stat(levelFile)
	if err != nil {
		t.Fail()
	}
	if !checkFile.Mode().IsRegular() {
		t.Fail()
	}

	// is that level.dat file a compressed NBT file?
	levelTag, err := nbt.ReadCompressedFile(levelFile)
	if err != nil {
		t.Fail()
	}

	// That's it!  contents of file are tested elsewhere...
	_ = levelTag
}

func Test_Two(t *testing.T) {
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

	if err := nw.Write(newSaveDir); err != nil {
		t.Fail()
	}

	// read it back
	sw, serr := ReadWorld(newSaveDir, newWorldName)
	if serr != nil {
		t.Fail()
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

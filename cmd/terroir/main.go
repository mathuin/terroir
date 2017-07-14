package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/mathuin/terroir/pkg/world"
)

func main() {

	log.Print("Begin!")

	saveDir := "/worlds"
	worldName := "TerroirTest"
	newWorldName := "TerroirTwo"

	w, err := world.ReadWorld(saveDir, worldName, true)
	if err != nil {
		log.Fatal(err)
	}

	nw := world.MakeWorld(newWorldName)

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
	obsidian, err := world.BlockNamed("Obsidian")
	if err != nil {
		log.Fatal(err)
	}
	pt := world.MakePoint(0, 60, 0)
	pt2 := world.MakePoint(1, 60, 0)
	nw.SetBlock(pt, *obsidian)
	nw.SetBlock(pt2, *obsidian)

	newSaveDir, nerr := ioutil.TempDir("", "")
	if nerr != nil {
		log.Panic(nerr)
	}
	defer os.RemoveAll(newSaveDir)

	nw.SetSaveDir(newSaveDir)
	if err := nw.Write(); err != nil {
		log.Panic(err)
	}

	// read it back
	sw, serr := world.ReadWorld(newSaveDir, newWorldName, true)
	if serr != nil {
		log.Panic(serr)
	}

	// check value of some particular block
	nbval, err := sw.Block(pt)
	if err != nil {
		log.Panic(err)
	}
	nb2val, err := sw.Block(pt2)
	if err != nil {
		log.Panic(err)
	}

	if *nbval != *obsidian {
		log.Fatalf("nbval %v is not equal to obsidian %v", nbval, obsidian)
	}
	if *nb2val != *obsidian {
		log.Fatalf("nb2val %v is not equal to obsidian %v", nb2val, obsidian)
	}

	log.Print("Success!")
}

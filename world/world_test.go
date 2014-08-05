// Tests for Minecraft world package.

package world

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/mathuin/terroir/nbt"
)

// test writeworld
// - check for:
//   - directory created under temp with correct name
//   - existence of level.dat in said directory

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

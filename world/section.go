package world

import (
	"fmt"
	"log"

	"github.com/mathuin/terroir/nbt"
)

type Section struct {
	Blocks     []byte
	Add        []byte
	Data       []byte
	BlockLight []byte
	SkyLight   []byte
}

func MakeSection() Section {
	if Debug {
		log.Printf("MAKE SECTION")
	}
	Blocks := make([]byte, 4096)
	Add := make([]byte, 2048)
	Data := make([]byte, 2048)
	BlockLight := make([]byte, 2048)
	SkyLight := make([]byte, 2048)
	return Section{Blocks: Blocks, Add: Add, Data: Data, BlockLight: BlockLight, SkyLight: SkyLight}
}

func (s Section) String() string {
	return fmt.Sprintf("Section{}")
}

func (w *World) Block(pt Point) int {
	s := w.Section(pt)
	base := int(s.Blocks[pt.Index()])
	add := int(Nibble(s.Add, pt.Index()))
	return base + add*256
}

func (w *World) SetBlock(pt Point, b int) {
	base := byte(b % 256)
	add := byte(b / 256)
	s := w.Section(pt)
	i := pt.Index()
	s.Blocks[i] = byte(base)
	WriteNibble(s.Add, i, add)
}

func (w *World) Data(pt Point) byte {
	return Nibble(w.Section(pt).Data, pt.Index())
}

func (w *World) SetData(pt Point, b byte) {
	WriteNibble(w.Section(pt).Data, pt.Index(), b)
}

func (w *World) BlockLight(pt Point) byte {
	return Nibble(w.Section(pt).BlockLight, pt.Index())
}

func (w *World) SetBlockLight(pt Point, b byte) {
	WriteNibble(w.Section(pt).BlockLight, pt.Index(), b)
}

func (w *World) SkyLight(pt Point) byte {
	return Nibble(w.Section(pt).SkyLight, pt.Index())
}

func (w *World) SetSkyLight(pt Point, b byte) {
	WriteNibble(w.Section(pt).SkyLight, pt.Index(), b)
}

func (w World) Section(pt Point) Section {
	return w.ChunkMap[pt.ChunkXZ()].Sections[int(floor(pt.Y, 16))]
}

func (s Section) write(y int) []nbt.Tag {
	sElems := []nbt.CompoundElem{
		{"Y", nbt.TAG_Byte, byte(y)},
		{"Blocks", nbt.TAG_Byte_Array, s.Blocks},
		{"Add", nbt.TAG_Byte_Array, s.Add},
		{"Data", nbt.TAG_Byte_Array, s.Data},
		{"BlockLight", nbt.TAG_Byte_Array, s.BlockLight},
		{"SkyLight", nbt.TAG_Byte_Array, s.SkyLight},
	}

	sTagPayload := nbt.MakeCompoundPayload(sElems)

	return sTagPayload
}

func ReadSection(tarr []nbt.Tag) Section {
	s := MakeSection()

	for _, tval := range tarr {
		switch tval.Name {
		case "Y":
		// Y tags are checked on the chunk level.
		case "Blocks":
			s.Blocks = tval.Payload.([]byte)
		case "Add":
			s.Add = tval.Payload.([]byte)
		case "Data":
			s.Data = tval.Payload.([]byte)
		case "BlockLight":
			s.BlockLight = tval.Payload.([]byte)
		case "SkyLight":
			s.SkyLight = tval.Payload.([]byte)
		default:
			log.Fatalf("tag name %s not required for section", tval.Name)
		}
	}

	return s
}

package world

import (
	"fmt"
	"log"

	"github.com/mathuin/terroir/pkg/nbt"
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

func (s Section) write(y int) []nbt.Tag {
	sElems := []nbt.CompoundElem{
		{Key: "Y", Tag: nbt.TAG_Byte, Value: byte(y)},
		{Key: "Blocks", Tag: nbt.TAG_Byte_Array, Value: s.Blocks},
		{Key: "Add", Tag: nbt.TAG_Byte_Array, Value: s.Add},
		{Key: "Data", Tag: nbt.TAG_Byte_Array, Value: s.Data},
		{Key: "BlockLight", Tag: nbt.TAG_Byte_Array, Value: s.BlockLight},
		{Key: "SkyLight", Tag: nbt.TAG_Byte_Array, Value: s.SkyLight},
	}

	sTagPayload := nbt.MakeCompoundPayload(sElems)

	return sTagPayload
}

func ReadSection(tarr []nbt.Tag) (*Section, error) {
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
			return nil, fmt.Errorf("tag name %s not required for section", tval.Name)
		}
	}

	return &s, nil
}

package world

import (
	"fmt"
	"log"

	"github.com/mathuin/terroir/pkg/nbt"
)

// PLEASE NOTE
// This is for future expansion.
// For now, generated worlds will have none of these features.

// Entity is a TAG_Compound
// Entities is a TAG_List of TAG_Compound
// (unless none exist, in which case it is a TAG_List of <nil>)

type Entity struct {
	tags []nbt.Tag
}

func MakeEntity() Entity {
	if Debug {
		log.Printf("MAKE ENTITY")
	}
	return Entity{}
}

func (e Entity) String() string {
	return fmt.Sprintf("Entity{tags: %v}", e.tags)
}

func ReadEntity(tags []nbt.Tag) Entity {
	e := MakeEntity()
	e.tags = tags
	return e
}

func (e Entity) write() []nbt.Tag {
	return e.tags
}

// TileEntity is a TAG_Compound
// TileEntities is a TAG_List of TAG_Compound
// (unless none exist, in which case it is a TAG_List of <nil>)

type TileEntity struct {
	tags []nbt.Tag
}

func MakeTileEntity() TileEntity {
	if Debug {
		log.Printf("MAKE TILEENTITY")
	}
	return TileEntity{}
}

func (e TileEntity) String() string {
	return fmt.Sprintf("TileEntity{tags: %v}", e.tags)
}

func ReadTileEntity(tags []nbt.Tag) TileEntity {
	e := MakeTileEntity()
	e.tags = tags
	return e
}

func (te TileEntity) write() []nbt.Tag {
	return te.tags
}

// TileTick is a TAG_Compound
// TileTicks TAG_List of TAG_Compound
// (unless none exist, in which case no tag is sent)
type TileTick struct {
	id    int32
	time  int32
	order int32
	point Point
}

func MakeTileTick() TileTick {
	if Debug {
		log.Printf("MAKE TILETICK")
	}
	return TileTick{}
}

func (t TileTick) String() string {
	return fmt.Sprintf("TileTick{ID: %d, Time: %d, Order: %d, Point: %v}", t.id, t.time, t.order, t.point)
}

func (tt TileTick) write() nbt.Tag {
	ttElems := []nbt.CompoundElem{
		{Key: "i", Tag: nbt.TAG_Int, Value: tt.id},
		{Key: "t", Tag: nbt.TAG_Int, Value: tt.time},
		{Key: "p", Tag: nbt.TAG_Int, Value: tt.order},
		{Key: "x", Tag: nbt.TAG_Int, Value: tt.point.X},
		{Key: "y", Tag: nbt.TAG_Int, Value: tt.point.Y},
		{Key: "z", Tag: nbt.TAG_Int, Value: tt.point.Z},
	}

	ttTag := nbt.MakeCompound("", ttElems)

	return ttTag
}

func ReadTileTick(tarr []nbt.Tag) TileTick {
	tt := MakeTileTick()
	for _, tval := range tarr {
		switch tval.Name {
		case "i":
			tt.id = tval.Payload.(int32)
		case "t":
			tt.time = tval.Payload.(int32)
		case "p":
			tt.order = tval.Payload.(int32)
		case "x":
			tt.point.X = tval.Payload.(int32)
		case "y":
			tt.point.Y = tval.Payload.(int32)
		case "z":
			tt.point.Z = tval.Payload.(int32)
		default:
			log.Fatalf("tag name %s not required for tiletick", tval.Name)
		}
	}
	// JMT: no check for incomplete tileticks
	return tt
}

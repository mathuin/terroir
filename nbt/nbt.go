package nbt

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	TAG_End = iota
	TAG_Byte
	TAG_Short
	TAG_Int
	TAG_Long
	TAG_Float
	TAG_Double
	TAG_Byte_Array
	TAG_String
	TAG_List
	TAG_Compound
	TAG_Int_Array
)

func tagToString(t byte) string {
	switch t {
	case TAG_End:
		return "TAG_End"
	case TAG_Byte:
		return "TAG_Byte"
	case TAG_Short:
		return "TAG_Short"
	case TAG_Int:
		return "TAG_Int"
	case TAG_Long:
		return "TAG_Long"
	case TAG_Float:
		return "TAG_Float"
	case TAG_Double:
		return "TAG_Double"
	case TAG_Byte_Array:
		return "TAG_Byte_Array"
	case TAG_String:
		return "TAG_String"
	case TAG_List:
		return "TAG_List"
	case TAG_Compound:
		return "TAG_Compound"
	case TAG_Int_Array:
		return "TAG_Int_Array"
	}
	return "Unknown tag!"
}

type Tag struct {
	Type    byte
	Name    string
	Payload interface{}
}

func ReadTag(r io.Reader) (t Tag, err error) {
	// read a byte
	ttype := make([]byte, 1)
	_, err = io.ReadFull(r, ttype)
	if err != nil {
		return
	}

	// TAG_End isn't really a tag
	if ttype[0] == byte(TAG_End) {
		return
	}

	// Real tags need types
	t.Type = ttype[0]

	// Now about that name
	lenbytes := make([]byte, 2)
	_, err = io.ReadFull(r, lenbytes)
	if err != nil {
		return
	}

	strlen := int(lenbytes[0])*256 + int(lenbytes[1])

	strbytes := make([]byte, strlen)

	_, err = io.ReadFull(r, strbytes)
	if err != nil {
		return
	}

	t.Name = string(strbytes)

	// Payload-specific code goes here
	switch t.Type {
	case TAG_Byte:
		payload := make([]byte, 1)
		_, err = io.ReadFull(r, payload)
		if err != nil {
			return
		}
		t.Payload = payload
	case TAG_Int:
		var payload int32
		err = binary.Read(r, binary.BigEndian, &payload)
		if err != nil {
			return
		}
		t.Payload = payload
	case TAG_Compound:
		payload := []Tag{}
		var newtag, emptytag Tag
		for newtag, err = ReadTag(r); newtag != emptytag; newtag, err = ReadTag(r) {
			payload = append(payload, newtag)
		}
		t.Payload = payload
	default:
		err = fmt.Errorf("unknown tag")
	}
	return
}

// other way

func WriteTag(w io.Writer, t Tag) (err error) {
	w.Write([]byte{t.Type})

	if t.Type == TAG_End {
		return
	}

	binary.Write(w, binary.BigEndian, int16(len(t.Name)))
	w.Write([]byte(t.Name))

	switch t.Type {
	case TAG_Byte:
		w.Write(t.Payload.([]byte))
	case TAG_Int:
		binary.Write(w, binary.BigEndian, t.Payload.(int32))
	case TAG_Compound:
		tags := append(t.Payload.([]Tag), Tag{Type: TAG_End})
		for _, tag := range tags {
			WriteTag(w, tag)
		}
	default:
		err = fmt.Errorf("unknown tag")
	}
	return
}

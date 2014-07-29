package nbt

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	TAG_End        byte = 0
	TAG_Byte       byte = 1
	TAG_Short      byte = 2
	TAG_Int        byte = 3
	TAG_Long       byte = 4
	TAG_Float      byte = 5
	TAG_Double     byte = 6
	TAG_Byte_Array byte = 7
	TAG_String     byte = 8
	TAG_List       byte = 9
	TAG_Compound   byte = 10
	TAG_Int_Array  byte = 11
)

type Tag struct {
	Type    byte
	Name    string
	Payload interface{}
}

type PayloadReader func(io.Reader) interface{}
type PayloadWriter func(io.Writer, interface{})

var Tags = map[byte]struct {
	String  string
	PReader PayloadReader
	PWriter PayloadWriter
}{
	TAG_End: {"TAG_End",
		nil,
		nil,
	},
	TAG_Byte: {"TAG_Byte",
		func(r io.Reader) interface{} {
			payload := make([]byte, 1)
			if _, err := io.ReadFull(r, payload); err != nil {
				panic(err)
			}
			return payload
		},
		func(w io.Writer, i interface{}) {
			w.Write(i.([]byte))
		},
	},
	TAG_Short: {"TAG_Short",
		func(r io.Reader) interface{} {
			var payload int16
			if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
				panic(err)
			}
			return payload
		},
		func(w io.Writer, i interface{}) {
			binary.Write(w, binary.BigEndian, i.(int16))
		},
	},
	TAG_Int: {"TAG_Int",
		func(r io.Reader) interface{} {
			var payload int32
			if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
				panic(err)
			}
			return payload
		},
		func(w io.Writer, i interface{}) {
			binary.Write(w, binary.BigEndian, i.(int32))
		},
	},
	TAG_Long: {"TAG_Long",
		func(r io.Reader) interface{} {
			var payload int64
			if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
				panic(err)
			}
			return payload
		},
		func(w io.Writer, i interface{}) {
			binary.Write(w, binary.BigEndian, i.(int64))
		},
	},
	TAG_Float: {"TAG_Float",
		func(r io.Reader) interface{} {
			var payload float32
			if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
				panic(err)
			}
			return payload
		},
		func(w io.Writer, i interface{}) {
			binary.Write(w, binary.BigEndian, i.(float32))
		},
	},
	TAG_Double: {"TAG_Double",
		func(r io.Reader) interface{} {
			var payload float64
			if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
				panic(err)
			}
			return payload
		},
		func(w io.Writer, i interface{}) {
			binary.Write(w, binary.BigEndian, i.(float64))
		},
	},
	TAG_Byte_Array: {"TAG_Byte_Array",
		func(r io.Reader) interface{} {
			var strlen int32
			if err := binary.Read(r, binary.BigEndian, &strlen); err != nil {
				panic(err)
			}

			strbytes := make([]byte, strlen)

			if _, err := io.ReadFull(r, strbytes); err != nil {
				panic(err)
			}

			return strbytes
		},
		func(w io.Writer, i interface{}) {
			binary.Write(w, binary.BigEndian, int32(len(i.([]byte))))
			for _, value := range i.([]byte) {
				binary.Write(w, binary.BigEndian, value)
			}
		},
	},
	TAG_String: {"TAG_String",
		func(r io.Reader) interface{} {
			var strlen int16
			if err := binary.Read(r, binary.BigEndian, &strlen); err != nil {
				panic(err)
			}

			strbytes := make([]byte, strlen)

			if _, err := io.ReadFull(r, strbytes); err != nil {
				panic(err)
			}

			return string(strbytes)
		},
		func(w io.Writer, i interface{}) {
			binary.Write(w, binary.BigEndian, int16(len(i.(string))))
			w.Write([]byte(i.(string)))
		},
	},
	TAG_List: {"TAG_List",
		func(r io.Reader) interface{} {
			var i interface{}
			return i
		},
		func(w io.Writer, i interface{}) {
			return
		},
	},
	TAG_Compound: {"TAG_Compound",
		// JMT: figure out how to break loop
		nil,
		nil,
	},
	TAG_Int_Array: {"TAG_Int_Array",
		func(r io.Reader) interface{} {
			var strlen int32
			if err := binary.Read(r, binary.BigEndian, &strlen); err != nil {
				panic(err)
			}

			ints := make([]int32, strlen)
			for key := range ints {
				if err := binary.Read(r, binary.BigEndian, &ints[key]); err != nil {
					panic(err)
				}
			}

			return ints
		},
		func(w io.Writer, i interface{}) {
			binary.Write(w, binary.BigEndian, int32(len(i.([]int32))))
			for _, value := range i.([]int32) {
				binary.Write(w, binary.BigEndian, value)
			}
		},
	},
}

// JMT: The compound reader and writer cause an initialization loop
// when added to the Tags variable.

func readCompound(r io.Reader) interface{} {
	payload := []Tag{}
	var emptytag Tag
	for newtag, err := ReadTag(r); newtag != emptytag; newtag, err = ReadTag(r) {
		if err != nil {
			panic(err)
		}

		payload = append(payload, newtag)
	}
	return payload
}

func writeCompound(w io.Writer, i interface{}) {
	tags := append(i.([]Tag), Tag{Type: TAG_End})
	for _, tag := range tags {
		WriteTag(w, tag)
	}
}

func ReadTag(r io.Reader) (t Tag, err error) {
	// read a byte
	ttype := Tags[TAG_Byte].PReader(r).([]byte)[0]

	// TAG_End isn't really a tag
	if ttype == TAG_End {
		return
	}

	// Real tags need types
	t.Type = ttype

	// Now about that name
	t.Name = Tags[TAG_String].PReader(r).(string)

	// Putting this in the widget causes an initialization loop issue
	// (Tags refers to readCompound refers to ReadTag refers to Tags)
	if t.Type == TAG_Compound {
		t.Payload = readCompound(r)
	} else if val, ok := Tags[t.Type]; ok {
		t.Payload = val.PReader(r)
	} else {
		err = fmt.Errorf("unknown tag")
	}
	return
}

func WriteTag(w io.Writer, t Tag) (err error) {
	// JMT: this []byte{} bit feels wrong
	Tags[TAG_Byte].PWriter(w, []byte{t.Type})

	if t.Type == TAG_End {
		return
	}

	Tags[TAG_String].PWriter(w, t.Name)

	// JMT: initialization loop issue here too
	if t.Type == TAG_Compound {
		writeCompound(w, t.Payload)
	} else if val, ok := Tags[t.Type]; ok {
		val.PWriter(w, t.Payload)
	} else {
		err = fmt.Errorf("unknown tag")
	}
	return
}

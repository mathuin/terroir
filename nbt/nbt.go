package nbt

import (
	"bytes"
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

type PayloadReader func(io.Reader) (interface{}, error)
type PayloadWriter func(io.Writer, interface{}) error

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
		func(r io.Reader) (interface{}, error) {
			payload := make([]byte, 1)
			if _, err := io.ReadFull(r, payload); err != nil {
				return nil, err
			}
			return payload[0], nil
		},
		func(w io.Writer, i interface{}) error {
			b := make([]byte, 1)
			b[0] = i.(byte)
			_, err := w.Write(b)
			return err
		},
	},
	TAG_Short: {"TAG_Short",
		func(r io.Reader) (interface{}, error) {
			var payload int16
			if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
				return nil, err
			}
			return payload, nil
		},
		func(w io.Writer, i interface{}) error {
			return binary.Write(w, binary.BigEndian, i.(int16))
		},
	},
	TAG_Int: {"TAG_Int",
		func(r io.Reader) (interface{}, error) {
			var payload int32
			if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
				return nil, err
			}
			return payload, nil
		},
		func(w io.Writer, i interface{}) error {
			return binary.Write(w, binary.BigEndian, i.(int32))
		},
	},
	TAG_Long: {"TAG_Long",
		func(r io.Reader) (interface{}, error) {
			var payload int64
			if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
				return nil, err
			}
			return payload, nil
		},
		func(w io.Writer, i interface{}) error {
			return binary.Write(w, binary.BigEndian, i.(int64))
		},
	},
	TAG_Float: {"TAG_Float",
		func(r io.Reader) (interface{}, error) {
			var payload float32
			if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
				return nil, err
			}
			return payload, nil
		},
		func(w io.Writer, i interface{}) error {
			return binary.Write(w, binary.BigEndian, i.(float32))
		},
	},
	TAG_Double: {"TAG_Double",
		func(r io.Reader) (interface{}, error) {
			var payload float64
			if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
				return nil, err
			}
			return payload, nil
		},
		func(w io.Writer, i interface{}) error {
			return binary.Write(w, binary.BigEndian, i.(float64))
		},
	},
	TAG_Byte_Array: {"TAG_Byte_Array",
		func(r io.Reader) (interface{}, error) {
			var strlen int32
			if err := binary.Read(r, binary.BigEndian, &strlen); err != nil {
				return nil, err
			}

			strbytes := make([]byte, strlen)

			if _, err := io.ReadFull(r, strbytes); err != nil {
				return nil, err
			}

			return strbytes, nil
		},
		func(w io.Writer, i interface{}) error {
			if err := binary.Write(w, binary.BigEndian, int32(len(i.([]byte)))); err != nil {
				return err
			}
			for _, value := range i.([]byte) {
				if err := binary.Write(w, binary.BigEndian, value); err != nil {
					return err
				}
			}
			return nil
		},
	},
	TAG_String: {"TAG_String",
		func(r io.Reader) (interface{}, error) {
			var strlen int16
			if err := binary.Read(r, binary.BigEndian, &strlen); err != nil {
				return nil, err
			}

			strbytes := make([]byte, strlen)

			if _, err := io.ReadFull(r, strbytes); err != nil {
				return nil, err
			}

			return string(strbytes), nil
		},
		func(w io.Writer, i interface{}) error {
			if err := binary.Write(w, binary.BigEndian, int16(len(i.(string)))); err != nil {
				return err
			}
			_, err := w.Write([]byte(i.(string)))
			return err
		},
	},
	TAG_List: {"TAG_List",
		// JMT: figure out how to break loop
		nil,
		nil,
	},
	TAG_Compound: {"TAG_Compound",
		// JMT: figure out how to break loop
		nil,
		nil,
	},
	TAG_Int_Array: {"TAG_Int_Array",
		func(r io.Reader) (interface{}, error) {
			var strlen int32
			if err := binary.Read(r, binary.BigEndian, &strlen); err != nil {
				return nil, err
			}

			ints := make([]int32, strlen)
			for key := range ints {
				if err := binary.Read(r, binary.BigEndian, &ints[key]); err != nil {
					return nil, err
				}
			}

			return ints, nil
		},
		func(w io.Writer, i interface{}) error {
			if err := binary.Write(w, binary.BigEndian, int32(len(i.([]int32)))); err != nil {
				return err
			}
			for _, value := range i.([]int32) {
				if err := binary.Write(w, binary.BigEndian, value); err != nil {
					return err
				}
			}
			return nil
		},
	},
}

// JMT: The compound reader and writer cause an initialization loop
// when added to the Tags variable.

func readCompound(r io.Reader) (interface{}, error) {
	payload := []Tag{}
	var emptytag Tag
	for newtag, err := ReadTag(r); newtag != emptytag; newtag, err = ReadTag(r) {
		if err != nil {
			return nil, err
		}
		payload = append(payload, newtag)
	}
	return payload, nil
}

func writeCompound(w io.Writer, i interface{}) error {
	tags := append(i.([]Tag), Tag{Type: TAG_End})
	for _, tag := range tags {
		if err := WriteTag(w, tag); err != nil {
			return err
		}
	}
	return nil
}

// same for the list

func readList(r io.Reader) (i interface{}, err error) {
	var tsubi interface{}
	if tsubi, err = Tags[TAG_Byte].PReader(r); err != nil {
		return
	}
	tsub := tsubi.(byte)
	var tleni interface{}
	if tleni, err = Tags[TAG_Int].PReader(r); err != nil {
		return
	}
	tlen := tleni.(int32)
	var iarrj interface{}
	var j int32
	switch tsub {
	case TAG_Byte:
		// JMT: so much duplicated code!
		iarr := make([]byte, tlen)
		for j = 0; j < tlen; j++ {
			if iarrj, err = Tags[tsub].PReader(r); err != nil {
				return
			}
			iarr[j] = iarrj.(byte)
		}
		return iarr, nil
	case TAG_Short:
		// JMT: so much duplicated code!
		iarr := make([]int16, tlen)
		for j = 0; j < tlen; j++ {
			if iarrj, err = Tags[tsub].PReader(r); err != nil {
				return
			}
			iarr[j] = iarrj.(int16)
		}
		return iarr, nil
	case TAG_Int:
		// JMT: so much duplicated code!
		iarr := make([]int32, tlen)
		for j = 0; j < tlen; j++ {
			if iarrj, err = Tags[tsub].PReader(r); err != nil {
				return
			}
			iarr[j] = iarrj.(int32)
		}
		return iarr, nil
	case TAG_Long:
		// JMT: so much duplicated code!
		iarr := make([]int64, tlen)
		for j = 0; j < tlen; j++ {
			if iarrj, err = Tags[tsub].PReader(r); err != nil {
				return
			}
			iarr[j] = iarrj.(int64)
		}
		return iarr, nil
	case TAG_Float:
		// JMT: so much duplicated code!
		iarr := make([]float32, tlen)
		for j = 0; j < tlen; j++ {
			if iarrj, err = Tags[tsub].PReader(r); err != nil {
				return
			}
			iarr[j] = iarrj.(float32)
		}
		return iarr, nil
	case TAG_Double:
		// JMT: so much duplicated code!
		iarr := make([]float64, tlen)
		for j = 0; j < tlen; j++ {
			if iarrj, err = Tags[tsub].PReader(r); err != nil {
				return
			}
			iarr[j] = iarrj.(float64)
		}
		return iarr, nil
	case TAG_Byte_Array:
		// JMT: so much duplicated code!
		iarr := make([][]byte, tlen)
		for j = 0; j < tlen; j++ {
			if iarrj, err = Tags[tsub].PReader(r); err != nil {
				return
			}
			iarr[j] = iarrj.([]byte)
		}
		return iarr, nil
	case TAG_String:
		// JMT: so much duplicated code!
		iarr := make([]string, tlen)
		for j = 0; j < tlen; j++ {
			if iarrj, err = Tags[tsub].PReader(r); err != nil {
				return
			}
			iarr[j] = iarrj.(string)
		}
		return iarr, nil
	case TAG_List:
		// JMT: so much duplicated code!
		iarr := make([][]interface{}, tlen)
		for j = 0; j < tlen; j++ {
			if iarrj, err = readList(r); err != nil {
				return
			}
			iarr[j] = iarrj.([]interface{})
		}
		return iarr, nil
	case TAG_Compound:
		// JMT: so much duplicated code!
		iarr := make([][]Tag, tlen)
		for j = 0; j < tlen; j++ {
			if iarrj, err = readCompound(r); err != nil {
				return
			}
			iarr[j] = iarrj.([]Tag)
		}
		return iarr, nil
	case TAG_Int_Array:
		// JMT: so much duplicated code!
		iarr := make([][]int32, tlen)
		for j = 0; j < tlen; j++ {
			if iarrj, err = Tags[tsub].PReader(r); err != nil {
				return
			}
			iarr[j] = iarrj.([]int32)
		}
		return iarr, nil
	}
	return
}

func writeList(w io.Writer, i interface{}) error {
	var tsub byte
	var tlen int32
	var tout bytes.Buffer
	switch arr := i.(type) {
	case []byte:
		// JMT: why must this code be repeated
		tsub = TAG_Byte
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := Tags[tsub].PWriter(&tout, value); err != nil {
				return err
			}
		}
	case []int16:
		// JMT: why must this code be repeated
		tsub = TAG_Short
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := Tags[tsub].PWriter(&tout, value); err != nil {
				return err
			}
		}
	case []int32:
		// JMT: why must this code be repeated
		tsub = TAG_Int
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := Tags[tsub].PWriter(&tout, value); err != nil {
				return err
			}
		}
	case []int64:
		// JMT: why must this code be repeated
		tsub = TAG_Long
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := Tags[tsub].PWriter(&tout, value); err != nil {
				return err
			}
		}
	case []float32:
		// JMT: why must this code be repeated
		tsub = TAG_Float
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := Tags[tsub].PWriter(&tout, value); err != nil {
				return err
			}
		}
	case []float64:
		// JMT: why must this code be repeated
		tsub = TAG_Double
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := Tags[tsub].PWriter(&tout, value); err != nil {
				return err
			}
		}
	case [][]byte:
		// JMT: why must this code be repeated
		tsub = TAG_Byte_Array
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := Tags[tsub].PWriter(&tout, value); err != nil {
				return err
			}
		}
	case []string:
		// JMT: why must this code be repeated
		tsub = TAG_String
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := Tags[tsub].PWriter(&tout, value); err != nil {
				return err
			}
		}
	case []interface{}:
		// JMT: why must this code be repeated
		tsub = TAG_List
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := writeList(&tout, value); err != nil {
				return err
			}
		}
	case [][]Tag:
		// JMT: why must this code be repeated
		tsub = TAG_Compound
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := writeCompound(&tout, value); err != nil {
				return err
			}
		}
	case [][]int32:
		// JMT: why must this code be repeated
		tsub = TAG_Int_Array
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := Tags[tsub].PWriter(&tout, value); err != nil {
				return err
			}
		}
	}
	Tags[TAG_Byte].PWriter(w, tsub)
	Tags[TAG_Int].PWriter(w, tlen)
	tout.WriteTo(w)
	return nil
}

func ReadTag(r io.Reader) (t Tag, err error) {
	// read a byte
	var ttypei, tnamei interface{}
	if ttypei, err = Tags[TAG_Byte].PReader(r); err != nil {
		return t, err
	}
	ttype := ttypei.(byte)

	// TAG_End isn't really a tag
	if ttype == TAG_End {
		return
	}

	// Now about that name
	if tnamei, err = Tags[TAG_String].PReader(r); err != nil {
		return
	}
	tname := tnamei.(string)

	// Putting this in the widget causes an initialization loop issue
	// (Tags refers to readCompound refers to ReadTag refers to Tags)
	switch ttype {
	case TAG_List:
		if payload, err := readList(r); err == nil {
			t = Tag{Type: ttype, Name: tname, Payload: payload}
		}
	case TAG_Compound:
		if payload, err := readCompound(r); err == nil {
			t = Tag{Type: ttype, Name: tname, Payload: payload}
		}
	default:
		if val, ok := Tags[ttype]; ok {
			if payload, err := val.PReader(r); err == nil {
				t = Tag{Type: ttype, Name: tname, Payload: payload}
			}
		} else {
			err = fmt.Errorf("unknown tag")
		}
	}
	return
}

func WriteTag(w io.Writer, t Tag) (err error) {
	Tags[TAG_Byte].PWriter(w, t.Type)

	if t.Type == TAG_End {
		return
	}

	Tags[TAG_String].PWriter(w, t.Name)

	// JMT: initialization loop issue here too
	switch t.Type {
	case TAG_List:
		err = writeList(w, t.Payload)
	case TAG_Compound:
		err = writeCompound(w, t.Payload)
	default:
		if val, ok := Tags[t.Type]; ok {
			err = val.PWriter(w, t.Payload)
		} else {
			err = fmt.Errorf("unknown tag")
		}
	}
	return
}

package nbt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

var Debug = false

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

func NewTag(Type byte, Name string) *Tag {
	if Debug {
		if Type != TAG_End {
			log.Printf("NEW TAG: type %s and name %#v", Names[Type], Name)
		}
	}
	return &Tag{Type: Type, Name: Name}
}

func MakeTag(Type byte, Name string) Tag {
	if Debug {
		if Type != TAG_End {
			log.Printf("MAKE TAG: type %s and name %#v", Names[Type], Name)
		}
	}
	return Tag{Type: Type, Name: Name}
}

func (t Tag) String() string {
	return fmt.Sprintf("Tag{Type: %s, Name: %#v, Payload: %v}", Names[t.Type], t.Name, t.Payload)
}

func (t *Tag) SetPayload(newp interface{}) error {
	if Debug {
		log.Printf("SET PAYLOAD: type %s, name %#v, payload type %T", Names[t.Type], t.Name, newp)
	}
	switch newp.(type) {
	case byte:
		if t.Type == TAG_Byte {
			t.Payload = newp
			return nil
		}
	case int16:
		if t.Type == TAG_Short {
			t.Payload = newp
			return nil
		}
	case int32:
		if t.Type == TAG_Int {
			t.Payload = newp
			return nil
		}
	case int64:
		if t.Type == TAG_Long {
			t.Payload = newp
			return nil
		}
	case float32:
		if t.Type == TAG_Float {
			t.Payload = newp
			return nil
		}
	case float64:
		if t.Type == TAG_Double {
			t.Payload = newp
			return nil
		}
	case []byte:
		if t.Type == TAG_Byte_Array {
			t.Payload = newp
			return nil
		}
	case string:
		if t.Type == TAG_String {
			t.Payload = newp
			return nil
		}
	case []Tag:
		if t.Type == TAG_Compound {
			t.Payload = newp
			return nil
		}
	case []int32:
		if t.Type == TAG_Int_Array {
			t.Payload = newp
			return nil
		}
	// JMT: this must be at the bottom because it's a wildcard effectively
	case interface{}:
		if t.Type == TAG_List {
			t.Payload = newp
			return nil
		}
	case nil:
		if t.Type == TAG_List {
			t.Payload = newp
			return nil
		}
	}
	// only set payload if appropriate for this type
	return fmt.Errorf("type %s does not match payload %v (%T)", Names[t.Type], newp, newp)
}

var Names = map[byte]string{
	TAG_End:        "TAG_End",
	TAG_Byte:       "TAG_Byte",
	TAG_Short:      "TAG_Short",
	TAG_Int:        "TAG_Int",
	TAG_Long:       "TAG_Long",
	TAG_Float:      "TAG_Float",
	TAG_Double:     "TAG_Double",
	TAG_Byte_Array: "TAG_Byte_Array",
	TAG_String:     "TAG_String",
	TAG_List:       "TAG_List",
	TAG_Compound:   "TAG_Compound",
	TAG_Int_Array:  "TAG_Int_Array",
}

type PayloadReader func(io.Reader) (interface{}, error)

var PReaders = map[byte]PayloadReader{
	TAG_Byte: func(r io.Reader) (interface{}, error) {
		payload := make([]byte, 1)
		if _, err := io.ReadFull(r, payload); err != nil {
			return nil, err
		}
		return payload[0], nil
	},
	TAG_Short: func(r io.Reader) (interface{}, error) {
		var payload int16
		if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
			return nil, err
		}
		return payload, nil
	},
	TAG_Int: func(r io.Reader) (interface{}, error) {
		var payload int32
		if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
			return nil, err
		}
		return payload, nil
	},
	TAG_Long: func(r io.Reader) (interface{}, error) {
		var payload int64
		if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
			return nil, err
		}
		return payload, nil
	},
	TAG_Float: func(r io.Reader) (interface{}, error) {
		var payload float32
		if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
			return nil, err
		}
		return payload, nil
	},
	TAG_Double: func(r io.Reader) (interface{}, error) {
		var payload float64
		if err := binary.Read(r, binary.BigEndian, &payload); err != nil {
			return nil, err
		}
		return payload, nil
	},
	TAG_Byte_Array: func(r io.Reader) (interface{}, error) {
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
	TAG_String: func(r io.Reader) (interface{}, error) {
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
	TAG_Int_Array: func(r io.Reader) (interface{}, error) {
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
}

type PayloadWriter func(io.Writer, interface{}) error

var PWriters = map[byte]PayloadWriter{
	TAG_Byte: func(w io.Writer, i interface{}) error {
		b := make([]byte, 1)
		b[0] = i.(byte)
		_, err := w.Write(b)
		return err
	},
	TAG_Short: func(w io.Writer, i interface{}) error {
		return binary.Write(w, binary.BigEndian, i.(int16))
	},
	TAG_Int: func(w io.Writer, i interface{}) error {
		return binary.Write(w, binary.BigEndian, i.(int32))
	},
	TAG_Long: func(w io.Writer, i interface{}) error {
		return binary.Write(w, binary.BigEndian, i.(int64))
	},
	TAG_Float: func(w io.Writer, i interface{}) error {
		return binary.Write(w, binary.BigEndian, i.(float32))
	},
	TAG_Double: func(w io.Writer, i interface{}) error {
		return binary.Write(w, binary.BigEndian, i.(float64))
	},
	TAG_Byte_Array: func(w io.Writer, i interface{}) error {
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
	TAG_String: func(w io.Writer, i interface{}) error {
		if err := binary.Write(w, binary.BigEndian, int16(len(i.(string)))); err != nil {
			return err
		}
		_, err := w.Write([]byte(i.(string)))
		return err
	},
	TAG_Int_Array: func(w io.Writer, i interface{}) error {
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
}

type ListReader func(io.Reader, int) (interface{}, error)

var LReaders = map[byte]ListReader{
	TAG_Byte: func(r io.Reader, tlen int) (interface{}, error) {
		iarr := make([]byte, tlen)
		for j := 0; j < tlen; j++ {
			if iarrj, err := PReaders[TAG_Byte](r); err != nil {
				return nil, err
			} else {
				iarr[j] = iarrj.(byte)
			}
		}
		return iarr, nil
	},
	TAG_Short: func(r io.Reader, tlen int) (interface{}, error) {
		iarr := make([]int16, tlen)
		for j := 0; j < tlen; j++ {
			if iarrj, err := PReaders[TAG_Short](r); err != nil {
				return nil, err
			} else {
				iarr[j] = iarrj.(int16)
			}
		}
		return iarr, nil
	},
	TAG_Int: func(r io.Reader, tlen int) (interface{}, error) {
		iarr := make([]int32, tlen)
		for j := 0; j < tlen; j++ {
			if iarrj, err := PReaders[TAG_Int](r); err != nil {
				return nil, err
			} else {
				iarr[j] = iarrj.(int32)
			}
		}
		return iarr, nil
	},
	TAG_Long: func(r io.Reader, tlen int) (interface{}, error) {
		iarr := make([]int64, tlen)
		for j := 0; j < tlen; j++ {
			if iarrj, err := PReaders[TAG_Long](r); err != nil {
				return nil, err
			} else {
				iarr[j] = iarrj.(int64)
			}
		}
		return iarr, nil
	},
	TAG_Float: func(r io.Reader, tlen int) (interface{}, error) {
		iarr := make([]float32, tlen)
		for j := 0; j < tlen; j++ {
			if iarrj, err := PReaders[TAG_Float](r); err != nil {
				return nil, err
			} else {
				iarr[j] = iarrj.(float32)
			}
		}
		return iarr, nil
	},
	TAG_Double: func(r io.Reader, tlen int) (interface{}, error) {
		iarr := make([]float64, tlen)
		for j := 0; j < tlen; j++ {
			if iarrj, err := PReaders[TAG_Double](r); err != nil {
				return nil, err
			} else {
				iarr[j] = iarrj.(float64)
			}
		}
		return iarr, nil
	},
	TAG_Byte_Array: func(r io.Reader, tlen int) (interface{}, error) {
		iarr := make([][]byte, tlen)
		for j := 0; j < tlen; j++ {
			if iarrj, err := PReaders[TAG_Byte_Array](r); err != nil {
				return nil, err
			} else {
				iarr[j] = iarrj.([]byte)
			}
		}
		return iarr, nil
	},
	TAG_String: func(r io.Reader, tlen int) (interface{}, error) {
		iarr := make([]string, tlen)
		for j := 0; j < tlen; j++ {
			if iarrj, err := PReaders[TAG_String](r); err != nil {
				return nil, err
			} else {
				iarr[j] = iarrj.(string)
			}
		}
		return iarr, nil
	},
	TAG_Int_Array: func(r io.Reader, tlen int) (interface{}, error) {
		iarr := make([][]int32, tlen)
		for j := 0; j < tlen; j++ {
			if iarrj, err := PReaders[TAG_Int_Array](r); err != nil {
				return nil, err
			} else {
				iarr[j] = iarrj.([]int32)
			}
		}
		return iarr, nil
	},
}

func readCompound(r io.Reader) (i interface{}, err error) {
	// must return interface{} since it may be called from a list
	payload := []Tag{}
	var newtag Tag
	endtag := MakeTag(TAG_End, "")
	// read the first one
	for newtag, err = ReadTag(r); newtag != endtag; newtag, err = ReadTag(r) {
		if err != nil {
			break
		}
		payload = append(payload, newtag)
	}
	i = payload
	return
}

func writeCompound(w io.Writer, i interface{}) error {
	tags := append(i.([]Tag), MakeTag(TAG_End, ""))
	for _, tag := range tags {
		if err := tag.Write(w); err != nil {
			return err
		}
	}
	return nil
}

func readList(r io.Reader) (i interface{}, err error) {
	var tsubi, tleni interface{}
	if tsubi, err = PReaders[TAG_Byte](r); err != nil {
		return
	}
	tsub := tsubi.(byte)

	if tleni, err = PReaders[TAG_Int](r); err != nil {
		return
	}
	tlen := int(tleni.(int32))

	switch tsub {
	case TAG_List:
		iarr := make([][]interface{}, tlen)
		for j := 0; j < tlen; j++ {
			if iarrj, err := readList(r); err != nil {
				return nil, err
			} else {
				iarr[j] = iarrj.([]interface{})
			}
		}
		return iarr, nil
	case TAG_Compound:
		iarr := make([][]Tag, tlen)
		for j := 0; j < tlen; j++ {
			if iarrj, err := readCompound(r); err != nil {
				return nil, err
			} else {
				iarr[j] = iarrj.([]Tag)
			}
		}
		return iarr, nil
	default:
		if val, ok := LReaders[tsub]; ok {
			return val(r, tlen)
		} else {
			return
		}
	}
}

// JMT: not sure if this can be normalized due to type thing
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
			if err := PWriters[tsub](&tout, value); err != nil {
				return err
			}
		}
	case []int16:
		// JMT: why must this code be repeated
		tsub = TAG_Short
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := PWriters[tsub](&tout, value); err != nil {
				return err
			}
		}
	case []int32:
		// JMT: why must this code be repeated
		tsub = TAG_Int
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := PWriters[tsub](&tout, value); err != nil {
				return err
			}
		}
	case []int64:
		// JMT: why must this code be repeated
		tsub = TAG_Long
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := PWriters[tsub](&tout, value); err != nil {
				return err
			}
		}
	case []float32:
		// JMT: why must this code be repeated
		tsub = TAG_Float
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := PWriters[tsub](&tout, value); err != nil {
				return err
			}
		}
	case []float64:
		// JMT: why must this code be repeated
		tsub = TAG_Double
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := PWriters[tsub](&tout, value); err != nil {
				return err
			}
		}
	case [][]byte:
		// JMT: why must this code be repeated
		tsub = TAG_Byte_Array
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := PWriters[tsub](&tout, value); err != nil {
				return err
			}
		}
	case []string:
		// JMT: why must this code be repeated
		tsub = TAG_String
		tlen = int32(len(arr))
		for _, value := range arr {
			if err := PWriters[tsub](&tout, value); err != nil {
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
			if err := PWriters[tsub](&tout, value); err != nil {
				return err
			}
		}
	}
	if err := PWriters[TAG_Byte](w, tsub); err != nil {
		return err
	}
	if err := PWriters[TAG_Int](w, tlen); err != nil {
		return err
	}
	if _, err := tout.WriteTo(w); err != nil {
		return err
	}
	return nil
}

func ReadTag(r io.Reader) (Tag, error) {
	// defaults
	t := MakeTag(TAG_End, "")
	err := fmt.Errorf("unknown tag")

	var ttypei, tnamei interface{}
	if ttypei, err = PReaders[TAG_Byte](r); err != nil {
		return t, err
	}
	ttype := ttypei.(byte)

	if ttype == TAG_End {
		return t, err
	}

	if tnamei, err = PReaders[TAG_String](r); err != nil {
		return t, err
	}
	tname := tnamei.(string)

	if Debug {
		log.Printf("ReadTag: type %s name %s", Names[ttype], tname)
	}

	var payload interface{}
	switch ttype {
	case TAG_List:
		if payload, err = readList(r); err == nil {
			t = MakeTag(ttype, tname)
			err = t.SetPayload(payload)
		}
	case TAG_Compound:
		if payload, err = readCompound(r); err == nil {
			t = MakeTag(ttype, tname)
			err = t.SetPayload(payload)
		}
	default:
		if val, ok := PReaders[ttype]; ok {
			if payload, err = val(r); err == nil {
				t = MakeTag(ttype, tname)
				err = t.SetPayload(payload)
			}
		} else {
			err = fmt.Errorf("no PReader found for type %s", Names[ttype])
		}
	}
	return t, err
}

func (t Tag) Write(w io.Writer) error {
	PWriters[TAG_Byte](w, t.Type)

	if t.Type == TAG_End {
		return nil
	}

	PWriters[TAG_String](w, t.Name)

	// JMT: initialization loop issue here too
	switch t.Type {
	case TAG_List:
		return writeList(w, t.Payload)
	case TAG_Compound:
		return writeCompound(w, t.Payload)
	default:
		if val, ok := PWriters[t.Type]; ok {
			return val(w, t.Payload)
		} else {
			return fmt.Errorf("unknown tag")
		}
	}
	return nil
}

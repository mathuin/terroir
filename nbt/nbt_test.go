package nbt

import (
	"bytes"
	"reflect"
	"testing"
)

var tag_tests = []struct {
	input  Tag
	output []byte
}{
	{Tag{Type: TAG_Byte, Name: "lookbyte", Payload: []byte{'0'}}, []byte{1, 0, 8, 'l', 'o', 'o', 'k', 'b', 'y', 't', 'e', '0'}},
	{Tag{Type: TAG_Short, Name: "lookshort", Payload: int16(2)}, []byte{2, 0, 9, 'l', 'o', 'o', 'k', 's', 'h', 'o', 'r', 't', 0, 2}},
	{Tag{Type: TAG_Int, Name: "lookint", Payload: int32(4)}, []byte{3, 0, 7, 'l', 'o', 'o', 'k', 'i', 'n', 't', 0, 0, 0, 4}},
	{Tag{Type: TAG_Long, Name: "looklong", Payload: int64(8)}, []byte{4, 0, 8, 'l', 'o', 'o', 'k', 'l', 'o', 'n', 'g', 0, 0, 0, 0, 0, 0, 0, 8}},
	{Tag{Type: TAG_Float, Name: "lookfloat", Payload: float32(3.14)}, []byte{0x5, 0x0, 0x9, 0x6c, 0x6f, 0x6f, 0x6b, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x40, 0x48, 0xf5, 0xc3}},
	{Tag{Type: TAG_Double, Name: "lookdouble", Payload: float64(6.28)}, []byte{0x6, 0x0, 0xa, 0x6c, 0x6f, 0x6f, 0x6b, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x40, 0x19, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f}},
	{Tag{Type: TAG_Byte_Array, Name: "lookbytes", Payload: []byte{'1', '2', '3'}}, []byte{0x7, 0x0, 0x9, 0x6c, 0x6f, 0x6f, 0x6b, 0x62, 0x79, 0x74, 0x65, 0x73, 0x0, 0x0, 0x0, 0x3, 0x31, 0x32, 0x33}},
	{Tag{Type: TAG_String, Name: "lookstring", Payload: "string"}, []byte{8, 0, 10, 'l', 'o', 'o', 'k', 's', 't', 'r', 'i', 'n', 'g', 0, 6, 's', 't', 'r', 'i', 'n', 'g'}},
	{Tag{Type: TAG_Compound, Name: "", Payload: []Tag{Tag{Type: TAG_Byte, Name: "lookbyte", Payload: []byte{'0'}}, Tag{Type: TAG_Int, Name: "lookint", Payload: int32(4)}}}, []byte{10, 0, 0, 1, 0, 8, 'l', 'o', 'o', 'k', 'b', 'y', 't', 'e', '0', 3, 0, 7, 'l', 'o', 'o', 'k', 'i', 'n', 't', 0, 0, 0, 4, 0}},
	{Tag{Type: TAG_Compound, Name: "", Payload: []Tag{Tag{Type: TAG_Compound, Name: "lookcompound", Payload: []Tag{Tag{Type: TAG_Byte, Name: "lookbyte", Payload: []byte{'0'}}, Tag{Type: TAG_Int, Name: "lookint", Payload: int32(4)}}}}}, []byte{10, 0, 0, 10, 0, 12, 'l', 'o', 'o', 'k', 'c', 'o', 'm', 'p', 'o', 'u', 'n', 'd', 1, 0, 8, 'l', 'o', 'o', 'k', 'b', 'y', 't', 'e', '0', 3, 0, 7, 'l', 'o', 'o', 'k', 'i', 'n', 't', 0, 0, 0, 4, 0, 0}},
	{Tag{Type: TAG_Int_Array, Name: "lookints", Payload: []int32{1, 2, 3}}, []byte{0xb, 0x0, 0x8, 0x6c, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x74, 0x73, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x3}},
}

func Test_WriteTag(t *testing.T) {
	for _, tt := range tag_tests {
		var b bytes.Buffer
		err := WriteTag(&b, tt.input)
		if err != nil {
			t.Fail()
		}
		if bytes.Compare(b.Bytes(), tt.output) != 0 {
			t.Errorf("Given %+#v, wanted %+#v, got %+#v", tt.input, tt.output, b.Bytes())
		}
	}
}

func Test_ReadTag(t *testing.T) {
	for _, tt := range tag_tests {
		b := bytes.NewBuffer(tt.output)
		readtags, err := ReadTag(b)
		if err != nil {
			t.Fail()
		}
		if !reflect.DeepEqual(readtags, tt.input) {
			t.Errorf("Given %+#v, wanted %+#v, got %+#v", tt.output, tt.input, readtags)
		}
	}
}

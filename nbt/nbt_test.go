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
	{Tag{Type: TAG_Int, Name: "lookint", Payload: int32(4)}, []byte{3, 0, 7, 'l', 'o', 'o', 'k', 'i', 'n', 't', 0, 0, 0, 4}},
	{Tag{Type: TAG_Compound, Name: "", Payload: []Tag{Tag{Type: TAG_Byte, Name: "lookbyte", Payload: []byte{'0'}}, Tag{Type: TAG_Int, Name: "lookint", Payload: int32(4)}}}, []byte{10, 0, 0, 1, 0, 8, 'l', 'o', 'o', 'k', 'b', 'y', 't', 'e', '0', 3, 0, 7, 'l', 'o', 'o', 'k', 'i', 'n', 't', 0, 0, 0, 4, 0}},
	{Tag{Type: TAG_Compound, Name: "", Payload: []Tag{Tag{Type: TAG_Compound, Name: "lookcompound", Payload: []Tag{Tag{Type: TAG_Byte, Name: "lookbyte", Payload: []byte{'0'}}, Tag{Type: TAG_Int, Name: "lookint", Payload: int32(4)}}}}}, []byte{10, 0, 0, 10, 0, 12, 'l', 'o', 'o', 'k', 'c', 'o', 'm', 'p', 'o', 'u', 'n', 'd', 1, 0, 8, 'l', 'o', 'o', 'k', 'b', 'y', 't', 'e', '0', 3, 0, 7, 'l', 'o', 'o', 'k', 'i', 'n', 't', 0, 0, 0, 4, 0, 0}},
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

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
	{Tag{Type: TAG_Byte, Name: "lookbyte", Payload: byte('0')}, []byte{1, 0, 8, 'l', 'o', 'o', 'k', 'b', 'y', 't', 'e', '0'}},
	{Tag{Type: TAG_Short, Name: "lookshort", Payload: int16(2)}, []byte{2, 0, 9, 'l', 'o', 'o', 'k', 's', 'h', 'o', 'r', 't', 0, 2}},
	{Tag{Type: TAG_Int, Name: "lookint", Payload: int32(4)}, []byte{3, 0, 7, 'l', 'o', 'o', 'k', 'i', 'n', 't', 0, 0, 0, 4}},
	{Tag{Type: TAG_Long, Name: "looklong", Payload: int64(8)}, []byte{4, 0, 8, 'l', 'o', 'o', 'k', 'l', 'o', 'n', 'g', 0, 0, 0, 0, 0, 0, 0, 8}},
	{Tag{Type: TAG_Float, Name: "lookfloat", Payload: float32(3.14)}, []byte{0x5, 0x0, 0x9, 0x6c, 0x6f, 0x6f, 0x6b, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x40, 0x48, 0xf5, 0xc3}},
	{Tag{Type: TAG_Double, Name: "lookdouble", Payload: float64(6.28)}, []byte{0x6, 0x0, 0xa, 0x6c, 0x6f, 0x6f, 0x6b, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x40, 0x19, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f}},
	{Tag{Type: TAG_Byte_Array, Name: "lookbytes", Payload: []byte{'1', '2', '3'}}, []byte{0x7, 0x0, 0x9, 0x6c, 0x6f, 0x6f, 0x6b, 0x62, 0x79, 0x74, 0x65, 0x73, 0x0, 0x0, 0x0, 0x3, 0x31, 0x32, 0x33}},
	{Tag{Type: TAG_String, Name: "lookstring", Payload: string("string")}, []byte{8, 0, 10, 'l', 'o', 'o', 'k', 's', 't', 'r', 'i', 'n', 'g', 0, 6, 's', 't', 'r', 'i', 'n', 'g'}},
	{Tag{Type: TAG_List, Name: "looklist", Payload: []float32{1.23, 4.56, 7.89}}, []byte{0x9, 0x0, 0x8, 0x6c, 0x6f, 0x6f, 0x6b, 0x6c, 0x69, 0x73, 0x74, 0x5, 0x0, 0x0, 0x0, 0x3, 0x3f, 0x9d, 0x70, 0xa4, 0x40, 0x91, 0xeb, 0x85, 0x40, 0xfc, 0x7a, 0xe1}},
	{Tag{Type: TAG_Compound, Name: "", Payload: []Tag{Tag{Type: TAG_Byte, Name: "lookbyte", Payload: byte('0')}, Tag{Type: TAG_Int, Name: "lookint", Payload: int32(4)}}}, []byte{10, 0, 0, 1, 0, 8, 'l', 'o', 'o', 'k', 'b', 'y', 't', 'e', '0', 3, 0, 7, 'l', 'o', 'o', 'k', 'i', 'n', 't', 0, 0, 0, 4, 0}},
	{Tag{Type: TAG_Compound, Name: "", Payload: []Tag{Tag{Type: TAG_Compound, Name: "lookcompound", Payload: []Tag{Tag{Type: TAG_Byte, Name: "lookbyte", Payload: byte('0')}, Tag{Type: TAG_Int, Name: "lookint", Payload: int32(4)}}}}}, []byte{10, 0, 0, 10, 0, 12, 'l', 'o', 'o', 'k', 'c', 'o', 'm', 'p', 'o', 'u', 'n', 'd', 1, 0, 8, 'l', 'o', 'o', 'k', 'b', 'y', 't', 'e', '0', 3, 0, 7, 'l', 'o', 'o', 'k', 'i', 'n', 't', 0, 0, 0, 4, 0, 0}},
	{Tag{Type: TAG_Int_Array, Name: "lookints", Payload: []int32{1, 2, 3}}, []byte{0xb, 0x0, 0x8, 0x6c, 0x6f, 0x6f, 0x6b, 0x69, 0x6e, 0x74, 0x73, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x3}},
	// file tests -- helloworld.nbt
	{Tag{Type: TAG_Compound, Name: "hello world", Payload: []Tag{Tag{Type: TAG_String, Name: "name", Payload: string("Bananrama")}}}, []byte{0xa, 0x0, 0xb, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x8, 0x0, 0x4, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0x9, 0x42, 0x61, 0x6e, 0x61, 0x6e, 0x72, 0x61, 0x6d, 0x61, 0x0}},
	// bigtest.nbt
	{Tag{Type: 0xa, Name: "Level", Payload: []Tag{Tag{Type: 0x4, Name: "longTest", Payload: int64(9223372036854775807)}, Tag{Type: 0x2, Name: "shortTest", Payload: int16(32767)}, Tag{Type: 0x8, Name: "stringTest", Payload: "HELLO WORLD THIS IS A TEST STRING ÅÄÖ!"}, Tag{Type: 0x5, Name: "floatTest", Payload: float32(0.49823147)}, Tag{Type: 0x3, Name: "intTest", Payload: int32(2147483647)}, Tag{Type: 0xa, Name: "nested compound test", Payload: []Tag{Tag{Type: 0xa, Name: "ham", Payload: []Tag{Tag{Type: 0x8, Name: "name", Payload: "Hampus"}, Tag{Type: 0x5, Name: "value", Payload: float32(0.75)}}}, Tag{Type: 0xa, Name: "egg", Payload: []Tag{Tag{Type: 0x8, Name: "name", Payload: "Eggbert"}, Tag{Type: 0x5, Name: "value", Payload: float32(0.5)}}}}}, Tag{Type: 0x9, Name: "listTest (long)", Payload: []int64{11, 12, 13, 14, 15}}, Tag{Type: 0x9, Name: "listTest (compound)", Payload: [][]Tag{[]Tag{Tag{Type: 0x8, Name: "name", Payload: "Compound tag #0"}, Tag{Type: 0x4, Name: "created-on", Payload: int64(1264099775885)}}, []Tag{Tag{Type: 0x8, Name: "name", Payload: "Compound tag #1"}, Tag{Type: 0x4, Name: "created-on", Payload: int64(1264099775885)}}}}, Tag{Type: 0x1, Name: "byteTest", Payload: byte(0x7f)}, Tag{Type: 0x7, Name: "byteArrayTest (the first 1000 values of (n*n*255+n*7)%100, starting with n=0 (0, 62, 34, 16, 8, ...))", Payload: []uint8{0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30}}, Tag{Type: 0x6, Name: "doubleTest", Payload: 0.4931287132182315}}}, []byte{0xa, 0x0, 0x5, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x4, 0x0, 0x8, 0x6c, 0x6f, 0x6e, 0x67, 0x54, 0x65, 0x73, 0x74, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x2, 0x0, 0x9, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x73, 0x74, 0x7f, 0xff, 0x8, 0x0, 0xa, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x54, 0x65, 0x73, 0x74, 0x0, 0x29, 0x48, 0x45, 0x4c, 0x4c, 0x4f, 0x20, 0x57, 0x4f, 0x52, 0x4c, 0x44, 0x20, 0x54, 0x48, 0x49, 0x53, 0x20, 0x49, 0x53, 0x20, 0x41, 0x20, 0x54, 0x45, 0x53, 0x54, 0x20, 0x53, 0x54, 0x52, 0x49, 0x4e, 0x47, 0x20, 0xc3, 0x85, 0xc3, 0x84, 0xc3, 0x96, 0x21, 0x5, 0x0, 0x9, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x54, 0x65, 0x73, 0x74, 0x3e, 0xff, 0x18, 0x32, 0x3, 0x0, 0x7, 0x69, 0x6e, 0x74, 0x54, 0x65, 0x73, 0x74, 0x7f, 0xff, 0xff, 0xff, 0xa, 0x0, 0x14, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x20, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x75, 0x6e, 0x64, 0x20, 0x74, 0x65, 0x73, 0x74, 0xa, 0x0, 0x3, 0x68, 0x61, 0x6d, 0x8, 0x0, 0x4, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0x6, 0x48, 0x61, 0x6d, 0x70, 0x75, 0x73, 0x5, 0x0, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3f, 0x40, 0x0, 0x0, 0x0, 0xa, 0x0, 0x3, 0x65, 0x67, 0x67, 0x8, 0x0, 0x4, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0x7, 0x45, 0x67, 0x67, 0x62, 0x65, 0x72, 0x74, 0x5, 0x0, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9, 0x0, 0xf, 0x6c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x73, 0x74, 0x20, 0x28, 0x6c, 0x6f, 0x6e, 0x67, 0x29, 0x4, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xb, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xd, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xe, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf, 0x9, 0x0, 0x13, 0x6c, 0x69, 0x73, 0x74, 0x54, 0x65, 0x73, 0x74, 0x20, 0x28, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x75, 0x6e, 0x64, 0x29, 0xa, 0x0, 0x0, 0x0, 0x2, 0x8, 0x0, 0x4, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0xf, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x75, 0x6e, 0x64, 0x20, 0x74, 0x61, 0x67, 0x20, 0x23, 0x30, 0x4, 0x0, 0xa, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x2d, 0x6f, 0x6e, 0x0, 0x0, 0x1, 0x26, 0x52, 0x37, 0xd5, 0x8d, 0x0, 0x8, 0x0, 0x4, 0x6e, 0x61, 0x6d, 0x65, 0x0, 0xf, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x75, 0x6e, 0x64, 0x20, 0x74, 0x61, 0x67, 0x20, 0x23, 0x31, 0x4, 0x0, 0xa, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x2d, 0x6f, 0x6e, 0x0, 0x0, 0x1, 0x26, 0x52, 0x37, 0xd5, 0x8d, 0x0, 0x1, 0x0, 0x8, 0x62, 0x79, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x7f, 0x7, 0x0, 0x65, 0x62, 0x79, 0x74, 0x65, 0x41, 0x72, 0x72, 0x61, 0x79, 0x54, 0x65, 0x73, 0x74, 0x20, 0x28, 0x74, 0x68, 0x65, 0x20, 0x66, 0x69, 0x72, 0x73, 0x74, 0x20, 0x31, 0x30, 0x30, 0x30, 0x20, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x20, 0x6f, 0x66, 0x20, 0x28, 0x6e, 0x2a, 0x6e, 0x2a, 0x32, 0x35, 0x35, 0x2b, 0x6e, 0x2a, 0x37, 0x29, 0x25, 0x31, 0x30, 0x30, 0x2c, 0x20, 0x73, 0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x6e, 0x3d, 0x30, 0x20, 0x28, 0x30, 0x2c, 0x20, 0x36, 0x32, 0x2c, 0x20, 0x33, 0x34, 0x2c, 0x20, 0x31, 0x36, 0x2c, 0x20, 0x38, 0x2c, 0x20, 0x2e, 0x2e, 0x2e, 0x29, 0x29, 0x0, 0x0, 0x3, 0xe8, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x0, 0x3e, 0x22, 0x10, 0x8, 0xa, 0x16, 0x2c, 0x4c, 0x12, 0x46, 0x20, 0x4, 0x56, 0x4e, 0x50, 0x5c, 0xe, 0x2e, 0x58, 0x28, 0x2, 0x4a, 0x38, 0x30, 0x32, 0x3e, 0x54, 0x10, 0x3a, 0xa, 0x48, 0x2c, 0x1a, 0x12, 0x14, 0x20, 0x36, 0x56, 0x1c, 0x50, 0x2a, 0xe, 0x60, 0x58, 0x5a, 0x2, 0x18, 0x38, 0x62, 0x32, 0xc, 0x54, 0x42, 0x3a, 0x3c, 0x48, 0x5e, 0x1a, 0x44, 0x14, 0x52, 0x36, 0x24, 0x1c, 0x1e, 0x2a, 0x40, 0x60, 0x26, 0x5a, 0x34, 0x18, 0x6, 0x62, 0x0, 0xc, 0x22, 0x42, 0x8, 0x3c, 0x16, 0x5e, 0x4c, 0x44, 0x46, 0x52, 0x4, 0x24, 0x4e, 0x1e, 0x5c, 0x40, 0x2e, 0x26, 0x28, 0x34, 0x4a, 0x6, 0x30, 0x6, 0x0, 0xa, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x54, 0x65, 0x73, 0x74, 0x3f, 0xdf, 0x8f, 0x6b, 0xbb, 0xff, 0x6a, 0x5e, 0x0}},
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

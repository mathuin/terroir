package world

import "fmt"

type Block struct {
	block int
	data  int
}

func MakeBlock(block int, data int) Block {
	return Block{block: block, data: data}
}

func (b Block) String() string {
	return fmt.Sprintf("Block{block: %d, data: %d}", b.block, b.data)
}

func (w *World) Block(pt Point) (*Block, error) {
	s, err := w.Section(pt)
	if err != nil {
		return nil, err
	}
	base := int(s.Blocks[pt.Index()])
	add := int(Nibble(s.Add, pt.Index()))
	data := int(Nibble(s.Data, pt.Index()))
	retval := base + add*256
	b := MakeBlock(retval, data)
	return &b, nil
}

func (w *World) SetBlock(pt Point, b Block) error {
	base := byte(b.block % 256)
	add := byte(b.block / 256)
	s, err := w.Section(pt)
	if err != nil {
		return err
	}
	i := pt.Index()
	s.Blocks[i] = byte(base)
	WriteNibble(s.Add, i, add)
	WriteNibble(s.Data, i, byte(b.data))
	return nil
}

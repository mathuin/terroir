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

func (w World) Block(pt Point) (*Block, error) {
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

var BlockNames = map[string]Block{}

func BlockNamed(name string) (*Block, error) {
	if val, ok := BlockNames[name]; ok {
		return &val, nil
	}
	return nil, fmt.Errorf("block with name %s does not exist!", name)
}

func (b Block) BlockName() (string, error) {
	for name := range BlockNames {
		if BlockNames[name] == b {
			return name, nil
		}
	}
	return "", fmt.Errorf("block %v has no name!", b)
}

func init() {
	for _, bd := range BlockData {
		b := MakeBlock(bd.Block, bd.Data)
		for _, name := range bd.Names {
			BlockNames[name] = b
		}
	}
}

var BlockData = []struct {
	Block int
	Data  int
	Names []string
}{
	{0, 0, []string{"Air", "Empty"}},
	{1, 0, []string{"Stone"}},
}

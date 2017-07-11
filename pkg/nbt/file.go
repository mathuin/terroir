package nbt

import (
	"compress/gzip"
	"io"
	"os"
)

// TODO: consider merging compressed and uncompressed methods together somehow?

func ReadCompressedFile(filename string) (Tag, error) {
	inf, err := os.Open(filename)
	defer inf.Close()
	if err != nil {
		return Tag{}, err
	}
	f, err := gzip.NewReader(inf)
	defer f.Close()
	if err != nil {
		return Tag{}, err
	}
	return ReadTag(f)
}

func ReadUncompressedFile(filename string) (Tag, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return Tag{}, err
	}
	return ReadTag(f)
}

func defClose(f io.Closer) {
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func WriteCompressedFile(filename string, t Tag) error {
	outf, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer defClose(outf)
	f := gzip.NewWriter(outf)
	if err := t.Write(f); err != nil {
		return err
	}
	f.Close()
	return nil
}

func WriteUncompressedFile(filename string, t Tag) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer defClose(f)
	return t.Write(f)
}

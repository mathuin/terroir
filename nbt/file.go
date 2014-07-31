package nbt

import (
	"compress/gzip"
	"os"
)

// readfile
// input: filename, compressed
// output: tag

// big question: can a file's compressed state be detected on the fly?
// bigger question: SHOULD IT?

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

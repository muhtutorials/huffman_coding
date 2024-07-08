package main

import (
	"huffman_coding/bits"
	"io"
)

type Reader struct {
	*symbols
	br *bits.Reader
}

func NewReader(in io.Reader) *Reader {
	return &Reader{
		symbols: newSymbols(),
		br:      bits.NewReader(in),
	}
}

// Read decompresses up to len(p) bytes from the source
func (r *Reader) Read(p []byte) (n int, err error) {
	for i := range p {
		if p[i], err = r.ReadByte(); err != nil {
			return i, err
		}
	}
	return len(p), nil
}

// ReadByte decompresses a single byte
func (r *Reader) ReadByte() (b byte, err error) {
	node := r.root
	for node.Left != nil { // read until we reach a leaf
		var right bool
		if right, err = r.br.ReadOneBit(); err != nil {
			return 0, err
		}
		if right {
			node = node.Right
		} else {
			node = node.Left
		}
	}

	switch node.Char {
	case newChar:
		if b, err = r.br.ReadByte(); err != nil {
			return 0, err
		}
		r.insert(rune(b))
		return b, nil
	case eof:
		return 0, io.EOF
	default:
		r.update(node)
		return byte(node.Char), nil
	}
}

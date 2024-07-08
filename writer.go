package main

import (
	"huffman_coding/bits"
	"io"
)

// Writer is the Huffman writer implementation.
// Must be closed in order to properly send EOF.
type Writer struct {
	*symbols
	bw *bits.Writer
}

func NewWriter(out io.Writer) *Writer {
	return &Writer{
		symbols: newSymbols(),
		bw:      bits.NewWriter(out),
	}
}

// Write writes the compressed form of p to the underlying io.Writer.
// The compressed byte(s) are not necessarily flushed until the Writer is closed.
func (w *Writer) Write(p []byte) (n int, err error) {
	for i, b := range p {
		if err = w.WriteByte(b); err != nil {
			return i, err
		}
	}
	return len(p), nil
}

// WriteByte writes the compressed form of b to the underlying io.Writer.
// The compressed byte(s) are not necessarily flushed until the Writer is closed.
func (w *Writer) WriteByte(b byte) error {
	char := rune(b)
	node := w.chars[char]

	if node == nil {
		// Character is encountered the first time.
		// So we write the "new character" character and then the character itself.
		if err := w.bw.WriteBits(w.chars[newChar].Code()); err != nil {
			return err
		}
		if err := w.bw.WriteByte(b); err != nil {
			return err
		}
		w.insert(char)
	} else {
		// Character has been encountered already.
		// So we write its code and update the tree.
		if err := w.bw.WriteBits(node.Code()); err != nil {
			return err
		}
		w.update(node)
	}
	return nil
}

// Close closes the Huffman writer properly, sending EOF.
// If the underlying io.Writer implements io.Closer
// it will be closed after sending EOF.
func (w *Writer) Close() error {
	if len(w.leaves) > 2 {
		if err := w.bw.WriteBits(w.chars[eof].Code()); err != nil {
			return err
		}
	}
	return w.bw.Close()
}

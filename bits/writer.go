package bits

import (
	"bufio"
	"io"
)

type Writer struct {
	out *bufio.Writer
	// bits buffer
	buf byte
	// number of bits written to buffer
	count uint8
}

func NewWriter(out io.Writer) *Writer {
	return &Writer{out: bufio.NewWriter(out)}
}

// Write implements io.Writer and gives a byte-level interface to the bit stream.
// This will give the best performance if the underlying io.Writer is aligned
// to a byte boundary (else all the individual bytes are spread to multiple bytes).
// Byte boundary can be ensured by calling Align().
func (w *Writer) Write(p []byte) (n int, err error) {
	if w.count == 0 {
		return w.out.Write(p)
	}
	for i, b := range p {
		if err = w.writeUnalignedByte(b); err != nil {
			return i, err
		}
	}
	return len(p), nil
}

func (w *Writer) WriteByte(b byte) error {
	if w.count == 0 {
		return w.out.WriteByte(b)
	}
	return w.writeUnalignedByte(b)
}

// writeUnalignedByte writes 8 bits which are unaligned
func (w *Writer) writeUnalignedByte(b byte) error {
	count := w.count // 3
	//                     11100000 | 11111111>>3
	//                     11100000 | 00011111
	//                     11111111
	if err := w.out.WriteByte(w.buf | b>>count); err != nil {
		return err
	}
	//       (11111111 & (1<<3 - 1)) << (8 - 3)
	//       (11111111 & (00001000 - 1)) << 5
	//       (11111111 & 00000111) << 5
	//       00000111 << 5
	//       11100000
	w.buf = (b & (1<<count - 1)) << (8 - count)
	return nil
}

// WriteBits writes the n lowest bits of r.
// Bits of r in positions higher than n are ignored.
//
// For example:
//
//	00010010 00110100
//
// w.WriteBits(0x1234, 8)
// is equivalent to:
//
//	00110100
//
// w.WriteBits(0x34, 8)
func (w *Writer) WriteBits(r uint64, n uint8) error {
	// if r had bits set at higher positions than n,
	// WriteBitsUnsafe implementation could "corrupt" bits in byte.
	// That is not acceptable. To be on the safe side, mask out higher bits.
	//      00010010 00110100 & (1<<8-1)
	//      00010010 00110100 & (00000001 00000000-1)
	//      00010010 00110100 & 00000000 111111111
	//               00110100
	return w.WriteBitsUnsafe(r&(1<<n-1), n)
}

// WriteBitsUnsafe writes the n lowest bits of r.
// r must not have bits set at higher positions than n.
// If r doesn't satisfy this, a mask must be explicitly applied before passing it to WriteBitsUnsafe(),
// or WriteBits() should be used instead.
//
// WriteBitsUnsafe() offers slightly better performance than WriteBits() because
// the input r is not masked. Calling WriteBitsUnsafe() with an r that does
// not satisfy this is undefined behavior (might corrupt previously written bits).
//
// E.g. if you want to write 8 bits:
//
// w.WriteBitsUnsafe(0x34, 8) // This is OK, 0x34 has no bits set higher than the 8th.
//
// w.WriteBitsUnsafe(0x1234&0xff, 8) // &0xff (111111111) masks out bits higher than the 8th.
//
// Or:
//
// w.WriteBits(0x1234, 8) // bits higher than the 8th are ignored here.
func (w *Writer) WriteBitsUnsafe(r uint64, n uint8) error {
	// w.count = 2
	// n = 4
	// 6 = 2 + 4
	totalCount := w.count + n
	if totalCount < 8 {
		// 11000000 |= byte(... 0000 0000 1111) << (8 - 6)
		// 11000000 |= 00001111 << (8 - 6)
		// 11000000 |= 00001111 << 2
		// 11000000 |= 00111100
		// 11111100
		w.buf |= byte(r) << (8 - totalCount)
		w.count = totalCount
		return nil
	}

	if totalCount > 8 {
		// n = 22
		// w.count = 3
		// 5 = 8 - 3
		free := 8 - w.count
		// 11100000 | byte(0000 0011 1111 1111 1111 1111 1111>>(22-5))
		// 11100000 | byte(000[0 0011 111]1 1111 1111 1111 1111>>17) removes 17 excess bits
		// 11100000 | byte(000[0 0011 111])
		// 11100000 | 00011111
		// 11111111
		if err := w.out.WriteByte(w.buf | byte(r>>(n-free))); err != nil {
			return err
		}
		// n = 22
		// free = 5
		// 17 = 22 - 5
		n -= free
		for n >= 8 {
			// 9 = 17 - 8
			n -= 8
			// No need to mask r, converting to byte will mask out higher bits
			//             byte(0011 111[1 1111 111]1 1111 1111) >> 9 removes 9 excess bits
			//                    byte(0001 1111 [1111 1111])
			//                             11111111
			// 1 = 9 - 8
			// n -= 8
			// No need to mask r, converting to byte will mask out higher bits
			//             byte(0011 1111 1111 111[1 1111 111]1) >> 1 removes 1 excess bit
			//                    byte(0001 1111 1111 1111 [1111 1111])
			//                             11111111
			if err := w.out.WriteByte(byte(r >> n)); err != nil {
				return err
			}
		}
		// Put remaining bits into buffer
		if n > 0 {
			// Note: n < 8 (in case of n=8, 1<<n would overflow byte)
			// (byte(0011 1111 1111 1111 1111 1111)&(1<<1-1))<<(8-1)
			// (byte(0011 1111 1111 1111 1111 1111)&1)<<7
			// 1111 1111&1<<7
			// 0000 0001<<7
			// 10000000
			w.buf, w.count = (byte(r)&(1<<n-1))<<(8-n), n
		} else {
			w.buf, w.count = 0, 0
		}
		return nil
	}

	// buffer will be filled exactly with the bits to be written
	// w.count = 4
	// n = 4
	// totalCount = 8
	// 11110000 | byte(... 0000 0000 1111)
	// 11110000 | 00001111
	// 11111111
	b := w.buf | byte(r)
	w.buf, w.count = 0, 0
	return w.out.WriteByte(b)
}

// WriteOneBit writes one bit to buffer from left to right (10000000 -> 1[1]000000)
func (w *Writer) WriteOneBit(b bool) error {
	if w.count == 7 {
		if b {
			// 11111110 | 1
			// 11111111
			if err := w.out.WriteByte(w.buf | 1); err != nil {
				return err
			}
		} else {
			// the last bit is zero, so we write buffer as it is
			// 11111110
			if err := w.out.WriteByte(w.buf); err != nil {
				return err
			}
		}
		w.buf, w.count = 0, 0
		return nil
	}

	// count = 1
	w.count++ // 2
	if b {
		// 10000000 | 1 << (8-2)
		// 10000000 | 1 << 6
		// 10000000 | 01000000
		// 11000000
		w.buf |= 1 << (8 - w.count)
	}

	return nil
}

// Align aligns the bit stream to a byte boundary,
// so next write will go into a new byte.
// If there are buffered bits, they are first written to the output.
// Returns the number of unset but still written bits.
func (w *Writer) Align() (unset uint8, err error) {
	if w.count > 0 {
		if err = w.out.WriteByte(w.buf); err != nil {
			return 0, err
		}
		unset = 8 - w.count
		w.buf, w.count = 0, 0
	}
	return unset, nil
}

func (w *Writer) Close() error {
	if _, err := w.Align(); err != nil {
		return err
	}
	return w.out.Flush()
}

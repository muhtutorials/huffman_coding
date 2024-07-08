package bits

import (
	"bufio"
	"io"
)

type Reader struct {
	in    *bufio.Reader
	buf   byte
	count uint8
}

func NewReader(in io.Reader) *Reader {
	return &Reader{in: bufio.NewReader(in)}
}

// Read implements io.Reader and gives a byte-level view of the bit stream.
// This will give the best performance if the underlying io.Reader is aligned
// to a byte boundary, else all the individual bytes are assembled from multiple bytes.
// Byte boundary can be ensured by calling Align().
func (r *Reader) Read(p []byte) (n int, err error) {
	if r.count == 0 {
		return r.in.Read(p)
	}
	for n = range p {
		if p[n], err = r.readUnalignedByte(); err != nil {
			return 0, err
		}
	}
	return n, nil
}

func (r *Reader) ReadByte() (b byte, err error) {
	if r.count == 0 {
		return r.in.ReadByte()
	}
	return r.readUnalignedByte()
}

// readUnalignedByte reads the next 8 bits which are unaligned and returns them as a byte
// bits are read from left to right:
//
//	buffer      read byte    buffer   returned byte
//
// 000000bb and 12345678 -> 00000078 and bb123456
func (r *Reader) readUnalignedByte() (b byte, err error) {
	count := r.count // 2
	// 00000011 << (8 - 2)
	// 00000011 << 6
	// 11000000
	b = r.buf << (8 - count)
	r.buf, err = r.in.ReadByte()
	if err != nil {
		return 0, err
	}
	// 11000000 | 11111111 >> 2
	// 11000000 | 00111111
	// 11111111
	b |= r.buf >> count
	// 11111111 & (1 << 2 - 1)
	// 11111111 & (00000100 - 1)
	// 11111111 & 00000011
	// 00000011
	r.buf &= 1<<count - 1
	return b, nil
}

// ReadBits reads n bits and returns them as the lowest n bits of u
func (r *Reader) ReadBits(n uint8) (u uint64, err error) {
	// enough bits in buffer to fill u
	// n = 4
	// r.count = 6
	// 4 < 6
	if n < r.count {
		// 2 = 6 - 4
		shift := r.count - n
		// 00111111 >> 2
		// 00001111
		u = uint64(r.buf >> shift)
		// 00111111 & 1<<2 - 1
		// 00111111 & 00000100 - 1
		// 00111111 & 00000011
		// 00000011
		r.buf &= 1<<shift - 1
		return u, nil
	}

	// n = 21
	// r.count = 7
	// 21 > 7
	if n > r.count {
		if r.count > 0 {
			// 01111111
			// 00000000 00000000 ... 01111111 (8 bytes or 64 bits)
			u = uint64(r.buf)
			// 14 = 21 - 7
			n -= r.count
		}
		for n >= 8 {
			b, err := r.in.ReadByte()
			if err != nil {
				return 0, err
			}
			// 00000000 ... 01111111<<8 + 0000000 ... 11111111
			// 00000000 ... 01111111 00000000 + 0000000 ... 11111111
			// 00000000 ... 01111111 11111111
			u = u<<8 + uint64(b)
			// 6 = 14 - 8
			n -= 8
		}
		// read last bits if any
		// 6 > 0
		if n > 0 {
			if r.buf, err = r.in.ReadByte(); err != nil {
				// 2 = 8 - 6
				shift := 8 - n
				// 00000000 ... 01111111 11111111 << 6 + uint64(11111111>>2)
				// 00000000 ... 01111111 11111111 << 6 + uint64(00111111)
				// 00000000 ... 01111111 11111111 << 6 + 00000000 ... 00111111
				// 00000000 ... 00011111 11111111 11000000 + 00000000 ... 00111111
				// 00000000 ... 00011111 11111111 11111111
				u = u<<n + uint64(r.buf>>shift)
				// 11111111 & 1<<2 - 1
				// 11111111 & 00000100 - 1
				// 11111111 & 00000011
				// 00000011
				r.buf &= 1<<shift - 1
				r.count = shift
			}
		} else {
			r.count = 0
		}
		return u, nil
	}

	// buffer has exactly as many bits as needed
	// no need to clear buffer since it will be overwritten on the next read
	r.count = 0
	return uint64(r.buf), nil
}

// ReadOneBit reads one bit from buffer from left to right ([1]1111111 -> 0[1]111111 -> 00111111)
func (r *Reader) ReadOneBit() (b bool, err error) {
	if r.count == 0 {
		r.buf, err = r.in.ReadByte()
		if err != nil {
			return false, err
		}
		// (11111111 & (1 << 7)) != 0
		// (11111111 & 10000000) != 0
		// 10000000 != 0
		// true
		b = (r.buf & (1 << 7)) != 0
		// 11111111 & (1<<7 - 1)
		// 11111111 & (10000000 - 1)
		// 11111111 & 01111111
		// 01111111
		r.buf, r.count = r.buf&(1<<7-1), 7
		return b, nil
	}

	// r.count = 6
	r.count--
	// r.count = 5
	// 00111111 & (1 << 5) != 0
	// 00111111 & 00100000 != 0
	// true
	b = (r.buf & (1 << r.count)) != 0
	// 00111111 & (1 << 5 - 1)
	// 00111111 & (00100000 - 1)
	// 00111111 & 00011111
	// 00011111
	r.buf &= 1<<r.count - 1
	return b, nil
}

// Align aligns the bit stream to a byte boundary,
// so next read will read data from the next byte.
// Returns the number of unread bits.
func (r *Reader) Align() (unread uint8) {
	unread = r.count
	r.count = 0
	return unread
}

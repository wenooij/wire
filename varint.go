package wire

import (
	"encoding/binary"
	"io"
)

var Uvarint64 = proto[uint64]{
	read:  readUvarint64,
	write: writeUvarint64,
	size:  sizeUvarint64,
}

func readUvarint64(r Reader) (uint64, error) {
	var (
		x uint64
		s uint
	)
	for i := 0; i < binary.MaxVarintLen64; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		if b < 0x80 {
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, io.ErrShortBuffer
}

func writeUvarint64(w Writer, x uint64) error {
	for i := 0; x >= 0x80; i++ {
		w.WriteByte(byte(x) | 0x80)
		x >>= 7
	}
	return w.WriteByte(byte(x))
}

func sizeUvarint64(x uint64) uint64 {
	n := uint64(1)
	for i := 0; x >= 0x80; i++ {
		n++
		x >>= 7
	}
	return n
}

var Varint64 = proto[int64]{
	read:  readVarint64,
	write: writeVarint64,
	size:  sizeVarint64,
}

func readVarint64(r Reader) (int64, error) {
	ux, err := readUvarint64(r)
	if err != nil {
		return 0, err
	}
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, nil
}

func writeVarint64(w Writer, x int64) error {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return writeUvarint64(w, ux)
}

func sizeVarint64(x int64) uint64 {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return sizeUvarint64(ux)
}

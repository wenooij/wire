package wire

import (
	"fmt"
	"io"
	"math"
)

var Fixed8 Proto[uint8] = proto[uint8]{
	read:  readFixed8,
	write: func(w Writer, x uint8) error { return w.WriteByte(byte(x)) },
	size:  func(uint8) uint64 { return 1 },
}

func readFixed8(r Reader) (uint8, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	return b, nil
}

var Fixed32 Proto[uint32] = proto[uint32]{
	read:  readFixed32,
	write: writeFixed32,
	size:  func(uint32) uint64 { return 4 },
}

func readFixed32(r Reader) (uint32, error) {
	var b [4]byte
	if _, err := r.Read(b[:]); err != nil {
		return 0, err
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24, nil
}

func writeFixed32(w Writer, x uint32) error {
	if _, err := w.Write([]byte{byte(x), byte(x >> 8), byte(x >> 16), byte(x >> 24)}); err != nil {
		return fmt.Errorf("Fixed32.Write: %w", err)
	}
	return nil
}

var Fixed64 Proto[uint64] = proto[uint64]{
	read:  readFixed64,
	write: writeFixed64,
	size:  func(uint64) uint64 { return 8 },
}

func readFixed64(r Reader) (uint64, error) {
	var b [8]byte
	if _, err := r.Read(b[:]); err != nil {
		return 0, err
	}
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56, nil
}

func writeFixed64(w Writer, x uint64) error {
	if _, err := w.Write([]byte{
		byte(x), byte(x >> 8), byte(x >> 16), byte(x >> 24),
		byte(x >> 32), byte(x >> 40), byte(x >> 48), byte(x >> 56)}); err != nil {
		return fmt.Errorf("Fixed64.Write: %w", err)
	}
	return nil
}

var Float32 Proto[float32] = proto[float32]{
	read:  readFloat32,
	write: writeFloat32,
	size:  func(float32) uint64 { return 4 },
}

func readFloat32(r Reader) (float32, error) {
	ux, err := readFixed32(r)
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(uint32(ux)), nil
}

func writeFloat32(w Writer, x float32) error { return writeFixed32(w, uint32(math.Float32bits(x))) }

var Float64 Proto[float64] = proto[float64]{
	read:  readFloat64,
	write: writeFloat64,
	size:  func(float64) uint64 { return 8 },
}

func readFloat64(r Reader) (float64, error) {
	ux, err := readFixed64(r)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(uint64(ux)), nil
}

func writeFloat64(w Writer, x float64) error { return writeFixed64(w, math.Float64bits(x)) }

func Fixed[T any](p Proto[T]) func(size uint64) Proto[T] {
	return func(size uint64) Proto[T] {
		return proto[T]{
			read:  readFixed(p, size),
			write: writeFixed(p, size),
			size:  func(T) uint64 { return size },
		}
	}
}

func readFixed[T any](proto Proto[T], size uint64) func(Reader) (T, error) {
	return func(r Reader) (T, error) {
		rf := newFixedReader(r, size)
		elem, err := proto.Read(rf)
		if err != nil {
			var t T
			return t, nil
		}
		if rf.n != 0 {
			var t T
			return t, fmt.Errorf("fixed reader expected %d more bytes: %w", rf.n, io.ErrUnexpectedEOF)
		}
		return elem, nil
	}
}

func writeFixed[T any](proto Proto[T], size uint64) func(Writer, T) error {
	return func(w Writer, e T) error {
		wf := newFixedWriter(w, size)
		if err := proto.Write(wf, e); err != nil {
			return fmt.Errorf("Fixed.Write: %w", err)
		}
		if wf.n != 0 {
			return fmt.Errorf("Fixed.Write: expected %d more bytes: %w", wf.n, io.ErrUnexpectedEOF)
		}
		return nil
	}
}

type fixedReader struct {
	Reader
	n uint64
}

func newFixedReader(r Reader, n uint64) *fixedReader { return &fixedReader{r, n} }

func (r *fixedReader) ReadByte() (b byte, err error) {
	if r.n <= 0 {
		return 0, io.ErrUnexpectedEOF
	}
	b, err = r.Reader.ReadByte()
	if err != nil {
		return 0, err
	}
	r.n--
	return b, nil
}

func (r *fixedReader) Read(p []byte) (n int, err error) {
	if r.n <= 0 {
		return 0, io.EOF
	}
	if uint64(len(p)) > r.n {
		p = p[0:r.n]
	}
	n, err = r.Reader.Read(p)
	r.n -= uint64(n)
	return
}

type fixedWriter struct {
	Writer
	n uint64
}

func newFixedWriter(w Writer, n uint64) *fixedWriter { return &fixedWriter{w, n} }

func (w *fixedWriter) WriteByte(b byte) error {
	if w.n <= 0 {
		return fmt.Errorf("fixed writer got extra byte")
	}
	if err := w.Writer.WriteByte(b); err != nil {
		return err
	}
	w.n--
	return nil
}

func (w *fixedWriter) Write(p []byte) (n int, err error) {
	if w.n < uint64(len(p)) {
		return 0, fmt.Errorf("fixed writer got %d extra byte(s)", uint64(len(p))-w.n)
	}
	n, err = w.Writer.Write(p)
	w.n -= uint64(n)
	return
}

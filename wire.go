package wire

import (
	"io"
)

type Reader interface {
	io.Reader
	io.ByteReader
}

type Writer interface {
	io.Writer
	io.ByteWriter
}

type Proto[T any] interface {
	Read(Reader) (T, error)
	Write(Writer, T) error
	Size(T) uint64
}

type proto[T any] struct {
	read  func(Reader) (T, error)
	write func(Writer, T) error
	size  func(T) uint64
}

func (f proto[T]) Read(r Reader) (T, error)     { return f.read(r) }
func (f proto[T]) Write(w Writer, elem T) error { return f.write(w, elem) }
func (f proto[T]) Size(elem T) uint64           { return f.size(elem) }

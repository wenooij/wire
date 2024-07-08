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

func (p proto[T]) Read(r Reader) (T, error)     { return p.read(r) }
func (p proto[T]) Write(w Writer, elem T) error { return p.write(w, elem) }
func (p proto[T]) Size(elem T) uint64           { return p.size(elem) }

type ProtoMaker[E, T any] interface {
	Proto[T]
	Make(E) T
}

type protoMaker[E, T any] struct {
	proto[T]
	makeFunc func(E) T
}

func (p protoMaker[E, T]) Make(e E) T { return p.makeFunc(e) }

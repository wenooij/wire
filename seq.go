package wire

import (
	"fmt"
	"io"
)

func Seq[T any](p Proto[T]) Proto[[]T] {
	return proto[[]T]{
		read:  readSeq(p),
		write: writeSeq(p),
		size:  sizeSeq(p),
	}
}

func readSeq[T any](proto Proto[T]) func(Reader) ([]T, error) {
	return func(r Reader) ([]T, error) {
		for res := make([]T, 0, 8); ; {
			elem, err := proto.Read(r)
			if err != nil {
				if err == io.EOF {
					return res, nil
				}
				return nil, err
			}
			res = append(res, elem)
		}
	}
}

func writeSeq[T any](proto Proto[T]) func(w Writer, seq []T) error {
	return func(w Writer, seq []T) error {
		for _, e := range seq {
			if err := proto.Write(w, e); err != nil {
				return fmt.Errorf("Seq.Write: %w", err)
			}
		}
		return nil
	}
}

func sizeSeq[T any](proto Proto[T]) func(seq []T) uint64 {
	return func(seq []T) uint64 {
		var n uint64
		for _, e := range seq {
			n += proto.Size(e)
		}
		return n
	}
}

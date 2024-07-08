package wire

import (
	"fmt"
	"io"
)

func RawSeq[T any](p Proto[T]) ProtoRanger[[]T, T] {
	return protoRanger[[]T, T]{
		proto: proto[[]T]{
			read:  readRawSeq(p),
			write: writeRawSeq(p),
			size:  sizeRawSeq(p),
		},
		rangeFunc: rangeRawSeq(p),
	}
}

func readRawSeq[T any](proto Proto[T]) func(Reader) ([]T, error) {
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

func writeRawSeq[T any](proto Proto[T]) func(w Writer, seq []T) error {
	return func(w Writer, seq []T) error {
		for _, e := range seq {
			if err := proto.Write(w, e); err != nil {
				return fmt.Errorf("Seq.Write: %w", err)
			}
		}
		return nil
	}
}

func sizeRawSeq[T any](proto Proto[T]) func(seq []T) uint64 {
	return func(seq []T) uint64 {
		var n uint64
		for _, e := range seq {
			n += proto.Size(e)
		}
		return n
	}
}

func rangeRawSeq[T any](proto Proto[T]) func(Reader, func(T) error) error {
	return func(r Reader, f func(T) error) error {
		for {
			elem, err := proto.Read(r)
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			if err := f(elem); err != nil {
				if err == ErrStop {
					return nil
				}
				return err
			}
		}
	}
}

func Seq[T any](proto Proto[T]) ProtoMakeRanger[[]T, SpanElem[[]T], T] {
	return spanMakeRanger[[]T, T](RawSeq(proto))
}

package wire

import (
	"fmt"
)

type SpanElem[T any] Tup2Val[uint64, T]

func (s SpanElem[T]) Any() SpanElem[any] { return SpanElem[any]{s.E0, s.E1} }

func (s SpanElem[T]) Size() uint64 { return s.E0 }
func (s SpanElem[T]) Elem() T      { return s.E1 }

func Span[T any](p Proto[T]) ProtoMaker[T, SpanElem[T]] {
	return protoMaker[T, SpanElem[T]]{
		proto: proto[SpanElem[T]]{
			read:  readSpan(p),
			write: writeSpan(p),
			size:  sizeSpan(p),
		},
		makeFunc: func(e T) SpanElem[T] { return SpanElem[T]{p.Size(e), e} },
	}
}

func readSpan[T any](proto Proto[T]) func(Reader) (SpanElem[T], error) {
	return func(r Reader) (SpanElem[T], error) {
		ux, err := readUvarint64(r)
		if err != nil {
			return SpanElem[T]{}, err
		}
		elem, err := Fixed(proto)(ux).Read(r)
		if err != nil {
			return SpanElem[T]{}, err
		}
		return SpanElem[T]{ux, elem}, nil
	}
}

func writeSpan[T any](proto Proto[T]) func(Writer, SpanElem[T]) error {
	return func(w Writer, span SpanElem[T]) error {
		if err := writeUvarint64(w, span.Size()); err != nil {
			return fmt.Errorf("Span.Write: %w", err)
		}
		if err := Fixed(proto)(span.Size()).Write(w, span.Elem()); err != nil {
			return fmt.Errorf("Span.Write: %w", err)
		}
		return nil
	}
}

func sizeSpan[T any](proto Proto[T]) func(SpanElem[T]) uint64 {
	size := sizeTup2(Uvarint64, proto)
	return func(span SpanElem[T]) uint64 { return size(Tup2Val[uint64, T](span)) }
}

func spanMakeRanger[T, R any](p Proto[T]) ProtoMakeRanger[T, SpanElem[T], R] {
	protoMaker := Span(p).(protoMaker[T, SpanElem[T]])
	return protoMakeRanger[T, SpanElem[T], R]{
		protoRanger: protoRanger[SpanElem[T], R]{
			proto: protoMaker.proto,
			rangeFunc: func(r Reader, f func(R) error) error {
				ranger, ok := p.(ProtoRanger[T, R])
				if !ok {
					panic(fmt.Errorf("not a ProtoRanger: %T", p))
				}
				return ranger.Range(r, f)
			},
		},
		makeFunc: protoMaker.makeFunc,
	}
}

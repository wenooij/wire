package wire

type FieldVal[T any] Tup2Val[uint64, SpanElem[T]]

func MakeField[T any](proto Proto[T]) func(num uint64, val T) FieldVal[T] {
	makeSpan := MakeSpan(proto)
	return func(num uint64, val T) FieldVal[T] { return FieldVal[T]{E0: num, E1: makeSpan(val)} }
}

// MakeAnyField creates a field by converting proto to an Any.
// It can prevent accidential type errors when working with Any(T) directly.
func MakeAnyField[T any](proto Proto[T]) func(num uint64, val T) FieldVal[any] {
	return func(num uint64, val T) FieldVal[any] {
		return FieldVal[any]{E0: num, E1: SpanElem[any]{proto.Size(val), val}}
	}
}

func (f FieldVal[T]) Num() uint64 { return f.E0 }
func (f FieldVal[T]) Val() T      { return f.E1.Elem() }

func Field[T any](p Proto[T]) Proto[FieldVal[T]] {
	return proto[FieldVal[T]]{
		read:  readField(p),
		write: writeField(p),
		size:  sizeField(p),
	}
}

func rawField[T any](p Proto[T]) Proto[Tup2Val[uint64, SpanElem[T]]] { return Tup2(Uvarint64, Span(p)) }

func readField[T any](proto Proto[T]) func(r Reader) (FieldVal[T], error) {
	readRaw := rawField(proto).Read
	return func(r Reader) (FieldVal[T], error) {
		tup, err := readRaw(r)
		if err != nil {
			return FieldVal[T]{}, err
		}
		return FieldVal[T](tup), nil
	}
}

func writeField[T any](proto Proto[T]) func(Writer, FieldVal[T]) error {
	writeRaw := rawField(proto).Write
	return func(w Writer, field FieldVal[T]) error { return writeRaw(w, Tup2Val[uint64, SpanElem[T]](field)) }
}

func sizeField[T any](proto Proto[T]) func(FieldVal[T]) uint64 {
	sizeRaw := rawField(proto).Size
	return func(field FieldVal[T]) uint64 { return sizeRaw(Tup2Val[uint64, SpanElem[T]](field)) }
}

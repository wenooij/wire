package wire

type FieldVal[T any] Tup2Val[uint64, SpanElem[T]]

func (f FieldVal[T]) Num() uint64        { return f.E0 }
func (f FieldVal[T]) Val() T             { return f.E1.Elem() }
func (f FieldVal[T]) Any() FieldVal[any] { return FieldVal[any]{f.Num(), f.E1.Any()} }

func Field[T any](p Proto[T]) ProtoMaker[Tup2Val[uint64, T], FieldVal[T]] {
	makeSpan := Span(p).Make
	return protoMaker[Tup2Val[uint64, T], FieldVal[T]]{
		proto: proto[FieldVal[T]]{
			read:  readField(p),
			write: writeField(p),
			size:  sizeField(p),
		},
		makeFunc: func(t Tup2Val[uint64, T]) FieldVal[T] { return FieldVal[T]{t.E0, makeSpan(t.E1)} },
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

package wire

import "fmt"

type FieldVal[T any] Tup2Val[uint64, SpanElem[T]]

func MakeField[T any](proto Proto[T]) func(num uint64, val T) FieldVal[T] {
	makeSpan := MakeSpan(proto)
	return func(num uint64, val T) FieldVal[T] { return FieldVal[T]{E0: num, E1: makeSpan(val)} }
}

func makeAnyField[T any](proto Proto[T]) func(num uint64, val T) FieldVal[any] {
	makeSpan := MakeSpan(proto)
	return func(num uint64, val T) FieldVal[any] {
		span := makeSpan(val)
		return FieldVal[any]{E0: num, E1: SpanElem[any]{span.Size(), span.Elem()}}
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

func Fields(p map[uint64]Proto[any]) Proto[FieldVal[any]] {
	return proto[FieldVal[any]]{
		read:  readFields(p),
		write: writeFields(p),
		size:  sizeFields(p),
	}
}

var anyRawSpan = Span(anyProto[[]byte]{Raw})

func readFields(protos map[uint64]Proto[any]) func(Reader) (FieldVal[any], error) {
	numReader := Uvarint64.Read
	spanReaders := make(map[uint64]func(Reader) (SpanElem[any], error), len(protos))
	for num, proto := range protos {
		spanReaders[num] = readSpan(proto)
	}
	anyRawSpanReader := anyRawSpan.Read
	return func(r Reader) (FieldVal[any], error) {
		num, err := numReader(r)
		if err != nil {
			return FieldVal[any]{}, err
		}
		readField, ok := spanReaders[num]
		if !ok {
			readField = anyRawSpanReader // Unknown fields can be read as Raw.
		}
		span, err := readField(r)
		if err != nil {
			return FieldVal[any]{}, err
		}
		return FieldVal[any]{E0: num, E1: span}, nil
	}
}

func writeFields(proto map[uint64]Proto[any]) func(Writer, FieldVal[any]) error {
	writeSpans := make(map[uint64]func(Writer, SpanElem[any]) error, len(proto))
	for num, fns := range proto {
		writeSpans[num] = writeSpan(fns)
	}
	return func(w Writer, field FieldVal[any]) (err error) {
		writeSpan, ok := writeSpans[field.Num()]
		if !ok {
			return fmt.Errorf("Fields.Write: unknown field: %d", field.Num())
		}
		if err := writeUvarint64(w, field.Num()); err != nil {
			return fmt.Errorf("Fields.Write: %w", err)
		}
		return writeSpan(w, field.E1)
	}
}

func sizeFields(protos map[uint64]Proto[any]) func(FieldVal[any]) uint64 {
	return func(field FieldVal[any]) uint64 {
		p, ok := protos[field.Num()]
		if !ok {
			panic(fmt.Errorf("invalid field: %d", field.Num()))
		}
		return Field(p).Size(field)
	}
}

package wire

import (
	"fmt"
)

var anyRawSpan = Span(Any(Raw))

func RawStruct(fields map[uint64]Proto[any]) ProtoRanger[[]FieldVal[any], FieldVal[any]] {
	return RawSeq(structField(fields))
}

// Struct enables coding based on a field-numbers-to-Proto mapping.
//
// Reading unknown fields will result in a Field of type Raw ([]byte).
// To appease the generics, the Field type T must be made any.
// See Any for help working with Any Protos.
func Struct(fields map[uint64]Proto[any]) ProtoRanger[SpanElem[[]FieldVal[any]], FieldVal[any]] {
	return spanMakeRanger[[]FieldVal[any], FieldVal[any]](RawStruct(fields))
}

func MakeStruct(fields map[uint64]Proto[any]) func([]FieldVal[any]) SpanElem[[]FieldVal[any]] {
	rawStruct := RawStruct(fields)
	return func(fields []FieldVal[any]) SpanElem[[]FieldVal[any]] { return Span(rawStruct).Make(fields) }
}

func structField(p map[uint64]Proto[any]) Proto[FieldVal[any]] {
	return proto[FieldVal[any]]{
		read:  readStructField(p),
		write: writeStructField(p),
		size:  sizeStructField(p),
	}
}

func readStructField(protos map[uint64]Proto[any]) func(Reader) (FieldVal[any], error) {
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

func writeStructField(proto map[uint64]Proto[any]) func(Writer, FieldVal[any]) error {
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

func sizeStructField(protos map[uint64]Proto[any]) func(FieldVal[any]) uint64 {
	return func(field FieldVal[any]) uint64 {
		p, ok := protos[field.Num()]
		if !ok {
			panic(fmt.Errorf("invalid field: %d", field.Num()))
		}
		return Field(p).Size(field)
	}
}

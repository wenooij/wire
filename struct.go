package wire

import (
	"fmt"
	"slices"
)

var anyRawSpan = Span(Any(Raw))

type StructVal map[uint64]SpanElem[any]

// Struct enables coding based on a field-numbers-to-Proto mapping.
//
// Reading unknown fields will result in a Field of type Raw ([]byte).
// To appease the generics, the Field type T must be made any.
// See Any for help working with Any Protos.
func Struct(fields map[uint64]Proto[any]) Proto[StructVal] {
	return proto[StructVal]{
		read:  readStruct(fields),
		write: writeStruct(fields),
		size:  sizeStruct(fields),
	}
}

func readStruct(protos map[uint64]Proto[any]) func(Reader) (StructVal, error) {
	visitSeq := VisitSeq(fieldMapper(protos))
	return func(r Reader) (StructVal, error) {
		s := make(StructVal, 4)
		visitSeq(r, func(field FieldVal[any]) error {
			s[field.Num()] = field.E1
			return nil
		})
		return s, nil
	}
}

func writeStruct(protos map[uint64]Proto[any]) func(Writer, StructVal) error {
	fieldWriter := fieldMapper(protos)
	return func(w Writer, s StructVal) error {
		keys := make([]uint64, 0, len(s))
		for k := range s {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, k := range keys {
			if err := fieldWriter.Write(w, FieldVal[any]{k, s[k]}); err != nil {
				return err
			}
		}
		return nil
	}
}

func sizeStruct(protos map[uint64]Proto[any]) func(StructVal) uint64 {
	return func(s StructVal) uint64 {
		var size uint64
		for k, v := range s {
			proto, ok := protos[k]
			if !ok {
				panic(fmt.Errorf("unknown field: %d", k))
			}
			size += proto.Size(v)
		}
		return size
	}
}

func VisitStruct(fields map[uint64]Proto[any]) func(Reader, func(FieldVal[any]) error) error {
	visitSeq := VisitSeq(fieldMapper(fields))
	return func(r Reader, f func(FieldVal[any]) error) error {
		return visitSeq(r, func(field FieldVal[any]) error {
			if err := f(FieldVal[any]{field.Num(), field.E1}); err != nil {
				if err == ErrStop {
					return nil
				}
				return err
			}
			return nil
		})
	}
}

func fieldMapper(p map[uint64]Proto[any]) Proto[FieldVal[any]] {
	return proto[FieldVal[any]]{
		read:  readFields(p),
		write: writeFields(p),
		size:  sizeFields(p),
	}
}

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

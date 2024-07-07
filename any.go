package wire

import "fmt"

type anyProto[T any] struct{ Proto[T] }

func (p anyProto[T]) Read(r Reader) (any, error)     { return p.Proto.Read(r) }
func (p anyProto[T]) Write(w Writer, elem any) error { return anyWrite(p.Proto.Write)(w, elem) }
func (p anyProto[T]) Size(elem any) uint64           { return anySize(p.Proto.Size)(elem) }

func anyWrite[T any](writeFn func(Writer, T) error) func(Writer, any) error {
	return func(w Writer, v any) error {
		t, ok := v.(T)
		if !ok {
			var t T
			return fmt.Errorf("type mismatch for Any: expected %T but found %T", t, v)
		}
		return writeFn(w, t)
	}
}

func anySize[T any](sizeFn func(T) uint64) func(any) uint64 {
	return func(v any) uint64 {
		t, ok := v.(T)
		if !ok {
			panic(fmt.Errorf("type mismatch for Any: expected %T but found %T", t, v))
		}
		return sizeFn(t)
	}
}

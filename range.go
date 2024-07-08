package wire

import "errors"

// ErrStop can be returned from a ProtoRanger handler to stop the Range as soon as possible.
var ErrStop = errors.New("stop")

type ProtoRanger[T, E any] interface {
	Proto[T]
	Range(Reader, func(E) error) error
}

type protoRanger[T, E any] struct {
	proto[T]
	rangeFunc func(Reader, func(E) error) error
}

func (p protoRanger[T, E]) Range(r Reader, f func(E) error) error {
	if err := p.rangeFunc(r, f); err != nil {
		if err == ErrStop {
			return nil
		}
		return err
	}
	return nil
}

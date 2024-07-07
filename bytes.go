package wire

import (
	"fmt"
	"io"
)

var Raw Proto[[]byte] = RawBufferSize(8)

func RawBufferSize(size uint64) Proto[[]byte] {
	return proto[[]byte]{
		read: readRawBufferSize(size),
		write: func(w Writer, b []byte) error {
			_, err := w.Write(b)
			if err != nil {
				return fmt.Errorf("Raw.Write: %w", err)
			}
			return nil
		},
		size: func(b []byte) uint64 { return uint64(len(b)) },
	}
}

func readRawBufferSize(size uint64) func(Reader) ([]byte, error) {
	return func(r Reader) ([]byte, error) {
		b := make([]byte, 0, size)
		for {
			n, err := r.Read(b[len(b):cap(b)])
			b = b[:len(b)+n]
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
			if len(b) == cap(b) {
				// Add more capacity (let append pick how much).
				b = append(b, 0)[:len(b)]
			}
		}
		// Trim extra capacity and return.
		return b[0:len(b):len(b)], nil
	}
}

var Bytes = Span(Raw)

var makeBytes = MakeSpan(Raw)

func MakeBytes(bs []byte) SpanElem[[]byte] { return makeBytes(bs) }

var RawString Proto[string] = proto[string]{
	read: func(r Reader) (string, error) {
		bs, err := readRawBufferSize(8)(r)
		if err != nil {
			return "", err
		}
		return string(bs), nil
	},
	write: func(w Writer, s string) error {
		if _, err := w.Write([]byte(s)); err != nil {
			return fmt.Errorf("RawString.Write: %w", err)
		}
		return nil
	},
	size: func(s string) uint64 { return uint64(len(s)) },
}

var String = Span(RawString)

var makeString = MakeSpan(RawString)

func MakeString(s string) SpanElem[string] { return makeString(s) }

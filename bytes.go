package wire

import (
	"fmt"
	"io"
)

var Raw Proto[[]byte] = proto[[]byte]{
	read: func(r Reader) ([]byte, error) { return ReadRawBuffer(r, make([]byte, 0, 8)) },
	write: func(w Writer, b []byte) error {
		_, err := w.Write(b)
		if err != nil {
			return fmt.Errorf("Raw.Write: %w", err)
		}
		return nil
	},
	size: func(b []byte) uint64 { return uint64(len(b)) },
}

// ReadRawBuffer reads raw contents from the Reader using buf.
func ReadRawBuffer(r Reader, buf []byte) ([]byte, error) {
	for {
		n, err := r.Read(buf[len(buf):cap(buf)])
		buf = buf[:len(buf)+n]
		if err != nil {
			if err == io.EOF {
				// Trim extra capacity and return.
				return buf[0:len(buf):len(buf)], nil
			}
			return nil, err
		}
		if len(buf) == cap(buf) {
			// Add more capacity (let append pick how much).
			buf = append(buf, 0)[:len(buf)]
		}
	}
}

var Bytes = Span(Raw)

var makeBytes = MakeSpan(Raw)

func MakeBytes(bs []byte) SpanElem[[]byte] { return makeBytes(bs) }

var RawString Proto[string] = proto[string]{
	read: func(r Reader) (string, error) {
		bs, err := ReadRawBuffer(r, make([]byte, 0, 8))
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

package wire

import (
	"fmt"
	"io"
)

var Raw Proto[[]byte] = proto[[]byte]{
	read: func(r Reader) ([]byte, error) { return readRawBuffer(r, make([]byte, 0, 8)) },
	write: func(w Writer, b []byte) error {
		_, err := w.Write(b)
		if err != nil {
			return fmt.Errorf("Raw.Write: %w", err)
		}
		return nil
	},
	size: func(b []byte) uint64 { return uint64(len(b)) },
}

// readRawBuffer reads raw contents from the Reader using buf.
func readRawBuffer(r Reader, buf []byte) ([]byte, error) {
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

var Rune Proto[rune] = proto[rune]{
	read:  readRune,
	write: writeRune,
	size:  sizeRune,
}

func readRune(r Reader) (rune, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	if b < utf8.RuneSelf {
		return rune(b), nil // ASCII
	}
	if !utf8.RuneStart(b) {
		return utf8.RuneError, nil
	}
	// FIXME: Reduce calls to DecodeRune somehow.
	s := make([]byte, 0, 4)
	b1, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	s = append(s, b1)
	if r, _ := utf8.DecodeRune(s); r != utf8.RuneError {
		return r, nil
	}
	b2, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	s = append(s, b2)
	if r, _ := utf8.DecodeRune(s); r != utf8.RuneError {
		return r, nil
	}
	b3, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	s = append(s, b3)
	if r, _ := utf8.DecodeRune(s); r != utf8.RuneError {
		return r, nil
	}
	return utf8.RuneError, nil
}
func writeRune(w Writer, r rune) error {
	if r < utf8.RuneSelf {
		w.WriteByte(byte(r))
	}
	b := make([]byte, utf8.RuneLen(r))
	_, err := w.Write(b[:utf8.EncodeRune(b, r)])
	return err
}
func sizeRune(r rune) uint64 {
	n := utf8.RuneLen(r)
	if n < 0 {
		panic("invalid rune")
	}
	return uint64(n)
}

var RawString Proto[string] = proto[string]{
	read: func(r Reader) (string, error) {
		bs, err := readRawBuffer(r, make([]byte, 0, 8))
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

package wire

import (
	"fmt"
	"io"
	"unicode/utf8"
)

var Raw ProtoMakeRanger[[]byte, []byte, byte] = protoMakeRanger[[]byte, []byte, byte]{
	protoRanger: protoRanger[[]byte, byte]{
		proto: proto[[]byte]{
			read:  readRawBytes,
			write: writeRawBytes,
			size:  func(b []byte) uint64 { return uint64(len(b)) },
		},
		rangeFunc: func(r Reader, f func(byte) error) error { return RawSeq(Fixed8).Range(r, f) },
	},
	makeFunc: func(b []byte) []byte { return b },
}

func readRawBytes(r Reader) ([]byte, error) {
	for buf := make([]byte, 0, 8); ; {
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

func writeRawBytes(w Writer, b []byte) error {
	_, err := w.Write(b)
	if err != nil {
		return fmt.Errorf("Raw.Write: %w", err)
	}
	return nil
}

var Bytes ProtoMakeRanger[[]byte, SpanElem[[]byte], byte] = spanMakeRanger[[]byte, byte](Raw)

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

var RawString ProtoRanger[string, rune] = protoRanger[string, rune]{
	proto: proto[string]{
		read:  readRawString,
		write: func(w Writer, s string) error { return Raw.Write(w, []byte(s)) },
		size:  func(s string) uint64 { return uint64(len(s)) },
	},
	rangeFunc: Seq(Rune).Range,
}

func readRawString(r Reader) (string, error) {
	bs, err := Raw.Read(r)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

var String ProtoMakeRanger[string, SpanElem[string], rune] = spanMakeRanger[string, rune](RawString)

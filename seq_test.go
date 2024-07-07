package wire

import (
	"bytes"
	"testing"
)

func TestWriteSeqFieldUvarint64(t *testing.T) {
	var b bytes.Buffer
	makeField := MakeField(Uvarint64)
	err := Seq(Field(Uvarint64)).Write(&b, []FieldVal[uint64]{
		makeField(1, 1),
		makeField(2, 2),
		makeField(3, 3),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%x\n", b.Bytes())
}

func TestWriteSeqFixedRaw(t *testing.T) {
	var b bytes.Buffer
	err := Seq(Fixed(Raw)(4)).Write(&b, [][]byte{
		[]byte("\x00\x01\x02\x03"),
		[]byte("\x00\x01\x02\x03"),
		[]byte("\x00\x01\x02\x03"),
		[]byte("\x00\x01\x02\x03"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%x\n", b)
}

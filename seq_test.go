package wire

import (
	"bytes"
	"testing"
)

func TestWriteSeqFieldUvarint64(t *testing.T) {
	var b bytes.Buffer
	err := RawSeq(Field(Uvarint64)).Write(&b, []FieldVal[uint64]{
		Field(Uvarint64).Make(Tup2Val[uint64, uint64]{1, 1}),
		Field(Uvarint64).Make(Tup2Val[uint64, uint64]{2, 2}),
		Field(Uvarint64).Make(Tup2Val[uint64, uint64]{3, 3}),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%x\n", b.Bytes())
}

func TestWriteSeqFixedRaw(t *testing.T) {
	var b bytes.Buffer
	err := RawSeq(Fixed(Raw)(4)).Write(&b, [][]byte{
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

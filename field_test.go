package wire

import (
	"bytes"
	"testing"
)

func TestWriteFieldUvarint64(t *testing.T) {
	var b bytes.Buffer
	if err := Field(Uvarint64).Write(&b, Field(Uvarint64).Make(Tup2Val[uint64, uint64]{1, 100})); err != nil {
		t.Fatal(err)
	}
	t.Logf("%x\n", b.Bytes())
}

func TestWriteFieldSeqUvarint64(t *testing.T) {
	var b bytes.Buffer
	if err := Field(Seq(Uvarint64)).Write(&b, Field(Seq(Uvarint64)).Make(Tup2Val[uint64, SpanElem[[]uint64]]{1, Seq(Uvarint64).Make([]uint64{1, 2, 3})})); err != nil {
		t.Fatal(err)
	}
	t.Logf("%x\n", b.Bytes())
}

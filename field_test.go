package wire

import (
	"bytes"
	"testing"
)

func TestWriteFieldUvarint64(t *testing.T) {
	var b bytes.Buffer
	if err := Field(Uvarint64).Write(&b, MakeField(Uvarint64)(100, 100)); err != nil {
		t.Fatal(err)
	}
	t.Logf("%x\n", b.Bytes())
}

func TestWriteFieldSeqUvarint64(t *testing.T) {
	var b bytes.Buffer
	if err := Field(Seq(Uvarint64)).Write(&b, MakeField(Seq(Uvarint64))(100, []uint64{1, 2, 3})); err != nil {
		t.Fatal(err)
	}
	t.Logf("%x\n", b.Bytes())
}

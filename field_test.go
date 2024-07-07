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

func TestWriteFields(t *testing.T) {
	var b bytes.Buffer
	fields := Fields(
		map[uint64]Proto[any]{
			100: Any(RawString),
			200: Any(Uvarint64),
		},
	)
	if err := fields.Write(&b, MakeAnyField(RawString)(100, "abc")); err != nil {
		t.Fatal(err)
	}
	if err := fields.Write(&b, MakeAnyField(Uvarint64)(200, uint64(123))); err != nil {
		t.Fatal(err)
	}

	t.Logf("%x\n", b.Bytes())
}

package wire

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteStruct(t *testing.T) {
	var b bytes.Buffer
	fields := map[uint64]Proto[any]{
		1: Any(RawString),
		2: Any(Uvarint64),
	}
	if err := Struct(fields).Write(&b, MakeStruct(fields)([]FieldVal[any]{
		MakeAnyField(RawString)(1, "a"),
		MakeAnyField(Uvarint64)(2, 123),
	})); err != nil {
		t.Fatal(err)
	}
	got := b.String()
	want := "\x06\x01\x01a\x02\x01{"
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestWriteStruct(): got diff:\n%s", diff)
	}
}

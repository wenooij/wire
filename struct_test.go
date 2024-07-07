package wire

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteStruct(t *testing.T) {
	var b bytes.Buffer
	s := Struct(
		map[uint64]Proto[any]{
			1: Any(RawString),
			2: Any(Uvarint64),
		},
	)
	if err := s.Write(&b, StructVal{
		1: MakeAnySpan(RawString)("a"),
		2: MakeAnySpan(Uvarint64)(uint64(123)),
	}); err != nil {
		t.Fatal(err)
	}
	got := b.String()
	want := "\x01\x01a\x02\x01{"
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestWriteStruct(): got diff:\n%s", diff)
	}
}

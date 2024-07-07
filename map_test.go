package wire

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadMapSet(t *testing.T) {
	m := Map(String, Empty)
	var b bytes.Buffer
	b.WriteString("\x06\x01a\x01b\x01c")
	got, err := m.Read(&b)
	if err != nil {
		t.Fatal(err)
	}
	want := MakeDeterministicMap(String, Empty)(
		map[SpanElem[string]]struct{}{MakeString("a"): {}, MakeString("b"): {}, MakeString("c"): {}},
		[]SpanElem[string]{MakeString("a"), MakeString("b"), MakeString("c")},
	)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestReadMapSet(): got diff (-want, +got):\n%s", diff)
	}
}

func TestWriteMapSet(t *testing.T) {
	m := Map(String, Empty)
	var b bytes.Buffer
	err := m.Write(&b, MakeDeterministicMap(String, Empty)(
		map[SpanElem[string]]struct{}{MakeString("a"): {}, MakeString("b"): {}, MakeString("c"): {}},
		[]SpanElem[string]{MakeString("a"), MakeString("b"), MakeString("c")},
	))
	if err != nil {
		t.Fatal(err)
	}
	want := "\x06\x01a\x01b\x01c"
	got := b.String()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestWriteMapSet(): got diff (-want, +got):\n%s", diff)
	}
}

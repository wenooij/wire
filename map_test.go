package wire

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadMapSet(t *testing.T) {
	var b bytes.Buffer
	b.WriteString("\x06\x01a\x01b\x01c")
	got, err := Map(String, Empty).Read(&b)
	if err != nil {
		t.Fatal(err)
	}
	want := Map(String, Empty).Make([]Tup2Val[SpanElem[string], struct{}]{
		{String.Make("a"), struct{}{}},
		{String.Make("b"), struct{}{}},
		{String.Make("c"), struct{}{}},
	})
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestReadMapSet(): got diff (-want, +got):\n%s", diff)
	}
}

func TestWriteMapSet(t *testing.T) {
	var b bytes.Buffer
	mv := Map(String, Empty).Make([]Tup2Val[SpanElem[string], struct{}]{
		{String.Make("a"), struct{}{}},
		{String.Make("b"), struct{}{}},
		{String.Make("c"), struct{}{}},
	})
	err := Map(String, Empty).Write(&b, mv)
	if err != nil {
		t.Fatal(err)
	}
	want := "\x06\x01a\x01b\x01c"
	got := b.String()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestWriteMapSet(): got diff (-want, +got):\n%s", diff)
	}
}

package wire

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFixedEmpty(t *testing.T) {
	var b bytes.Buffer
	got, err := Empty.Read(&b)
	if err != nil {
		t.Fatal(err)
	}
	want := struct{}{}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestWriteFixedEmpty(): got diff:\n%s", diff)
	}
}

func TestWriteFixedEmpty(t *testing.T) {
	var b bytes.Buffer
	if err := Empty.Write(&b, struct{}{}); err != nil {
		t.Fatal(err)
	}
	if b.Bytes() != nil {
		t.Errorf("TestWriteFixedEmpty(): got %v, want %v", b.Bytes(), nil)
	}
}

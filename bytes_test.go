package wire

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadBytes(t *testing.T) {
	var b bytes.Buffer
	b.WriteString("\x06\x61\x62\x63\x31\x32\x33")
	got, err := Bytes.Read(&b)
	if err != nil {
		t.Fatal(err)
	}
	want := SpanElem[[]byte]{6, []byte("abc123")}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestReadBytes(): got diff (-want, +got):\n%s", diff)
	}
}

func TestWriteBytes(t *testing.T) {
	var b bytes.Buffer
	if err := Bytes.Write(&b, MakeSpan(Raw)([]byte("abc123"))); err != nil {
		t.Fatal(err)
	}
	got := b.Bytes()
	want := []byte("\x06\x61\x62\x63\x31\x32\x33")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestWriteBytes(): got diff (-want, +got):\n%s", diff)
	}
}

func TestSizeBytes(t *testing.T) {
	got := Bytes.Size(MakeBytes([]byte("abc123")))
	want := uint64(7)
	if want != got {
		t.Errorf("TestSizeBytes(): got %d, want %d", got, want)
	}
}

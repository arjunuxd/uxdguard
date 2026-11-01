package gitutil

import (
	"strings"
	"testing"
)

func TestCappedBufferKeepsDataUnderLimit(t *testing.T) {
	b := &cappedBuffer{limit: 32}
	if n, err := b.Write([]byte("hello")); err != nil || n != 5 {
		t.Fatalf("Write = %d, %v", n, err)
	}
	if b.overflow {
		t.Fatal("small write must not set overflow")
	}
	if got := string(b.buf.Bytes()); got != "hello" {
		t.Fatalf("buffer = %q, want %q", got, "hello")
	}
}

func TestCappedBufferDiscardsPastLimit(t *testing.T) {
	b := &cappedBuffer{limit: 4}
	n, err := b.Write([]byte("abcdefgh"))
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if n != 8 {
		t.Fatalf("Write = %d, want full input length 8", n)
	}
	if !b.overflow {
		t.Fatal("overflow must be reported")
	}
	if got := string(b.buf.Bytes()); got != "abcd" {
		t.Fatalf("buffer = %q, want first 4 bytes", got)
	}
	if strings.Contains(string(b.buf.Bytes()), "efgh") {
		t.Fatal("buffer retained bytes past its limit")
	}
}

func TestCappedBufferZeroLimitIsUnlimited(t *testing.T) {
	b := &cappedBuffer{}
	if _, err := b.Write([]byte("unbounded")); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if b.overflow {
		t.Fatal("zero limit must not overflow")
	}
	if got := string(b.buf.Bytes()); got != "unbounded" {
		t.Fatalf("buffer = %q", got)
	}
}

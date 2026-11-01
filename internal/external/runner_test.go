package external

import (
	"strings"
	"testing"
	"time"
)

func TestTimeoutFromConfig(t *testing.T) {
	def := 60 * time.Second
	cases := map[string]time.Duration{
		"":      def,
		"30s":   30 * time.Second,
		"2m":    2 * time.Minute,
		"bogus": def,
		"-10s":  def,
		"0s":    def,
		" 5s ":  5 * time.Second,
	}
	for in, want := range cases {
		if got := timeoutFromConfig(in, def); got != want {
			t.Fatalf("timeoutFromConfig(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestLimitedBufferCapsOutput(t *testing.T) {
	b := &limitedBuffer{limit: 8}
	n, err := b.Write([]byte("0123456789abcdef"))
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if n != 16 {
		t.Fatalf("Write returned %d, want 16 (full input length)", n)
	}
	got := string(b.Bytes())
	if !strings.HasPrefix(got, "01234567") {
		t.Fatalf("buffer did not keep the first 8 bytes: %q", got)
	}
	if !strings.Contains(got, "truncated") {
		t.Fatalf("truncated output should be marked: %q", got)
	}
}

func TestLimitedBufferUnderLimitUnchanged(t *testing.T) {
	b := &limitedBuffer{limit: 32}
	if _, err := b.Write([]byte("hello")); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if got := string(b.Bytes()); got != "hello" {
		t.Fatalf("Bytes() = %q, want %q", got, "hello")
	}
}

func TestLimitedBufferStopsAtZeroRemaining(t *testing.T) {
	b := &limitedBuffer{limit: 4}
	if _, err := b.Write([]byte("abcd")); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if _, err := b.Write([]byte("efgh")); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	got := string(b.Bytes())
	if !strings.HasPrefix(got, "abcd") || strings.Contains(got, "efgh") {
		t.Fatalf("buffer wrote past its limit: %q", got)
	}
}

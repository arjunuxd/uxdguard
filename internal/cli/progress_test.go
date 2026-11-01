package root

import (
	"bytes"
	"strings"
	"testing"
)

func TestProgressSilentForNonTerminalWriter(t *testing.T) {
	var buf bytes.Buffer
	p := newProgress(&buf, true)
	p.Status("Discovering files... %d found", 10)
	p.End()
	if buf.Len() != 0 {
		t.Fatalf("non-terminal progress must stay silent, got %q", buf.String())
	}
}

func TestProgressDisabledProducesNoOutput(t *testing.T) {
	var buf bytes.Buffer
	p := newProgress(&buf, false)
	p.Status("scanning")
	p.End()
	if buf.Len() != 0 {
		t.Fatalf("disabled progress must stay silent, got %q", buf.String())
	}
}

func TestIsTerminalRejectsBuffers(t *testing.T) {
	var buf bytes.Buffer
	if isTerminal(&buf) {
		t.Fatal("bytes.Buffer must not be treated as a terminal")
	}
}

func TestProgressStatusOverwritesPreviousLine(t *testing.T) {
	var buf bytes.Buffer
	p := &progress{out: &buf, interactive: true}
	first := "Discovering files... 1234 found"
	second := "Scanning 3 files"
	p.Status("%s", first)
	p.Status("%s", second)
	p.End()

	out := buf.String()
	if !strings.Contains(out, "\r"+first) {
		t.Fatalf("missing first status line: %q", out)
	}
	// The shorter second line must be padded to erase the previous text, and
	// End must clear the widest line written.
	padding := strings.Repeat(" ", len(first)-len(second))
	if !strings.Contains(out, "\r"+second+padding) {
		t.Fatalf("shorter status did not pad to previous width: %q", out)
	}
	if !strings.Contains(out, "\r"+strings.Repeat(" ", len(first))+"\r") {
		t.Fatalf("End did not clear the status line: %q", out)
	}
}

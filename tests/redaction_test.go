package tests

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arjunuxd/uxdguard/internal/config"
	"github.com/arjunuxd/uxdguard/internal/project"
	"github.com/arjunuxd/uxdguard/internal/report"
	"github.com/arjunuxd/uxdguard/internal/scan"
)

func TestRedactionInJSON(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.txt")
	// Use synthetic AWS key
	if err := os.WriteFile(f, []byte("AKIAIOSFODNN7EXAMPLE\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Defaults()
	files, err := project.Discover(tmp, cfg)
	if err != nil {
		t.Fatal(err)
	}
	ctx := scan.NewContext(tmp, cfg, files)
	res := scan.DefaultRegistry().Run(ctx)
	var buf bytes.Buffer
	jr := &report.JSONReporter{Out: &buf}
	if err := jr.Report(res, tmp, "test"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "AKIAIOSFODNN7EXAMPLE") {
		t.Fatalf("secret should not appear in JSON output")
	}
}

func TestRedactionInTerminal(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.txt")
	if err := os.WriteFile(f, []byte("AKIAIOSFODNN7EXAMPLE\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Defaults()
	files, err := project.Discover(tmp, cfg)
	if err != nil {
		t.Fatal(err)
	}
	ctx := scan.NewContext(tmp, cfg, files)
	res := scan.DefaultRegistry().Run(ctx)
	var outbuf, errbuf bytes.Buffer
	tr := &report.TerminalReporter{Out: &outbuf, Err: &errbuf}
	tr.Report(res, tmp, "test")
	out := outbuf.String() + errbuf.String()
	if strings.Contains(out, "AKIAIOSFODNN7EXAMPLE") {
		t.Fatalf("secret should not appear in terminal output")
	}
}

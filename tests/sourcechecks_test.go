package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arjunuxd/uxdguard/internal/config"
	"github.com/arjunuxd/uxdguard/internal/project"
	"github.com/arjunuxd/uxdguard/internal/scan"
)

func TestSourceCheckTLSInsecure(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "main.go")
	if err := os.WriteFile(f, []byte("x := &tls.Config{InsecureSkipVerify: true}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Defaults()
	files, err := project.Discover(tmp, cfg)
	if err != nil {
		t.Fatal(err)
	}
	ctx := scan.NewContext(tmp, cfg, files)
	res := scan.DefaultRegistry().Run(ctx)
	found := false
	for _, fd := range res.Findings {
		if fd.RuleID == "UXD-SEC-201" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected TLS insecure finding")
	}
}

func TestSourceCheckEval(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "app.js")
	if err := os.WriteFile(f, []byte("eval(userInput)\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Defaults()
	files, err := project.Discover(tmp, cfg)
	if err != nil {
		t.Fatal(err)
	}
	ctx := scan.NewContext(tmp, cfg, files)
	res := scan.DefaultRegistry().Run(ctx)
	found := false
	for _, fd := range res.Findings {
		if fd.RuleID == "UXD-SEC-200" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected eval finding")
	}
}

func TestSourceCheckHTTPInsecureLocalhost(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "main.go")
	if err := os.WriteFile(f, []byte("url := \"http://localhost:8080\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Defaults()
	files, err := project.Discover(tmp, cfg)
	if err != nil {
		t.Fatal(err)
	}
	ctx := scan.NewContext(tmp, cfg, files)
	res := scan.DefaultRegistry().Run(ctx)
	for _, fd := range res.Findings {
		if fd.RuleID == "UXD-SEC-203" {
			t.Fatalf("localhost should not be reported")
		}
	}
}

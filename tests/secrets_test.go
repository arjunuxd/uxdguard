package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arjunuxd/uxdguard/internal/config"
	"github.com/arjunuxd/uxdguard/internal/project"
	"github.com/arjunuxd/uxdguard/internal/scan"
)

func TestSecretDetectionAWS(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.txt")
	if err := os.WriteFile(f, []byte("key=AKIAIOSFODNN7EXAMPLE\n"), 0644); err != nil {
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
		if fd.RuleID == "UXD-SEC-100" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected AWS key finding")
	}
}

func TestSecretDetectionPlaceholder(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.txt")
	if err := os.WriteFile(f, []byte("api_key = \"YOUR_API_KEY\"\n"), 0644); err != nil {
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
		if fd.RuleID == "UXD-SEC-105" {
			t.Fatalf("placeholder should not be reported")
		}
	}
}

func TestSecretDetectionPrivateKey(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "key.pem")
	if err := os.WriteFile(f, []byte("-----BEGIN RSA PRIVATE KEY-----\nabc\n"), 0644); err != nil {
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
		if fd.RuleID == "UXD-SEC-103" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected private key finding")
	}
}

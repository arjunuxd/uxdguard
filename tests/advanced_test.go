package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arjunuxd/uxdguard/internal/config"
	"github.com/arjunuxd/uxdguard/internal/project"
	"github.com/arjunuxd/uxdguard/internal/scan"
)

func TestAdvancedCommandInjection(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.js")
	if err := os.WriteFile(f, []byte("exec('ls', {shell: true})\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Defaults()
	files, _ := project.Discover(tmp, cfg)
	ctx := scan.NewContext(tmp, cfg, files)
	res := scan.DefaultRegistry().Run(ctx)
	found := false
	for _, fd := range res.Findings {
		if fd.RuleID == "UXD-SEC-300" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected command injection finding")
	}
}

func TestAdvancedJWTVerify(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "auth.go")
	if err := os.WriteFile(f, []byte("jwt.verify(token, {algorithms: ['none']})\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Defaults()
	files, _ := project.Discover(tmp, cfg)
	ctx := scan.NewContext(tmp, cfg, files)
	res := scan.DefaultRegistry().Run(ctx)
	found := false
	for _, fd := range res.Findings {
		if fd.RuleID == "UXD-SEC-400" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected JWT verify finding")
	}
}

func TestAdvancedWeakHash(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "crypto.go")
	if err := os.WriteFile(f, []byte("hash := md5.Sum(data)\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := config.Defaults()
	files, _ := project.Discover(tmp, cfg)
	ctx := scan.NewContext(tmp, cfg, files)
	res := scan.DefaultRegistry().Run(ctx)
	found := false
	for _, fd := range res.Findings {
		if fd.RuleID == "UXD-SEC-500" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected weak hash finding")
	}
}

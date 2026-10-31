package tests

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/arjunuxd/uxdguard/internal/config"
	"github.com/arjunuxd/uxdguard/internal/scan"
)

func TestStagedScanUsesIndexContent(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	repo := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	run("init")
	path := filepath.Join(repo, "secret.txt")
	if err := os.WriteFile(path, []byte("AKIAIOSFODNN7EXAMPLE\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run("add", "secret.txt")
	// Replace the working tree with a safe value. Staged scanning must still
	// inspect the indexed content.
	if err := os.WriteFile(path, []byte("redacted\n"), 0644); err != nil {
		t.Fatal(err)
	}

	res, err := scan.ScanStaged(context.Background(), repo, config.Defaults())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, f := range res.Findings {
		if f.RuleID == "UXD-SEC-100" {
			found = true
			if f.File != "secret.txt" {
				t.Fatalf("staged finding path = %q, want secret.txt", f.File)
			}
		}
	}
	if !found {
		t.Fatalf("expected staged secret finding: %+v", res.Findings)
	}
}

package root

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsePositiveInt(t *testing.T) {
	for _, in := range []string{"1", "30", " 7 "} {
		n, err := parsePositiveInt("--page", in)
		if err != nil {
			t.Fatalf("parsePositiveInt(%q) returned error: %v", in, err)
		}
		if n < 1 {
			t.Fatalf("parsePositiveInt(%q) = %d, want positive", in, n)
		}
	}
}

func TestParsePositiveIntRejectsInvalid(t *testing.T) {
	for _, in := range []string{"", "0", "-1", "abc", "1.5", "999999999999999999999999"} {
		if _, err := parsePositiveInt("--limit", in); err == nil {
			t.Fatalf("parsePositiveInt(%q) should fail", in)
		}
	}
}

func TestValidSeverity(t *testing.T) {
	for _, s := range []string{"critical", "high", "medium", "low", "info"} {
		if !validSeverity(s) {
			t.Fatalf("%q should be valid", s)
		}
	}
	for _, s := range []string{"", "HIGH", "severe", "warning"} {
		if validSeverity(s) {
			t.Fatalf("%q should be invalid", s)
		}
	}
}

func TestInitRefusesDanglingConfigSymlink(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.yaml")
	link := filepath.Join(root, ".uxdguard.yaml")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if code := cmdInit([]string{"--project-root", root}); code != 2 {
		t.Fatalf("init exit code = %d, want 2", code)
	}
	if _, err := os.Stat(outside); !os.IsNotExist(err) {
		t.Fatalf("init wrote through the symlink: %v", err)
	}
}

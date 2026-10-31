package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arjunuxd/uxdguard/internal/config"
)

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func containsPath(files []string, want string) bool {
	want = filepath.Clean(want)
	for _, f := range files {
		if filepath.Clean(f) == want {
			return true
		}
	}
	return false
}

func TestLoadWithConfigRejectsMissingProjectRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "missing")
	if _, _, err := LoadWithConfig(root, ""); err == nil {
		t.Fatal("expected error for missing project root")
	}
}

func TestDiscoverSymlinkNotFollowedByDefault(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "real.go"), "package real\n")
	link := filepath.Join(root, "link.go")
	if err := os.Symlink(filepath.Join(root, "real.go"), link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	files, err := Discover(root, config.Defaults())
	if err != nil {
		t.Fatal(err)
	}
	if containsPath(files, link) {
		t.Fatalf("symlink should be skipped when follow_symlinks is false: %v", files)
	}
	if !containsPath(files, filepath.Join(root, "real.go")) {
		t.Fatalf("regular file missing: %v", files)
	}
}

func TestDiscoverFollowsSymlinkInsideRoot(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target.go")
	writeTestFile(t, target, "package target\n")
	link := filepath.Join(root, "link.go")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	cfg := config.Defaults()
	cfg.Scan.FollowSymlinks = true
	files, err := Discover(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !containsPath(files, link) {
		t.Fatalf("expected symlink to be followed: %v", files)
	}
}

func TestDiscoverDoesNotFollowSymlinkOutsideRoot(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "root")
	outside := filepath.Join(base, "outside")
	writeTestFile(t, filepath.Join(root, "main.go"), "package main\n")
	writeTestFile(t, filepath.Join(outside, "secret.go"), "package secret\n")
	link := filepath.Join(root, "external")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	cfg := config.Defaults()
	cfg.Scan.FollowSymlinks = true
	files, err := Discover(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if containsPath(files, filepath.Join(link, "secret.go")) {
		t.Fatalf("symlink outside root must not be scanned: %v", files)
	}
}

func TestDiscoverHandlesSymlinkCycle(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "main.go"), "package main\n")
	if err := os.Symlink(root, filepath.Join(root, "loop")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	cfg := config.Defaults()
	cfg.Scan.FollowSymlinks = true
	files, err := Discover(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !containsPath(files, filepath.Join(root, "main.go")) {
		t.Fatalf("regular file missing: %v", files)
	}
}

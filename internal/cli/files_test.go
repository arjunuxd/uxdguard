package root

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateExclusiveCreatesNewFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new.txt")
	f, err := createExclusive(path, 0600)
	if err != nil {
		t.Fatalf("createExclusive returned error: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close returned error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file was not created: %v", err)
	}
}

func TestCreateExclusiveRefusesExistingFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "existing.txt")
	if err := os.WriteFile(path, []byte("original"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := createExclusive(path, 0600); !os.IsExist(err) {
		t.Fatalf("createExclusive error = %v, want ErrExist", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "original" {
		t.Fatalf("existing file was modified: %q", got)
	}
}

func TestCreateExclusiveRefusesDanglingSymlink(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.txt")
	link := filepath.Join(root, "link.txt")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := createExclusive(link, 0600); !os.IsExist(err) {
		t.Fatalf("createExclusive error = %v, want ErrExist", err)
	}
	if _, err := os.Stat(outside); !os.IsNotExist(err) {
		t.Fatalf("createExclusive wrote through the symlink: %v", err)
	}
}

package scan

import (
	"path/filepath"
	"testing"
)

func TestValidStagedPathAcceptsRepositoryPaths(t *testing.T) {
	for _, in := range []string{"app.js", "src/main.go", "docs/readme.md"} {
		got, ok := validStagedPath(in)
		if !ok {
			t.Fatalf("validStagedPath(%q) rejected a normal path", in)
		}
		if got != filepath.ToSlash(in) {
			t.Fatalf("validStagedPath(%q) = %q", in, got)
		}
	}
}

func TestValidStagedPathRejectsEscapes(t *testing.T) {
	for _, in := range []string{"", ".", "../secret", "../../etc/passwd", "/etc/passwd", ":evil"} {
		if _, ok := validStagedPath(in); ok {
			t.Fatalf("validStagedPath(%q) should be rejected", in)
		}
	}
}

func TestValidStagedPathRejectsWindowsVolume(t *testing.T) {
	if _, ok := validStagedPath(`C:\Windows\system32\drivers\etc\hosts`); ok {
		t.Fatal("absolute Windows path should be rejected")
	}
	if _, ok := validStagedPath(`C:relative`); ok {
		t.Fatal("volume-relative Windows path should be rejected")
	}
	if _, ok := validStagedPath(`src\main.go`); ok {
		t.Fatal("backslash-separated path should be rejected")
	}
}

func TestValidStagedPathRejectsWindowsDeviceNames(t *testing.T) {
	for _, in := range []string{
		"CON", "con", "CON.txt", "CON.", "CON ", "dir/CON",
		"PRN", "AUX", "NUL", "NUL.txt", "COM1", "LPT9", "src/COM1.txt",
		"CONIN$", "CONOUT$", "CON:stream", "NUL:stream", "COM1:stream",
	} {
		if _, ok := validStagedPath(in); ok {
			t.Fatalf("validStagedPath(%q) should reject a Windows device name", in)
		}
	}
	for _, in := range []string{"console.js", "auxiliary.py", "com10.txt", "lpt0.txt", "src/nullable.go"} {
		if _, ok := validStagedPath(in); !ok {
			t.Fatalf("validStagedPath(%q) should accept a normal name", in)
		}
	}
}

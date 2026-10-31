package root

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHookScriptQuotesExecutable(t *testing.T) {
	script := hookScript(`C:\Program Files\UXDGuard\uxdguard.exe`)
	if !strings.Contains(script, "'C:/Program Files/UXDGuard/uxdguard.exe' check --staged") {
		t.Fatalf("executable path is not quoted correctly:\n%s", script)
	}
	if !strings.HasPrefix(script, "#!/bin/sh\n") {
		t.Fatalf("hook must use a POSIX shell shebang:\n%s", script)
	}
}

func TestGetHooksDirHonorsCoreHooksPath(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	repo := t.TempDir()
	runGit := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	runGit("init")
	runGit("config", "core.hooksPath", ".githooks")

	hooks, err := getHooksDir(repo)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(repo, ".githooks")
	if filepath.Clean(hooks) != filepath.Clean(want) {
		t.Fatalf("hooks dir = %q, want %q", hooks, want)
	}
}

package gitutil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// RepoRoot finds git repository root.
func RepoRoot(start string) (string, error) {
	if start == "" {
		var err error
		start, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "-C", start, "rev-parse", "--show-toplevel")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("not a git repository: %s", strings.TrimSpace(stderr.String()))
	}
	root := strings.TrimSpace(out.String())
	if root == "" {
		return "", errors.New("empty git root")
	}
	// Normalize
	root, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	return root, nil
}

// StagedFiles returns NUL-delimited staged file paths.
func StagedFiles(repoRoot string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Use -z for NUL delimiters
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "diff", "--cached", "--name-only", "-z", "--diff-filter=ACDMRTUXB")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git diff failed: %s", strings.TrimSpace(stderr.String()))
	}
	if out.Len() == 0 {
		return []string{}, nil
	}
	// Split by NUL
	parts := bytes.Split(out.Bytes(), []byte{0})
	var files []string
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		// git may output paths
		files = append(files, string(p))
	}
	return files, nil
}

// GetStagedContent reads staged blob content for path (HEAD :path or :path)
func GetStagedContent(repoRoot, relPath string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Use git show :path to read from index
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "show", ":"+relPath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to read staged content: %s", strings.TrimSpace(stderr.String()))
	}
	return out.Bytes(), nil
}

// IsBinaryContent checks for null bytes
func IsBinaryContent(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	// check first 8KB
	check := b
	if len(check) > 8192 {
		check = check[:8192]
	}
	for _, c := range check {
		if c == 0 {
			return true
		}
	}
	return false
}

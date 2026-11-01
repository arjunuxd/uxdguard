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

// HooksDir returns the effective Git hooks directory for a repository,
// honoring core.hooksPath and linked worktrees.
func HooksDir(repoRoot string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "rev-parse", "--git-path", "hooks")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to locate git hooks: %s", strings.TrimSpace(stderr.String()))
	}
	p := strings.TrimSpace(out.String())
	if p == "" {
		return "", errors.New("empty git hooks path")
	}
	if !filepath.IsAbs(p) {
		p = filepath.Join(repoRoot, p)
	}
	return filepath.Abs(p)
}

// StagedFiles returns NUL-delimited staged file paths.
func StagedFiles(repoRoot string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Use -z for NUL delimiters
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "diff", "--cached", "--name-only", "-z", "--diff-filter=ACMRTUXB")
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

// ErrStagedBlobTooLarge is returned when a staged file exceeds the requested
// size limit.
var ErrStagedBlobTooLarge = errors.New("staged blob exceeds size limit")

// GetStagedContent reads staged blob content for path (HEAD :path or :path)
func GetStagedContent(repoRoot, relPath string) ([]byte, error) {
	return getStagedContent(repoRoot, relPath, 0)
}

// GetStagedContentLimited reads staged blob content while keeping at most limit
// bytes in memory. A limit of zero means unlimited. Large blobs are rejected
// with ErrStagedBlobTooLarge instead of exhausting memory.
func GetStagedContentLimited(repoRoot, relPath string, limit int64) ([]byte, error) {
	return getStagedContent(repoRoot, relPath, limit)
}

func getStagedContent(repoRoot, relPath string, limit int64) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Use git show :path to read from index
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "show", ":"+relPath)
	out := &cappedBuffer{limit: limit}
	var stderr bytes.Buffer
	cmd.Stdout = out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to read staged content: %s", strings.TrimSpace(stderr.String()))
	}
	if out.overflow {
		return nil, ErrStagedBlobTooLarge
	}
	return out.buf.Bytes(), nil
}

// cappedBuffer accumulates up to limit bytes and discards the rest so a large
// subprocess stream cannot exhaust memory.
type cappedBuffer struct {
	buf      bytes.Buffer
	limit    int64
	overflow bool
}

func (c *cappedBuffer) Write(p []byte) (int, error) {
	n := len(p)
	if c.limit > 0 {
		remaining := c.limit - int64(c.buf.Len())
		if remaining <= 0 {
			c.overflow = true
			return n, nil
		}
		if int64(len(p)) > remaining {
			_, _ = c.buf.Write(p[:remaining])
			c.overflow = true
			return n, nil
		}
	}
	_, _ = c.buf.Write(p)
	return n, nil
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

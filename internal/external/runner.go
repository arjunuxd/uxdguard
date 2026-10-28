package external

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

// Runner executes external commands safely
type Runner struct {
	Timeout time.Duration
	MaxOut  int64
}

func NewRunner(timeout time.Duration) *Runner {
	if timeout == 0 {
		timeout = 60 * time.Second
	}
	return &Runner{Timeout: timeout, MaxOut: 10 * 1024 * 1024}
}

func (r *Runner) Run(ctx context.Context, name string, args []string) ([]byte, []byte, error) {
	cctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()
	cmd := exec.CommandContext(cctx, name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if errors.Is(cctx.Err(), context.DeadlineExceeded) {
			return stdout.Bytes(), stderr.Bytes(), errors.New("timeout")
		}
		// Return output even on error for some tools
		return stdout.Bytes(), stderr.Bytes(), err
	}
	return stdout.Bytes(), stderr.Bytes(), nil
}

func IsInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "... (truncated)"
}

func RedactOutput(s string) string {
	// Basic redaction - avoid leaking long tokens
	// This is conservative
	lines := strings.Split(s, "\n")
	var out []string
	for _, l := range lines {
		// If line looks like it contains a secret assignment, redact
		ll := strings.ToLower(l)
		if strings.Contains(ll, "token") || strings.Contains(ll, "secret") || strings.Contains(ll, "password") || strings.Contains(ll, "key") {
			if len(l) > 60 {
				out = append(out, l[:20]+"... [REDACTED]")
				continue
			}
		}
		out = append(out, l)
	}
	return strings.Join(out, "\n")
}

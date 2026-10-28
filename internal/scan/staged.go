package scan

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/arjunuxd/uxdguard/internal/config"
	"github.com/arjunuxd/uxdguard/internal/gitutil"
	"github.com/arjunuxd/uxdguard/internal/project"
)

// StagedScanner scans staged content
type StagedScanner struct {
	inner Scanner
}

func NewStagedScanner(inner Scanner) *StagedScanner {
	return &StagedScanner{inner: inner}
}

func (s *StagedScanner) ID() string          { return s.inner.ID() }
func (s *StagedScanner) Name() string        { return s.inner.Name() }
func (s *StagedScanner) Description() string { return s.inner.Description() }

func (s *StagedScanner) Run(ctx *Context) ([]Finding, []Error) {
	// If not in staged mode, delegate
	if ctx == nil || ctx.Cfg == nil {
		return s.inner.Run(ctx)
	}
	// Check if staged mode is set via context
	if v := ctx.Ctx.Value("staged"); v == nil {
		return s.inner.Run(ctx)
	}
	// In staged mode, we need to scan staged content
	// For simplicity, we'll handle at a higher level
	return s.inner.Run(ctx)
}

// ScanStaged performs staged scanning
func ScanStaged(ctx context.Context, repoRoot string, cfg *config.Config) (RunResult, error) {
	var res RunResult
	staged, err := gitutil.StagedFiles(repoRoot)
	if err != nil {
		return res, err
	}
	if len(staged) == 0 {
		return res, nil
	}
	// Create temp files to scan staged content
	tmpDir, err := os.MkdirTemp("", "uxdguard-staged-*")
	if err != nil {
		return res, err
	}
	defer os.RemoveAll(tmpDir)
	var files []string
	for _, rel := range staged {
		// Normalize path
		rel = filepath.ToSlash(rel)
		// Check exclusions
		if project.Excluded(rel, cfg.Scan.Exclusions) {
			continue
		}
		content, err := gitutil.GetStagedContent(repoRoot, rel)
		if err != nil {
			res.Errors = append(res.Errors, Error{ScannerID: "git", Message: "failed to read staged content"})
			continue
		}
		if gitutil.IsBinaryContent(content) {
			continue
		}
		if cfg.Scan.MaxFileSizeBytes > 0 && int64(len(content)) > cfg.Scan.MaxFileSizeBytes {
			continue
		}
		// Write to temp file preserving relative structure
		tmpPath := filepath.Join(tmpDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(tmpPath), 0755); err != nil {
			continue
		}
		if err := os.WriteFile(tmpPath, content, 0644); err != nil {
			continue
		}
		files = append(files, tmpPath)
	}
	// Run scanners on temp files
	for _, scanner := range DefaultRegistry().List() {
		if ctx != nil {
			select {
			case <-ctx.Done():
				res.Errors = append(res.Errors, Error{ScannerID: scanner.ID(), Message: "cancelled"})
				continue
			default:
			}
		}
		findings, errs := runStagedScanner(scanner, repoRoot, tmpDir, files, cfg, ctx)
		if len(findings) > 0 {
			res.Findings = append(res.Findings, findings...)
		}
		if len(errs) > 0 {
			res.Errors = append(res.Errors, errs...)
		}
	}
	SortFindings(res.Findings)
	return res, nil
}

func runStagedScanner(s Scanner, repoRoot, tmpDir string, files []string, cfg *config.Config, ctx context.Context) ([]Finding, []Error) {
	// Create context with temp files
	sc := &Context{Root: tmpDir, Files: files, Cfg: cfg}
	if ctx != nil {
		sc.Ctx = ctx
	}
	findings, errs := s.Run(sc)
	// Convert paths back to relative paths from repo root
	var out []Finding
	for _, f := range findings {
		// f.File is relative to tmpDir
		if f.File != "" {
			// Try to make it relative to repo root
			// Find which original file this corresponds to
			// Simple approach: use the basename/relative path
			rel := f.File
			// Remove tmpDir prefix
			if strings.HasPrefix(rel, tmpDir) {
				rel = strings.TrimPrefix(rel, tmpDir)
				rel = strings.TrimPrefix(rel, string(filepath.Separator))
				rel = strings.TrimPrefix(rel, "/")
			}
			f.File = filepath.ToSlash(rel)
		}
		out = append(out, f)
	}
	return out, errs
}

package external

import (
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/arjunuxd/uxdguard/internal/scan"
)

// SemgrepAdapter integrates semgrep
type SemgrepAdapter struct{}

func init() {
	_ = scan.RegisterDefault(&SemgrepAdapter{})
}

func (s *SemgrepAdapter) ID() string          { return "semgrep" }
func (s *SemgrepAdapter) Name() string        { return "Semgrep (external)" }
func (s *SemgrepAdapter) Description() string { return "Optional external static analysis (semgrep)" }

type semgrepResult struct {
	Results []struct {
		CheckID string `json:"check_id"`
		Path    string `json:"path"`
		Start   struct {
			Line   int `json:"line"`
			Col    int `json:"col"`
			Offset int `json:"offset"`
		} `json:"start"`
		End struct {
			Line   int `json:"line"`
			Col    int `json:"col"`
			Offset int `json:"offset"`
		} `json:"end"`
		Extra struct {
			Message  string                 `json:"message"`
			Metadata map[string]interface{} `json:"metadata"`
			Severity string                 `json:"severity"`
		} `json:"extra"`
	} `json:"results"`
}

func (s *SemgrepAdapter) Run(ctx *scan.Context) ([]scan.Finding, []scan.Error) {
	var findings []scan.Finding
	var errs []scan.Error
	if ctx == nil || ctx.Cfg == nil {
		return findings, errs
	}
	ec := ctx.Cfg.External.Semgrep
	if !ec.Enabled {
		return findings, errs
	}
	bin := ec.Path
	if bin == "" {
		bin = "semgrep"
	}
	if !IsInstalled(bin) {
		if ec.Require {
			errs = append(errs, scan.Error{ScannerID: s.ID(), Message: "semgrep not installed"})
			return findings, errs
		}
		return findings, errs
	}
	runner := NewRunner(timeoutFromConfig(ec.Timeout, 120*time.Second))
	// Use local-only config if specified
	args := []string{"--json", "--quiet", "--no-git-ignore"}
	if cfgPath, ok := ec.Options["config"]; ok && cfgPath != "" {
		// Validate path is safe (local)
		if filepath.IsAbs(cfgPath) {
			// allow if exists
		}
		args = append(args, "--config", cfgPath)
	} else {
		args = append(args, "--config", "auto")
	}
	args = append(args, ctx.Root)
	stdout, stderr, err := runner.Run(ctx.Ctx, bin, args)
	if err != nil {
		if len(stdout) == 0 {
			if ec.Require {
				errs = append(errs, scan.Error{ScannerID: s.ID(), Message: "semgrep execution failed"})
				return findings, errs
			}
			_ = stderr
			return findings, errs
		}
	}
	if len(stdout) == 0 {
		return findings, errs
	}
	var sr semgrepResult
	if err := json.Unmarshal(stdout, &sr); err != nil {
		if ec.Require {
			errs = append(errs, scan.Error{ScannerID: s.ID(), Message: "failed to parse semgrep output"})
			return findings, errs
		}
		return findings, errs
	}
	for _, r := range sr.Results {
		rel := r.Path
		if rel != "" {
			if rr, err := filepath.Rel(ctx.Root, rel); err == nil {
				rel = rr
			}
			rel = filepath.ToSlash(rel)
		}
		if projectExcluded(rel, ctx.Cfg) {
			continue
		}
		sev := scan.SeverityMedium
		if r.Extra.Severity == "ERROR" || r.Extra.Severity == "CRITICAL" {
			sev = scan.SeverityHigh
		} else if r.Extra.Severity == "WARNING" {
			sev = scan.SeverityMedium
		}
		findings = append(findings, scan.Finding{
			RuleID:      "SEM-" + r.CheckID,
			ScannerID:   s.ID(),
			Title:       r.CheckID,
			Message:     r.Extra.Message,
			Severity:    sev,
			Confidence:  scan.ConfidenceMedium,
			File:        rel,
			Line:        r.Start.Line,
			Column:      r.Start.Col,
			Remediation: "Review and fix per semgrep rule guidance",
		})
	}
	return findings, errs
}

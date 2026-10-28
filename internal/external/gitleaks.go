package external

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/arjunuxd/uxdguard/internal/config"
	"github.com/arjunuxd/uxdguard/internal/scan"
)

// GitleaksAdapter integrates gitleaks
type GitleaksAdapter struct{}

func init() {
	_ = scan.RegisterDefault(&GitleaksAdapter{})
}

func (g *GitleaksAdapter) ID() string          { return "gitleaks" }
func (g *GitleaksAdapter) Name() string        { return "Gitleaks (external)" }
func (g *GitleaksAdapter) Description() string { return "Optional external secret scanner (gitleaks)" }

type gitleaksFinding struct {
	Description string   `json:"Description"`
	StartLine   int      `json:"StartLine"`
	EndLine     int      `json:"EndLine"`
	StartColumn int      `json:"StartColumn"`
	EndColumn   int      `json:"EndColumn"`
	Match       string   `json:"Match"`
	Secret      string   `json:"Secret"`
	File        string   `json:"File"`
	SymlinkFile string   `json:"SymlinkFile"`
	Commit      string   `json:"Commit"`
	Entropy     float64  `json:"Entropy"`
	Author      string   `json:"Author"`
	Email       string   `json:"Email"`
	Date        string   `json:"Date"`
	Message     string   `json:"Message"`
	Tags        []string `json:"Tags"`
	RuleID      string   `json:"RuleID"`
	Fingerprint string   `json:"Fingerprint"`
}

func (g *GitleaksAdapter) Run(ctx *scan.Context) ([]scan.Finding, []scan.Error) {
	var findings []scan.Finding
	var errs []scan.Error
	if ctx == nil || ctx.Cfg == nil {
		return findings, errs
	}
	ec := ctx.Cfg.External.Gitleaks
	if !ec.Enabled {
		return findings, errs
	}
	bin := ec.Path
	if bin == "" {
		bin = "gitleaks"
	}
	if !IsInstalled(bin) {
		if ec.Require {
			errs = append(errs, scan.Error{ScannerID: g.ID(), Message: "gitleaks not installed"})
			return findings, errs
		}
		return findings, errs
	}
	runner := NewRunner(60)
	// Run gitleaks detect in repo
	args := []string{"detect", "--source", ctx.Root, "--report-format", "json", "--report-path", "-", "--no-git", "--no-banner"}
	stdout, stderr, err := runner.Run(ctx.Ctx, bin, args)
	if err != nil {
		// Some versions might behave differently; check if we got output
		if len(stdout) == 0 {
			if ec.Require {
				errs = append(errs, scan.Error{ScannerID: g.ID(), Message: "gitleaks execution failed"})
				return findings, errs
			}
			// Don't fail if optional
			_ = stderr
			return findings, errs
		}
	}
	if len(stdout) == 0 {
		return findings, errs
	}
	var glFindings []gitleaksFinding
	if err := json.Unmarshal(stdout, &glFindings); err != nil {
		if ec.Require {
			errs = append(errs, scan.Error{ScannerID: g.ID(), Message: "failed to parse gitleaks output"})
			return findings, errs
		}
		return findings, errs
	}
	for _, gf := range glFindings {
		rel := gf.File
		if rel != "" {
			if r, err := filepath.Rel(ctx.Root, rel); err == nil {
				rel = r
			}
			rel = filepath.ToSlash(rel)
		}
		// Skip if outside repo or excluded
		if projectExcluded(rel, ctx.Cfg) {
			continue
		}
		findings = append(findings, scan.Finding{
			RuleID:      "GL-" + gf.RuleID,
			ScannerID:   g.ID(),
			Title:       gf.Description,
			Message:     "Secret detected by gitleaks",
			Severity:    scan.SeverityHigh,
			Confidence:  scan.ConfidenceHigh,
			File:        rel,
			Line:        gf.StartLine,
			Column:      gf.StartColumn,
			Remediation: "Remove hardcoded secret and rotate if exposed",
		})
	}
	return findings, errs
}

func projectExcluded(rel string, cfg *config.Config) bool {
	if cfg == nil {
		return false
	}
	// Use project.Excluded if possible
	// Import cycle? project imports config, scan imports config - we need to check
	for _, p := range cfg.Scan.Exclusions {
		pn := strings.ReplaceAll(p, "**", "")
		if pn == "" {
			continue
		}
		if strings.Contains(rel, pn) {
			return true
		}
	}
	return false
}

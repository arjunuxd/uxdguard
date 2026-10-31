package scan

import (
	"fmt"
	"path/filepath"
	"sort"
	"sync"

	"github.com/arjunuxd/uxdguard/internal/config"
)

// Registry holds scanners.
type Registry struct {
	mu       sync.RWMutex
	scanners map[string]Scanner
}

// NewRegistry creates a new registry.
func NewRegistry() *Registry {
	return &Registry{
		scanners: make(map[string]Scanner),
	}
}

// Register registers a scanner.
func (r *Registry) Register(s Scanner) error {
	if s == nil {
		return fmt.Errorf("nil scanner")
	}
	id := s.ID()
	if id == "" {
		return fmt.Errorf("empty scanner id")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.scanners[id]; exists {
		return fmt.Errorf("duplicate scanner id: %s", id)
	}
	r.scanners[id] = s
	return nil
}

// List returns all scanners in deterministic order by ID.
func (r *Registry) List() []Scanner {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Scanner, 0, len(r.scanners))
	for _, s := range r.scanners {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// RunResult contains results from a scan run.
type RunResult struct {
	Findings []Finding
	Errors   []Error
	Skipped  []Skipped
}

// Skipped represents a skipped check.
type Skipped struct {
	ScannerID string `json:"scanner_id"`
	Reason    string `json:"reason"`
}

// Run executes all enabled scanners and applies rule configuration.
func (r *Registry) Run(ctx *Context) RunResult {
	res := RunResult{
		Findings: []Finding{},
		Errors:   []Error{},
		Skipped:  []Skipped{},
	}
	if ctx == nil {
		res.Errors = append(res.Errors, Error{ScannerID: "registry", Message: "nil scan context"})
		return res
	}

	for _, s := range r.List() {
		if ctx.Ctx != nil {
			select {
			case <-ctx.Ctx.Done():
				res.Errors = append(res.Errors, Error{ScannerID: s.ID(), Message: "cancelled"})
				continue
			default:
			}
		}
		if !isRuleEnabled(s.ID(), ctx.Cfg) {
			res.Skipped = append(res.Skipped, Skipped{ScannerID: s.ID(), Reason: "disabled by configuration"})
			continue
		}
		findings, errs := s.Run(ctx)
		if len(findings) > 0 {
			res.Findings = append(res.Findings, findings...)
		}
		if len(errs) > 0 {
			res.Errors = append(res.Errors, errs...)
		}
	}

	res.Findings = applyRuleConfig(res.Findings, ctx.Cfg)
	SortFindings(res.Findings)
	res.Findings = DeduplicateFindings(res.Findings)
	SortFindings(res.Findings)
	for i := range res.Findings {
		res.Findings[i].File = filepath.ToSlash(res.Findings[i].File)
	}
	return res
}

var defaultRegistry = NewRegistry()

// DefaultRegistry returns the global registry.
func DefaultRegistry() *Registry {
	return defaultRegistry
}

// RegisterDefault registers a scanner in the default registry.
func RegisterDefault(s Scanner) error {
	return defaultRegistry.Register(s)
}

// isRuleEnabled reports whether a scanner is enabled. Rule-level enable/disable
// filtering is applied to findings in applyRuleConfig.
func isRuleEnabled(scannerID string, cfg *config.Config) bool {
	if cfg == nil {
		return true
	}
	for _, disabled := range cfg.Rules.Disable {
		if disabled == scannerID {
			return false
		}
	}
	return true
}

// applyRuleConfig filters findings by rules.enable/rules.disable and applies
// rules.severity overrides. Entries may be scanner IDs or individual rule IDs.
func applyRuleConfig(findings []Finding, cfg *config.Config) []Finding {
	if cfg == nil {
		return findings
	}

	disabled := make(map[string]bool, len(cfg.Rules.Disable))
	for _, id := range cfg.Rules.Disable {
		disabled[id] = true
	}
	enabled := make(map[string]bool, len(cfg.Rules.Enable))
	for _, id := range cfg.Rules.Enable {
		enabled[id] = true
	}

	out := make([]Finding, 0, len(findings))
	for _, f := range findings {
		if disabled[f.ScannerID] || disabled[f.RuleID] {
			continue
		}
		if len(enabled) > 0 && !enabled[f.ScannerID] && !enabled[f.RuleID] {
			continue
		}
		if severity, ok := cfg.Rules.Severity[f.RuleID]; ok {
			f.Severity = Severity(severity)
		}
		out = append(out, f)
	}
	return out
}

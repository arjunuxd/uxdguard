package scan

import (
	"fmt"
	"sync"
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
	var out []Scanner
	for _, s := range r.scanners {
		out = append(out, s)
	}
	// sort
	// simple sort by id
	type pair struct {
		id string
		s  Scanner
	}
	var pairs []pair
	for _, s := range out {
		pairs = append(pairs, pair{id: s.ID(), s: s})
	}
	// sort
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[i].id > pairs[j].id {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}
	res := make([]Scanner, len(pairs))
	for i, p := range pairs {
		res[i] = p.s
	}
	return res
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

// Run executes all scanners.
func (r *Registry) Run(ctx *Context) RunResult {
	var res RunResult
	for _, s := range r.List() {
		if ctx.Ctx != nil {
			select {
			case <-ctx.Ctx.Done():
				res.Errors = append(res.Errors, Error{ScannerID: s.ID(), Message: "cancelled"})
				continue
			default:
			}
		}
		findings, errs := s.Run(ctx)
		if len(findings) > 0 {
			res.Findings = append(res.Findings, findings...)
		}
		if len(errs) > 0 {
			res.Errors = append(res.Errors, errs...)
		}
	}
	SortFindings(res.Findings)
	res.Findings = DeduplicateFindings(res.Findings)
	SortFindings(res.Findings)
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

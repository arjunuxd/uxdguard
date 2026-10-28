package scan

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

type Baseline struct {
	Version  string          `json:"version"`
	Findings []BaselineEntry `json:"findings"`
}

type BaselineEntry struct {
	Fingerprint string   `json:"fingerprint"`
	RuleID      string   `json:"rule_id"`
	File        string   `json:"file"`
	Title       string   `json:"title"`
	Severity    Severity `json:"severity"`
}

func Fingerprint(f Finding) string {
	h := sha256.New()
	h.Write([]byte(f.RuleID + "|" + filepath.ToSlash(f.File) + "|" + f.Title))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

func LoadBaseline(path string) (*Baseline, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var bl Baseline
	if err := json.Unmarshal(b, &bl); err != nil {
		return nil, err
	}
	return &bl, nil
}

func SaveBaseline(path string, res RunResult) error {
	bl := Baseline{Version: "1.0", Findings: []BaselineEntry{}}
	for _, f := range res.Findings {
		bl.Findings = append(bl.Findings, BaselineEntry{
			Fingerprint: Fingerprint(f),
			RuleID:      f.RuleID,
			File:        filepath.ToSlash(f.File),
			Title:       f.Title,
			Severity:    f.Severity,
		})
	}
	b, _ := json.MarshalIndent(bl, "", "  ")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func FilterBaseline(res RunResult, bl *Baseline) RunResult {
	if bl == nil {
		return res
	}
	seen := map[string]bool{}
	for _, e := range bl.Findings {
		seen[e.Fingerprint] = true
	}
	var out RunResult
	for _, f := range res.Findings {
		if seen[Fingerprint(f)] {
			continue
		}
		out.Findings = append(out.Findings, f)
	}
	out.Errors = res.Errors
	out.Skipped = res.Skipped
	return out
}

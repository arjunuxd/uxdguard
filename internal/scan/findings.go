package scan

import "sort"

// Severity represents finding severity.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// Confidence represents confidence level.
type Confidence string

const (
	ConfidenceHigh   Confidence = "high"
	ConfidenceMedium Confidence = "medium"
	ConfidenceLow    Confidence = "low"
)

// Finding represents a scan finding.
type Finding struct {
	RuleID      string     `json:"rule_id"`
	ScannerID   string     `json:"scanner_id"`
	Title       string     `json:"title"`
	Message     string     `json:"message"`
	Severity    Severity   `json:"severity"`
	Confidence  Confidence `json:"confidence,omitempty"`
	File        string     `json:"file,omitempty"`
	Line        int        `json:"line,omitempty"`
	Column      int        `json:"column,omitempty"`
	Remediation string     `json:"remediation,omitempty"`
}

// SortFindings sorts findings deterministically.
func SortFindings(findings []Finding) {
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Severity != findings[j].Severity {
			order := map[Severity]int{SeverityCritical: 0, SeverityHigh: 1, SeverityMedium: 2, SeverityLow: 3, SeverityInfo: 4}
			if order[findings[i].Severity] != order[findings[j].Severity] {
				return order[findings[i].Severity] < order[findings[j].Severity]
			}
		}
		if findings[i].File != findings[j].File {
			return findings[i].File < findings[j].File
		}
		if findings[i].Line != findings[j].Line {
			return findings[i].Line < findings[j].Line
		}
		if findings[i].RuleID != findings[j].RuleID {
			return findings[i].RuleID < findings[j].RuleID
		}
		return findings[i].Title < findings[j].Title
	})
}

// DeduplicateFindings removes duplicate findings based on stable attributes
func DeduplicateFindings(findings []Finding) []Finding {
	seen := map[string]bool{}
	var out []Finding
	for _, f := range findings {
		key := string(f.Severity) + "|" + f.RuleID + "|" + f.File + "|" + strconvItoa(f.Line) + "|" + strconvItoa(f.Column) + "|" + f.Title
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, f)
	}
	return out
}

func strconvItoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var buf [32]byte
	b := len(buf)
	for i > 0 {
		b--
		buf[b] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		b--
		buf[b] = '-'
	}
	return string(buf[b:])
}

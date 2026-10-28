package report

import (
	"encoding/json"
	"io"

	"github.com/arjunuxd/uxdguard/internal/scan"
)

// JSONReport is the versioned JSON output.
type JSONReport struct {
	Tool     ToolInfo       `json:"tool"`
	Project  string         `json:"project"`
	Status   string         `json:"status"`
	Findings []scan.Finding `json:"findings"`
	Skipped  []scan.Skipped `json:"skipped"`
	Errors   []scan.Error   `json:"errors"`
	Summary  Summary        `json:"summary"`
}

// ToolInfo contains tool metadata.
type ToolInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Summary contains counts.
type Summary struct {
	Findings int `json:"findings"`
	Skipped  int `json:"skipped"`
	Errors   int `json:"errors"`
}

// JSONReporter writes JSON output.
type JSONReporter struct {
	Out io.Writer
}

// Report writes JSON report.
func (j *JSONReporter) Report(res scan.RunResult, root string, version string) error {
	status := "ok"
	if len(res.Errors) > 0 {
		status = "error"
	} else if len(res.Findings) > 0 {
		status = "findings"
	}
	rep := JSONReport{
		Tool:     ToolInfo{Name: "uxdguard", Version: version},
		Project:  root,
		Status:   status,
		Findings: res.Findings,
		Skipped:  res.Skipped,
		Errors:   res.Errors,
		Summary: Summary{
			Findings: len(res.Findings),
			Skipped:  len(res.Skipped),
			Errors:   len(res.Errors),
		},
	}
	enc := json.NewEncoder(j.Out)
	enc.SetIndent("", "  ")
	return enc.Encode(rep)
}

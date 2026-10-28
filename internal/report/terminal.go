package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/arjunuxd/uxdguard/internal/scan"
)

// TerminalReporter prints human readable output.
type TerminalReporter struct {
	Out io.Writer
	Err io.Writer
}

// Report writes terminal output.
func (t *TerminalReporter) Report(res scan.RunResult, root string, version string) {
	fmt.Fprintf(t.Out, "uxdguard %s\n", version)
	fmt.Fprintf(t.Out, "Project: %s\n\n", root)
	if len(res.Findings) == 0 && len(res.Errors) == 0 && len(res.Skipped) == 0 {
		fmt.Fprintf(t.Out, "No issues found.\n")
		return
	}
	if len(res.Findings) > 0 {
		fmt.Fprintf(t.Out, "Findings:\n")
		for _, f := range res.Findings {
			loc := f.File
			if f.Line > 0 {
				loc = fmt.Sprintf("%s:%d", f.File, f.Line)
			}
			fmt.Fprintf(t.Out, "- [%s] %s (%s) %s\n", strings.ToUpper(string(f.Severity)), f.Title, f.RuleID, loc)
			if f.Message != "" {
				fmt.Fprintf(t.Out, "  %s\n", f.Message)
			}
			if f.Remediation != "" {
				fmt.Fprintf(t.Out, "  Remediation: %s\n", f.Remediation)
			}
			fmt.Fprintf(t.Out, "\n")
		}
	}
	if len(res.Skipped) > 0 {
		fmt.Fprintf(t.Out, "Skipped:\n")
		for _, s := range res.Skipped {
			fmt.Fprintf(t.Out, "- %s: %s\n", s.ScannerID, s.Reason)
		}
		fmt.Fprintf(t.Out, "\n")
	}
	if len(res.Errors) > 0 {
		fmt.Fprintf(t.Err, "Errors:\n")
		for _, e := range res.Errors {
			fmt.Fprintf(t.Err, "- %s: %s\n", e.ScannerID, e.Message)
		}
		fmt.Fprintf(t.Err, "\n")
	}
	fmt.Fprintf(t.Out, "Summary: %d findings, %d skipped, %d errors\n", len(res.Findings), len(res.Skipped), len(res.Errors))
}

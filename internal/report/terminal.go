package report

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/arjunuxd/uxdguard/internal/scan"
)

var (
	red    = "\033[31m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	gray   = "\033[90m"
	reset  = "\033[0m"
	bold   = "\033[1m"
)

// sanitizeOutput removes terminal control characters from untrusted values
// such as repository file names and tool messages, preventing escape-sequence
// injection into the user's terminal.
func sanitizeOutput(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\t' {
			return r
		}
		if unicode.IsControl(r) {
			return '?'
		}
		return r
	}, s)
}

func colorEnabled(w io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if f, ok := w.(*os.File); ok {
		fi, err := f.Stat()
		if err != nil {
			return false
		}
		return (fi.Mode() & os.ModeCharDevice) != 0
	}
	return false
}

// ReportOptions controls pagination and filtering for terminal output.
type ReportOptions struct {
	Page     int
	Limit    int
	Severity string
	Rule     string
}

// TerminalReporter prints human readable output.
type TerminalReporter struct {
	Out io.Writer
	Err io.Writer
}

// Report writes terminal output.
func (t *TerminalReporter) Report(res scan.RunResult, root string, version string) {
	t.ReportWithOptions(res, root, version, ReportOptions{Page: 1, Limit: 30})
}

// ReportWithOptions writes terminal output with pagination.
func (t *TerminalReporter) ReportWithOptions(res scan.RunResult, root string, version string, opts ReportOptions) {
	if opts.Limit <= 0 {
		opts.Limit = 30
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}
	useColor := colorEnabled(t.Out)

	if useColor {
		fmt.Fprintf(t.Out, "%sUXDGuard%s\n", bold, reset)
		fmt.Fprintf(t.Out, "%sSecurity Scanner · %s%s\n\n", gray, version, reset)
	} else {
		fmt.Fprintf(t.Out, "UXDGuard\nSecurity Scanner · %s\n\n", version)
	}
	fmt.Fprintf(t.Out, "Project: %s\n\n", sanitizeOutput(root))

	if len(res.Findings) == 0 && len(res.Errors) == 0 {
		if useColor {
			fmt.Fprintf(t.Out, "%sNo findings detected.%s\n", gray, reset)
		} else {
			fmt.Fprintf(t.Out, "No findings detected.\n")
		}
		fmt.Fprintf(t.Out, "Scan completed with no issues found.\n")
		if len(res.Skipped) > 0 {
			fmt.Fprintf(t.Out, "\nSkipped checks: %d\n", len(res.Skipped))
		}
		fmt.Fprintf(t.Out, "\nStatus: OK\n")
		return
	}

	if len(res.Findings) > 0 {
		// Apply filtering
		filtered := make([]scan.Finding, 0, len(res.Findings))
		for _, f := range res.Findings {
			if opts.Severity != "" {
				if !strings.EqualFold(string(f.Severity), opts.Severity) {
					continue
				}
			}
			if opts.Rule != "" {
				if !strings.EqualFold(f.RuleID, opts.Rule) {
					continue
				}
			}
			filtered = append(filtered, f)
		}
		total := len(filtered)
		if total == 0 {
			if useColor {
				fmt.Fprintf(t.Out, "%sFindings%s (0 of %d total match the selected filters)\n\n", bold, reset, len(res.Findings))
			} else {
				fmt.Fprintf(t.Out, "Findings: 0 of %d total match the selected filters\n\n", len(res.Findings))
			}
		} else {
			totalPages := (total + opts.Limit - 1) / opts.Limit
			requestedPage := opts.Page
			if opts.Page > totalPages {
				opts.Page = totalPages
			}
			start := (opts.Page - 1) * opts.Limit
			end := start + opts.Limit
			if end > total {
				end = total
			}
			if useColor {
				fmt.Fprintf(t.Out, "%sFindings%s (showing %d-%d of %d total)%s\n", bold, reset, start+1, end, total, reset)
				fmt.Fprintf(t.Out, "%sPage %d of %d (%d per page)%s\n\n", gray, opts.Page, totalPages, opts.Limit, reset)
			} else {
				fmt.Fprintf(t.Out, "Findings: (showing %d-%d of %d total)\n", start+1, end, total)
				fmt.Fprintf(t.Out, "Page %d of %d (%d per page)\n\n", opts.Page, totalPages, opts.Limit)
			}
			if requestedPage > totalPages {
				fmt.Fprintf(t.Out, "Requested page %d, but only %d page(s) available. Showing page %d of %d.\n\n", requestedPage, totalPages, opts.Page, totalPages)
			}
			t.renderFindings(filtered[start:end], useColor)
			if totalPages > 1 {
				hint := ""
				if opts.Page > 1 {
					hint = fmt.Sprintf("uxdguard check --page %d --limit %d", opts.Page-1, opts.Limit)
					if opts.Page < totalPages {
						hint += "; next: " + fmt.Sprintf("uxdguard check --page %d --limit %d", opts.Page+1, opts.Limit)
					}
				} else {
					hint = fmt.Sprintf("uxdguard check --page %d --limit %d", opts.Page+1, opts.Limit)
				}
				hint += pageFilterSuffix(opts)
				fmt.Fprintf(t.Out, "Pages: %d total. Navigate: %s\n\n", totalPages, hint)
			}
		}
	}

	if len(res.Skipped) > 0 {
		if useColor {
			fmt.Fprintf(t.Out, "%sSkipped%s\n", bold, reset)
		} else {
			fmt.Fprintf(t.Out, "Skipped:\n")
		}
		for _, s := range res.Skipped {
			fmt.Fprintf(t.Out, "- %s: %s\n", sanitizeOutput(s.ScannerID), sanitizeOutput(s.Reason))
		}
		fmt.Fprintf(t.Out, "\n")
	}
	if len(res.Errors) > 0 {
		if useColor {
			fmt.Fprintf(t.Err, "%sErrors%s\n", bold+red, reset)
		} else {
			fmt.Fprintf(t.Err, "Errors:\n")
		}
		for _, e := range res.Errors {
			fmt.Fprintf(t.Err, "- %s: %s\n", sanitizeOutput(e.ScannerID), sanitizeOutput(e.Message))
		}
		fmt.Fprintf(t.Err, "\n")
	}

	crit, high, med, low := countSeverities(res.Findings)
	if useColor {
		fmt.Fprintf(t.Out, "%sSummary:%s %d critical, %d high, %d medium, %d low | %d findings, %d skipped, %d errors\n",
			bold, reset, crit, high, med, low, len(res.Findings), len(res.Skipped), len(res.Errors))
	} else {
		fmt.Fprintf(t.Out, "Summary: %d critical, %d high, %d medium, %d low | %d findings, %d skipped, %d errors\n",
			crit, high, med, low, len(res.Findings), len(res.Skipped), len(res.Errors))
	}

	if len(res.Errors) > 0 {
		fmt.Fprintf(t.Out, "Status: Scan completed with errors\n")
	} else if len(res.Findings) > 0 {
		fmt.Fprintf(t.Out, "Status: Findings detected\n")
	} else {
		fmt.Fprintf(t.Out, "Status: OK\n")
	}
}

// pageFilterSuffix appends the active filters to a navigation hint.
func pageFilterSuffix(opts ReportOptions) string {
	suffix := ""
	if opts.Severity != "" {
		suffix += " --severity " + sanitizeOutput(opts.Severity)
	}
	if opts.Rule != "" {
		suffix += " --rule " + sanitizeOutput(opts.Rule)
	}
	return suffix
}

// renderFindings prints one page of findings.
func (t *TerminalReporter) renderFindings(findings []scan.Finding, useColor bool) {
	for _, f := range findings {
		loc := f.File
		if f.Line > 0 {
			loc = fmt.Sprintf("%s:%d", f.File, f.Line)
		}
		loc = sanitizeOutput(loc)
		ruleID := sanitizeOutput(f.RuleID)
		title := sanitizeOutput(f.Title)
		message := sanitizeOutput(f.Message)
		remediation := sanitizeOutput(f.Remediation)
		sev := strings.ToUpper(sanitizeOutput(string(f.Severity)))
		if useColor {
			switch sev {
			case "CRITICAL":
				fmt.Fprintf(t.Out, "%s%-8s%s %s%s%s  %s\n", red+bold, sev, reset, bold, ruleID, reset, title)
			case "HIGH":
				fmt.Fprintf(t.Out, "%s%-8s%s %s%s%s  %s\n", red, sev, reset, bold, ruleID, reset, title)
			case "MEDIUM":
				fmt.Fprintf(t.Out, "%s%-8s%s %s%s%s  %s\n", yellow, sev, reset, bold, ruleID, reset, title)
			case "LOW":
				fmt.Fprintf(t.Out, "%s%-8s%s %s%s%s  %s\n", blue, sev, reset, bold, ruleID, reset, title)
			default:
				fmt.Fprintf(t.Out, "%s%-8s%s %s%s%s  %s\n", gray, sev, reset, bold, ruleID, reset, title)
			}
		} else {
			fmt.Fprintf(t.Out, "%-8s %s  %s\n", sev, ruleID, title)
		}
		if useColor {
			fmt.Fprintf(t.Out, "         %s%s%s\n", gray, loc, reset)
		} else {
			fmt.Fprintf(t.Out, "         %s\n", loc)
		}
		if message != "" {
			fmt.Fprintf(t.Out, "         %s\n", message)
		}
		if remediation != "" {
			fmt.Fprintf(t.Out, "         Remediation: %s\n", remediation)
		}
		fmt.Fprintf(t.Out, "\n")
	}
}

func countSeverities(findings []scan.Finding) (crit, high, med, low int) {
	for _, f := range findings {
		switch strings.ToLower(string(f.Severity)) {
		case "critical":
			crit++
		case "high":
			high++
		case "medium":
			med++
		case "low":
			low++
		}
	}
	return
}

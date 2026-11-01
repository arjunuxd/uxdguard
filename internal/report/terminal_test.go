package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/arjunuxd/uxdguard/internal/scan"
)

func findings(n int) []scan.Finding {
	out := make([]scan.Finding, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, scan.Finding{
			RuleID:   "UXD-SEC-100",
			Title:    "Potential AWS access key detected",
			Severity: scan.SeverityHigh,
			File:     "file.js",
			Line:     i + 1,
		})
	}
	return out
}

func reportText(t *testing.T, res scan.RunResult, opts ReportOptions) string {
	t.Helper()
	var buf bytes.Buffer
	tr := &TerminalReporter{Out: &buf, Err: &buf}
	tr.ReportWithOptions(res, "/project", "1.1.0", opts)
	return buf.String()
}

func TestDefaultPageLimitIs30(t *testing.T) {
	out := reportText(t, scan.RunResult{Findings: findings(65)}, ReportOptions{Page: 1})
	if !strings.Contains(out, "showing 1-30 of 65 total") {
		t.Fatalf("default limit is not 30:\n%s", out)
	}
	if !strings.Contains(out, "Page 1 of 3 (30 per page)") {
		t.Fatalf("missing page count:\n%s", out)
	}
}

func TestPaginationShowsTotalPagesAndNavigation(t *testing.T) {
	out := reportText(t, scan.RunResult{Findings: findings(65)}, ReportOptions{Page: 2, Limit: 30})
	if !strings.Contains(out, "showing 31-60 of 65 total") {
		t.Fatalf("wrong slice for page 2:\n%s", out)
	}
	if !strings.Contains(out, "Page 2 of 3 (30 per page)") {
		t.Fatalf("missing page count:\n%s", out)
	}
	if !strings.Contains(out, "Pages: 3 total.") {
		t.Fatalf("missing total pages hint:\n%s", out)
	}
	if !strings.Contains(out, "uxdguard check --page 3 --limit 30") {
		t.Fatalf("missing next-page command:\n%s", out)
	}
	if !strings.Contains(out, "uxdguard check --page 1 --limit 30") {
		t.Fatalf("missing previous-page command:\n%s", out)
	}
}

func TestPaginationReportsClampedPage(t *testing.T) {
	out := reportText(t, scan.RunResult{Findings: findings(65)}, ReportOptions{Page: 99, Limit: 30})
	if !strings.Contains(out, "Requested page 99, but only 3 page(s) available. Showing page 3 of 3.") {
		t.Fatalf("missing clamp notice:\n%s", out)
	}
}

func TestPaginationKeepsFiltersInNavigation(t *testing.T) {
	out := reportText(t, scan.RunResult{Findings: findings(65)}, ReportOptions{
		Page:     1,
		Limit:    30,
		Severity: "high",
		Rule:     "UXD-SEC-100",
	})
	if !strings.Contains(out, "--severity high --rule UXD-SEC-100") {
		t.Fatalf("navigation lost filters:\n%s", out)
	}
}

func TestPaginationZeroMatchesShowsTotal(t *testing.T) {
	out := reportText(t, scan.RunResult{Findings: findings(5)}, ReportOptions{Page: 1, Limit: 30, Severity: "critical"})
	if !strings.Contains(out, "0 of 5 total match the selected filters") {
		t.Fatalf("missing filtered-count message:\n%s", out)
	}
}

func TestTerminalOutputRemovesControlCharacters(t *testing.T) {
	res := scan.RunResult{Findings: []scan.Finding{{
		RuleID:   "UXD-SEC-100",
		Title:    "evil\x1b[31m title",
		Message:  "message\x00 with null",
		Severity: scan.SeverityHigh,
		File:     "evil\x1b[2Jname.js",
		Line:     1,
	}}}
	out := reportText(t, res, ReportOptions{Page: 1, Limit: 30})
	if strings.ContainsRune(out, '\x1b') {
		t.Fatalf("terminal output contains an escape character: %q", out)
	}
	if strings.ContainsRune(out, '\x00') {
		t.Fatalf("terminal output contains a NUL byte: %q", out)
	}
}

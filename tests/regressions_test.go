package tests

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arjunuxd/uxdguard/internal/config"
	"github.com/arjunuxd/uxdguard/internal/project"
	"github.com/arjunuxd/uxdguard/internal/scan"
)

func runScan(t *testing.T, root string, cfg *config.Config) scan.RunResult {
	t.Helper()
	files, err := project.Discover(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	return scan.DefaultRegistry().Run(scan.NewContext(root, cfg, files))
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func hasRule(res scan.RunResult, ruleID string) bool {
	for _, f := range res.Findings {
		if f.RuleID == ruleID {
			return true
		}
	}
	return false
}

func countRule(res scan.RunResult, ruleID string) int {
	n := 0
	for _, f := range res.Findings {
		if f.RuleID == ruleID {
			n++
		}
	}
	return n
}

func TestRegExpExecIsNotCommandExecution(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "app.js"), "const m = regex.exec(line);\n")
	res := runScan(t, root, config.Defaults())
	if hasRule(res, "UXD-SEC-207") {
		t.Fatalf("RegExp.exec must not be reported as command execution: %+v", res.Findings)
	}
}

func TestCommandExecutionReportedOnce(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "app.js"), "exec(userInput);\n")
	res := runScan(t, root, config.Defaults())
	if got := countRule(res, "UXD-SEC-207"); got != 1 {
		t.Fatalf("expected one UXD-SEC-207 finding, got %d", got)
	}
	if hasRule(res, "UXD-SEC-300") {
		t.Fatalf("generic process execution without a shell must not be reported as shell injection: %+v", res.Findings)
	}
}

func TestDocumentationExamplesSkippedByCodeScanners(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "eval(userInput)\nurl := \"http://example.com\"\n")
	res := runScan(t, root, config.Defaults())
	if hasRule(res, "UXD-SEC-200") || hasRule(res, "UXD-SEC-203") {
		t.Fatalf("documentation examples must not be reported as code findings: %+v", res.Findings)
	}
}

func TestSecureJWTVerifyNotReported(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "auth.go"), "jwt.verify(token, secret)\n")
	res := runScan(t, root, config.Defaults())
	if hasRule(res, "UXD-SEC-400") {
		t.Fatalf("plain jwt.verify call is secure and must not be reported: %+v", res.Findings)
	}
}

func TestJWTSecretPlaceholderNotReported(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "config.js"), "jwt_secret: \"your_jwt_secret_here\"\n")
	res := runScan(t, root, config.Defaults())
	if hasRule(res, "UXD-SEC-401") {
		t.Fatalf("placeholder JWT secret must not be reported: %+v", res.Findings)
	}
}

func TestHardcodedKeyPlaceholderNotReported(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "config.go"), "encryption_key := \"your_encryption_key_here\"\n")
	res := runScan(t, root, config.Defaults())
	if hasRule(res, "UXD-SEC-502") {
		t.Fatalf("placeholder encryption key must not be reported: %+v", res.Findings)
	}
}

func TestRelativeImportNotPathTraversal(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "app.js"), "import data from \"../lib/file\"\n")
	res := runScan(t, root, config.Defaults())
	if hasRule(res, "UXD-SEC-301") || hasRule(res, "UXD-SEC-205") {
		t.Fatalf("relative import must not be reported as path traversal: %+v", res.Findings)
	}
}

func TestAdvancedDocExamplesSkipped(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "guide.md"), "jwt.verify(token, secret, {algorithms: ['none']})\n")
	res := runScan(t, root, config.Defaults())
	if hasRule(res, "UXD-SEC-400") {
		t.Fatalf("documentation example must not be reported by advanced scanner: %+v", res.Findings)
	}
}

func TestRuleSeverityOverrideAndDisable(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "main.go"), "cfg := &tls.Config{InsecureSkipVerify: true}\n")

	cfg := config.Defaults()
	cfg.Rules.Severity = map[string]string{"UXD-SEC-201": "low"}
	res := runScan(t, root, cfg)
	found := false
	for _, f := range res.Findings {
		if f.RuleID == "UXD-SEC-201" {
			found = true
			if f.Severity != scan.SeverityLow {
				t.Fatalf("severity override not applied: %s", f.Severity)
			}
		}
	}
	if !found {
		t.Fatal("expected UXD-SEC-201 finding")
	}

	cfg.Rules.Disable = []string{"UXD-SEC-201"}
	res = runScan(t, root, cfg)
	if hasRule(res, "UXD-SEC-201") {
		t.Fatalf("disabled rule was still reported: %+v", res.Findings)
	}
}

func TestDisabledScannerProducesSkip(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "main.go"), "cfg := &tls.Config{InsecureSkipVerify: true}\n")
	cfg := config.Defaults()
	cfg.Rules.Disable = []string{"sourcechecks"}
	res := runScan(t, root, cfg)
	if hasRule(res, "UXD-SEC-201") {
		t.Fatalf("findings from disabled scanner: %+v", res.Findings)
	}
	skipped := false
	for _, s := range res.Skipped {
		if s.ScannerID == "sourcechecks" {
			skipped = true
		}
	}
	if !skipped {
		t.Fatalf("expected sourcechecks skip: %+v", res.Skipped)
	}
}

func TestCancelledScanStops(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "main.go"), "cfg := &tls.Config{InsecureSkipVerify: true}\n")
	files, err := project.Discover(root, config.Defaults())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	res := scan.DefaultRegistry().Run(scan.NewContext(root, config.Defaults(), files).WithContext(ctx))
	cancelled := false
	for _, e := range res.Errors {
		if strings.Contains(e.Message, "cancelled") {
			cancelled = true
		}
	}
	if !cancelled {
		t.Fatalf("expected cancellation error: %+v", res.Errors)
	}
}

func TestFindingPathsUseForwardSlashes(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "sub", "secret.txt"), "AKIAIOSFODNN7EXAMPLE\n")
	res := runScan(t, root, config.Defaults())
	found := false
	for _, f := range res.Findings {
		if f.RuleID == "UXD-SEC-100" {
			found = true
			if f.File != "sub/secret.txt" {
				t.Fatalf("finding path = %q, want sub/secret.txt", f.File)
			}
		}
	}
	if !found {
		t.Fatal("expected AWS key finding")
	}
}

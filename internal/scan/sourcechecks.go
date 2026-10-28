package scan

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SourceCheckScanner performs source code security checks.
type SourceCheckScanner struct{}

func init() {
	_ = RegisterDefault(&SourceCheckScanner{})
}

func (s *SourceCheckScanner) ID() string   { return "sourcechecks" }
func (s *SourceCheckScanner) Name() string { return "Source Code Security Checks" }
func (s *SourceCheckScanner) Description() string {
	return "Detects common insecure patterns in source code"
}

var (
	reEval                = regexp.MustCompile(`\beval\s*\(`)
	reTLSInsecure         = regexp.MustCompile(`(?i)InsecureSkipVerify\s*:\s*true`)
	reHTTPInsecure        = regexp.MustCompile(`http://`)
	reSQLConcat           = regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE).*(\+|\%s|fmt\.Sprintf).*(SELECT|INSERT|UPDATE|DELETE)|["']\s*(SELECT|INSERT|UPDATE|DELETE)[^"']*["']\s*\+`)
	reUnsafeDeserPyPickle = regexp.MustCompile(`pickle\.loads\s*\(`)
	reUnsafeDeserJava     = regexp.MustCompile(`ObjectInputStream\s*\(`)
)

func (s *SourceCheckScanner) Run(ctx *Context) ([]Finding, []Error) {
	var findings []Finding
	var errs []Error
	if ctx == nil || ctx.Cfg == nil {
		return findings, errs
	}
	if !isRuleEnabled("sourcechecks", ctx.Cfg) {
		return findings, errs
	}
	for _, f := range ctx.Files {
		file, err := os.Open(f)
		if err != nil {
			errs = append(errs, Error{ScannerID: s.ID(), Message: "failed to open file"})
			continue
		}
		scanner := bufio.NewScanner(file)
		const maxTokenSize = 1024 * 64
		scanner.Buffer(make([]byte, maxTokenSize), maxTokenSize)
		scanner.Split(bufio.ScanLines)
		lineNum := 0
		rel, _ := filepath.Rel(ctx.Root, f)
		if rel == "" {
			rel = f
		}
		ext := strings.ToLower(filepath.Ext(f))
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			trim := strings.TrimSpace(line)
			if trim == "" || strings.HasPrefix(trim, "//") || strings.HasPrefix(trim, "#") {
				continue
			}
			// eval
			if reEval.MatchString(line) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-200",
					ScannerID:   s.ID(),
					Title:       "Use of eval() detected",
					Message:     "eval() can lead to code injection if used with untrusted input.",
					Severity:    Severity(getSeverity("UXD-SEC-200", string(SeverityHigh), ctx.Cfg)),
					Confidence:  ConfidenceMedium,
					File:        rel,
					Line:        lineNum,
					Remediation: "Avoid eval(); use safe parsers or explicit logic.",
				})
				continue
			}
			// TLS insecure skip verify (Go/Python common)
			if reTLSInsecure.MatchString(line) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-201",
					ScannerID:   s.ID(),
					Title:       "TLS certificate verification disabled",
					Message:     "Disabling certificate verification weakens TLS security.",
					Severity:    Severity(getSeverity("UXD-SEC-201", string(SeverityHigh), ctx.Cfg)),
					Confidence:  ConfidenceHigh,
					File:        rel,
					Line:        lineNum,
					Remediation: "Enable certificate verification and use proper trust stores.",
				})
				continue
			}
			// SQL concat (basic heuristic)
			if reSQLConcat.MatchString(line) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-202",
					ScannerID:   s.ID(),
					Title:       "Possible SQL string concatenation",
					Message:     "Building SQL queries via string concatenation/interpolation can lead to SQL injection.",
					Severity:    Severity(getSeverity("UXD-SEC-202", string(SeverityHigh), ctx.Cfg)),
					Confidence:  ConfidenceMedium,
					File:        rel,
					Line:        lineNum,
					Remediation: "Use parameterized queries/prepared statements.",
				})
				continue
			}
			// Insecure http (skip localhost/127.0.0.1/::1)
			if reHTTPInsecure.MatchString(line) {
				l := strings.ToLower(line)
				if !strings.Contains(l, "localhost") && !strings.Contains(l, "127.0.0.1") && !strings.Contains(l, "[::1]") {
					findings = append(findings, Finding{
						RuleID:      "UXD-SEC-203",
						ScannerID:   s.ID(),
						Title:       "Insecure HTTP endpoint detected",
						Message:     "Use of http:// in non-local contexts may transmit data in plaintext.",
						Severity:    Severity(getSeverity("UXD-SEC-203", string(SeverityMedium), ctx.Cfg)),
						Confidence:  ConfidenceMedium,
						File:        rel,
						Line:        lineNum,
						Remediation: "Use HTTPS where possible.",
					})
					continue
				}
			}
			// Deserialization
			if (ext == ".py" && reUnsafeDeserPyPickle.MatchString(line)) || (ext == ".java" && reUnsafeDeserJava.MatchString(line)) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-204",
					ScannerID:   s.ID(),
					Title:       "Potential unsafe deserialization",
					Message:     "Deserializing untrusted data can lead to remote code execution.",
					Severity:    Severity(getSeverity("UXD-SEC-204", string(SeverityHigh), ctx.Cfg)),
					Confidence:  ConfidenceMedium,
					File:        rel,
					Line:        lineNum,
					Remediation: "Avoid deserializing untrusted data or use safe deserialization practices.",
				})
				continue
			}
		}
		if err := scanner.Err(); err != nil {
			errs = append(errs, Error{ScannerID: s.ID(), Message: "error reading file"})
		}
		file.Close()
	}
	return findings, errs
}

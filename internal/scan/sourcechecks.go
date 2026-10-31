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
	reFunctionConstructor = regexp.MustCompile(`new\s+Function\s*\(`)
	reExec                = regexp.MustCompile(`\b(exec|spawn|execSync|spawnSync|system|popen|child_process\.exec|child_process\.spawn)\s*\(`)
	reTLSInsecure         = regexp.MustCompile(`(?i)InsecureSkipVerify\s*:\s*true|rejectUnauthorized\s*:\s*false|NODE_TLS_REJECT_UNAUTHORIZED\s*.*0`)
	reHTTPInsecure        = regexp.MustCompile(`http://`)
	reSQLConcat           = regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE).*(\+|\%s|fmt\.Sprintf|\$\{).*(SELECT|INSERT|UPDATE|DELETE)|["']\s*(SELECT|INSERT|UPDATE|DELETE)[^"']*["']\s*(\+|concat)|query\s*\(\s*["'][^"']+\$\{`)
	reUnsafeDeserPyPickle = regexp.MustCompile(`pickle\.loads\s*\(|yaml\.load\s*\(|marshal\.load\s*\(`)
	reUnsafeDeserJava     = regexp.MustCompile(`ObjectInputStream\s*\(|readObject\s*\(`)
	rePathTraversal       = regexp.MustCompile(`\.\./|\.\.\\`)
	reWeakHashSrc         = regexp.MustCompile(`\b(md5|sha1)\s*\(`)
	reSSRFSrc             = regexp.MustCompile(`(?i)(fetch|axios\.|http\.request|requests\.get|urllib\.request|curl|http\.Get)\s*\([^)]*(\$\{|user|req\.|params|query|body)`)
	reAuthBypass          = regexp.MustCompile(`(?i)(if\s*\([^)]*auth[^)]*\)\s*\{?\s*return\s+next|skipAuth|bypass.*auth|disable.*auth)`)
	reBase64Eval          = regexp.MustCompile(`(?i)(atob|Buffer\.from\([^)]*,\s*['"]base64['"]\)|base64.*decode).*(eval|Function|exec)`)
	reOpenRedirectSrc     = regexp.MustCompile(`(?i)(res\.redirect|window\.location|location\.href)\s*=\s*[^;]*(req\.|params|query|body)`)
	rePlainPwd            = regexp.MustCompile(`(?i)function\s+\w*[Pp]ass\w*\s*\([^)]*\)\s*\{[^}]*return\s+.*(String\(|password|pwd)[^}]*\}`)
	rePlainPwdAlt         = regexp.MustCompile(`(?i)return\s+String\s*\(\s*\w*[Pp]ass\w*\s*\)`)
	rePwdCompare          = regexp.MustCompile(`(?i)(input|password|pwd|pass)\s*===?\s*(stored|hash|password|pwd|pass)|\w*[Pp]ass\w*\s*===?\s*\w*[Pp]ass\w*`)
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
		if ctx.Cancelled() {
			errs = append(errs, Error{ScannerID: s.ID(), Message: "cancelled"})
			return findings, errs
		}
		rel, _ := filepath.Rel(ctx.Root, f)
		if rel == "" {
			rel = f
		}
		// Documentation examples are not executable source code. Secrets are
		// still handled by the dedicated secret scanner.
		if isDocFile(rel) {
			continue
		}
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
		ext := strings.ToLower(filepath.Ext(f))
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if strings.TrimSpace(line) == "" || isCommentLine(line) {
				continue
			}
			// eval / Function constructor
			if reEval.MatchString(line) || reFunctionConstructor.MatchString(line) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-200",
					ScannerID:   s.ID(),
					Title:       "Use of dynamic code evaluation detected",
					Message:     "eval() or Function() can lead to code injection if used with untrusted input.",
					Severity:    Severity(getSeverity("UXD-SEC-200", string(SeverityHigh), ctx.Cfg)),
					Confidence:  ConfidenceMedium,
					File:        rel,
					Line:        lineNum,
					Remediation: "Avoid dynamic code evaluation; use safe parsers or explicit logic.",
				})
				continue
			}
			// command exec (basic)
			if reExec.MatchString(line) && !isRegExpExec(line) {
				l := strings.ToLower(line)
				if !strings.Contains(l, "console.log") {
					findings = append(findings, Finding{
						RuleID:      "UXD-SEC-207",
						ScannerID:   s.ID(),
						Title:       "Potential command execution detected",
						Message:     "Executing shell commands can lead to command injection if input is untrusted.",
						Severity:    Severity(getSeverity("UXD-SEC-207", string(SeverityHigh), ctx.Cfg)),
						Confidence:  ConfidenceMedium,
						File:        rel,
						Line:        lineNum,
						Remediation: "Use parameterized APIs and avoid passing untrusted input to shells.",
					})
					continue
				}
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
			// Path traversal
			if rePathTraversal.MatchString(line) && !isImportLikeLine(line) && (ext == ".js" || ext == ".ts" || ext == ".py" || ext == ".php" || ext == ".java" || ext == ".go") {
				l := strings.ToLower(line)
				if strings.Contains(l, "read") || strings.Contains(l, "open") || strings.Contains(l, "path") || strings.Contains(l, "file") {
					findings = append(findings, Finding{
						RuleID:      "UXD-SEC-205",
						ScannerID:   s.ID(),
						Title:       "Potential path traversal pattern",
						Message:     "Path traversal sequences may indicate unsafe file access.",
						Severity:    Severity(getSeverity("UXD-SEC-205", string(SeverityHigh), ctx.Cfg)),
						Confidence:  ConfidenceMedium,
						File:        rel,
						Line:        lineNum,
						Remediation: "Validate and normalize paths; use allowlists and avoid concatenating user input.",
					})
					continue
				}
			}
			// Weak hash
			if reWeakHashSrc.MatchString(line) {
				l := strings.ToLower(line)
				if strings.Contains(l, "password") || strings.Contains(l, "hash") || strings.Contains(l, "crypto") {
					findings = append(findings, Finding{
						RuleID:      "UXD-SEC-206",
						ScannerID:   s.ID(),
						Title:       "Weak cryptographic hash in security context",
						Message:     "MD5/SHA1 are weak for security-sensitive use cases.",
						Severity:    Severity(getSeverity("UXD-SEC-206", string(SeverityMedium), ctx.Cfg)),
						Confidence:  ConfidenceMedium,
						File:        rel,
						Line:        lineNum,
						Remediation: "Use strong algorithms like SHA-256, bcrypt, argon2, or PBKDF2.",
					})
					continue
				}
			}
			// SSRF
			if reSSRFSrc.MatchString(line) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-208",
					ScannerID:   s.ID(),
					Title:       "Potential SSRF pattern detected",
					Message:     "Making network requests with user-controlled input can lead to SSRF.",
					Severity:    Severity(getSeverity("UXD-SEC-208", string(SeverityHigh), ctx.Cfg)),
					Confidence:  ConfidenceMedium,
					File:        rel,
					Line:        lineNum,
					Remediation: "Allowlist destinations, block private/internal IPs, and disable redirects where possible.",
				})
				continue
			}
			// Auth bypass
			if reAuthBypass.MatchString(line) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-209",
					ScannerID:   s.ID(),
					Title:       "Potential authentication bypass detected",
					Message:     "Authentication logic appears to be skipped or bypassed.",
					Severity:    Severity(getSeverity("UXD-SEC-209", string(SeverityCritical), ctx.Cfg)),
					Confidence:  ConfidenceMedium,
					File:        rel,
					Line:        lineNum,
					Remediation: "Review authentication/authorization logic and ensure all routes are properly protected.",
				})
				continue
			}
			// Base64 + eval
			if reBase64Eval.MatchString(line) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-210",
					ScannerID:   s.ID(),
					Title:       "Potential obfuscated code execution",
					Message:     "Decoding base64 and passing to eval/Function/exec is often malicious or unsafe.",
					Severity:    Severity(getSeverity("UXD-SEC-210", string(SeverityCritical), ctx.Cfg)),
					Confidence:  ConfidenceHigh,
					File:        rel,
					Line:        lineNum,
					Remediation: "Avoid executing decoded content. If required, use strict allowlists and safe parsers.",
				})
				continue
			}
			// Open redirect
			if reOpenRedirectSrc.MatchString(line) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-211",
					ScannerID:   s.ID(),
					Title:       "Potential open redirect",
					Message:     "Redirecting based on user input can lead to phishing attacks.",
					Severity:    Severity(getSeverity("UXD-SEC-211", string(SeverityMedium), ctx.Cfg)),
					Confidence:  ConfidenceMedium,
					File:        rel,
					Line:        lineNum,
					Remediation: "Use allowlists for redirect destinations or relative paths only.",
				})
				continue
			}
			// Plaintext password handling
			if rePlainPwd.MatchString(line) || rePlainPwdAlt.MatchString(line) {
				l := strings.ToLower(line)
				if strings.Contains(l, "bcrypt") || strings.Contains(l, "argon2") || strings.Contains(l, "pbkdf2") || strings.Contains(l, "scrypt") {
					continue
				}
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-212",
					ScannerID:   s.ID(),
					Title:       "Plaintext password handling detected",
					Message:     "Returning or storing passwords in plaintext is insecure.",
					Severity:    Severity(getSeverity("UXD-SEC-212", string(SeverityHigh), ctx.Cfg)),
					Confidence:  ConfidenceHigh,
					File:        rel,
					Line:        lineNum,
					Remediation: "Use a strong password hashing algorithm (e.g., bcrypt, argon2, PBKDF2).",
				})
				continue
			}
			// Direct password comparison
			if rePwdCompare.MatchString(line) {
				l := strings.ToLower(line)
				if strings.Contains(l, "bcrypt") || strings.Contains(l, "argon2") || strings.Contains(l, "pbkdf2") || strings.Contains(l, "scrypt") {
					goto skipPwd2
				}
				if (strings.Contains(l, "confirmpassword") || strings.Contains(l, "confirmpass") || strings.Contains(l, "passwordconfirm") || strings.Contains(l, "password") && strings.Contains(l, "confirm")) && strings.Contains(line, "===") {
					goto skipPwd2
				}
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-213",
					ScannerID:   s.ID(),
					Title:       "Insecure password comparison",
					Message:     "Comparing plaintext passwords with === is insecure and timing-vulnerable.",
					Severity:    Severity(getSeverity("UXD-SEC-213", string(SeverityHigh), ctx.Cfg)),
					Confidence:  ConfidenceMedium,
					File:        rel,
					Line:        lineNum,
					Remediation: "Use constant-time comparison and proper password hash verification.",
				})
				continue
			}
		skipPwd2:
		}
		if err := scanner.Err(); err != nil {
			errs = append(errs, Error{ScannerID: s.ID(), Message: "error reading file"})
		}
		file.Close()
	}
	return findings, errs
}

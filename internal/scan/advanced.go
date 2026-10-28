package scan

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type AdvancedScanner struct{}

func init() {
	_ = RegisterDefault(&AdvancedScanner{})
}

func (a *AdvancedScanner) ID() string   { return "advanced" }
func (a *AdvancedScanner) Name() string { return "Advanced Security Checks" }
func (a *AdvancedScanner) Description() string {
	return "Extended security detection for injection, auth, crypto, SSRF, and more"
}

var (
	reExecCmd       = regexp.MustCompile(`(?i)(exec\.Command|os\.Exec|child_process\.exec|subprocess\.call|subprocess\.run|system\(|popen\(|exec\()`)
	reShellTrue     = regexp.MustCompile(`(?i)(shell\s*[:=]\s*true|Shell\s*=\s*true|{shell)`)
	rePathJoin      = regexp.MustCompile(`(?i)(filepath\.Join|path\.join).*\+|concat.*path|\.\./`)
	reInnerHTML     = regexp.MustCompile(`(?i)(innerHTML|dangerouslySetInnerHTML|v-html|outerHTML)`)
	reSQLConcatAdv  = regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE).*(\+|concat|%s|\$\{)`)
	reOpenRedirect  = regexp.MustCompile(`(?i)(redirect|Location).*http|res\.redirect\(.*req\.|header.*location`)
	reDebugEndpoint = regexp.MustCompile(`(?i)(/debug|/admin|/swagger|/api-docs|/actuator)`)
	reJWTVerify     = regexp.MustCompile(`(?i)(jwt\.verify|verify.*false|algorithms.*\[.*none.*\]|"alg".*"none")`)
	reJWTSecret     = regexp.MustCompile(`(?i)(jwt.*secret|JWT_SECRET)\s*[:=]\s*["'][^"']{4,}["']`)
	reCookieInsec   = regexp.MustCompile(`(?i)(Set-Cookie|cookie).*HttpOnly\s*=\s*false|Secure\s*=\s*false`)
	reWeakHash      = regexp.MustCompile(`(?i)(md5\.|sha1\.|MD5|SHA1)\(|\bmd5\b|\bsha1\b`)
	reWeakRand      = regexp.MustCompile(`(?i)(math/rand|rand\.Int|random\.random|Math\.random\(\)).*(token|session|password|secret)`)
	reHardcodeKey   = regexp.MustCompile(`(?i)(encryption.*key|aes.*key|secret.*key)\s*[:=]\s*["'][A-Za-z0-9+/=]{16,}["']`)
	reLogSecret     = regexp.MustCompile(`(?i)(log.*password|print.*secret|console\.log.*token|logger.*key)`)
	reSSRF          = regexp.MustCompile(`(?i)(http\.Get|axios\.get|requests\.get|fetch\().*req\.|url.*=.*req\.`)
	reCSRFMissing   = regexp.MustCompile(`(?i)(app\.post|router\.post).*/(delete|update|create|transfer)`)
	reFileUpload    = regexp.MustCompile(`(?i)(multipart|form-data|uploaded.*file|MoveUploadedFile)`)
	reDeserPickle   = regexp.MustCompile(`(?i)pickle\.loads|yaml\.load\(|unserialize\(|ObjectInputStream`)
)

func (a *AdvancedScanner) Run(ctx *Context) ([]Finding, []Error) {
	var findings []Finding
	var errs []Error
	if ctx == nil || ctx.Cfg == nil {
		return findings, errs
	}
	if !isRuleEnabled("advanced", ctx.Cfg) {
		return findings, errs
	}

	for _, f := range ctx.Files {
		file, err := os.Open(f)
		if err != nil {
			errs = append(errs, Error{ScannerID: a.ID(), Message: "failed to open file"})
			continue
		}
		scanner := bufio.NewScanner(file)
		const maxTokenSize = 1024 * 64
		scanner.Buffer(make([]byte, maxTokenSize), maxTokenSize)
		scanner.Split(bufio.ScanLines)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			rel, _ := filepath.Rel(ctx.Root, f)
			if rel == "" {
				rel = f
			}
			relSlash := filepath.ToSlash(rel)

			// Command injection
			if reExecCmd.MatchString(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-300", ScannerID: a.ID(), Title: "Potential command injection with shell execution",
					Message: "Process execution with shell=true and dynamic input can lead to command injection.", Severity: SeverityHigh, Confidence: ConfidenceMedium,
					File: relSlash, Line: lineNum, Remediation: "Avoid shell execution; pass arguments as arrays and validate/sanitize inputs.",
				})
				continue
			}
			// Path traversal
			if rePathJoin.MatchString(line) && strings.Contains(line, "..") {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-301", ScannerID: a.ID(), Title: "Potential path traversal",
					Message: "Path construction with user input may allow directory traversal.", Severity: SeverityHigh, Confidence: ConfidenceMedium,
					File: relSlash, Line: lineNum, Remediation: "Validate and normalize paths; use filepath.Clean and restrict to allowed directories.",
				})
				continue
			}
			// XSS
			if reInnerHTML.MatchString(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-302", ScannerID: a.ID(), Title: "Potential XSS via unsafe HTML rendering",
					Message: "Direct HTML injection can lead to cross-site scripting.", Severity: SeverityHigh, Confidence: ConfidenceMedium,
					File: relSlash, Line: lineNum, Remediation: "Use safe templating with auto-escaping; sanitize user input.",
				})
				continue
			}
			// SQL injection broader
			if reSQLConcatAdv.MatchString(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-303", ScannerID: a.ID(), Title: "Potential SQL injection via string concatenation",
					Message: "SQL queries built with string concatenation are vulnerable to injection.", Severity: SeverityHigh, Confidence: ConfidenceMedium,
					File: relSlash, Line: lineNum, Remediation: "Use parameterized queries/prepared statements.",
				})
				continue
			}
			// Open redirect
			if reOpenRedirect.MatchString(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-304", ScannerID: a.ID(), Title: "Potential open redirect",
					Message: "Redirecting to user-controlled URLs can lead to phishing.", Severity: SeverityMedium, Confidence: ConfidenceMedium,
					File: relSlash, Line: lineNum, Remediation: "Validate redirect destinations against an allowlist.",
				})
				continue
			}
			// Debug endpoints
			if reDebugEndpoint.MatchString(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-305", ScannerID: a.ID(), Title: "Potential debug/admin endpoint exposed",
					Message: "Debug or administrative endpoints may be exposed in production.", Severity: SeverityMedium, Confidence: ConfidenceLow,
					File: relSlash, Line: lineNum, Remediation: "Disable debug endpoints in production and protect admin routes with auth.",
				})
				continue
			}
			// JWT issues
			if reJWTVerify.MatchString(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-400", ScannerID: a.ID(), Title: "Insecure JWT verification",
					Message: "JWT verification may be disabled or using insecure algorithms.", Severity: SeverityCritical, Confidence: ConfidenceHigh,
					File: relSlash, Line: lineNum, Remediation: "Always verify JWT signatures with secure algorithms (e.g., RS256).",
				})
				continue
			}
			if reJWTSecret.MatchString(line) && !isPlaceholder(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-401", ScannerID: a.ID(), Title: "Hardcoded JWT secret",
					Message: "JWT secret is hardcoded in source code.", Severity: SeverityCritical, Confidence: ConfidenceHigh,
					File: relSlash, Line: lineNum, Remediation: "Store JWT secrets in environment variables or secure secret management.",
				})
				continue
			}
			// Cookie security
			if reCookieInsec.MatchString(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-402", ScannerID: a.ID(), Title: "Insecure cookie configuration",
					Message: "Cookie missing HttpOnly/Secure flags or has insecure settings.", Severity: SeverityMedium, Confidence: ConfidenceMedium,
					File: relSlash, Line: lineNum, Remediation: "Set HttpOnly, Secure, and appropriate SameSite attributes.",
				})
				continue
			}
			// Crypto issues
			if reWeakHash.MatchString(line) && !strings.Contains(line, "//") && !strings.Contains(line, "#") {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-500", ScannerID: a.ID(), Title: "Weak cryptographic hash",
					Message: "MD5/SHA1 may be weak for security purposes.", Severity: SeverityMedium, Confidence: ConfidenceMedium,
					File: relSlash, Line: lineNum, Remediation: "Use SHA-256 or stronger algorithms for security-sensitive operations.",
				})
				continue
			}
			if reWeakRand.MatchString(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-501", ScannerID: a.ID(), Title: "Weak random number generation",
					Message: "Predictable RNG used for security-sensitive values.", Severity: SeverityHigh, Confidence: ConfidenceMedium,
					File: relSlash, Line: lineNum, Remediation: "Use cryptographically secure random generators (crypto/rand).",
				})
				continue
			}
			if reHardcodeKey.MatchString(line) && !isPlaceholder(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-502", ScannerID: a.ID(), Title: "Hardcoded encryption key",
					Message: "Encryption key is hardcoded in source code.", Severity: SeverityCritical, Confidence: ConfidenceHigh,
					File: relSlash, Line: lineNum, Remediation: "Store encryption keys in secure key management systems.",
				})
				continue
			}
			if reLogSecret.MatchString(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-503", ScannerID: a.ID(), Title: "Potential secret logging",
					Message: "Logging sensitive data may expose secrets.", Severity: SeverityHigh, Confidence: ConfidenceMedium,
					File: relSlash, Line: lineNum, Remediation: "Avoid logging passwords, tokens, or secrets.",
				})
				continue
			}
			// SSRF
			if reSSRF.MatchString(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-600", ScannerID: a.ID(), Title: "Potential SSRF",
					Message: "Making requests with user-controlled URLs can lead to SSRF.", Severity: SeverityHigh, Confidence: ConfidenceMedium,
					File: relSlash, Line: lineNum, Remediation: "Validate and allowlist URLs; block internal/private IP ranges.",
				})
				continue
			}
			// File upload
			if reFileUpload.MatchString(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-601", ScannerID: a.ID(), Title: "File upload handling detected",
					Message: "File uploads require validation for type, size, and filename.", Severity: SeverityMedium, Confidence: ConfidenceLow,
					File: relSlash, Line: lineNum, Remediation: "Validate file types, restrict extensions, use safe filenames, and scan uploads.",
				})
				continue
			}
			// Unsafe deserialization
			if reDeserPickle.MatchString(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-602", ScannerID: a.ID(), Title: "Unsafe deserialization",
					Message: "Deserializing untrusted data can lead to remote code execution.", Severity: SeverityCritical, Confidence: ConfidenceHigh,
					File: relSlash, Line: lineNum, Remediation: "Avoid deserializing untrusted data; use safe serialization formats.",
				})
				continue
			}
		}
		if err := scanner.Err(); err != nil {
			errs = append(errs, Error{ScannerID: a.ID(), Message: "error reading file"})
		}
		file.Close()
	}
	return findings, errs
}

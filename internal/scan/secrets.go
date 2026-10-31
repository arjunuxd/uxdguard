package scan

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/arjunuxd/uxdguard/internal/config"
)

// SecretScanner detects common secrets.
type SecretScanner struct{}

func init() {
	_ = RegisterDefault(&SecretScanner{})
}

func (s *SecretScanner) ID() string   { return "secrets" }
func (s *SecretScanner) Name() string { return "Secret Detection" }
func (s *SecretScanner) Description() string {
	return "Detects likely hardcoded secrets and credentials"
}

var (
	reAWSKey      = regexp.MustCompile(`\b(AKIA|ASIA)[A-Z0-9]{16}\b`)
	reGitHubToken = regexp.MustCompile(`\bghp_[A-Za-z0-9]{36}\b|\bgho_[A-Za-z0-9]{36}\b|\bghu_[A-Za-z0-9]{36}\b|\bghs_[A-Za-z0-9]{36}\b|\bghr_[A-Za-z0-9]{36}\b|\bgithub_pat_[A-Za-z0-9_]{22,}\b`)
	reSlackToken  = regexp.MustCompile(`\bxox[baprs]-[0-9a-zA-Z-]{10,48}\b`)
	rePrivateKey  = regexp.MustCompile(`-----BEGIN (RSA|EC|DSA|OPENSSH|PGP|ENCRYPTED) PRIVATE KEY-----`)
	reJWT         = regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\b`)
	reDBConn      = regexp.MustCompile(`(?i)(postgres|mysql|mongodb(\+srv)?|redis|mssql)://[^:\s]+:[^@\s]+@[^\s/]+`)
	reAssign      = regexp.MustCompile(`(?i)^\s*(?:export\s+)?([A-Z0-9_]+)\s*=\s*(.+)$`)
	reGoogleKey   = regexp.MustCompile(`\bAIza[0-9A-Za-z_-]{35}\b`)
	reResendKey   = regexp.MustCompile(`\bre_[A-Za-z0-9]{24,}\b`)
	reRazorpayKey = regexp.MustCompile(`\brzp_(live|test)_[A-Za-z0-9]{14,}\b`)
	rePosthogKey  = regexp.MustCompile(`\bphc_[A-Za-z0-9]{32,}\b`)
	reStripeKey   = regexp.MustCompile(`\bsk_(live|test)_[A-Za-z0-9]{24,}\b|\brk_(live|test)_[A-Za-z0-9]{24,}\b`)
	reSendGridKey = regexp.MustCompile(`\bSG\.[A-Za-z0-9_-]{22}\.[A-Za-z0-9_-]{43}\b`)
	reDiscordHook = regexp.MustCompile(`https://discord(?:app)?\.com/api/webhooks/[0-9]+/[A-Za-z0-9_-]+`)
	reCloudinary  = regexp.MustCompile(`cloudinary://[0-9]+:[A-Za-z0-9_-]+@[A-Za-z0-9_-]+`)
)

var placeholders = map[string]bool{
	"your_api_key": true, "replace_me": true, "example": true, "changeme": true, "password": true, "secret": true, "token": true,
	"xxxx": true, "test": true, "dummy": true, "sample": true, "placeholder": true,
}

func isPlaceholder(v string) bool {
	lv := strings.ToLower(strings.TrimSpace(v))
	if lv == "" {
		return true
	}
	// exact matches
	if placeholders[lv] {
		return true
	}
	// Common placeholder markers. These are intentionally conservative and
	// only cover values that are clearly not production credentials.
	if strings.HasPrefix(lv, "xxx") || strings.HasPrefix(lv, "yyy") {
		return true
	}
	if strings.HasPrefix(lv, "your_") || strings.HasPrefix(lv, "your-") {
		return true
	}
	for _, marker := range []string{
		"your_api_key", "your_secret", "your-secret", "your_key_here",
		"insert_key_here", "replace_me", "placeholder", "not_a_real",
		"notreal", "changeme", "redacted",
	} {
		if strings.Contains(lv, marker) {
			return true
		}
	}
	return false
}

// looksLikeExample reports whether a line appears to contain a documentation
// or placeholder example rather than a live secret.
func looksLikeExample(line string) bool {
	l := strings.ToLower(line)
	for _, marker := range []string{
		"example", "placeholder", "dummy", "sample", "your_", "your-",
		"changeme", "xxxx", "yyyy", "redacted", "not-a-real", "not_a_real",
	} {
		if strings.Contains(l, marker) {
			return true
		}
	}
	return false
}

func stripQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	// remove inline comments for shell/env style
	if idx := strings.Index(s, " #"); idx > 0 {
		s = strings.TrimSpace(s[:idx])
	}
	if idx := strings.Index(s, " //"); idx > 0 {
		s = strings.TrimSpace(s[:idx])
	}
	return strings.TrimSpace(s)
}

func getSeverity(id string, def string, cfg *config.Config) string {
	if cfg != nil && cfg.Rules.Severity != nil {
		if v, ok := cfg.Rules.Severity[id]; ok {
			return v
		}
	}
	return def
}

func (s *SecretScanner) Run(ctx *Context) ([]Finding, []Error) {
	var findings []Finding
	var errs []Error
	if ctx == nil || ctx.Cfg == nil {
		return findings, errs
	}
	// Check if secrets scanner enabled
	if !isRuleEnabled("secrets", ctx.Cfg) {
		return findings, errs
	}
	for _, f := range ctx.Files {
		if ctx.Cancelled() {
			errs = append(errs, Error{ScannerID: s.ID(), Message: "cancelled"})
			return findings, errs
		}
		// skip very large files already handled by discovery
		file, err := os.Open(f)
		if err != nil {
			errs = append(errs, Error{ScannerID: s.ID(), Message: "failed to open file"})
			continue
		}
		scanner := bufio.NewScanner(file)
		// increase buffer for long lines
		const maxTokenSize = 1024 * 64
		scanner.Buffer(make([]byte, maxTokenSize), maxTokenSize)
		scanner.Split(bufio.ScanLines)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if strings.TrimSpace(line) == "" || isCommentLine(line) {
				continue
			}
			rel, _ := filepath.Rel(ctx.Root, f)
			if rel == "" {
				rel = f
			}
			// AWS key
			if reAWSKey.MatchString(line) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-100",
					ScannerID:   s.ID(),
					Title:       "Potential AWS access key detected",
					Message:     "A string matching AWS access key pattern was found.",
					Severity:    Severity(getSeverity("UXD-SEC-100", string(SeverityHigh), ctx.Cfg)),
					Confidence:  ConfidenceMedium,
					File:        rel,
					Line:        lineNum,
					Remediation: "Remove hardcoded credentials and use environment variables or a secure secret manager.",
				})
				continue
			}
			// GitHub tokens
			if reGitHubToken.MatchString(line) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-101",
					ScannerID:   s.ID(),
					Title:       "Potential GitHub token detected",
					Message:     "A string matching GitHub token pattern was found.",
					Severity:    Severity(getSeverity("UXD-SEC-101", string(SeverityCritical), ctx.Cfg)),
					Confidence:  ConfidenceHigh,
					File:        rel,
					Line:        lineNum,
					Remediation: "Revoke the token if exposed and use secure credential storage.",
				})
				continue
			}
			// Slack tokens
			if reSlackToken.MatchString(line) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-102",
					ScannerID:   s.ID(),
					Title:       "Potential Slack token detected",
					Message:     "A string matching Slack token pattern was found.",
					Severity:    Severity(getSeverity("UXD-SEC-102", string(SeverityHigh), ctx.Cfg)),
					Confidence:  ConfidenceHigh,
					File:        rel,
					Line:        lineNum,
					Remediation: "Remove hardcoded tokens and rotate if exposed.",
				})
				continue
			}
			// Private key
			if rePrivateKey.MatchString(line) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-103",
					ScannerID:   s.ID(),
					Title:       "Private key detected",
					Message:     "A private key header was found in the file.",
					Severity:    Severity(getSeverity("UXD-SEC-103", string(SeverityCritical), ctx.Cfg)),
					Confidence:  ConfidenceHigh,
					File:        rel,
					Line:        lineNum,
					Remediation: "Remove private keys from source control and use secure key management.",
				})
				continue
			}
			// DB connection string
			if reDBConn.MatchString(line) {
				low := strings.ToLower(line)
				if !strings.Contains(low, "user:pass") && !strings.Contains(low, "example") && !strings.Contains(low, "your_") {
					findings = append(findings, Finding{
						RuleID:      "UXD-SEC-104",
						ScannerID:   s.ID(),
						Title:       "Database connection string with credentials detected",
						Message:     "A connection string with embedded credentials was found.",
						Severity:    Severity(getSeverity("UXD-SEC-104", string(SeverityHigh), ctx.Cfg)),
						Confidence:  ConfidenceHigh,
						File:        rel,
						Line:        lineNum,
						Remediation: "Use environment variables or a secret manager for database credentials.",
					})
					continue
				}
			}
			// JWT
			if reJWT.MatchString(line) && !looksLikeExample(line) {
				findings = append(findings, Finding{
					RuleID:      "UXD-SEC-106",
					ScannerID:   s.ID(),
					Title:       "Potential hardcoded JWT detected",
					Message:     "A JSON Web Token was found hardcoded in source.",
					Severity:    Severity(getSeverity("UXD-SEC-106", string(SeverityHigh), ctx.Cfg)),
					Confidence:  ConfidenceMedium,
					File:        rel,
					Line:        lineNum,
					Remediation: "Use secure token generation and storage; avoid hardcoding tokens.",
				})
				continue
			}
			// Provider-specific patterns
			if reGoogleKey.MatchString(line) {
				findings = append(findings, Finding{
					RuleID: "UXD-SEC-107", ScannerID: s.ID(), Title: "Potential Google API key detected",
					Message: "A string matching Google API key pattern was found.", Severity: Severity(getSeverity("UXD-SEC-107", string(SeverityHigh), ctx.Cfg)),
					Confidence: ConfidenceHigh, File: rel, Line: lineNum, Remediation: "Remove hardcoded credentials and use environment variables or a secure secret manager.",
				})
				continue
			}
			if reResendKey.MatchString(line) {
				findings = append(findings, Finding{RuleID: "UXD-SEC-108", ScannerID: s.ID(), Title: "Potential Resend API key detected", Message: "A string matching Resend API key pattern was found.", Severity: Severity(getSeverity("UXD-SEC-108", string(SeverityHigh), ctx.Cfg)), Confidence: ConfidenceHigh, File: rel, Line: lineNum, Remediation: "Remove hardcoded credentials and rotate if exposed."})
				continue
			}
			if reRazorpayKey.MatchString(line) {
				findings = append(findings, Finding{RuleID: "UXD-SEC-109", ScannerID: s.ID(), Title: "Potential Razorpay key detected", Message: "A string matching Razorpay key pattern was found.", Severity: Severity(getSeverity("UXD-SEC-109", string(SeverityHigh), ctx.Cfg)), Confidence: ConfidenceHigh, File: rel, Line: lineNum, Remediation: "Remove hardcoded credentials and rotate if exposed."})
				continue
			}
			if rePosthogKey.MatchString(line) {
				findings = append(findings, Finding{RuleID: "UXD-SEC-110", ScannerID: s.ID(), Title: "Potential PostHog key detected", Message: "A string matching PostHog key pattern was found.", Severity: Severity(getSeverity("UXD-SEC-110", string(SeverityHigh), ctx.Cfg)), Confidence: ConfidenceHigh, File: rel, Line: lineNum, Remediation: "Remove hardcoded credentials and rotate if exposed."})
				continue
			}
			if reStripeKey.MatchString(line) {
				findings = append(findings, Finding{RuleID: "UXD-SEC-111", ScannerID: s.ID(), Title: "Potential Stripe key detected", Message: "A string matching Stripe key pattern was found.", Severity: Severity(getSeverity("UXD-SEC-111", string(SeverityCritical), ctx.Cfg)), Confidence: ConfidenceHigh, File: rel, Line: lineNum, Remediation: "Remove hardcoded credentials and rotate if exposed."})
				continue
			}
			if reSendGridKey.MatchString(line) {
				findings = append(findings, Finding{RuleID: "UXD-SEC-112", ScannerID: s.ID(), Title: "Potential SendGrid API key detected", Message: "A string matching SendGrid API key pattern was found.", Severity: Severity(getSeverity("UXD-SEC-112", string(SeverityHigh), ctx.Cfg)), Confidence: ConfidenceHigh, File: rel, Line: lineNum, Remediation: "Remove hardcoded credentials and rotate if exposed."})
				continue
			}
			if reDiscordHook.MatchString(line) {
				findings = append(findings, Finding{RuleID: "UXD-SEC-113", ScannerID: s.ID(), Title: "Potential Discord webhook URL detected", Message: "A Discord webhook URL was found.", Severity: Severity(getSeverity("UXD-SEC-113", string(SeverityHigh), ctx.Cfg)), Confidence: ConfidenceHigh, File: rel, Line: lineNum, Remediation: "Remove webhook URLs from source control and use secure storage."})
				continue
			}
			if reCloudinary.MatchString(line) {
				findings = append(findings, Finding{RuleID: "UXD-SEC-114", ScannerID: s.ID(), Title: "Potential Cloudinary credentials detected", Message: "Cloudinary URL with credentials was found.", Severity: Severity(getSeverity("UXD-SEC-114", string(SeverityHigh), ctx.Cfg)), Confidence: ConfidenceHigh, File: rel, Line: lineNum, Remediation: "Remove hardcoded credentials and use environment variables."})
				continue
			}
			// Assignment pattern (quoted and unquoted)
			if m := reAssign.FindStringSubmatch(line); len(m) == 3 {
				key := strings.ToUpper(m[1])
				val := stripQuotes(m[2])
				if isPlaceholder(val) {
					continue
				}
				if len(val) < 8 {
					continue
				}
				// skip common non-secret keys and UI copy fields
				if key == "NODE_ENV" || key == "PORT" || key == "HOST" {
					continue
				}
				if strings.Contains(key, "LABEL") || strings.Contains(key, "HINT") ||
					strings.Contains(key, "MESSAGE") || strings.Contains(key, "TITLE") ||
					strings.Contains(key, "DESCRIPTION") {
					continue
				}
				// check if key suggests secret
				if strings.Contains(key, "SECRET") || strings.Contains(key, "TOKEN") || strings.Contains(key, "KEY") || strings.Contains(key, "PASSWORD") || strings.Contains(key, "API_KEY") || strings.Contains(key, "PRIVATE") {
					findings = append(findings, Finding{
						RuleID:      "UXD-SEC-105",
						ScannerID:   s.ID(),
						Title:       "Hardcoded credential assignment detected",
						Message:     "A suspicious hardcoded credential assignment was found.",
						Severity:    Severity(getSeverity("UXD-SEC-105", string(SeverityHigh), ctx.Cfg)),
						Confidence:  ConfidenceMedium,
						File:        rel,
						Line:        lineNum,
						Remediation: "Move credentials to environment variables or a secure secret manager.",
					})
					continue
				}
			}
		}
		if err := scanner.Err(); err != nil {
			errs = append(errs, Error{ScannerID: s.ID(), Message: "error reading file"})
		}
		file.Close()
	}
	return findings, errs
}

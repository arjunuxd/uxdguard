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
	reAssign      = regexp.MustCompile(`(?i)\b(password|passwd|pwd|api[_-]?key|apikey|secret|client[_-]?secret|access[_-]?token|auth[_-]?token|refresh[_-]?token|private[_-]?key)\s*[:=]\s*["']([^"']{8,})["']`)
)

var placeholders = map[string]bool{
	"your_api_key": true, "replace_me": true, "example": true, "changeme": true, "password": true, "secret": true, "token": true,
	"xxxx": true, "test": true, "dummy": true, "sample": true,
}

func isPlaceholder(v string) bool {
	lv := strings.ToLower(strings.TrimSpace(v))
	if lv == "" {
		return true
	}
	if placeholders[lv] {
		return true
	}
	if strings.Contains(lv, "your_") || strings.Contains(lv, "replace") || strings.Contains(lv, "example") || strings.Contains(lv, "changeme") {
		return true
	}
	if strings.HasPrefix(lv, "xxx") || strings.HasPrefix(lv, "yyy") {
		return true
	}
	return false
}

func isRuleEnabled(id string, cfg *config.Config) bool {
	if cfg == nil {
		return true
	}
	// check disable first
	for _, d := range cfg.Rules.Disable {
		if d == id {
			return false
		}
	}
	if len(cfg.Rules.Enable) > 0 {
		for _, e := range cfg.Rules.Enable {
			if e == id {
				return true
			}
		}
		return false
	}
	return true
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
			trim := strings.TrimSpace(line)
			if trim == "" || strings.HasPrefix(trim, "//") || strings.HasPrefix(trim, "#") || strings.HasPrefix(trim, "*") {
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
			if reJWT.MatchString(line) {
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
			// Assignment pattern
			if m := reAssign.FindStringSubmatch(line); len(m) == 3 {
				val := m[2]
				if isPlaceholder(val) {
					continue
				}
				// skip very short
				if len(val) < 8 {
					continue
				}
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
		if err := scanner.Err(); err != nil {
			errs = append(errs, Error{ScannerID: s.ID(), Message: "error reading file"})
		}
		file.Close()
	}
	return findings, errs
}

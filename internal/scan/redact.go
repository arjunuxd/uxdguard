package scan

import (
	"strings"
)

// RedactSecret attempts to redact secret-like values from strings.
// This is best-effort and avoids logging full secrets.
func RedactSecret(s string) string {
	if s == "" {
		return s
	}
	// Replace long token-like strings
	// AWS keys
	if len(s) > 20 {
		// simple redaction: keep first 4 and last 4 if looks sensitive
		if strings.Contains(s, "AKIA") || strings.Contains(s, "ASIA") {
			if len(s) >= 16 {
				return s[:4] + "..." + s[len(s)-4:]
			}
		}
		// generic long hex/base64-like tokens
		if len(s) > 32 {
			return s[:4] + "..." + s[len(s)-4:]
		}
	}
	return s
}

// RedactLine removes obvious secret assignments' values partially
func RedactLine(line string) string {
	// For safety, if line looks like assignment with long value, redact
	// This is conservative
	l := strings.ToLower(line)
	if strings.Contains(l, "password") || strings.Contains(l, "secret") || strings.Contains(l, "api_key") || strings.Contains(l, "token") || strings.Contains(l, "private_key") {
		// Try to redact after = or :
		for _, sep := range []string{"=", ":"} {
			if idx := strings.Index(line, sep); idx >= 0 {
				// keep key part
				prefix := line[:idx+1]
				// replace rest with redaction if long
				if len(line)-idx-1 > 8 {
					return prefix + " [REDACTED]"
				}
			}
		}
	}
	return line
}

package project

import (
	"path/filepath"
	"strings"
)

// Excluded checks if path matches any exclusion pattern (simple glob-like).
func Excluded(rel string, patterns []string) bool {
	// Normalize to forward slashes style
	norm := filepath.ToSlash(rel)
	for _, p := range patterns {
		pn := filepath.ToSlash(p)
		// Handle **
		if strings.Contains(pn, "**") {
			// prefix match if ends with /**
			if strings.HasSuffix(pn, "/**") {
				prefix := strings.TrimSuffix(pn, "/**")
				if prefix == "" {
					continue
				}
				if norm == prefix || strings.HasPrefix(norm+"/", prefix+"/") || strings.HasPrefix(norm, prefix+"/") {
					return true
				}
				// also if it starts with prefix and has / in between
				if strings.HasPrefix(norm, prefix+"/") {
					return true
				}
				continue
			}
		}
		// directory name match
		base := filepath.Base(norm)
		if pn == base || pn == norm {
			return true
		}
		if strings.HasSuffix(pn, "/**") {
			continue
		}
		// simple prefix match for directories
		if strings.HasSuffix(pn, "/") {
			if strings.HasPrefix(norm, strings.TrimSuffix(pn, "/")+"/") {
				return true
			}
			continue
		}
		// exact component match
		parts := strings.Split(norm, "/")
		for _, part := range parts {
			if matched, _ := filepath.Match(pn, part); matched {
				return true
			}
		}
		if matched, _ := filepath.Match(pn, norm); matched {
			return true
		}
		if matched, _ := filepath.Match(pn, base); matched {
			return true
		}
	}
	return false
}

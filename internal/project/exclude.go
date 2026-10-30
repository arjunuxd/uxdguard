package project

import (
	"path/filepath"
	"strings"
)

// Excluded checks if path matches any exclusion pattern (simple glob-like).
func Excluded(rel string, patterns []string) bool {
	// Normalize to forward slashes style
	norm := filepath.ToSlash(filepath.Clean(rel))
	if norm == "." {
		return false
	}
	for _, p := range patterns {
		if p == "" {
			continue
		}
		pn := filepath.ToSlash(filepath.Clean(p))
		// Handle **
		if strings.Contains(pn, "**") {
			if strings.HasSuffix(pn, "/**") {
				prefix := strings.TrimSuffix(pn, "/**")
				if prefix == "" {
					continue
				}
				if norm == prefix || strings.HasPrefix(norm, prefix+"/") || strings.Contains("/"+norm+"/", "/"+prefix+"/") {
					return true
				}
				continue
			}
			if strings.Contains(pn, "/**/") {
				parts := strings.Split(pn, "/**/")
				if len(parts) == 2 {
					prefix := parts[0]
					suffix := parts[1]
					if prefix != "" && strings.HasPrefix(norm, prefix+"/") {
						rest := strings.TrimPrefix(norm, prefix+"/")
						if strings.HasSuffix(rest, suffix) || strings.Contains(rest, "/"+suffix) {
							return true
						}
					} else if prefix == "" {
						if strings.HasSuffix(norm, suffix) || strings.Contains("/"+norm+"/", "/"+suffix) {
							return true
						}
					}
				}
				continue
			}
		}
		// directory name match
		base := filepath.Base(norm)
		if pn == base || pn == norm {
			return true
		}
		if strings.HasSuffix(pn, "/") {
			pn2 := strings.TrimSuffix(pn, "/")
			if norm == pn2 || strings.HasPrefix(norm, pn2+"/") {
				return true
			}
			continue
		}
		// exact component match for directories
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
		if strings.HasPrefix(norm, pn+"/") {
			return true
		}
	}
	return false
}

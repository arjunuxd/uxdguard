package scan

import (
	"path/filepath"
	"strings"
)

// isDocFile reports whether a path is documentation rather than executable
// source. Secrets are still scanned in documentation; code-pattern rules are
// not applied there because examples are not runnable code.
func isDocFile(rel string) bool {
	switch strings.ToLower(filepath.Ext(rel)) {
	case ".md", ".markdown", ".rst", ".adoc", ".txt":
		return true
	}
	return false
}

// isCommentLine reports whether a trimmed line is a whole-line comment.
func isCommentLine(line string) bool {
	trim := strings.TrimSpace(line)
	return strings.HasPrefix(trim, "//") ||
		strings.HasPrefix(trim, "#") ||
		trim == "*" ||
		strings.HasPrefix(trim, "* ") ||
		strings.HasPrefix(trim, "*/") ||
		strings.HasPrefix(trim, "<!--")
}

// isRegExpExec reports whether a line is a regular-expression match using
// RegExp.prototype.exec rather than a process/spawn call. The generic exec
// pattern intentionally matches both, so this removes the common false
// positive without weakening process-execution detection.
func isRegExpExec(line string) bool {
	l := strings.ToLower(line)
	if !strings.Contains(l, ".exec(") {
		return false
	}
	for _, marker := range []string{
		"regexp", "regex", "pattern", "matcher", "expression",
		"re.exec", "match.exec", "reg.exec",
	} {
		if strings.Contains(l, marker) {
			return true
		}
	}
	return false
}

// isImportLikeLine reports whether a line imports or requires another module.
// Relative import specifiers such as "../lib/file" are normal and must not be
// treated as path traversal in file-access code.
func isImportLikeLine(line string) bool {
	l := strings.ToLower(strings.TrimSpace(line))
	return strings.Contains(l, "import ") ||
		strings.Contains(l, "import(") ||
		strings.Contains(l, "require(") ||
		strings.Contains(l, "#include") ||
		strings.HasPrefix(l, "from ")
}

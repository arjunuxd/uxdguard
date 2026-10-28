package scan

import (
	"github.com/arjunuxd/uxdguard/internal/config"
	"github.com/arjunuxd/uxdguard/internal/exitcode"
)

// EvaluatePolicy evaluates policy and returns exit code.
func EvaluatePolicy(cfg *config.Config, res RunResult) int {
	// If there are operational errors from required scans, it's operational failure
	if len(res.Errors) > 0 {
		return exitcode.OperationalFail
	}
	// Check if any finding is blocking
	failSet := map[string]bool{}
	for _, f := range cfg.Policy.FailOn {
		failSet[f] = true
	}
	for _, f := range res.Findings {
		if failSet[string(f.Severity)] {
			return exitcode.Findings
		}
	}
	return exitcode.OK
}

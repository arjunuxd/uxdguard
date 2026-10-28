package app

import "github.com/arjunuxd/uxdguard/internal/exitcode"

const (
	ExitOK              = exitcode.OK
	ExitFindings        = exitcode.Findings
	ExitUsageOrConfig   = exitcode.UsageOrConfig
	ExitOperationalFail = exitcode.OperationalFail
)

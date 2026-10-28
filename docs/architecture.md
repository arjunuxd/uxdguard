# Architecture

## Overview

UXDGuard is a local-first, modular monolithic CLI. Phase 1 establishes the core scanning engine and configuration system.

## Modules

- `cmd/uxdguard/`: CLI entrypoint
- `internal/app/`: Application orchestration and exit codes
- `internal/cli/`: Command/flag parsing and command handlers
- `internal/config/`: Configuration schema, defaults, loading, validation
- `internal/project/`: Project root discovery, file traversal, exclusion logic
- `internal/scan/`: Scanner interface, registry, context, finding model, policy engine
- `internal/report/`: Terminal and JSON reporters
- `internal/exitcode/`: Shared exit code constants

## Data Flow

1. CLI parses args and resolves project root/config
2. `project.LoadWithConfig` loads/merges config (defaults if missing)
3. `project.Discover` enumerates eligible files with exclusions, size limits, symlink rules
4. `scan.NewContext` creates context with root/files/config
5. `scan.DefaultRegistry().Run` executes all registered scanners
6. `scan.EvaluatePolicy` applies policy to determine exit code
7. Reporters output results (JSON to stdout, diagnostics to stderr)

## Scanner Interface

```go
type Scanner interface {
    ID() string
    Name() string
    Description() string
    Run(ctx *Context) ([]Finding, []Error)
}
```

Scanners receive only the scan context (root, files, config). No unrestricted filesystem/process access is granted. Scanner failures are collected as operational errors.

## Findings and Policy

Findings have stable rule/scanner IDs, severity/confidence, file/line/column, remediation. Policy maps configured `fail_on` severities to blocking behavior; operational errors result in exit code 3.

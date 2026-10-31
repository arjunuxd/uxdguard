# UXDGuard - Architecture

**Version:** 1.1.0

**Type:** Local-first CLI (Go)

**Scope:** Security and code-quality scanning for source files and staged Git changes

## Overview

UXDGuard is a modular, local-first Go CLI. It discovers files in a project, runs a registry of scanners, applies rule configuration, evaluates policy against findings, and reports results to the terminal or as JSON. Built-in scanning involves no telemetry, network calls, or AI services.

## High-Level Flow

```text
cmd/uxdguard/main.go
  -> internal/cli/root.go (command dispatch)
     -> internal/cli/check.go / hook.go / etc.
        -> internal/project (discovery + config loading)
        -> internal/scan (registry + scanners + policy)
        -> internal/report (terminal/json)
```

## Package Structure

| Package | Purpose |
| --- | --- |
| `cmd/uxdguard` | Binary entrypoint. Receives version/commit/build-date linker values, registers external adapters, and calls `cli.Execute`. |
| `internal/cli` | Command parsing and handlers (`check`, `init`, `config validate`, `rules`, `hook`, `uninstall`, `version`). Includes version resolution and pagination/filter flags. |
| `internal/app` | Reusable scan orchestration and exit code mapping. |
| `internal/exitcode` | Exit code constants (`OK`, `Findings`, `UsageOrConfig`, `OperationalFail`). |
| `internal/config` | Config schema (`.uxdguard.yaml`), defaults, loading, merging, and validation. |
| `internal/project` | Project root discovery, file traversal, exclusion matching, symlink handling, size checks, and path containment. |
| `internal/scan` | Core engine: scanner interface, registry, context, findings, policy, deduplication, secrets, source checks, advanced checks, staged scanning, and redaction. |
| `internal/gitutil` | Git helpers using explicit argument lists (repo root, hooks path, staged files, staged content, binary detection). |
| `internal/report` | Terminal (paginated) and JSON reporters. |
| `internal/external` | Optional Gitleaks and Semgrep adapters with configurable timeouts and safe process execution. |

## Key Components

### CLI and Commands

- `uxdguard --version` / `--help`
- `uxdguard version` - version, commit, and build date
- `uxdguard init` - create `.uxdguard.yaml` without overwriting an existing file
- `uxdguard check` - scan with `--json`, `--staged`, `--project-root`, `--page`, `--limit`, `--severity`, `--rule`, and `--no-progress`
- `uxdguard config validate` - validate configuration
- `uxdguard rules` - list registered scanners
- `uxdguard hook install|uninstall|status` - manage the pre-commit hook
- `uxdguard uninstall` - remove the installed executable

### Version Resolution

`internal/cli/version.go` resolves build metadata in this order:

1. Version: linker-injected `main.version`, then Go module metadata (`debug.ReadBuildInfo`) when the main module path matches and the version is not `(devel)`, then the truthful fallback `dev`.
2. Commit and build date: VCS settings (`vcs.revision`, `vcs.time`), then any linker-injected `main.commit` / `main.buildDate` value, then the truthful fallback `unknown`.

Reporters use the same resolved version as the `version` command, so JSON output and terminal output cannot diverge.

### Configuration

- Defaults are applied when no configuration file is present.
- User configuration is merged with defaults and validated.
- Supports `scan` (size limit, symlink policy, exclusions), `policy.fail_on`, `rules.enable` / `rules.disable` / `rules.severity`, and optional `external` adapters.
- Rule filtering is centralized in `internal/scan/registry.go`; every finding, including staged and external findings, passes through it.

### Discovery and Files

- Walks the project root with deterministic ordering.
- Skips default and user-configured exclusions.
- Enforces `max_file_size_bytes` and processes regular files only.
- With `follow_symlinks: false` (default), symlinks are ignored.
- With `follow_symlinks: true`, symlinks are resolved only when the target stays inside the project root; resolved directories are tracked to prevent cycles.
- Rejects paths that escape the project root.
- Produces forward-slash finding paths on every platform.

### Scanning Engine

- `Scanner` interface: `ID()`, `Name()`, `Description()`, `Run(*Context) ([]Finding, []Error)`.
- Built-in scanners: `secrets`, `sourcechecks`, and `advanced`.
- Optional external scanners: `gitleaks` and `semgrep`, only active when enabled in configuration.
- Context carries `Root`, `Files`, `Cfg`, and an optional `context.Context`.
- Context cancellation is checked between files and before scanners run.
- Findings are filtered by rule configuration, sorted deterministically, deduplicated, and sorted again.
- Policy evaluates configured `fail_on` severities.
- Secret values are redacted in terminal and JSON output.

### Reporting

- Terminal reporter with pagination (20 per page by default), color support, and severity/rule filters.
- Empty filtered results are reported explicitly rather than as an invalid page range.
- JSON reporter includes all findings, skipped scanners, errors, and summary counts.
- JSON output is not truncated by pagination options.

### Staged Scanning and Hooks

- `--staged` reads indexed content through `git show :path`, so the scan reflects exactly what would be committed.
- Staged content is written to a private temporary directory (`0700`) with private files (`0600`) and removed afterward.
- Hook installation uses `git rev-parse --git-path hooks`, honoring `core.hooksPath` and linked worktrees.
- The generated hook quotes the executable path so installations under paths containing spaces work.

## Security Design

- **Local-first**: Built-in scanning uses no telemetry, network calls, or external uploads.
- **No code execution**: Scanned files are analyzed statically.
- **Path safety**: Discovery confines resolved symlinks to the project root and prevents traversal.
- **Secret safety**: Secret values are not printed in reports.
- **Minimal dependencies**: The only runtime dependency is `gopkg.in/yaml.v3`.
- **Explicit execution**: Git and optional external tools are invoked with argument lists, never through a shell.

## Limitations

- Pattern-based detection only; no AST parsing, data-flow, or taint analysis.
- False positives and false negatives are possible.
- No dependency CVE scanning.
- Language coverage varies by rule.
- External integrations depend on separately installed third-party tools; their network behavior is controlled by those tools.

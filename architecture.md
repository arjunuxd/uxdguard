# UXDGuard - Architecture

**Version:** v1.1.0 Flash  
**Type:** Local-first CLI (Go)  
**Scope:** Security and code-quality scanning for source files and staged Git changes

## Overview

UXDGuard is a modular, local-first Go CLI. It discovers files in a project, runs a registry of scanners, evaluates policy against findings, and reports results to the terminal or as JSON. No telemetry, network calls, or AI API usage are involved in the built-in scanning flow.

## High-Level Flow

```text
cmd/uxdguard/main.go
  -> internal/cli/root.go (flag/command dispatch)
     -> internal/cli/check.go / hook.go / etc.
        -> internal/project (discovery + config loading)
        -> internal/scan (registry + scanners + policy)
        -> internal/report (terminal/json)
```

## Package Structure

| Package | Purpose |
| --- | --- |
| `cmd/uxdguard` | Binary entrypoint. Supports version injection via linker flags and calls `cli.SetVersion`. |
| `internal/cli` | Command parsing and handlers (`check`, `init`, `config validate`, `rules`, `hook`, `uninstall`, `version`). Supports pagination, filtering, and flags. |
| `internal/app` | Orchestration helpers and exit code mapping. |
| `internal/exitcode` | Exit code constants (`OK`, `Findings`, `UsageOrConfig`, `OperationalFail`). |
| `internal/config` | Config schema (`.uxdguard.yaml`), defaults, loading, merging, and validation. |
| `internal/project` | Project root discovery, file traversal, exclusion matching, symlink handling, size checks, and path-traversal guards. |
| `internal/scan` | Core engine: scanner interface, registry, context, findings, policy, deduplication, secrets, source checks, advanced checks, staged scanning, redaction. |
| `internal/gitutil` | Git helpers using explicit `git` args (repo root, staged files, staged content). Includes binary detection. |
| `internal/report` | Terminal (paginated) and JSON reporters. |
| `internal/external` | Optional adapters for `gitleaks` and `semgrep` with timeouts and safe process execution. |

## Key Components

### CLI & Commands
- `uxdguard --version` / `--help`
- `uxdguard version` - version information
- `uxdguard init` - create `.uxdguard.yaml` and `.uxdguard.yaml.example`
- `uxdguard check` - normal scan with `--json`, `--staged`, `--project-root`, `--page`, `--limit`, `--severity`, `--rule`, `--no-progress`
- `uxdguard config validate` - validate configuration
- `uxdguard rules` - list registered scanners/rules
- `uxdguard hook install|uninstall|status` - manage pre-commit hook
- `uxdguard uninstall` - remove global executable

### Configuration
- Defaults applied if no config present
- Merged with user config
- Validated via `internal/config/validate.go`
- Supports `scan` (size limits, symlink following, exclusions), `policy.fail_on`, `rules.enable/disable/severity`, `external` (opt-in)

### Discovery & Files
- Walks from resolved root with symlink safety
- Respects default exclusions (dependencies, build outputs, caches)
- Custom exclusions supported
- Enforces `max_file_size_bytes`
- Filters regular files only
- Blocks paths escaping project root
- Sorted results for deterministic behavior

### Scanning Engine
- `Scanner` interface: `ID()`, `Name()`, `Description()`, `Run(*Context) ([]Finding, []Error)`
- Built-in scanners: `secrets`, `sourcechecks`, `advanced`
- Context carries `Root`, `Files`, `Cfg`, and optional `context.Context`
- Findings deduplicated and sorted deterministically
- Policy evaluates `fail_on` severities
- Secret redaction applied in outputs

### Reporting
- Terminal reporter with pagination (default 20 per page), severity/color support, filters
- JSON reporter with complete results (no truncation)
- Machine-readable output safe from progress pollution

## Security Design

- **Local-first**: No telemetry, network calls, or external uploads
- **No code execution**: Static analysis only
- **Path safety**: Traversal and symlink protections
- **Secret safety**: Redaction in outputs
- **Minimal dependencies**: Only `gopkg.in/yaml.v3`
- **Explicit execution**: Git and external tools use argument lists only

## Limitations

- Pattern-based detection only - no AST parsing or taint analysis
- May produce false positives or negatives
- No comprehensive dependency CVE scanning
- Coverage varies by language and pattern availability
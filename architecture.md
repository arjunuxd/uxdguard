# UXDGuard – Current Architecture

**Version:** 0.1.0-dev  
**Type:** Local-first CLI (Go)  
**Scope:** Security and code-quality scanning for source files and staged Git changes

## 1. Overview

UXDGuard is a modular, monolithic Go CLI designed to run entirely locally. It discovers files in a project, runs a registry of scanners, evaluates policy against findings, and reports results to the terminal or as JSON. No telemetry, network calls, or AI API usage are involved in the built-in flow.

## 2. High-Level Flow

```text
cmd/uxdguard/main.go
  -> internal/cli/root.go (flag/command dispatch)
     -> internal/cli/check.go / hook.go / etc.
        -> internal/project (discovery + config loading)
        -> internal/scan (registry + scanners + policy)
        -> internal/report (terminal/json)
```

## 3. Package Structure

| Package | Purpose |
| --- | --- |
| `cmd/uxdguard` | Binary entrypoint. Supports version injection via `-ldflags -X main.version=...` and calls `cli.SetVersion` if set. |
| `internal/cli` | Command parsing and handlers (`check`, `init`, `config validate`, `rules`, `hook install/uninstall/status`). Handles `--json`, `--staged`, `--project-root`. |
| `internal/app` | Application orchestration helpers and exit code mapping (`exitcodes.go`). |
| `internal/exitcode` | Central exit code constants (`OK`, `Findings`, `UsageOrConfig`, `OperationalFail`). |
| `internal/config` | Config schema (`.uxdguard.yaml`), defaults, loading, merging, and validation. |
| `internal/project` | Project root discovery, file traversal (`filepath.WalkDir`), exclusion matching, symlink handling, size checks, and path-traversal guards. |
| `internal/scan` | Core engine: `Scanner` interface, registry, context, findings model, policy evaluation, secret/source scanners, redaction helpers, and staged scanning. |
| `internal/gitutil` | Git helpers using explicit `git` args (repo root, staged files via `diff --cached -z`, staged content via `git show :<path>`). Includes binary content detection. |
| `internal/report` | Output formatters: human-readable terminal reporter and versioned JSON reporter. |
| `internal/external` | Optional adapters for `gitleaks` and `semgrep`, plus a safe process runner with timeouts, output handling, and redaction helpers. |

## 4. Key Components

### 4.1 CLI & Commands
- `uxdguard --version` / `--help`
- `uxdguard init` – creates `.uxdguard.yaml` and `.uxdguard.yaml.example` if missing
- `uxdguard check` – normal scan of discovered files
- `uxdguard check --staged` – scans staged Git index content only
- `uxdguard check --json` – machine-readable output
- `uxdguard config validate` – validates config against schema
- `uxdguard rules` – lists registered scanners/rules
- `uxdguard hook install|uninstall|status` – manages pre-commit hook in `.git/hooks/` (only manages hooks containing the UXDGuard marker)

### 4.2 Configuration
- Defaults applied if no config present
- Merged with user config
- Validated via `internal/config/validate.go`
- Supports `scan` (size limits, symlink following, exclusions), `policy.fail_on` (severity-based blocking), `rules.enable/disable/severity`, `external` (opt-in)

### 4.3 Discovery & Files
- Walks from resolved root
- Skips `.git` directories
- Respects exclusions (supports `**`, directory/pattern matching)
- Skips symlinks unless `follow_symlinks: true`
- Enforces `max_file_size_bytes`
- Filters to regular files
- Blocks paths escaping project root (traversal guards)
- Sorted results for deterministic behavior

### 4.4 Scanning Engine
- `Scanner` interface: `ID()`, `Name()`, `Description()`, `Run(*Context) ([]Finding, []Error)`
- `DefaultRegistry()` holds built-in scanners: `secrets` and `sourcechecks`
- Context carries `Root`, `Files`, `Cfg`, and optional `context.Context` for cancellation
- Findings include `RuleID`, `ScannerID`, `Title`, `Message`, `Severity`, `Confidence`, `File`, `Line`, `Remediation`
- Policy evaluates `fail_on` severities; operational errors from scans return exit code `3`

### 4.5 Built-in Rules
**Secrets (`internal/scan/secrets.go`):**
- `UXD-SEC-100` AWS access keys (`AKIA|ASIA...`)
- `UXD-SEC-101` GitHub tokens (`ghp_/gho_/ghu_/ghs_/ghr_`)
- `UXD-SEC-102` Slack tokens (`xox*`)
- `UXD-SEC-103` Private key headers
- `UXD-SEC-104` DB connection strings with credentials
- `UXD-SEC-105` Hardcoded credential assignments (with placeholder filtering)

**Source checks (`internal/scan/sourcechecks.go`):**
- `UXD-SEC-200` `eval()` usage
- `UXD-SEC-201` `InsecureSkipVerify: true`
- `UXD-SEC-202` SQL string concatenation patterns
- `UXD-SEC-203` Insecure `http://` endpoints (excluding localhost)
- `UXD-SEC-204` Unsafe deserialization patterns

All pattern-based. Secret values are not emitted in findings/messages/reports.

### 4.6 Staged Scanning
- `scan.ScanStaged` uses `gitutil.StagedFiles` (`git diff --cached --name-only -z --diff-filter=ACDMRTUXB`)
- Reads index content via `git show :<path>`
- Skips binaries, respects exclusions and `max_file_size_bytes`
- Writes staged content to a temp directory preserving relative paths
- Runs scanners against temp files, rewrites findings back to original repo-relative paths
- Cleans up temp directory

### 4.7 Git Hooks
- Installs a POSIX shell pre-commit hook that calls `uxdguard check --staged`
- Hook marked with `# UXDGuard pre-commit hook`
- Will not overwrite existing non-UXDGuard hooks
- Status/uninstall only act on UXDGuard-managed hooks
- Uses `os.Executable()` path in generated hook

### 4.8 External Scanners
- `internal/external/runner.go`: explicit `exec.CommandContext`, timeouts, stdout/stderr capture, basic redaction helpers
- `gitleaks` and `semgrep` adapters (opt-in via config)
- Never auto-installed
- Failures treated as operational/scan errors surfaced clearly

### 4.9 Reporting
- `TerminalReporter`: concise human output with rule IDs, locations, remediation, summary
- `JSONReporter`: structured output with `tool`, `project`, `status`, `findings`, `skipped`, `errors`, `summary`
- Diagnostics go to stderr where appropriate

## 5. Security Design

- **Local-first:** No uploads, telemetry, or AI calls
- **No code execution:** Scans by pattern/regex only
- **Safe process execution:** Explicit args, no shell string interpolation
- **Path safety:** Root boundary checks, traversal guards, symlink control
- **Resource limits:** File size caps, timeouts/cancellation
- **Secret safety:** Redaction logic and tests; findings avoid leaking raw secrets
- **Hook safety:** Marker-based ownership, refusal to overwrite foreign hooks
- **Index-only staged scan:** Reads from Git index, not working tree
- **Least privilege:** CI uses minimal permissions

## 6. Build, Test & Release

- **Go:** `go 1.22.9`
- **Build:** `go build ./cmd/uxdguard`
- **Cross-compile targets:** Windows amd64, Linux amd64, macOS arm64
- **CI:** `.github/workflows/ci.yml` – format, vet, build, test, coverage, race (non-Windows), cross-compile (actions pinned by SHA)
- **Release:** `.github/workflows/release.yml` – tag-triggered (`v*`), version injection via `-ldflags "-X main.version=$VERSION"`, SHA256 checksums, draft release only
- **Quality checks:** `gofmt`, `go vet`, `go test`, `go mod verify`

## 7. File Structure

```text
F:.
   .gitattributes
   .gitignore
   .uxdguard.yaml
   .uxdguard.yaml.example
   architecture.md
   go.mod
   go.sum
   LICENSE
   README.md
   uxdguard.exe

+---.github
    +---workflows
            ci.yml
            release.yml

+---cmd
    +---uxdguard
            main.go

+---dist
        uxdguard.exe

+---docs
        architecture.md
        configuration.md
        external-scanners.md
        hooks.md
        privacy.md
        PROJECT-STATUS.md
        rules.md
        security-policy.md
        threat-model.md

+---internal
    +---app
        app.go
        exitcodes.go
    +---cli
        check.go
        hook.go
        root.go
    +---config
        config.go
        defaults.go
        external.go
        validate.go
    +---exitcode
        exitcode.go
    +---external
        gitleaks.go
        runner.go
        semgrep.go
    +---gitutil
        git.go
    +---project
        discover.go
        exclude.go
        files.go
    +---report
        json.go
        terminal.go
    +---scan
        context.go
        findings.go
        noop.go
        policy.go
        redact.go
        registry.go
        scanner.go
        secrets.go
        sourcechecks.go
        staged.go

+---tests
    redaction_test.go
    secrets_test.go
    sourcechecks_test.go
    +---fixtures
        +---clean
        +---secrets
            aws.txt
            db.txt
            github.txt
            placeholder.txt
            privatekey.pem
        +---source
            insecure.go
            sql.go
    +---integration
```

## 8. Current State

- Core functionality implemented
- Security fixes applied (path traversal, hooks, redaction)
- Staged scanning and optional external adapters implemented
- CI/release workflows in place
- Tests cover secrets, source checks, and redaction
- Pre-release (`0.1.0-dev`), pending human review

**Note:** Detection is pattern-based. Git hooks can be bypassed with `git commit --no-verify`. External tools require manual installation.


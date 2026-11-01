# Contributing to UXDGuard

Thank you for helping improve UXDGuard. Contributions of code, detection rules,
tests, documentation, platform testing, and bug reports are all welcome.

UXDGuard is a local-first security scanner. Contributions must preserve that
model: no telemetry, no uploading scanned code, no executing scanned code, and no
network dependency in the built-in scanners.

## Before You Start

- Search [existing issues](https://github.com/arjunuxd/uxdguard/issues) before
  opening a new one.
- For a large feature, refactor, or new dependency, open an issue first so the
  design can be discussed before you invest time.
- For a suspected vulnerability in UXDGuard itself, follow
  [SECURITY.md](SECURITY.md). Do not open a public issue with exploit details.
- Never include real credentials, personal data, or third-party data in an
  issue, test, fixture, or commit. Use synthetic placeholders such as
  `AKIAIOSFODNN7EXAMPLE`.
- Be respectful and constructive. Harassment, discrimination, and personal
  attacks are not acceptable in any project space.

## Ways to Contribute

- Report reproducible bugs, false positives, and false negatives
- Improve detection rules without weakening existing coverage
- Add regression tests for reported issues
- Improve documentation, examples, and error messages
- Test on Windows, Linux, and macOS
- Review pull requests
- Improve performance and memory usage on large repositories

## Development Setup

### Requirements

| Tool | Version / notes |
| --- | --- |
| Go | 1.22.9 or later |
| Git | Any recent version |
| C compiler | Optional; needed only for `go test -race` locally |
| Gitleaks / Semgrep | Optional; only for external adapter testing |

### Clone and build

```bash
git clone https://github.com/arjunuxd/uxdguard.git
cd uxdguard
go mod download
go build ./...
go test -count=1 ./...
```

Run the CLI from source:

```bash
go run ./cmd/uxdguard check
go run ./cmd/uxdguard check --help
```

On Windows PowerShell, the same commands work. If `go test -race` reports that
cgo is unavailable, use the normal test command; CI runs the race detector on
Linux and macOS.

## Project Layout

| Path | Purpose |
| --- | --- |
| `cmd/uxdguard` | Executable entry point and linker-injected build metadata |
| `internal/cli` | Command parsing, progress output, `init`, hooks, uninstall |
| `internal/app` | High-level scan orchestration and exit codes |
| `internal/config` | Configuration schema, defaults, loading, and validation |
| `internal/project` | Root discovery, file discovery, exclusions, symlink policy |
| `internal/scan` | Scanner interface, registry, built-in scanners, rules, policy |
| `internal/gitutil` | Git commands, staged files, staged blob reads, hooks path |
| `internal/report` | Terminal, JSON, and SARIF reporting |
| `internal/external` | Optional Gitleaks and Semgrep adapters |
| `internal/exitcode` | Stable process exit codes |
| `tests` | Cross-package regression and end-to-end tests |
| `docs` | User and maintainer documentation |

Architecture details are in [architecture.md](architecture.md). User-facing
configuration, rules, hooks, and privacy docs are under [`docs/`](docs).

## Quality Gates

Every change must pass the same checks CI runs:

```bash
gofmt -w .
gofmt -l .              # must print nothing
go build ./...
go vet ./...
go mod verify
go test -count=1 ./...
```

CI also runs coverage and `go test -race ./...` on Linux and macOS, plus
cross-compilation for Windows amd64, Linux amd64, and macOS arm64.

Useful optional checks:

```bash
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
go run github.com/securego/gosec/v2/cmd/gosec@latest ./...
```

Report new static-analysis findings in the pull request and either fix them or
explain precisely why they are intentional or a false positive. Do not silence a
finding without justification.

## Adding or Changing Detection Rules

1. Choose an unused `UXD-SEC-###` rule ID in the appropriate category and add it
   to [docs/rules.md](docs/rules.md).
2. Add the finding in the relevant scanner under `internal/scan`.
3. Populate the finding fields consistently:

   | Field | Guidance |
   | --- | --- |
   | `RuleID` | Stable `UXD-SEC-###` identifier |
   | `ScannerID` | Use the scanner's `ID()` value |
   | `Title` | Short, specific, user-facing description |
   | `Message` | Why the pattern is risky, without including secret values |
   | `Severity` | One of `critical`, `high`, `medium`, `low`, `info` |
   | `Confidence` | `high`, `medium`, or `low` |
   | `File`, `Line`, `Column` | Relative path with forward slashes; accurate location |
   | `Remediation` | Actionable fix, not generic advice |

4. Prefer a narrow, context-aware pattern over a broad one that creates noise.
5. Do not remove or weaken an existing detection merely to reduce findings.
   Fix a false positive by narrowing the match, excluding non-executable
   documentation, or adding a clearly justified exception.
6. Keep secret values out of titles, messages, JSON, and terminal output.
7. Add regression tests that cover both the positive case and the nearest
   realistic negative case. Put cross-package fixtures in `tests/` and
   package-level tests next to the code.
8. Update [docs/VERIFICATION-REPORT.md](docs/VERIFICATION-REPORT.md) when the
   security posture or validation evidence changes.

## Adding a Scanner

A scanner implements the `scan.Scanner` interface:

```go
type Scanner interface {
    ID() string
    Name() string
    Description() string
    Run(ctx *Context) ([]Finding, []Error)
}
```

Requirements:

- Register the scanner exactly once with `scan.RegisterDefault` from an `init`
  function, or explicitly from the CLI entry point.
- Use a short, stable, lower-case scanner ID. IDs are part of configuration and
  JSON output and must not change casually.
- Handle a nil context or nil configuration defensively.
- Honour cancellation through `ctx.Cancelled()` / `ctx.Ctx`.
- Read only files supplied in `ctx.Files`, close files, and use bounded line
  buffers.
- Return findings separately from operational errors. Do not turn a normal
  finding into an error or vice versa.
- Never execute scanned code, load scanned plugins, or make network calls from a
  built-in scanner.
- Keep output deterministic; central sorting and deduplication run after all
  scanners finish.

External adapters belong in `internal/external`. They must be disabled by
default, require explicit configuration, use argument lists rather than a shell,
apply a timeout, cap captured output, and document any network behaviour.

## Engineering Standards

### Security-sensitive code

Changes involving any of the following need extra tests and a short threat
analysis in the pull request:

- Path handling, symlinks, file discovery, or exclusions
- Temporary files, permissions, or cleanup
- Configuration parsing, defaults, or validation
- Git index reads, staged scanning, or hook installation
- External process execution, argument construction, or output capture
- Terminal, JSON, or SARIF rendering of untrusted values
- Any code that could expose a secret or follow a path outside the project root

Rules of thumb:

- Treat repository paths, file contents, tool output, and configuration as
  untrusted input.
- Keep followed paths inside the project root and reject traversal, absolute
  paths, and volume-relative paths.
- Use private temporary directories and remove them after use.
- Never interpolate untrusted values into a shell command.
- Bound memory and output when reading files or subprocess streams.
- Sanitize control characters before writing untrusted text to a terminal.
- Preserve `--json` output as clean, machine-readable JSON with no progress text.

### Compatibility

- Keep the Go module path `github.com/arjunuxd/uxdguard`.
- Support Go 1.22.9 and later unless a change is explicitly discussed.
- Support Windows, Linux, and macOS.
- Use `filepath` for filesystem work and forward slashes in reported paths.
- Preserve the documented exit codes: `0` success, `1` findings, `2` usage or
  configuration error, `3` operational failure.
- Treat the JSON report fields, rule IDs, CLI flags, and configuration schema as
  stable interfaces. Discuss breaking changes before implementing them.

### Dependencies

Prefer the Go standard library. A new dependency must be justified by a clear
benefit, actively maintained, and reviewed for supply-chain and licensing impact.
Do not add a dependency for a small helper that can be written safely in a few
lines.

### Code style

- Follow standard Go formatting and idioms.
- Return errors with context; do not ignore errors silently.
- Check and handle close errors for files being written.
- Keep functions focused and testable.
- Avoid hidden global state beyond the scanner registry.
- Keep comments useful and specific; explain why, not what.

## Tests

- Use table-driven tests where they improve clarity.
- Use `t.TempDir()` for filesystem tests.
- Skip symlink tests when the platform or user account cannot create symlinks,
  and rely on CI for Linux and macOS coverage.
- Do not depend on the network, a real home directory, or a pre-existing Git
  configuration.
- Test both success and failure paths, including invalid flags and malformed
  configuration.
- For bug fixes, add a regression test that fails before the fix.
- For false-positive fixes, keep a negative test so the pattern does not regress.

Run the full suite before opening a pull request:

```bash
go test -count=1 ./...
```

## Commit and Pull Request Guidelines

Use a clear commit subject with a conventional prefix where it fits:

```text
feat: add scan progress feedback
fix: reject Windows paths during staged scanning
docs: expand security policy
test: cover oversized staged blobs
refactor: simplify rule configuration
chore: update CI actions
```

In the pull request description, include:

- What changed and why
- User-visible behaviour changes
- Security or privacy implications
- Tests added or updated
- Documentation updated
- Platforms tested
- Any known limitations or follow-up work

Keep pull requests focused. Avoid unrelated formatting or refactoring changes in
the same commit. Do not rewrite shared history or force-push a branch that others
may be reviewing without coordination.

## Review and Release

Maintainers review pull requests for correctness, security, compatibility, and
documentation. A pull request is ready to merge when CI is green and review
comments are resolved.

Releases follow semantic versioning. Published tags and release assets are
immutable: never move, delete, or recreate a published tag, and never replace an
asset for an existing version. A change that ships to users requires a new
version and a new release.

## Documentation

Update documentation in the same pull request when behaviour changes:

- [README.md](README.md) for installation, usage, and feature summaries
- [docs/configuration.md](docs/configuration.md) for configuration fields
- [docs/rules.md](docs/rules.md) for rule IDs and severities
- [docs/hooks.md](docs/hooks.md) for staged scanning and hooks
- [docs/external-scanners.md](docs/external-scanners.md) for integrations
- [docs/privacy.md](docs/privacy.md) and [docs/threat-model.md](docs/threat-model.md)
  for data handling and security assumptions
- [docs/VERIFICATION-REPORT.md](docs/VERIFICATION-REPORT.md) for audit evidence

Do not claim capabilities the code does not have. Document known limitations and
installation caveats honestly.

## License

UXDGuard is released under the [MIT License](LICENSE). By submitting a
contribution, you agree that it may be distributed under the same license and
confirm that you have the right to submit it.

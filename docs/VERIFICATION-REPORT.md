# UXDGuard v1.1.1 Verification Report

**Date:** 2026-11-02
**Version:** 1.1.1
**Branch:** main
**Tag:** `v1.1.1` (published; see Release and Tag State)
**Go:** go1.22.9 windows/amd64
**Module:** `github.com/arjunuxd/uxdguard`

## Scope

This report covers the code, tests, documentation, CLI behaviour, staged scanning, security hardening, and release packaging for UXDGuard v1.1.1. It includes the 2026-11-01 audit and the scan-feedback, pagination, and hardening changes that were made after the `v1.1.0` tag, so the published `v1.1.1` source and binaries contain them. No finding categories were removed merely to reduce output; detection was not weakened and all remaining static-analysis findings are documented as intentional or false positives.

## Validation Summary

| Check | Result |
| --- | --- |
| `gofmt -w .` / `gofmt -l .` | Clean |
| `go test -count=1 ./...` | Pass: 64 tests passed, 8 Windows symlink tests skipped, 0 failures, 8 packages |
| `go build ./...` | Pass |
| `go vet ./...` | Pass |
| `go mod verify` | Pass |
| `govulncheck ./...` | No vulnerabilities found |
| `gosec ./...` | 16 findings, all triaged as intentional or false positives (see Static Analysis) |
| Cross-build Windows amd64 | Pass |
| Cross-build Linux amd64 | Pass |
| Cross-build macOS arm64 | Pass |
| Linker version injection (`main.version`, `main.commit`, `main.buildDate`) | Pass, reports `1.1.1` and supplied metadata |
| Untagged local build | Pass, reports `dev` with VCS commit and date |
| Go module version detection | Pass, verified with a file-based module proxy, reporting `1.1.1` |
| JSON output | Pass, valid JSON with non-null arrays and no progress text |
| Progress output | Pass, staged lines on the default terminal path; silent for `--json` and `--no-progress` |
| Terminal pagination | Pass, 30 per page default, page totals, navigation, filters, and out-of-range clamp |
| Invalid `--page` / `--limit` / `--severity` | Pass, exits with usage/config status 2 |
| Rule disable / policy | Pass |
| Invalid or missing explicit config | Pass, exits with usage/config status 2 |
| Missing project root | Pass, exits with operational status 3 |
| Staged scan | Pass, detects index-only secret and skips an oversized staged blob |
| Staged deletion | Pass, deleted index entry is not scanned |
| Windows path output | Pass, findings use forward slashes |
| Cancellation | Pass |
| False-positive regressions | Pass |

## Code and Behaviour Changes Verified

### Scan progress feedback

- `uxdguard check` prints `Scanning <project>` before discovery.
- On an interactive terminal, the discovered-file counter updates in place and is erased before the report.
- After discovery it prints `N files found; running scanners...`.
- After the scan it prints `Scan completed in <duration>`.
- `--no-progress` suppresses all progress lines.
- `--json` keeps stdout machine-readable with no progress text.
- Staged scans print `Scanning staged changes...` and the elapsed time unless progress is disabled.

### Pagination

- Terminal reports show 30 findings per page by default.
- Every paginated report states the total, for example `Findings: (showing 1-30 of 65 total)` and `Page 1 of 3 (30 per page)`.
- A `Pages: N total. Navigate: ...` line shows the exact previous and next commands, including active `--severity` and `--rule` filters.
- Requesting a page past the end prints how many pages exist and shows the last page instead of an empty result.
- `--page` and `--limit` must be strictly positive integers; zero, negative, and non-numeric values exit with status 2.
- `--severity` must be one of `critical`, `high`, `medium`, `low`, or `info`; invalid values exit with status 2. Values are matched case-insensitively.

### Version metadata

- Linker-injected version, commit, and build date take priority.
- A build without linker flags reports `dev`, the real VCS revision, and the VCS commit timestamp.
- Go module metadata is used only when the main module path matches and the version is not `(devel)`.
- A file-based Go module proxy install of `github.com/arjunuxd/uxdguard/cmd/uxdguard@v1.1.1` reports `Version: 1.1.1`; commit and build date remain `unknown` because module archives do not carry VCS metadata.
- The public proxy entry for `v1.1.0` is stale because an earlier artifact was cached for that version on 2026-10-29. `v1.1.1` is a fresh module version, so that entry does not apply to it.
- Release output format:

```text
UXDGuard
Security Scanner
Version: 1.1.1
Commit: <commit>
Build date: <date>
```

### Rule configuration

- `rules.enable`, `rules.disable`, and `rules.severity` are applied centrally after scanning.
- Staged and external findings are filtered by the same settings.
- Disabling a scanner records a skip reason.
- Rule-level severity overrides are reflected in terminal and JSON output.

### False positives and detection accuracy

- `RegExp.prototype.exec` is not reported as command execution.
- Generic process execution is reported once as `UXD-SEC-207`.
- `UXD-SEC-300` is limited to actual shell execution.
- A normal `jwt.verify(...)` call is not reported as insecure JWT verification.
- Placeholder JWT secrets and encryption keys are not reported.
- Relative imports such as `../lib/file` are not reported as path traversal.
- Code-pattern scanners skip documentation files; secret detection still scans them.
- UI label, hint, message, title, and description assignments are not treated as credentials.
- A file with a very long minified line is scanned without producing a false operational error; the line buffer is derived from the configured maximum file size.

### Staged scanning and Git hooks

- Staged scanning reads the Git index, not working-tree edits.
- Staged deletions are excluded from the file list.
- Index paths are validated before temporary files are created: absolute paths, Windows volume paths, backslash-separated paths, `..` traversal, leading separators, pathspec-like leading colons, and Windows reserved device names (`CON`, `NUL`, `COM1`, `LPT1`, `CONIN$`, and variants with extensions or trailing dots/spaces) are rejected. The Windows checks apply on every host OS.
- Staged blob reads are capped at `scan.max_file_size_bytes`; an oversized blob is skipped without exhausting memory.
- Staged content uses a private `0700` temporary directory, `0700` subdirectories, and `0600` files.
- Hook installation honours `core.hooksPath` and linked worktrees.
- Hook executable paths are quoted, including paths containing spaces and single quotes.
- Unmanaged hooks are never overwritten or removed.
- Hook creation refuses any existing path, including a dangling symlink, before using exclusive creation. On Windows, `O_EXCL` alone would follow a dangling link and create its target, so a non-following `Lstat` check is performed first.

### Discovery and configuration

- Symlinks are ignored by default.
- With `follow_symlinks: true`, targets must resolve inside the project root and cycles terminate.
- A config path outside the project root is rejected.
- A config file that is a symlink resolving outside the project root is rejected, even if its lexical path is inside the root.
- `uxdguard init` refuses any existing path, including a dangling `.uxdguard.yaml` symlink, so it cannot write through the link. On Windows this requires a non-following `Lstat` before exclusive creation because `O_EXCL` follows a dangling link.
- Missing config keys receive defaults.
- Explicit zero or empty values are validated instead of being replaced silently.
- An explicit missing config file is rejected.
- A missing project root is rejected.

### Reporting and cancellation

- Findings are sorted and deduplicated.
- Finding paths use forward slashes on Windows.
- JSON collections are empty arrays rather than `null`.
- Filtered terminal output with zero matches reports `0 of N total` without a misleading next page.
- Terminal output sanitizes control characters in untrusted values such as repository file names and tool messages, preventing ANSI escape-sequence injection.
- A cancelled scan stops and reports cancellation.

### External adapters

- Gitleaks and Semgrep are registered in the CLI binary but disabled by default.
- Missing optional tools do not fail a scan unless `require: true`.
- Configured timeouts are parsed as Go durations rather than nanoseconds.
- External tools are invoked with explicit argument lists, not through a shell.
- Captured stdout and stderr are capped at 10 MiB per stream so a misbehaving tool cannot exhaust memory.
- Semgrep runs with `--metrics=off` and `--disable-version-check`, preventing metrics and update checks during a scan.
- External adapters run only when enabled explicitly in configuration.

## Static Analysis

`govulncheck` reports no known vulnerabilities in the module graph.

`gosec` reports 16 findings. All are intentional behaviour or false positives:

| Category | Locations | Disposition |
| --- | --- | --- |
| TLS fixture with `InsecureSkipVerify` | `tests/fixtures/source/insecure.go` | Test fixture used to prove detection; never shipped as production logic |
| Subprocess with variable arguments | `internal/gitutil`, `internal/external/runner.go` | Argument lists are passed separately; no shell is used. The executable is the local `git` binary or an explicitly configured external tool |
| File inclusion via variable | scanners, config loader, hooks, exclusive file creation | These paths are the intended scan/config/hook/output targets, and discovery or containment checks constrain them |
| File/directory permissions | staged temp dir, hook script, generated config | `0700` is correct for a private temporary directory, `0755` is required for an executable hook file, and `0644` is correct for a shareable project config |

No finding indicated a vulnerability that could be exploited by running a scan without explicit user action.

## Manual CLI Evidence

| Scenario | Observed result |
| --- | --- |
| Default scan of 65 findings | `Findings: (showing 1-30 of 65 total)`, `Page 1 of 3 (30 per page)` |
| `--page 3` | `showing 61-65 of 65 total`, `Page 3 of 3`, previous-page hint |
| `--page 99` | `Requested page 99, but only 3 page(s) available. Showing page 3 of 3.` |
| `--page 0` | Error `--page requires a positive integer, got "0"`, exit 2 |
| `--limit abc` | Error `--limit requires a positive integer, got "abc"`, exit 2 |
| `--severity severe` | Error listing valid severities, exit 2 |
| `--no-progress` | Report starts at `UXDGuard`; no progress lines |
| `--json` | Starts with `{`; valid JSON; no progress lines |
| Long minified line | Finding reported, no `error reading file`, status findings |
| Staged index secret | `UXD-SEC-100` at `app.js:1`, exit 1 |
| Staged oversized file | Skipped without error; only the small staged secret reported |
| Hook install/status/uninstall | Installed, detected, removed; unmanaged hooks untouched |
| `init` then `init` again | First creates config, second prints `config already exists` and exits 2 |
| `config validate` | Valid before and after `init` |
| Self-scan of the repository | 71 findings in fixtures/docs, 0 errors, paginated to 3 pages |
| Linker build `version` | `Version: 1.1.1`, supplied commit and date |
| Untagged local build | `Version: dev` with real VCS commit/date |

## Release and Tag State

- The annotated tag `v1.1.1` points to the release commit on `main`; the tag and the release commit match.
- The GitHub release `v1.1.1` is published, latest, non-draft, and non-prerelease, with six assets.
- All six assets match their `.sha256` files and GitHub digests, and the embedded build metadata reports version `1.1.1` and the tagged commit.
- The published `v1.1.1` source and binaries include the audit hardening, progress feedback, pagination, and staged-path fixes described in this report.
- The `v1.1.0` tag and its release assets were left unchanged. Published module versions and release tags are immutable, so the historical `v1.1.0` caveats below still apply only to that version.

## Known Limitations

- Detection remains pattern-based; there is no AST, data-flow, or taint analysis.
- False positives and false negatives are possible and require human review.
- Dependency CVE and full SCA coverage are not included in the built-in scanners.
- External integrations depend on separately installed tools and their behaviour. With no configured Semgrep ruleset, `semgrep --config auto` may fetch registry rules over the network; this is opt-in and disabled by default.
- Eight symlink-dependent tests are skipped on this Windows machine when symlink creation is unavailable; CI exercises them on Linux and macOS.
- `go test -race` cannot run on this Windows machine because no C compiler is installed; CI runs the race detector on Linux and macOS.
- `go install github.com/arjunuxd/uxdguard/cmd/uxdguard@v1.1.0` and `@latest` at the time of that release used a stale `proxy.golang.org` entry: a pre-release artifact was cached for that version on 2026-10-29, before the final `v1.1.0` tag was created. That artifact reports `Version: 0.1.0-dev`, and forcing `GOPROXY=direct` alone fails with a checksum mismatch because the checksum database record is permanent. Set `GOPRIVATE='github.com/arjunuxd/*'` (or `GONOSUMDB`) to fetch the tag directly from GitHub and skip the stale checksum record; if `v1.1.0` was already downloaded once, run `go clean -modcache` first. This affects only `v1.1.0`; `v1.1.1` is a fresh module version.
- The `v1.1.0` release assets predate the post-tag commits described above and do not include the new progress or pagination behaviour. The `v1.1.1` assets do include them.

## Conclusion

The v1.1.1 codebase, CLI, JSON output, staged scanning, cancellation, configuration handling, and cross-platform builds pass the validation above. The GitHub release is published as the latest stable release, and its binaries, checksums, and tag match the release commit. The previously published `v1.1.0` tag and assets remain unchanged.

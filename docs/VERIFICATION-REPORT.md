# UXDGuard v1.1.0 Verification Report

**Date:** 2026-10-31
**Version:** 1.1.0
**Branch:** main
**Tag:** `v1.1.0`
**Go:** go1.22.9 windows/amd64
**Module:** `github.com/arjunuxd/uxdguard`

## Scope

This report covers the code, tests, documentation, CLI behaviour, and release packaging for UXDGuard v1.1.0. The release is a stable production release. No finding categories were removed merely to reduce output; false-positive fixes are limited to cases where the matched text is not executable code or is clearly a placeholder.

## Validation Summary

| Check | Result |
| --- | --- |
| `gofmt -w .` / `gofmt -l .` | Clean |
| `go test -count=1 ./...` | Pass: 40 tests, 36 passed, 4 Windows symlink tests skipped, 0 failures |
| `go build ./...` | Pass |
| `go vet ./...` | Pass |
| `go mod verify` | Pass |
| Cross-build Windows amd64 | Pass |
| Cross-build Linux amd64 | Pass |
| Cross-build macOS arm64 | Pass |
| Linker version injection | Pass |
| Untagged local build | Pass, reports `dev` with VCS commit and date |
| Go module version detection | Pass, verified with a file-based module proxy and with a direct `GOPRIVATE` install of `...@v1.1.0`, both reporting `1.1.0` |
| JSON output | Pass, valid JSON with non-null arrays |
| Terminal pagination | Pass, including zero filtered results |
| Rule disable / policy | Pass |
| Invalid or missing explicit config | Pass, exits with usage/config status 2 |
| Missing project root | Pass, exits with operational status 3 |
| Staged scan | Pass, detects index-only secret |
| Staged deletion | Pass, deleted index entry is not scanned |
| Windows path output | Pass, findings use forward slashes |
| Cancellation | Pass |
| False-positive regressions | Pass |

## Code and Behaviour Changes Verified

### Version metadata

- Linker-injected version, commit, and build date take priority.
- A build without linker flags reports `dev`, the real VCS revision, and the VCS commit timestamp.
- Go module metadata is used only when the main module path matches and the version is not `(devel)`.
- A file-based Go module proxy install of `github.com/arjunuxd/uxdguard/cmd/uxdguard@v1.1.0` reports `Version: 1.1.0`; commit and build date remain `unknown` when the module proxy does not carry VCS metadata.
- The public proxy entry for `v1.1.0` is stale because an earlier artifact was cached for that version on 2026-10-29. A direct install with `GOPRIVATE='github.com/arjunuxd/*'` bypasses the stale proxy and checksum entry, builds the current tag, and reports `Version: 1.1.0`; see Known Limitations.
- Release output format:

```text
UXDGuard
Security Scanner
Version: 1.1.0
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

### Staged scanning and Git hooks

- Staged scanning reads the Git index, not working-tree edits.
- Staged deletions are excluded from the file list.
- Index paths are validated before temporary files are created.
- Staged content uses a private `0700` temporary directory and `0600` files.
- Hook installation honours `core.hooksPath` and linked worktrees.
- Hook executable paths are quoted, including paths containing spaces.
- Unmanaged hooks are never overwritten or removed.

### Discovery and configuration

- Symlinks are ignored by default.
- With `follow_symlinks: true`, targets must resolve inside the project root and cycles terminate.
- Missing config keys receive defaults.
- Explicit zero or empty values are validated instead of being replaced silently.
- An explicit missing config file is rejected.
- A config path outside the project root is rejected.
- A missing project root is rejected.

### Reporting and cancellation

- Findings are sorted and deduplicated.
- Finding paths use forward slashes on Windows.
- JSON collections are empty arrays rather than `null`.
- Filtered terminal output with zero matches reports `0 of N total` without a misleading next page.
- A cancelled scan stops and reports cancellation.

### External adapters

- Gitleaks and Semgrep are registered in the CLI binary but disabled by default.
- Missing optional tools do not fail a scan unless `require: true`.
- Configured timeouts are parsed as Go durations rather than nanoseconds.
- External tools are invoked with explicit argument lists, not through a shell.

## Manual CLI Evidence

| Scenario | Observed result |
| --- | --- |
| Linker build `version` | `Version: 1.1.0`, supplied commit and date |
| Linker build `--version` | `uxdguard version 1.1.0` |
| Untagged local build | `Version: dev` with real VCS commit/date |
| JSON scan of fixture project | 3 expected findings, 0 errors, no documentation finding |
| Zero-match severity filter | `Findings: 0 of 3 total match the selected filters` |
| Disabled rule | Exit 0, 0 findings, no `null` arrays |
| Staged index secret | `UXD-SEC-100` at `secret.txt:1`, exit 1 |
| Staged secret deletion | 0 findings, 0 errors, exit 0 |
| Malformed config | Error, exit 2 |
| Missing explicit config | Error, exit 2 |
| Missing project root | Error, exit 3 |

## Release Packaging

The release workflow builds one binary for each target:

- `uxdguard-windows-amd64.exe`
- `uxdguard-linux-amd64`
- `uxdguard-darwin-arm64`

Each binary has a matching `.sha256` file. The workflow verifies module integrity, runs vet and tests, injects version/commit/date, verifies the Linux binary version, checks all checksums, and publishes a non-draft, non-prerelease release with `make_latest: true`.

## Known Limitations

- Detection remains pattern-based; there is no AST, data-flow, or taint analysis.
- False positives and false negatives are possible and require human review.
- Dependency CVE and full SCA coverage are not included.
- External integrations depend on separately installed tools and their behaviour.
- Four Windows symlink tests are skipped locally when symlink creation is unavailable; CI exercises them on Linux and macOS.
- `go install github.com/arjunuxd/uxdguard/cmd/uxdguard@v1.1.0` and `@latest` use a stale `proxy.golang.org` entry: a pre-release artifact was cached for that version on 2026-10-29, before the final `v1.1.0` tag was created. That artifact reports `Version: 0.1.0-dev`, and forcing `GOPROXY=direct` alone fails with a checksum mismatch because the checksum database record is permanent. Set `GOPRIVATE='github.com/arjunuxd/*'` (or `GONOSUMDB`) to fetch the tag directly from GitHub and skip the stale checksum record; this was verified to install the current code and report `Version: 1.1.0`. The published release binaries and source builds from the `v1.1.0` tag are unaffected.

## Conclusion

The v1.1.0 codebase, CLI, JSON output, staged scanning, cancellation, configuration handling, and cross-platform builds pass the validation above, and the GitHub release is published as the latest stable release. The default Go module proxy entry for `v1.1.0` is stale; the documented `GOPRIVATE` install command, the published binaries, and source builds from the tag all produce correct code.

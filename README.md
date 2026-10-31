# UXDGuard

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22.9%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-555555)](#supported-platforms)
[![Release](https://img.shields.io/badge/Release-v1.1.0-2ea44f)](https://github.com/arjunuxd/uxdguard/releases/latest)
[![Local First](https://img.shields.io/badge/Local--First-No%20Telemetry-orange)](docs/privacy.md)

**Build with AI. Commit with confidence.**

UXDGuard is a lightweight, local-first security scanner for developers. It detects exposed secrets, insecure coding patterns, and common security vulnerabilities in source code and staged Git changes before they are committed.

## Overview

UXDGuard runs entirely on your local machine. Built-in scanning requires no telemetry, no source uploads, and no network access. It reads files statically, reports actionable findings, and never executes scanned project code.

## Key Features

- **Secret detection** - API keys, access tokens, private keys, credentials, and database connection strings
- **Vulnerability scanning** - SQL injection, command injection, XSS, SSRF, path traversal, weak cryptography, insecure JWT handling, and more
- **Password safety** - Plaintext password handling and insecure password comparisons
- **Staged scanning** - Inspect the Git index with `uxdguard check --staged`
- **Git hooks** - Install a pre-commit hook that honors `core.hooksPath` and linked worktrees
- **Configurable rules** - Enable, disable, or change the severity of individual rules or scanners
- **Deterministic output** - Stable finding order, deduplication, and forward-slash paths on every platform
- **Pagination and filtering** - Navigate terminal output with `--page`, `--limit`, `--severity`, and `--rule`
- **JSON output** - Complete, machine-readable reports for CI pipelines
- **Cancellation** - Ctrl+C stops a scan promptly and reports cancellation
- **Symlink controls** - Optional symlink traversal that is confined to the project root and protected against cycles
- **Optional integrations** - Extend detection with Gitleaks and Semgrep when explicitly configured

## Supported Platforms

| Operating System | Architecture | Release binary |
|---|---|---|
| Windows 10/11 | amd64 | `uxdguard-windows-amd64.exe` |
| Linux | amd64 | `uxdguard-linux-amd64` |
| macOS | Apple Silicon (arm64) | `uxdguard-darwin-arm64` |

Other Go-supported targets can be built from source.

## Requirements

You do not need Go or any other dependency to run a release binary.

| Dependency | Needed for |
|---|---|
| [Go](https://go.dev/) 1.22.9 or later | Building from source or installing with `go install` |
| [Git](https://git-scm.com/) | `--staged` scans and Git hooks only |

## Install

### Option 1: Download a release binary (recommended)

Download from the [latest release](https://github.com/arjunuxd/uxdguard/releases/latest). Every binary has a matching `.sha256` file.

#### Windows (PowerShell)

```powershell
$v = "v1.1.0"
$u = "https://github.com/arjunuxd/uxdguard/releases/download/$v"

Invoke-WebRequest "$u/uxdguard-windows-amd64.exe" -OutFile uxdguard.exe
Invoke-WebRequest "$u/uxdguard-windows-amd64.exe.sha256" -OutFile uxdguard.exe.sha256

Get-Content .\uxdguard.exe.sha256
(Get-FileHash .\uxdguard.exe -Algorithm SHA256).Hash.ToLower()
```

The two hashes must match. Then run it:

```powershell
.\uxdguard.exe version
```

To call it from any directory, move `uxdguard.exe` into a folder on your `PATH`, for example:

```powershell
New-Item -ItemType Directory -Force "$env:LOCALAPPDATA\uxdguard" | Out-Null
Move-Item .\uxdguard.exe "$env:LOCALAPPDATA\uxdguard\uxdguard.exe"
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:LOCALAPPDATA\uxdguard", "User")
```

Open a new terminal, then check `uxdguard version`.

#### Linux (x86_64)

```bash
V=v1.1.0
U="https://github.com/arjunuxd/uxdguard/releases/download/$V"

curl -LO "$U/uxdguard-linux-amd64"
curl -LO "$U/uxdguard-linux-amd64.sha256"

sha256sum -c uxdguard-linux-amd64.sha256
chmod +x uxdguard-linux-amd64
sudo mv uxdguard-linux-amd64 /usr/local/bin/uxdguard

uxdguard version
```

#### macOS (Apple Silicon)

```bash
V=v1.1.0
U="https://github.com/arjunuxd/uxdguard/releases/download/$V"

curl -LO "$U/uxdguard-darwin-arm64"
curl -LO "$U/uxdguard-darwin-arm64.sha256"

shasum -a 256 -c uxdguard-darwin-arm64.sha256
chmod +x uxdguard-darwin-arm64
xattr -d com.apple.quarantine uxdguard-darwin-arm64 2>/dev/null || true
sudo mv uxdguard-darwin-arm64 /usr/local/bin/uxdguard

uxdguard version
```

### Option 2: Install with Go

This release needs one extra variable because Go's public proxy still holds an
older artifact for the `v1.1.0` version (see [Known Installation Issue](#known-installation-issue)).
`GOPRIVATE` makes Go fetch the tag directly from GitHub.

If you already ran the plain `go install` once, Go saved the stale copy in its
module cache; clear it first:

```bash
go clean -modcache
```

This removes all downloaded Go modules; Go re-downloads them as needed.

Linux and macOS:

```bash
GOPRIVATE='github.com/arjunuxd/*' go install github.com/arjunuxd/uxdguard/cmd/uxdguard@v1.1.0
```

Windows PowerShell:

```powershell
$env:GOPRIVATE = 'github.com/arjunuxd/*'
go install github.com/arjunuxd/uxdguard/cmd/uxdguard@v1.1.0
```

The binary is written to `$(go env GOPATH)/bin` (`%USERPROFILE%\go\bin` on
Windows). Add that directory to your `PATH`, then run:

```bash
uxdguard version
```

Check which executable your shell finds:

```powershell
where.exe uxdguard
```

If it points somewhere other than your Go bin directory, remove the old copy
or put the Go bin directory earlier on `PATH`.

Expected output:

```text
UXDGuard
Security Scanner
Version: 1.1.0
Commit: unknown
Build date: unknown
```

`Commit` and `Build date` are `unknown` for module installs because Go module
archives do not carry VCS metadata. Release binaries do include both.

### Option 3: Build from source

```bash
git clone https://github.com/arjunuxd/uxdguard.git
cd uxdguard
git checkout v1.1.0
go build ./cmd/uxdguard
```

On Windows the output is `uxdguard.exe`. A source build reports `dev` with the
real Git commit and build date.

## Use UXDGuard

### Step 1: Open a terminal in your project

```bash
cd /path/to/your-project
```

### Step 2: Run a scan

```bash
uxdguard check
```

UXDGuard scans the current directory recursively. No configuration is required.

### Step 3: Read the report

Example output:

```text
UXDGuard
Security Scanner · 1.1.0

Project: /path/to/your-project

Findings: (showing 1-2 of 2 total)
Page 1 of 1

HIGH     UXD-SEC-200  Use of dynamic code evaluation detected
         app.js:1
         eval() or Function() can lead to code injection if used with untrusted input.
         Remediation: Avoid dynamic code evaluation; use safe parsers or explicit logic.

HIGH     UXD-SEC-100  Potential AWS access key detected
         config.js:1
         A string matching AWS access key pattern was found.
         Remediation: Remove hardcoded credentials and use environment variables or a secure secret manager.

Summary: 0 critical, 2 high, 0 medium, 0 low | 2 findings, 0 skipped, 0 errors
Status: Findings detected
```

Each finding shows its severity, rule ID, `file:line`, description, and a
remediation hint. Secret values are redacted.

### Step 4: Fix the findings and rerun

```bash
uxdguard check
```

Repeat until the scan is clean. Exit codes make this usable in CI:

| Exit code | Meaning |
|---|---|
| `0` | No findings (or all findings below the policy threshold) |
| `1` | Findings at or above the policy threshold |
| `2` | Invalid usage or configuration |
| `3` | Operational failure (for example, an unreadable project root) |

### Step 5: Check changes before you commit

Scan only what is staged in Git:

```bash
git add .
uxdguard check --staged
```

Install a pre-commit hook so this happens automatically:

```bash
uxdguard hook install
uxdguard hook status
```

Remove the hook later with `uxdguard hook uninstall`.

### Other useful commands

```bash
uxdguard check --json                       # full machine-readable report
uxdguard check --staged                     # staged Git changes only
uxdguard check --severity high              # only high severity findings
uxdguard check --rule UXD-SEC-101           # one rule
uxdguard check --page 2 --limit 50          # paginate terminal output
uxdguard check --project-root /path/to/repo # scan another directory
uxdguard init                               # create .uxdguard.yaml
uxdguard config validate                    # validate configuration
uxdguard rules                              # list scanners
uxdguard version                            # version, commit, build date
```

### Command reference

| Command | Description |
|---|---|
| `uxdguard check` | Scan all files in the project root |
| `uxdguard check --staged` | Scan files staged in Git |
| `uxdguard check --json` | Output results in JSON format |
| `uxdguard init` | Create a starter `.uxdguard.yaml` |
| `uxdguard config validate` | Validate configuration |
| `uxdguard rules` | List registered scanners |
| `uxdguard hook [install\|uninstall\|status]` | Manage pre-commit hooks |
| `uxdguard version` | Display version information |
| `uxdguard uninstall` | Uninstall UXDGuard |

JSON output is always complete and is not affected by pagination options.

## Configuration

Generate a starter configuration:

```bash
uxdguard init
```

For the full schema, defaults, and rule controls, see [Configuration](docs/configuration.md).

## Default Exclusions

UXDGuard excludes dependency, build, cache, and editor directories by default, including `.git`, `node_modules`, `vendor`, `dist`, `build`, `.next`, `.nuxt`, `.svelte-kit`, `.turbo`, `.cache`, `coverage`, `target`, `out`, `.venv`, `venv`, `__pycache__`, `.terraform`, `.idea`, and `.vscode`. Dependency manifests and lockfiles are not excluded.

## Version Reporting

`uxdguard version` reports the version, commit, and build date. Values come from:

1. Linker flags used by the official release build;
2. Go module metadata when installed with `go install ...@version`;
3. The truthful fallback `dev` for an untagged local build.

## Known Installation Issue

A pre-release artifact was published under the `v1.1.0` version on 2026-10-29
and later replaced. `proxy.golang.org` and `sum.golang.org` keep that earlier
artifact permanently, so plain `go install ...@v1.1.0` builds old code, and
`GOPROXY=direct` alone fails checksum verification.

Use one of these instead:

- the release binaries (Option 1),
- the `GOPRIVATE` install command (Option 2),
- a source build from tag `v1.1.0` (Option 3).

This affects only the Go module path for `v1.1.0`. Full details are in
[docs/VERIFICATION-REPORT.md](docs/VERIFICATION-REPORT.md).

## Privacy and Security

- **Local-first execution** - All built-in analysis runs on your machine
- **No telemetry** - No usage data is collected
- **No code execution** - Scanned files are read statically
- **No data exfiltration** - No source code or findings are transmitted
- **Redacted output** - Secret values are not printed in reports
- **Optional integrations only** - External tools require explicit configuration

For details, see [Privacy](docs/privacy.md), [Threat Model](docs/threat-model.md), and [Security Policy](docs/security-policy.md).

## Documentation

| Document | Description |
|---|---|
| [Configuration](docs/configuration.md) | Schema, defaults, and rule controls |
| [Rules](docs/rules.md) | Built-in rule IDs and classifications |
| [Git Hooks](docs/hooks.md) | Staged scanning and pre-commit hook setup |
| [External Scanners](docs/external-scanners.md) | Optional Gitleaks and Semgrep integration |
| [Privacy](docs/privacy.md) | Data handling and privacy practices |
| [Threat Model](docs/threat-model.md) | Threat assumptions and scope |
| [Security Policy](docs/security-policy.md) | Vulnerability reporting guidelines |
| [Architecture](architecture.md) | Internal architecture and design |

## Development

```bash
gofmt -w .
go test -count=1 ./...
go build ./...
go vet ./...
go mod verify
```

## Limitations

- Detection is pattern-based. It does not perform AST, data-flow, or taint analysis.
- False positives and false negatives are possible; review findings in context.
- Built-in rules have different language coverage; see [Rules](docs/rules.md).
- Dependency CVE scanning is not included.
- Documentation files are skipped by code-pattern rules because examples are not executable, while secret detection still applies to them.

## Uninstall

```bash
uxdguard uninstall
```

If self-removal is blocked, delete the `uxdguard` executable you installed.

## License

UXDGuard is released under the [MIT License](LICENSE).

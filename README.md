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

| Operating System | Support |
|---|---|
| Windows | Supported |
| Linux | Supported |
| macOS | Supported |

Release binaries are published for Windows amd64, Linux amd64, and macOS arm64. Other Go-supported targets can be built from source.

## Requirements

| Dependency | Version | Purpose |
|---|---|---|
| [Go](https://go.dev/) | 1.22.9 or later | Installation from source |
| [Git](https://git-scm.com/) | Any recent version | Staged scanning and Git hooks |

## Installation

### Install with Go (Recommended)

```bash
go install github.com/arjunuxd/uxdguard/cmd/uxdguard@latest
```

Verify the installation:

```bash
uxdguard version
uxdguard rules
```

Ensure your Go bin directory (`$(go env GOPATH)/bin`) is on your `PATH`.

### Pre-built Binaries

Download official binaries and SHA-256 checksums from the [latest release](https://github.com/arjunuxd/uxdguard/releases/latest).

```bash
sha256sum -c uxdguard-linux-amd64.sha256
```

### Build from Source

```bash
git clone https://github.com/arjunuxd/uxdguard.git
cd uxdguard
go build ./cmd/uxdguard
```

On Windows the output is `uxdguard.exe`.

## Version Reporting

`uxdguard version` reports the version, commit, and build date. The value comes from:

1. Linker flags used by the official release build;
2. Go module metadata when installed with `go install ...@version`;
3. The truthful fallback `dev` for an untagged local build.

Release binaries therefore report `1.1.0`, while `go install ...@latest` reports the module version selected by Go.

## Quick Start

```bash
# Navigate to your project directory
cd /path/to/your-project

# Run a full scan
uxdguard check

# Scan only staged changes
uxdguard check --staged

# Generate JSON output
uxdguard check --json

# Filter and paginate results
uxdguard check --severity high --page 1 --limit 20

# Initialize and validate configuration
uxdguard init
uxdguard config validate

# Manage Git hooks
uxdguard hook status
uxdguard hook install
uxdguard hook uninstall
```

## Usage

### Command Reference

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

### Pagination and Filtering

```bash
uxdguard check                       # 20 findings per page by default
uxdguard check --limit 50
uxdguard check --page 2
uxdguard check --severity critical
uxdguard check --rule UXD-SEC-101
```

JSON output is always complete and is not affected by pagination options.

### Configuration

Generate a starter configuration:

```bash
uxdguard init
```

For the full schema, defaults, and rule controls, see [Configuration](docs/configuration.md).

### Default Exclusions

UXDGuard excludes dependency, build, cache, and editor directories by default, including `.git`, `node_modules`, `vendor`, `dist`, `build`, `.next`, `.nuxt`, `.svelte-kit`, `.turbo`, `.cache`, `coverage`, `target`, `out`, `.venv`, `venv`, `__pycache__`, `.terraform`, `.idea`, and `.vscode`. Dependency manifests and lockfiles are not excluded.

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

If self-removal is blocked, remove the `uxdguard` executable from your Go bin directory manually.

## License

UXDGuard is released under the [MIT License](LICENSE).

# UXDGuard

**Build with AI. Commit with confidence.**

UXDGuard is a lightweight, local-first security scanner for developers. It helps identify exposed secrets, insecure coding patterns, and common security issues in source code and staged Git changes before they are committed.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22.9+-00ADD8?logo=go&logoColor=white)](https://go.dev/)

## Overview

UXDGuard runs entirely on your local machine with no telemetry, no source uploads, and no mandatory network access. It uses pattern-based detection to provide fast, actionable feedback during development and in pre-commit workflows.

## Key Features

- **Secret Detection** - Identifies API keys, access tokens, private keys, credentials, and database connection strings
- **Security Checks** - Detects SQL injection, command injection, XSS, SSRF, path traversal, weak cryptography, and more
- **Password Safety** - Flags plaintext password handling and insecure password comparisons
- **Staged Scanning** - Scan only Git-staged changes with `uxdguard check --staged`
- **Git Hooks** - Integrate checks into your pre-commit workflow
- **Configurable Rules** - Enable, disable, or adjust severities via `.uxdguard.yaml`
- **Smart Exclusions** - Automatically skips dependencies, build outputs, and generated files while preserving manifests
- **Pagination & Filtering** - Navigate large result sets with page-based output (20 results per page by default)
- **JSON Output** - Machine-readable results for CI/CD integration
- **Optional Integrations** - Support for Gitleaks and Semgrep when explicitly enabled
- **Local-First** - No code execution, no telemetry, no external data transfer

## Supported Platforms

- Windows
- Linux
- macOS

## Requirements

- [Go](https://go.dev/) 1.22.9 or later (for installation from source)
- [Git](https://git-scm.com/) (for staged scanning and Git hooks)

## Installation

### Install with Go

```bash
go install github.com/arjunuxd/uxdguard/cmd/uxdguard@latest
```

Ensure your Go bin directory is in your `PATH`. Verify the installation:

```bash
uxdguard version
uxdguard rules
```

### Build from Source

```bash
git clone https://github.com/arjunuxd/uxdguard.git
cd uxdguard
go build -o uxdguard ./cmd/uxdguard
```

On Windows, use `uxdguard.exe`. Add the binary to your `PATH` or run it from the build directory.

### Pre-built Binaries

Pre-built binaries and checksums are available on the [GitHub Releases](https://github.com/arjunuxd/uxdguard/releases) page when published.

## Quick Start

```bash
# Scan the current project
uxdguard check

# Scan only staged changes
uxdguard check --staged

# Generate JSON output
uxdguard check --json

# View results with pagination and filters
uxdguard check --page 1 --limit 20 --severity high

# Initialize configuration
uxdguard init

# Validate configuration
uxdguard config validate

# Manage pre-commit hooks
uxdguard hook status
uxdguard hook install
uxdguard hook uninstall
```

## Usage

### Pagination and Filtering

By default, terminal output shows 20 findings per page. You can navigate and filter results:

```bash
uxdguard check --page 2 --limit 50
uxdguard check --severity critical
uxdguard check --rule UXD-SEC-101
```

JSON output always includes the complete result set without truncation.

### Default Exclusions

The following directories are excluded by default: `.git`, `node_modules`, `vendor`, `dist`, `build`, `.next`, `.nuxt`, `.svelte-kit`, `.turbo`, `.cache`, `coverage`, `target`, `out`, `.venv`, `venv`, `__pycache__`, `.terraform`, `.idea`, `.vscode`. Dependency manifests and lockfiles remain discoverable.

### Configuration

Create a starter configuration file:

```bash
uxdguard init
```

Edit `.uxdguard.yaml` in your project root as needed, then validate it:

```bash
uxdguard config validate
```

Refer to [docs/configuration.md](docs/configuration.md) for the full schema and options.

## Privacy and Security

- **Local-first**: Core scanning runs locally with no mandatory network access
- **No telemetry**: No analytics, tracking, or data collection
- **No code execution**: Scanned files are never executed
- **Safe output**: Findings are reported without exposing raw secret values
- **Optional integrations**: External tools are disabled by default
- **Pattern-based**: Detection is heuristic-based; false positives and negatives are possible

For details, see [docs/privacy.md](docs/privacy.md), [docs/threat-model.md](docs/threat-model.md), and [docs/security-policy.md](docs/security-policy.md).

## Documentation

| Document | Description |
| --- | --- |
| [docs/configuration.md](docs/configuration.md) | Configuration options and defaults |
| [docs/rules.md](docs/rules.md) | Built-in security rules and categories |
| [docs/hooks.md](docs/hooks.md) | Staged scanning and Git hooks |
| [docs/external-scanners.md](docs/external-scanners.md) | Gitleaks and Semgrep integrations |
| [docs/privacy.md](docs/privacy.md) | Privacy and data handling |
| [docs/threat-model.md](docs/threat-model.md) | Threat model and assumptions |
| [docs/security-policy.md](docs/security-policy.md) | Reporting security issues |
| [architecture.md](architecture.md) | System architecture overview |

## Development

```bash
# Build all packages
go build ./...

# Run tests
go test ./...

# Run static analysis
go vet ./...

# Format code
gofmt -l .
```

## Uninstall

```bash
uxdguard uninstall
```

If self-removal is blocked by your operating system, remove the `uxdguard` binary manually from your Go bin directory.

## License

[MIT License](LICENSE)
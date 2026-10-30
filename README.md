# UXDGuard

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22.9%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-555555)](#supported-platforms)
[![Status](https://img.shields.io/badge/Status-v1.1.0%20Flash-2ea44f)](https://github.com/arjunuxd/uxdguard)
[![Local First](https://img.shields.io/badge/Local--First-No%20Telemetry-orange)](docs/privacy.md)

**Build with AI. Commit with confidence.**

UXDGuard is a lightweight, local-first security scanner for developers. It identifies exposed secrets, insecure coding patterns, and common security vulnerabilities in source code and staged Git changes before they are committed.

## Overview

Designed for speed and privacy, UXDGuard runs entirely on your local machine with no telemetry, no source uploads, and no mandatory network access. It provides fast, actionable feedback to help enforce secure development practices without relying on cloud services.

## Key Features

- **Secret Detection** – Detects API keys, access tokens, private keys, credentials, and database connection strings
- **Vulnerability Scanning** – Covers SQL injection, command injection, XSS, SSRF, path traversal, weak cryptography, and more
- **Password Safety** – Identifies plaintext password handling and insecure password comparisons
- **Staged Scanning** – Scan only Git-staged changes with `uxdguard check --staged`
- **Git Hooks** – Automate security checks with pre-commit hooks
- **Configurable Rules** – Enable, disable, or adjust rule severity via `.uxdguard.yaml`
- **Smart Exclusions** – Automatically skips dependencies, build artifacts, and generated files while preserving manifests and lockfiles
- **Pagination & Filtering** – Navigate large result sets with paginated terminal output (20 results per page by default)
- **JSON Output** – Complete, machine-readable reports for CI/CD pipelines
- **Optional Integrations** – Extend detection with Gitleaks and Semgrep when explicitly configured
- **Local-First Architecture** – No code execution, no telemetry, and no external data transfer

## Supported Platforms

| Operating System | Support |
|---|---|
| Windows | Supported |
| Linux | Supported |
| macOS | Supported |

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

Ensure your Go bin directory (`$(go env GOPATH)/bin`) is added to your system `PATH`.

### Build from Source

```bash
git clone https://github.com/arjunuxd/uxdguard.git
cd uxdguard
go build ./cmd/uxdguard
```

For Windows, replace `uxdguard` with `uxdguard.exe`.

### Pre-built Binaries

Download official binaries and checksums from the [GitHub Releases](https://github.com/arjunuxd/uxdguard/releases) page.

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
| `uxdguard init` | Create `.uxdguard.yaml` configuration file |
| `uxdguard config validate` | Validate configuration file |
| `uxdguard rules` | List all active security rules |
| `uxdguard hook [install\|uninstall\|status]` | Manage pre-commit hooks |
| `uxdguard version` | Display version information |
| `uxdguard uninstall` | Uninstall UXDGuard |

### Pagination & Filtering

```bash
# Default pagination (20 per page)
uxdguard check

# Custom page size
uxdguard check --limit 50

# Navigate to a specific page
uxdguard check --page 2

# Filter by severity
uxdguard check --severity critical

# Filter by rule ID
uxdguard check --rule UXD-SEC-101
```

JSON output remains complete regardless of pagination settings.

### Configuration

Generate a starter configuration:

```bash
uxdguard init
```

For detailed options, see [Configuration Documentation](docs/configuration.md).

### Default Exclusions

UXDGuard excludes the following by default: `.git`, `node_modules`, `vendor`, `dist`, `build`, `.next`, `.nuxt`, `.svelte-kit`, `.turbo`, `.cache`, `coverage`, `target`, `out`, `.venv`, `venv`, `__pycache__`, `.terraform`, `.idea`, `.vscode`. Dependency manifests and lockfiles are always preserved.

## Privacy & Security

- **Local-first execution** – All analysis runs on your machine
- **No telemetry or tracking** – No usage data is collected
- **No code execution** – Scanned files are read statically
- **No data exfiltration** – No source code or findings are transmitted
- **Safe output** – Secret values are redacted in all reports
- **Optional integrations only** – External tools require explicit configuration

For more information, refer to [Privacy Policy](docs/privacy.md), [Threat Model](docs/threat-model.md), and [Security Policy](docs/security-policy.md).

## Documentation

| Document | Description |
|---|---|
| [Configuration](docs/configuration.md) | Schema, defaults, and configuration options |
| [Rules](docs/rules.md) | Built-in security rules and classifications |
| [Git Hooks](docs/hooks.md) | Staged scanning and pre-commit hook setup |
| [External Scanners](docs/external-scanners.md) | Integration with Gitleaks and Semgrep |
| [Privacy](docs/privacy.md) | Data handling and privacy practices |
| [Threat Model](docs/threat-model.md) | Threat assumptions and scope |
| [Security Policy](docs/security-policy.md) | Vulnerability reporting guidelines |
| [Architecture](architecture.md) | Internal architecture and design |

## Development

```bash
# Build
go build ./...

# Test
go test ./...

# Lint
go vet ./...

# Format
gofmt -l .
```

## Uninstall

```bash
uxdguard uninstall
```

If self-removal is blocked, manually remove the `uxdguard` executable from your Go bin directory.

## License

UXDGuard is released under the [MIT License](LICENSE).
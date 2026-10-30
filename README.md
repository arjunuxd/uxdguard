<div align="center">

# UXDGuard

**Build with AI. Commit with confidence.**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22.9+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](#supported-platforms)

*A lightweight, local-first security scanner that helps developers catch secrets, insecure patterns, and common vulnerabilities before they reach production.*

[Features](#-features) • [Installation](#-installation) • [Quick Start](#-quick-start) • [Usage](#-usage) • [Documentation](#-documentation) • [Contributing](#-contributing)

</div>

## ✨ Features

- **🔐 Secret Detection** – Catch API keys, tokens, private keys, database credentials, and connection strings
- **🛡️ Security Analysis** – Detect SQL injection, XSS, command injection, SSRF, weak crypto, and more
- **🔑 Password Safety** – Identify plaintext password handling and insecure comparisons
- **📦 Staged Scanning** – Review only Git-staged changes with `--staged` before committing
- **🪝 Git Hooks** – Automate checks with pre-commit hooks
- **⚙️ Configurable** – Fine-tune rules and severities via `.uxdguard.yaml`
- **🚫 Smart Exclusions** – Auto-skip `node_modules`, `.next`, `dist`, `build`, and other generated directories
- **📄 Pagination & Filtering** – Navigate large result sets with ease (20 results per page by default)
- **📊 JSON Output** – CI/CD ready with complete, machine-readable reports
- **🔌 Optional Integrations** – Extend with [Gitleaks](https://github.com/gitleaks/gitleaks) and [Semgrep](https://github.com/semgrep/semgrep)
- **🏠 Local-First** – Runs entirely offline, no telemetry, no uploads, no code execution

## 🎯 Supported Platforms

| OS | Status |
| --- | --- |
| Windows | ✅ Supported |
| Linux | ✅ Supported |
| macOS | ✅ Supported |

## 📋 Prerequisites

- [Go](https://go.dev/) 1.22.9 or higher (for building from source)
- [Git](https://git-scm.com/) (required for staged scanning and hooks)

## 🚀 Installation

### Option 1: Install with Go (Recommended)

```bash
go install github.com/arjunuxd/uxdguard/cmd/uxdguard@latest
```

Ensure your Go binary directory is in your `PATH`:

<details>
<summary><b>Windows (PowerShell)</b></summary>

```powershell
$env:PATH += ";$(go env GOPATH)\bin"
# Add to your profile for persistence
```

</details>

<details>
<summary><b>Linux / macOS (Bash/Zsh)</b></summary>

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
# Add to ~/.bashrc, ~/.zshrc, or ~/.profile
```

</details>

Verify your installation:

```bash
uxdguard version
uxdguard rules
```

### Option 2: Build from Source

```bash
git clone https://github.com/arjunuxd/uxdguard.git
cd uxdguard
go build ./cmd/uxdguard
```

On Windows, the binary will be `uxdguard.exe`. Add it to your `PATH` or run it directly.

### Option 3: Download Binaries

Pre-built binaries are available on the [Releases](https://github.com/arjunuxd/uxdguard/releases) page when published.

## ⚡ Quick Start

```bash
# Navigate to your project
cd /path/to/your-project

# Run a full scan
uxdguard check

# Scan only staged changes (ideal for pre-commit)
uxdguard check --staged

# Get JSON output for CI
uxdguard check --json

# View paginated results with filters
uxdguard check --page 1 --limit 20 --severity high

# List all security rules
uxdguard rules
```

## 💡 Usage

### Commands

| Command | Description |
| --- | --- |
| `uxdguard check` | Scan all files in the current project |
| `uxdguard check --staged` | Scan only Git-staged changes |
| `uxdguard check --json` | Output results in JSON format |
| `uxdguard init` | Create a starter `.uxdguard.yaml` config |
| `uxdguard config validate` | Validate your configuration |
| `uxdguard rules` | Display all active security rules |
| `uxdguard hook install` | Install pre-commit hook |
| `uxdguard hook status` | Check hook status |
| `uxdguard hook uninstall` | Remove pre-commit hook |
| `uxdguard version` | Show version information |
| `uxdguard uninstall` | Uninstall UXDGuard |

### Pagination & Filtering

```bash
# Default: 20 results per page
uxdguard check

# Custom page size
uxdguard check --limit 50

# Navigate pages
uxdguard check --page 2

# Filter by severity
uxdguard check --severity critical

# Filter by specific rule
uxdguard check --rule UXD-SEC-101

# Combine filters
uxdguard check --severity high --limit 10
```

**Note:** JSON output (`--json`) always returns complete results without pagination.

### Configuration

Create and customize your configuration:

```bash
uxdguard init
```

This generates `.uxdguard.yaml` in your project root. Validate it with:

```bash
uxdguard config validate
```

See [Configuration Guide](docs/configuration.md) for all available options.

### Default Exclusions

UXDGuard automatically skips common directories to improve performance:

```
.git, node_modules, vendor, dist, build, .next, .nuxt, .svelte-kit,
.turbo, .cache, coverage, target, out, .venv, venv, __pycache__,
.terraform, .idea, .vscode
```

Dependency manifests (`package.json`, `go.mod`, `requirements.txt`, `pom.xml`, etc.) are always preserved.

## 🔒 Privacy & Security

- **Local-First** – All scanning happens on your machine. No data leaves your system.
- **No Telemetry** – Zero tracking, analytics, or usage data collection.
- **No Code Execution** – Files are read statically, never executed.
- **Safe Reporting** – Secret values are redacted from all outputs.
- **Optional Integrations** – External tools are disabled by default and only used if explicitly enabled.
- **Open Source** – Fully auditable codebase.

> **Note:** UXDGuard uses pattern-based detection and may produce false positives or miss issues. It's a helpful guardrail, not a replacement for comprehensive security reviews.

## 📚 Documentation

| Document | Description |
| --- | --- |
| [Configuration](docs/configuration.md) | Schema, defaults, and options |
| [Rules](docs/rules.md) | Complete list of built-in security rules |
| [Git Hooks](docs/hooks.md) | Pre-commit hook setup and usage |
| [External Scanners](docs/external-scanners.md) | Gitleaks and Semgrep integration |
| [Privacy](docs/privacy.md) | Data handling and privacy policy |
| [Threat Model](docs/threat-model.md) | Security assumptions and scope |
| [Security Policy](docs/security-policy.md) | Reporting vulnerabilities |
| [Architecture](architecture.md) | Technical design overview |

## 🛠️ Development

```bash
# Build
go build ./...

# Run tests
go test ./...

# Static analysis
go vet ./...

# Format code
gofmt -l .
```

## 🗑️ Uninstall

```bash
uxdguard uninstall
```

If self-removal is blocked by your OS, manually delete the `uxdguard` binary from your Go bin directory.

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

<div align="center">

**Made with ❤️ for developers who build with AI**

</div>
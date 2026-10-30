# UXDGuard

**Build with AI. Commit with confidence.**

UXDGuard is a lightweight, local-first security scanner for developers. It helps identify potential secrets, insecure coding patterns, and common security issues in source code and staged Git changes before committing.

## Status

- **Version:** `1.1.0 Flash`
- **Release:** Release candidate (`1.1.0-rc.1`)
- **License:** MIT

## Key Features

- **Secret detection** – Identifies common patterns for tokens, keys, private keys, credentials, and connection strings.
- **Source-code security checks** – Detects selected insecure patterns (e.g., `eval()`, disabled TLS verification, unsafe deserialization, weak cryptography, and more).
- **Advanced checks** – Includes additional checks for injection, authentication, cookies, SSRF, and related issues.
- **Password handling checks** – Detects plaintext password handling and insecure password comparison patterns.
- **Staged Git scanning** – Scan only staged changes with `uxdguard check --staged`.
- **Git hooks** – Install, check, or remove a pre-commit hook.
- **Configurable rules** – Enable/disable rules and override severities via `.uxdguard.yaml`.
- **Large result handling** – Pagination, filtering by severity/rule, and complete JSON exports for large finding sets.
- **Expanded exclusions** – Common dependency and generated directories are excluded by default while dependency manifests are preserved.
- **Machine-readable output** – JSON output for CI and automation.
- **Optional external integrations** – Gitleaks and Semgrep when explicitly enabled and installed.
- **Local-first** – No source code uploads, no telemetry, and no mandatory network access for core scanning.

## Supported Platforms

- Windows
- Linux
- macOS

## Prerequisites

- [Go](https://go.dev/) 1.22.9+ (for building from source or `go install`)
- [Git](https://git-scm.com/) (required for staged scanning and hooks)

## Installation

### Option 1: Install with Go (recommended)

```bash
go install github.com/arjunuxd/uxdguard/cmd/uxdguard@latest
```

Ensure your Go bin directory is in your `PATH`:

- **Windows (PowerShell):** `$(go env GOPATH)\bin`
- **Linux/macOS:** `$(go env GOPATH)/bin` (or `$GOBIN`)

Verify installation:

```bash
uxdguard version
uxdguard --version
uxdguard rules
```

### Option 2: Build from source

```bash
git clone https://github.com/arjunuxd/uxdguard.git
cd uxdguard
go build ./cmd/uxdguard
```

On Windows, use `uxdguard.exe`. Add the binary to your `PATH` or run it from the build directory.

### Option 3: Release binaries

When published, download binaries and checksums from the [GitHub Releases](https://github.com/arjunuxd/uxdguard/releases) page for your platform. UXDGuard does not currently auto-update.

## Quick Start

```bash
# Check version and active scanners
uxdguard version
uxdguard rules

# Scan current project
uxdguard check

# Scan staged changes
uxdguard check --staged

# Output as JSON
uxdguard check --json

# Paginate/filter large results (terminal output)
uxdguard check --page 1 --limit 100 --severity high

# Initialize configuration
uxdguard init

# Validate configuration
uxdguard config validate

# Manage Git hooks
uxdguard hook status
uxdguard hook install
uxdguard hook uninstall
```

## Usage Notes

### Large result sets
For large scans, terminal output is paginated by default (100 findings per page). Use `--page`, `--limit`, `--severity`, and `--rule` to navigate and filter. JSON output (`--json`) always includes complete results with no truncation.

### Default exclusions
Common dependency and generated directories are excluded by default (e.g., `node_modules/`, `.git/`, `vendor/`, `dist/`, `build/`, `.next/`, `.nuxt/`, `coverage/`, `.turbo/`, `.cache/`, `out/`, `__pycache__/`, `.venv/`, `venv/`). Dependency manifests and lockfiles (e.g., `package.json`, `go.mod`, `requirements.txt`) remain discoverable. Exclusions can be configured in `.uxdguard.yaml`.

### Git hooks
Hooks can be bypassed with `git commit --no-verify`. Hooks are not a security boundary.

## Configuration

Create a starter config:

```bash
uxdguard init
```

Edit `.uxdguard.yaml` as needed and validate:

```bash
uxdguard config validate
```

See [docs/configuration.md](docs/configuration.md) for details.

## Privacy & Security

- **Local-first:** Core scanning runs locally with no mandatory network access.
- **No telemetry:** No telemetry, analytics, or tracking.
- **No code execution:** Scanned files are never executed.
- **Safe output:** Findings avoid exposing raw secret values.
- **Optional integrations:** External tools are disabled by default and only used if explicitly enabled.
- **Limitations:** Pattern-based scanning only; may produce false positives/negatives and is not a substitute for comprehensive security review.

See [docs/privacy.md](docs/privacy.md), [docs/threat-model.md](docs/threat-model.md), and [docs/security-policy.md](docs/security-policy.md) for more information.

## Documentation

| Document | Description |
| --- | --- |
| [docs/configuration.md](docs/configuration.md) | Configuration options and defaults |
| [docs/rules.md](docs/rules.md) | Built-in rules and categories |
| [docs/hooks.md](docs/hooks.md) | Staged scanning and Git hooks |
| [docs/external-scanners.md](docs/external-scanners.md) | Gitleaks and Semgrep integrations |
| [docs/privacy.md](docs/privacy.md) | Privacy and data handling |
| [docs/threat-model.md](docs/threat-model.md) | Threat model and assumptions |
| [docs/security-policy.md](docs/security-policy.md) | Reporting security issues |
| [architecture.md](architecture.md) | Architecture overview |

## Development

```bash
# Build
go build ./...

# Test
go test ./...

# Vet
go vet ./...

# Format
gofmt -l .
```

## Uninstall

```bash
uxdguard uninstall
```

If self-removal is blocked by your OS, remove the `uxdguard` binary from your Go bin directory manually.

## License

[MIT License](LICENSE)
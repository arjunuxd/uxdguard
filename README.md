# UXDGuard

**Build with AI. Commit with confidence.**

UXDGuard is a lightweight, local-first security scanner for developers. It helps detect potential secrets, insecure coding patterns, and common security issues in source code and staged Git changes before they are committed.

## Key Features

- **Secret Detection** – Identifies patterns for API keys, tokens, private keys, credentials, and database connection strings.
- **Source & Advanced Security Checks** – Covers common issues including `eval()` usage, disabled TLS verification, SQL concatenation, weak cryptography, injection patterns, SSRF, and more.
- **Password Handling Checks** – Detects plaintext password handling and insecure password comparisons.
- **Staged Scanning** – Scan only Git-staged changes with `uxdguard check --staged`.
- **Git Hooks** – Install, check, or remove a pre-commit hook to catch issues early.
- **Configurable Rules** – Enable, disable, or adjust rule severity via `.uxdguard.yaml`.
- **Smart Exclusions** – Skips common dependency and build directories by default while preserving manifests and lockfiles.
- **Pagination & Filtering** – Navigate large result sets in the terminal with `--page`, `--limit`, `--severity`, and `--rule`.
- **JSON Output** – Complete, machine-readable results for CI/CD pipelines.
- **Optional Integrations** – Works with Gitleaks and Semgrep when explicitly enabled and installed.
- **Local-First & Private** – Runs entirely on your machine with no telemetry, no source uploads, and no mandatory network access.

## Requirements

- [Go](https://go.dev/) 1.22.9+ (for source builds/install)
- [Git](https://git-scm.com/) (for staged scanning and hooks)

## Installation

### Install with Go (recommended)

```bash
go install github.com/arjunuxd/uxdguard/cmd/uxdguard@latest
```

Make sure your Go bin directory is in your `PATH`:

| Platform | Path |
| --- | --- |
| Windows (PowerShell) | `$(go env GOPATH)\bin` |
| Linux/macOS | `$(go env GOPATH)/bin` or `$GOBIN` |

Verify:

```bash
uxdguard version
uxdguard rules
```

### Build from source

```bash
git clone https://github.com/arjunuxd/uxdguard.git
cd uxdguard
go build -o uxdguard ./cmd/uxdguard
```

On Windows, use `uxdguard.exe`.

### Release binaries

Pre-built binaries and checksums are available on the [Releases](https://github.com/arjunuxd/uxdguard/releases) page when published.

## Quick Start

```bash
# Navigate to your project
cd /path/to/your/project

# Scan all files
uxdguard check

# Scan staged changes only
uxdguard check --staged

# Get JSON output
uxdguard check --json

# View paginated results with filters
uxdguard check --page 1 --limit 50 --severity high

# List available rules
uxdguard rules

# Initialize and validate configuration
uxdguard init
uxdguard config validate

# Manage pre-commit hooks
uxdguard hook status
uxdguard hook install
uxdguard hook uninstall
```

## Configuration

Create a starter configuration:

```bash
uxdguard init
```

This creates `.uxdguard.yaml` in your project root. Validate changes with `uxdguard config validate`. See [docs/configuration.md](docs/configuration.md) for all options.

**Default exclusions include:** `.git/`, `node_modules/`, `vendor/`, `dist/`, `build/`, `.next/`, `.nuxt/`, `coverage/`, `.turbo/`, `.cache/`, `out/`, `__pycache__/`, `.venv/`, `venv/`. Dependency manifests (e.g., `package.json`, `go.mod`, `requirements.txt`) are still discovered.

## Privacy & Security

- **Local-first** – All core scanning runs locally with no mandatory network access.
- **No telemetry** – No analytics, tracking, or data collection.
- **No code execution** – Scanned files are never executed.
- **Safe reporting** – Outputs avoid exposing raw secret values.
- **Optional integrations** – External tools are disabled by default.
- **Limitations** – Pattern-based only; may produce false positives/negatives and is not a substitute for a full security review.

See [docs/privacy.md](docs/privacy.md), [docs/threat-model.md](docs/threat-model.md), and [docs/security-policy.md](docs/security-policy.md).

## Documentation

| Document | Description |
| --- | --- |
| [docs/configuration.md](docs/configuration.md) | Configuration schema and defaults |
| [docs/rules.md](docs/rules.md) | Built-in rules and categories |
| [docs/hooks.md](docs/hooks.md) | Staged scanning and Git hooks |
| [docs/external-scanners.md](docs/external-scanners.md) | Gitleaks and Semgrep integration |
| [docs/privacy.md](docs/privacy.md) | Privacy and data handling |
| [docs/threat-model.md](docs/threat-model.md) | Threat model and assumptions |
| [docs/security-policy.md](docs/security-policy.md) | Reporting security issues |
| [architecture.md](architecture.md) | Architecture overview |

## Development

```bash
go build ./...
go test ./...
go vet ./...
gofmt -l .
```

## Uninstall

```bash
uxdguard uninstall
```

If self-removal is blocked, manually delete the `uxdguard` binary from your Go bin directory.

## License

[MIT License](LICENSE)
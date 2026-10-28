# UXDGuard

**A local-first security and code-quality CLI for developers.**

UXDGuard scans source code and staged Git changes for potential security issues, exposed secrets, and insecure coding patterns. It provides terminal and JSON reports, configurable rules, and optional integrations with Gitleaks and Semgrep.

Designed as a lightweight Go application, UXDGuard performs local scanning without uploading source code or sending telemetry.

**Version:** `0.1.0-dev`  
**Status:** Pre-release — human review pending

## Features

- **Secret detection:** Identifies patterns associated with cloud keys, access tokens, private keys, database credentials, and hardcoded credentials.
- **Source-code security checks:** Detects selected insecure coding patterns, including unsafe evaluation, disabled TLS verification, SQL string concatenation, insecure HTTP endpoints, and unsafe deserialization.
- **Staged Git scanning:** Scans staged changes using `uxdguard check --staged`.
- **Configurable rules:** Configure rule behavior and severity through `.uxdguard.yaml`.
- **JSON reporting:** Produce machine-readable scan results.
- **Git hooks:** Install, inspect, and uninstall a pre-commit hook.
- **Optional external scanners:** Integrate with Gitleaks and Semgrep when manually installed and explicitly configured.
- **Cross-platform support:** Built with Go for Windows, Linux, and macOS.
- **Local-first operation:** No source-code uploads, telemetry, or AI API calls.

## Built-in security rules

### Secret detection

| Rule ID | Detection |
|---|---|
| `UXD-SEC-100` | AWS access keys |
| `UXD-SEC-101` | GitHub tokens |
| `UXD-SEC-102` | Slack tokens |
| `UXD-SEC-103` | Private keys |
| `UXD-SEC-104` | Database connection strings containing credentials |
| `UXD-SEC-105` | Hardcoded credential assignments |

### Source-code checks

| Rule ID | Detection |
|---|---|
| `UXD-SEC-200` | `eval()` usage |
| `UXD-SEC-201` | Disabled TLS certificate verification |
| `UXD-SEC-202` | SQL string concatenation |
| `UXD-SEC-203` | Insecure HTTP endpoints, excluding localhost |
| `UXD-SEC-204` | Unsafe deserialization patterns |

These rules use pattern-based detection. They do not replace a full semantic security analysis or a comprehensive security review.

## Requirements

- Go `1.22.9` (development and verification version).
- Git for staged scanning and Git-hook functionality.
- Gitleaks and/or Semgrep only if their integrations are needed.

## Installation

### Build from source

Clone the repository:

```bash
git clone https://github.com/arjunuxd/uxdguard.git
cd uxdguard
```

Build the CLI:

```bash
go build -o uxdguard ./cmd/uxdguard
```

On Windows PowerShell, use:

```powershell
go build -o uxdguard.exe ./cmd/uxdguard
```

Run the executable from the project directory or add its location to your `PATH`.

> UXDGuard is currently in pre-release development. A stable public release and published installation packages are not yet assumed to be available.

## Usage

The examples below show the supported CLI commands.

### Display help

```bash
uxdguard --help
```

### Display version

```bash
uxdguard --version
```

### Initialize configuration

Create a starter configuration in the current project:

```bash
uxdguard init
```

### Scan a project

```bash
uxdguard check
```

### Scan staged Git changes

```bash
uxdguard check --staged
```

Use this before committing to check the content staged in Git.

### Generate JSON output

```bash
uxdguard check --json
```

### Validate configuration

```bash
uxdguard config validate
```

### List available rules

```bash
uxdguard rules
```

### Manage the pre-commit hook

Check hook status:

```bash
uxdguard hook status
```

Install the hook:

```bash
uxdguard hook install
```

Uninstall the hook:

```bash
uxdguard hook uninstall
```

Review the hook behavior and configuration before enabling it in an existing repository.

## Configuration

UXDGuard uses `.uxdguard.yaml` for configuration.

Initialize a starter configuration:

```bash
uxdguard init
```

Validate it after making changes:

```bash
uxdguard config validate
```

Configuration supports rule settings, including enabling or disabling rules and overriding severity. Refer to the project's configuration documentation for the exact supported fields and values.

## Optional external scanners

UXDGuard supports optional integrations with:

- **Gitleaks** — additional secret detection.
- **Semgrep** — additional static-analysis checks.

These tools are disabled by default and must be installed manually and configured explicitly. UXDGuard does not automatically install them.

See the external scanner documentation for configuration and usage details.

## Security and privacy

UXDGuard is designed to operate locally.

- Source code is not uploaded to a remote service by the built-in scanning workflow.
- No telemetry or AI API calls are used.
- The scanner does not execute the source code it analyzes.
- External scanner integrations are optional.
- Findings are designed to redact detected secret values.
- Filesystem boundaries, file-size limits, process timeouts, and cancellation are part of the security design.

These are project design goals and implemented safeguards, not a guarantee that every possible vulnerability or secret will be detected. Review the threat model and known limitations before relying on UXDGuard in security-sensitive workflows.

## Development and testing

Run the following commands from the repository root.

Format check:

```bash
gofmt -l .
```

Run tests:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

Build all packages:

```bash
go build ./...
```

Measure test coverage:

```bash
go test -cover ./...
```

Verify module checksums:

```bash
go mod verify
```

Run race detection where supported:

```bash
go test -race ./...
```

Race detection may not be available in every Windows development environment without the required C toolchain.

The GitHub Actions workflows provide CI checks for formatting, tests, static analysis, builds, coverage, and supported cross-platform targets.

## Project documentation

Additional documentation is available in the repository:

| Document | Description |
|---|---|
| `docs/architecture.md` | System architecture |
| `docs/configuration.md` | Configuration options |
| `docs/privacy.md` | Privacy and local-first behavior |
| `docs/rules.md` | Built-in security rules |
| `docs/hooks.md` | Git-hook functionality |
| `docs/external-scanners.md` | Gitleaks and Semgrep integrations |
| `docs/threat-model.md` | Threat analysis and security assumptions |
| `docs/security-policy.md` | Security policy |

## Known limitations

- Detection is pattern-based and may produce false positives or miss issues that require semantic analysis.
- Gitleaks and Semgrep require separate installation and explicit configuration.
- Git hooks can be bypassed, including with `git commit --no-verify`.
- Race detection depends on the host platform and toolchain.
- Test fixtures may intentionally contain synthetic insecure patterns or credentials to verify detection.
- A clean scan does not establish that a project is free of security vulnerabilities.

## Project status

UXDGuard's core implementation, security checks, configuration, staged scanning, hook integration, optional scanner adapters, documentation, and release workflows have been reported as complete.

The project remains at version `0.1.0-dev`, pending final human review and the first public release decision.

## License

MIT License. See [LICENSE](LICENSE) for details.


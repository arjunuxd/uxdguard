# Privacy

## Local-First Design

UXDGuard runs on your machine. Built-in scanning does not require network access and does not include telemetry, analytics, or tracking.

## Data Access

- **Files**: Reads files under the scanned project root, respecting exclusions, file size limits, and symlink settings.
- **Git**: When using `--staged` or Git hooks, reads Git index content for the affected files.
- **Configuration**: Reads `.uxdguard.yaml` from the project root if present.

## Network Communication

- **Built-in scanning**: No network requests.
- **External tools**: If optional Gitleaks or Semgrep integrations are enabled, those separate programs may have their own network behavior. For example, Semgrep with `config: auto` may retrieve rule packs. Configure a local rules path for offline use.
- **No telemetry**: UXDGuard does not send usage data or findings anywhere.

## Data Handling

- **No code execution**: Scanned files are never executed.
- **Redacted output**: Secret values are not printed in terminal or JSON reports.
- **Local storage only**: Scan results, reports, and configuration remain on your machine.
- **Staged content**: Written only to a private temporary directory (`0700` directory, `0600` files) and removed when the scan finishes.
- **User control**: You choose what to scan, how to configure it, and where to write reports.

## Recommendations

- Do not commit configuration files that contain sensitive values.
- Review output before sharing it.
- Understand the behavior of any external tool you enable.
- Keep dependencies and the UXDGuard binary up to date.

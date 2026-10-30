# Privacy

## Local-First Design
UXDGuard is designed to run locally on your machine. Core scanning functionality does not require network access.

## Data Access
- **Files**: Reads files under the scanned project root, respecting exclusions, file size limits, and symlink handling settings.
- **Git**: When using `--staged` or Git hooks, reads Git index information as needed for those features.
- **Configuration**: Reads `.uxdguard.yaml` from the project root if present.

## Network Communication
- **Core scanning**: No mandatory network requests for normal scanning operations.
- **External tools**: If optional external integrations are enabled, those tools may have their own network behavior. Refer to their documentation.
- **No telemetry**: UXDGuard does not include telemetry, analytics, or tracking by default.

## Data Handling
- **No code execution**: Scanned files are never executed.
- **Safe output**: Findings are reported without exposing raw secret values or unnecessary sensitive content.
- **Local storage only**: Scan results, reports, and configuration remain on your local machine.
- **User control**: You control what is scanned, how it is configured, and where reports are written.

## Recommendations
- Avoid committing configuration files that contain sensitive values
- Review outputs before sharing them
- Understand the behavior of any optional external tools you enable
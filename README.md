# UXDGuard

**Build with AI. Commit with confidence.**

UXDGuard is a local-first, open-source CLI tool designed to help developers detect security and code-quality issues before committing or pushing code. All core scanning runs locally with no uploads, telemetry, or AI API calls.

## Project Status

Phase 5: Hardening, Cross-platform Testing & Release Preparation (complete)

All phases 1-5 implemented. See [docs/PROJECT-STATUS.md](docs/PROJECT-STATUS.md) for details.

## Principles

- Local-first scanning
- Privacy by default
- Minimal dependencies
- Secure, readable, maintainable code
- Transparent findings and failures

## Prerequisites

- Go 1.22+ (for building from source)

## Build

```bash
go build ./cmd/uxdguard
```

## Commands

```bash
uxdguard --help
uxdguard --version
uxdguard init
uxdguard check
uxdguard check --json
uxdguard config validate
uxdguard rules
```

## Usage

1. Initialize configuration (optional):
```bash
uxdguard init
```
This creates `.uxdguard.yaml` and `.uxdguard.yaml.example`. Existing config is not overwritten.

2. Run a scan:
```bash
uxdguard check
uxdguard check --json
```

3. Validate configuration:
```bash
uxdguard config validate
```

4. List registered rules/scanners:
```bash
uxdguard rules
```

## Configuration

See [docs/configuration.md](docs/configuration.md). Defaults are used if no `.uxdguard.yaml` is present.

## Architecture

See [docs/architecture.md](docs/architecture.md).

## Privacy

See [docs/privacy.md](docs/privacy.md).

## Testing

```bash
go test ./...
go vet ./...
gofmt -l .
go build ./...
```

## Limitations (Phase 1)

- Phase 1 provides the scanning infrastructure only. Comprehensive vulnerability detection and the full secret-detection rule set are planned for Phase 2.
- No Git staged-index scanning (Phase 3).
- No dynamic plugin loading.

## License

MIT

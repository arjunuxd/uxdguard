# Optional External Scanner Integrations (Phase 4)

## Supported Integrations
- **gitleaks** (optional): Secret detection. Disabled by default. Requires gitleaks installed.
- **semgrep** (optional): Static analysis. Disabled by default. Requires semgrep installed. Uses local-only mode when possible.

## Configuration
```yaml
external:
  gitleaks:
    enabled: false
    require: false
    path: ""  # optional custom path
    timeout: "60s"
  semgrep:
    enabled: false
    require: false
    path: ""
    timeout: "120s"
    options:
      config: ""  # local config path if needed
```

## Behavior
- Missing optional tools do not break built-in scans
- If `require: true` and tool missing/fails, scan returns operational error
- Tools invoked with explicit args, no shell execution
- Outputs parsed defensively
- Secrets redacted
- No automatic installation or updates
- No network calls for credential validation

## Privacy
- Source code not uploaded by UXDGuard
- Adapters do not enable cloud features by default
- Users responsible for understanding external tool behavior

## Limitations
- No Gitleaks/Semgrep installed in test environment
- Integration tested via unit tests with mocked behavior patterns
- Dependency auditing deferred as not feasible with safe offline guarantees

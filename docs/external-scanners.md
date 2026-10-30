# External Scanners

UXDGuard supports optional integration with external security tools. These integrations are disabled by default and require separate installation of the external tools.

## Supported Integrations
| Tool | Purpose | Status |
| --- | --- | --- |
| [gitleaks](https://github.com/gitleaks/gitleaks) | Additional secret detection | Optional, disabled by default |
| [semgrep](https://github.com/semgrep/semgrep) | Additional static analysis | Optional, disabled by default |

## Configuration
Enable integrations in `.uxdguard.yaml`:

```yaml
external:
  gitleaks:
    enabled: false
    require: false
    path: ""
    timeout: "60s"
  semgrep:
    enabled: false
    require: false
    path: ""
    timeout: "120s"
    options:
      config: ""
```

## Usage Notes
- Missing optional tools do not prevent built-in scans from running.
- If `require: true` and the tool is missing or fails, the scan will return an error.
- External tools are invoked locally with explicit arguments.
- UXDGuard does not automatically install or update external tools.
- The behavior and privacy practices of external tools are governed by their own documentation.

## Limitations
- Integrations depend on third-party tools being installed and available
- Results depend on the configuration of external tools
- Network behavior, if any, is determined by the external tool
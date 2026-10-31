# External Scanners

UXDGuard supports optional integration with external security tools. Integrations are disabled by default, require explicit configuration, and require the tools to be installed separately.

## Supported Integrations

| Tool | Purpose | Default |
| --- | --- | --- |
| [gitleaks](https://github.com/gitleaks/gitleaks) | Additional secret detection | Disabled |
| [semgrep](https://github.com/semgrep/semgrep) | Additional static analysis | Disabled |

## Configuration

Enable integrations in `.uxdguard.yaml`:

```yaml
external:
  gitleaks:
    enabled: true
    require: false
    path: ""
    timeout: "60s"
  semgrep:
    enabled: true
    require: false
    path: ""
    timeout: "120s"
    options:
      config: ""
```

When `path` is empty, UXDGuard looks for `gitleaks` or `semgrep` on `PATH`. When `options.config` is empty, Semgrep uses its `auto` configuration.

## Behavior

- Disabled integrations do nothing and produce no findings.
- Missing optional tools are ignored when `require: false`.
- If `require: true` and the tool is missing, times out, fails, or returns unparsable output, the scan reports an error.
- External tools are invoked locally with explicit argument lists, never through a shell.
- Findings use the `GL-` (Gitleaks) or `SEM-` (Semgrep) rule prefixes.
- External findings pass through the same rule filtering and severity overrides as built-in findings.

## Privacy

Built-in UXDGuard scanning never sends data anywhere. External tools are separate programs with their own behavior:

- Gitleaks runs locally against the project directory.
- Semgrep with `config: auto` may download rule packs depending on the installed version and environment. Use a local configuration path if you require fully offline operation.
- Review the documentation and defaults of any external tool before enabling it.

## Limitations

- Integrations depend on third-party tools being installed and available.
- Output formats can vary between tool versions.
- Results depend on the external tool's configuration.
- UXDGuard does not install or update external tools.

# Configuration

## Configuration File

UXDGuard looks for `.uxdguard.yaml` in the project root. If no configuration file is found, built-in defaults are used. Generate a starter file with:

```bash
uxdguard init
```

`init` never overwrites an existing configuration file.

## Schema (v1)

```yaml
version: 1

scan:
  max_file_size_bytes: 1048576
  follow_symlinks: false
  exclusions:
    - ".git/**"
    - "vendor/**"
    - "node_modules/**"
    - "dist/**"
    - "build/**"
    - ".next/**"
    - ".nuxt/**"
    - ".svelte-kit/**"
    - ".turbo/**"
    - ".cache/**"
    - "coverage/**"
    - "target/**"
    - "out/**"
    - ".venv/**"
    - "venv/**"
    - "__pycache__/**"
    - ".terraform/**"
    - ".idea/**"
    - ".vscode/**"

policy:
  fail_on:
    - critical
    - high

rules:
  enable: []
  disable: []
  severity:
    UXD-SEC-203: low

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

## Options

| Setting | Type | Default | Description |
| --- | --- | --- | --- |
| `version` | number | `1` | Configuration schema version. |
| `scan.max_file_size_bytes` | number | `1048576` | Maximum file size to scan, in bytes. |
| `scan.follow_symlinks` | boolean | `false` | Follow symlinks whose targets stay inside the project root. |
| `scan.exclusions` | list | See above | Patterns for files and directories to exclude. |
| `policy.fail_on` | list | `["critical", "high"]` | Severities that cause a non-zero exit code. |
| `rules.enable` | list | empty | Scanner or rule IDs to report when non-empty. |
| `rules.disable` | list | empty | Scanner or rule IDs to suppress. Takes precedence over `enable`. |
| `rules.severity` | map | empty | Per-rule severity overrides. |
| `external.*.enabled` | boolean | `false` | Enable an optional external scanner. |
| `external.*.require` | boolean | `false` | Treat a missing or failing external tool as an error. |
| `external.*.path` | string | tool name | Custom executable path. |
| `external.*.timeout` | duration | `60s` / `120s` | Execution timeout, for example `30s` or `2m`. |
| `external.semgrep.options.config` | string | `auto` | Semgrep configuration path or `auto`. |

## Rule Controls

- Entries in `enable` and `disable` may be scanner IDs or individual rule IDs.
- When `enable` is non-empty, only listed scanner IDs or rule IDs are reported. Use it to narrow output, not to disable a scanner while leaving another enabled implicitly.
- When a scanner is disabled, it appears in the report's `skipped` list.
- Severity overrides accept `critical`, `high`, `medium`, `low`, and `info`.
- Rule filtering and severity overrides are applied centrally to built-in, staged, and external findings.

## Symlink Behavior

- `follow_symlinks: false` ignores all symlinks.
- `follow_symlinks: true` follows symlinks only when the resolved target is inside the project root.
- Traversal tracks resolved directories, so symlink cycles terminate.
- Symlinks pointing outside the project root are never scanned.

## Validation

`uxdguard config validate` checks:

- `version` is `1`;
- `max_file_size_bytes` is positive;
- exclusion patterns are non-empty;
- `fail_on` values are one of `critical`, `high`, `medium`, `low`, `info`;
- `rules.severity` keys are non-empty and values are valid severities.

Invalid configuration produces a clear error and a non-zero exit code.

## Precedence

1. CLI flags such as `--project-root`, `--config`, `--severity`, and `--rule`;
2. `.uxdguard.yaml` in the project root (or the file passed with `--config`);
3. Built-in defaults.

## Notes

- Exclusions are matched against paths relative to the project root.
- Dependency manifests and lockfiles are not excluded by default.
- A `--config` file must be inside the project root.
- `uxdguard init` creates only `.uxdguard.yaml`; it does not create a separate example file in your project.
- Terminal filters (`--severity`, `--rule`) affect display only. JSON output always contains the complete configured result set.

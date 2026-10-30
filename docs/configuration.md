# Configuration

## Configuration File
UXDGuard looks for `.uxdguard.yaml` in the project root. If no configuration file is found, built-in defaults are used. You can generate a starter file with:

```bash
uxdguard init
```

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
    - "coverage/**"
    - ".turbo/**"
    - ".cache/**"
    - "out/**"
    - "__pycache__/**"
    - ".venv/**"
    - "venv/**"

policy:
  fail_on:
    - critical
    - high
```

## Options

| Setting | Type | Default | Description |
| --- | --- | --- | --- |
| `version` | number | `1` | Configuration schema version. |
| `scan.max_file_size_bytes` | number | `1048576` | Maximum file size to scan (bytes). |
| `scan.follow_symlinks` | boolean | `false` | Whether to follow symbolic links. |
| `scan.exclusions` | list | See above | Patterns for files and directories to exclude. |
| `policy.fail_on` | list | `["critical", "high"]` | Severities that cause a non-zero exit code. |

## Validation
- `version` must be `1`
- `max_file_size_bytes` must be positive
- Exclusion patterns must be non-empty
- `fail_on` values must be one of: `critical`, `high`, `medium`, `low`, `info`
- Invalid configuration produces clear errors

## Precedence
1. CLI flags (e.g., `--project-root`)
2. `.uxdguard.yaml` in the project root
3. Built-in defaults

## Notes
- Creating configuration with `uxdguard init` does not overwrite existing files.
- Exclusions are matched against relative paths from the project root.
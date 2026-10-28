# Configuration

## File

Configuration is read from `.uxdguard.yaml` in the project root. If not present, defaults are used. A starter example is written by `uxdguard init` as `.uxdguard.yaml.example`.

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

policy:
  fail_on:
    - critical
    - high
```

## Defaults

- `version`: 1
- `scan.max_file_size_bytes`: 1048576 (1MB)
- `scan.follow_symlinks`: false
- `scan.exclusions`: includes `.git/**`, `vendor/**`, `node_modules/**`, `dist/**`, `build/**`
- `policy.fail_on`: `["critical", "high"]`

## Validation

- `version` must be `1`
- `max_file_size_bytes` must be positive
- Exclusion patterns must be non-empty
- `fail_on` severities must be one of `critical`, `high`, `medium`, `low`, `info`
- Malformed YAML produces a clear error with file context
- Missing optional config file is not an error (defaults used)

## Discovery and Precedence

1. CLI flags: `--project-root`, `--config` (where applicable)
2. Config file discovery: look for `.uxdguard.yaml` in resolved project root
3. Defaults applied for missing fields

`init` does not overwrite existing configuration.

# Built-in Rules (Phase 2)

## Secret Detection (`secrets`)

| Rule ID | Title | Severity (default) | Confidence | Notes |
| --- | --- | --- | --- | --- |
| UXD-SEC-100 | Potential AWS access key detected | high | medium | Matches `AKIA...` / `ASIA...` patterns. |
| UXD-SEC-101 | Potential GitHub token detected | critical | high | Matches common `ghp_`/`gho_`/`ghu_`/`ghs_`/`ghr_` prefixes. |
| UXD-SEC-102 | Potential Slack token detected | high | high | Matches `xox...` token patterns. |
| UXD-SEC-103 | Private key detected | critical | high | Detects PEM private key headers. |
| UXD-SEC-104 | Database connection string with credentials detected | high | high | Detects URLs with embedded credentials for common DBs. |
| UXD-SEC-105 | Hardcoded credential assignment detected | high | medium | Detects assignments like `password=...`, `api_key=...` etc. Skips placeholders (`YOUR_*`, `REPLACE_ME`, `example`, `changeme`, etc.). |

## Source Code Security Checks (`sourcechecks`)

| Rule ID | Title | Severity (default) | Confidence | Notes |
| --- | --- | --- | --- | --- |
| UXD-SEC-200 | Use of eval() detected | high | medium | Language-neutral pattern check. |
| UXD-SEC-201 | TLS certificate verification disabled | high | high | Detects `InsecureSkipVerify: true`. |
| UXD-SEC-202 | Possible SQL string concatenation | high | medium | Heuristic for SQL built via concatenation/interpolation. |
| UXD-SEC-203 | Insecure HTTP endpoint detected | medium | medium | Flags `http://` outside `localhost`/`127.0.0.1`/`::1`. |
| UXD-SEC-204 | Potential unsafe deserialization | high | medium | Flags `pickle.loads` in `.py` and `ObjectInputStream` in `.java`. |

## Configuration

Rules can be configured in `.uxdguard.yaml`:
```yaml
rules:
  enable:
    - secrets
  disable:
    - sourcechecks
  severity:
    UXD-SEC-100: critical
```

Unknown rule IDs in severity overrides are not specially rejected by schema validation except severity values are checked; invalid values produce validation errors.

## Redaction
Findings do not include raw secret values or full source snippets. The secret scanner and reporters are designed to avoid leaking secrets; tests verify no secrets appear in outputs.

## Limitations
- Pattern-based only - no semantic/data-flow analysis
- May produce false positives/negatives
- Not all token formats/providers covered
- Not a replacement for proper secret management
- Scanning is local only, no network validation

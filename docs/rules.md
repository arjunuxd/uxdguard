# Built-in Rules

UXDGuard uses pattern-based checks. The tables list every built-in rule in the current release. Code-pattern rules are not applied to documentation files (`.md`, `.markdown`, `.rst`, `.adoc`, `.txt`) because examples are not executable; secret detection still applies to those files.

## Secret Detection (`secrets`)

| Rule ID | Name | Severity |
| --- | --- | --- |
| UXD-SEC-100 | Potential AWS access key detected | High |
| UXD-SEC-101 | Potential GitHub token detected | Critical |
| UXD-SEC-102 | Potential Slack token detected | High |
| UXD-SEC-103 | Private key detected | Critical |
| UXD-SEC-104 | Database connection string with credentials detected | High |
| UXD-SEC-105 | Hardcoded credential assignment detected | High |
| UXD-SEC-106 | Potential hardcoded JWT detected | High |
| UXD-SEC-107 | Potential Google API key detected | High |
| UXD-SEC-108 | Potential Resend API key detected | High |
| UXD-SEC-109 | Potential Razorpay key detected | High |
| UXD-SEC-110 | Potential PostHog key detected | High |
| UXD-SEC-111 | Potential Stripe key detected | Critical |
| UXD-SEC-112 | Potential SendGrid API key detected | High |
| UXD-SEC-113 | Potential Discord webhook URL detected | High |
| UXD-SEC-114 | Potential Cloudinary credentials detected | High |

Provider-specific rules use unique IDs. Placeholder values such as `YOUR_API_KEY`, `your_jwt_secret_here`, `changeme`, `example`, and `replace_me` are suppressed for assignment-style checks.

## Source Code Security Checks (`sourcechecks`)

| Rule ID | Name | Severity |
| --- | --- | --- |
| UXD-SEC-200 | Use of dynamic code evaluation detected | High |
| UXD-SEC-201 | TLS certificate verification disabled | High |
| UXD-SEC-202 | Possible SQL string concatenation | High |
| UXD-SEC-203 | Insecure HTTP endpoint detected | Medium |
| UXD-SEC-204 | Potential unsafe deserialization | High |
| UXD-SEC-205 | Potential path traversal pattern | High |
| UXD-SEC-206 | Weak cryptographic hash in security context | Medium |
| UXD-SEC-207 | Potential command execution detected | High |
| UXD-SEC-208 | Potential SSRF pattern detected | High |
| UXD-SEC-209 | Potential authentication bypass detected | Critical |
| UXD-SEC-210 | Potential obfuscated code execution | Critical |
| UXD-SEC-211 | Potential open redirect | Medium |
| UXD-SEC-212 | Plaintext password handling detected | High |
| UXD-SEC-213 | Insecure password comparison | High |

`RegExp.prototype.exec` calls are not treated as process execution. Relative import specifiers (`../lib/file`) are not treated as path traversal.

## Advanced Security Checks (`advanced`)

| Rule ID | Name | Severity |
| --- | --- | --- |
| UXD-SEC-300 | Potential command injection with shell execution | High |
| UXD-SEC-301 | Potential path traversal | High |
| UXD-SEC-302 | Potential XSS via unsafe HTML rendering | High |
| UXD-SEC-303 | Potential SQL injection via string concatenation | High |
| UXD-SEC-304 | Potential open redirect | Medium |
| UXD-SEC-305 | Potential debug/admin endpoint exposed | Medium |
| UXD-SEC-400 | Insecure JWT verification | Critical |
| UXD-SEC-401 | Hardcoded JWT secret | Critical |
| UXD-SEC-402 | Insecure cookie configuration | Medium |
| UXD-SEC-500 | Weak cryptographic hash | Medium |
| UXD-SEC-501 | Weak random number generation | High |
| UXD-SEC-502 | Hardcoded encryption key | Critical |
| UXD-SEC-503 | Potential secret logging | High |
| UXD-SEC-600 | Potential SSRF | High |
| UXD-SEC-601 | File upload handling detected | Medium |
| UXD-SEC-602 | Unsafe deserialization | Critical |

`UXD-SEC-400` reports disabled verification, alg=none, and ignored expiration. A normal `jwt.verify(...)` call is not reported by itself.

## External Scanner Results

When explicitly enabled, external adapters contribute findings with their own IDs:

| Prefix | Source |
| --- | --- |
| `GL-<rule>` | Gitleaks |
| `SEM-<rule>` | Semgrep |

External findings are filtered by the same `rules.enable`, `rules.disable`, and `rules.severity` configuration as built-in findings.

## Configuration

Rules can be configured in `.uxdguard.yaml`:

```yaml
rules:
  # When enable is non-empty, only listed scanner IDs or rule IDs are reported.
  enable:
    - secrets
    - UXD-SEC-201
  disable:
    - UXD-SEC-203
  severity:
    UXD-SEC-100: critical
    UXD-SEC-203: low
```

- `disable` takes precedence over `enable`.
- Entries may be scanner IDs (`secrets`, `sourcechecks`, `advanced`, `gitleaks`, `semgrep`) or rule IDs.
- `severity` accepts `critical`, `high`, `medium`, `low`, or `info`.
- Disabling a scanner records it in the `skipped` list in reports.

## Limitations

- Pattern-based scanning only; no data-flow or semantic analysis.
- False positives and false negatives are possible.
- Coverage varies by language and framework.
- These rules do not replace a comprehensive security review.

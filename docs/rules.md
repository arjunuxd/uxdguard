# Built-in Rules

UXDGuard includes pattern-based security checks. Only rules that are registered and active in the current release are listed below.

## Secret Detection (`secrets`)
| Rule ID | Name | Severity | Description |
| --- | --- | --- | --- |
| UXD-SEC-100 | Potential AWS access key detected | High | Detects common AWS access key formats. |
| UXD-SEC-101 | Potential GitHub token detected | Critical | Detects common GitHub token formats. |
| UXD-SEC-102 | Potential Slack token detected | High | Detects common Slack token formats. |
| UXD-SEC-103 | Private key detected | Critical | Detects private key material. |
| UXD-SEC-104 | Database connection string with credentials detected | High | Detects database connection strings containing credentials. |
| UXD-SEC-105 | Hardcoded credential assignment detected | High | Detects likely hardcoded credential assignments. |

## Source Code Security Checks (`sourcechecks`)
| Rule ID | Name | Severity | Description |
| --- | --- | --- | --- |
| UXD-SEC-200 | Use of eval() detected | High | Detects use of `eval()`. |
| UXD-SEC-201 | TLS certificate verification disabled | High | Detects disabled TLS certificate verification. |
| UXD-SEC-202 | Possible SQL string concatenation | High | Detects potential SQL built through string concatenation. |
| UXD-SEC-203 | Insecure HTTP endpoint detected | Medium | Flags non-local HTTP endpoints. |
| UXD-SEC-204 | Potential unsafe deserialization | High | Detects patterns associated with unsafe deserialization. |
| UXD-SEC-212 | Plaintext password handling detected | High | Detects plaintext password handling patterns. |
| UXD-SEC-213 | Insecure password comparison | High | Detects insecure password comparison patterns. |

## Advanced Security Checks (`advanced`)
| Rule ID | Name | Severity | Description |
| --- | --- | --- | --- |
| UXD-SEC-300 | Potential command injection with shell execution | High | Detects patterns involving shell execution with dynamic input. |
| UXD-SEC-301 | Potential path traversal | High | Detects potential path traversal patterns. |
| UXD-SEC-302 | Potential XSS via unsafe HTML rendering | High | Detects potential cross-site scripting patterns. |
| UXD-SEC-303 | Potential SQL injection via string concatenation | High | Detects SQL injection patterns via concatenation. |
| UXD-SEC-304 | Potential open redirect | Medium | Detects potential open redirect patterns. |
| UXD-SEC-305 | Potential debug/admin endpoint exposed | Medium | Detects potential debug or admin endpoints. |
| UXD-SEC-400 | Insecure JWT verification | Critical | Detects insecure JWT verification patterns. |
| UXD-SEC-401 | Hardcoded JWT secret | Critical | Detects hardcoded JWT secrets. |
| UXD-SEC-402 | Insecure cookie configuration | Medium | Detects missing secure cookie attributes. |
| UXD-SEC-500 | Weak cryptographic hash | Medium | Detects use of weak cryptographic hashes. |
| UXD-SEC-501 | Weak random number generation | High | Detects use of weak random number generation for security purposes. |
| UXD-SEC-502 | Hardcoded encryption key | Critical | Detects hardcoded encryption keys. |
| UXD-SEC-503 | Potential secret logging | High | Detects potential logging of sensitive data. |
| UXD-SEC-600 | Potential SSRF | High | Detects potential server-side request forgery patterns. |
| UXD-SEC-601 | File upload handling detected | Medium | Detects file upload handling that may require validation. |
| UXD-SEC-602 | Unsafe deserialization | Critical | Detects unsafe deserialization patterns. |

## Configuration
Rules can be configured in `.uxdguard.yaml`:

```yaml
rules:
  enable:
    - secrets
    - sourcechecks
    - advanced
  disable: []
  severity:
    UXD-SEC-100: critical
```

## Limitations
- Pattern-based scanning only; no data-flow or semantic analysis
- May produce false positives or false negatives
- Coverage varies by language and framework
- Not a replacement for comprehensive security review
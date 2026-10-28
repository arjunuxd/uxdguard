# Threat Model

## Assets
- Source code and configuration
- Secrets/credentials
- Git index and hooks
- User environment

## Threats and Mitigations

| Threat | Impact | Mitigation | Residual Risk |
| --- | --- | --- | --- |
| Accidental secrets in commits | High | Built-in secret detection, staged scanning, pre-commit hook; redaction in outputs | Pattern-based, may miss obfuscated secrets |
| Path traversal | High | Path validation, stay within repo root, no shell execution | Complex symlink cases |
| Malicious repo content | Medium | No code execution, size limits, binary detection | Some parsers could be stressed |
| Malformed config | Medium | Validation, clear errors | Unknown fields handled per design |
| External tool output abuse | Medium | Defensive parsing, no shell interpolation, timeouts | Depends on external tool |
| Hook tampering | Medium | Marker-based install/uninstall, refuse to overwrite unknown hooks | Users can use --no-verify |
| Info leakage | High | Redaction, no secrets in logs/errors, relative paths by default | Debugging may need care |
| Resource exhaustion | Medium | Size limits, timeouts, bounded buffers | Very large repos still take time |

## Assumptions
- Local-first, no network required for core functionality
- Git available for staged/hook features
- Users review findings

## Non-goals
- Complete security guarantee
- Malware detection
- Full SAST replacement

# Threat Model

## Purpose
UXDGuard is a local-first security scanner that helps developers identify common security issues in source code before committing changes. This document describes its intended scope and general security considerations at a high level.

## Assets Protected
- Source code and project files under review
- Developer credentials and secrets present in repositories
- Git workflow integrity (when using staged scanning and hooks)

## Intended Usage
- Run locally on developer machines
- Scan project directories without executing scanned code
- Optionally scan staged Git changes
- Use as part of local development workflows

## General Risks Considered
- Accidental exposure of secrets or credentials
- Common insecure coding patterns
- Misconfiguration in project files
- Path traversal attempts in untrusted repositories

## Security Principles
- **Local-first**: Core scanning runs entirely on the user's machine.
- **No code execution**: Scanned files are never executed.
- **Read-only by default**: Normal scans do not modify the target repository.
- **Least privilege**: Only accesses files under the scanned project root.
- **Safe reporting**: Outputs are designed to avoid exposing sensitive values.

## User Responsibilities
- Review findings and remediation guidance
- Use appropriate secret management practices
- Keep UXDGuard up to date
- Understand the limitations of pattern-based scanning

## Known Limitations
- Pattern-based detection only; not a complete replacement for security review or dynamic analysis
- May produce false positives or false negatives
- Some risks require context that static pattern matching cannot determine
- External integrations, if used, have their own behaviors and limitations
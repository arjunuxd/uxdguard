# Security Policy

UXDGuard is a local-first security scanner. We take reports about the security of
the tool itself seriously and welcome good-faith research.

This policy covers vulnerabilities in UXDGuard's code, CLI, build, and release
process. It does not cover false positives, false negatives, or vulnerabilities
in code that UXDGuard scans.

## Supported Versions

| Version | Supported |
| --- | --- |
| `1.1.x` (latest stable) | :white_check_mark: |
| `main` (unreleased) | :white_check_mark: Security fixes land here first |
| `< 1.1` | :x: |

Only the latest stable release and the current `main` branch receive security
updates. Pre-release artifacts, forks, and modified binaries are not supported.

## Reporting a Vulnerability

**Do not open a public issue, discussion, or pull request for a suspected
vulnerability.** Public disclosure before a fix is available can put users at
risk.

### Preferred channel

Use GitHub's private vulnerability reporting:

**https://github.com/arjunuxd/uxdguard/security/advisories/new**

This creates a private advisory visible only to you and the maintainers. It also
lets us collaborate on a fix, request a CVE if appropriate, and publish a
coordinated advisory.

If the form is unavailable, contact the maintainer through
**https://github.com/arjunuxd** and ask for a private reporting channel. Do not
include exploit details in a public message.

### What to include

Please provide as much of the following as you can:

- Affected version, commit, or release asset
- Operating system and architecture
- How UXDGuard was installed or built
- A description of the vulnerability and its impact
- Reproduction steps or a minimal proof of concept
- Relevant configuration, command line, or repository layout
- Whether the issue is already public or known to others
- Any suggested fix or mitigation

Use test fixtures and synthetic data. **Never include real credentials, personal
data, or third-party data** in a report. Redact secrets from logs and command
output.

### What to expect

These are best-effort targets, not contractual guarantees:

| Stage | Target |
| --- | --- |
| Acknowledge receipt | 3 business days |
| Initial assessment and severity | 7 business days |
| Fix or mitigation plan | Based on severity and complexity |
| Coordinated advisory | After a fixed release is available |

We will keep you informed of progress and tell you if a report is out of scope or
requires more information. Complex issues, dependency coordination, and
multi-platform fixes may take longer.

UXDGuard does not currently operate a bug bounty program. We gratefully credit
reporters in advisories and release notes unless anonymity is requested.

## Scope

### In scope

Vulnerabilities in UXDGuard itself, including:

- CLI argument parsing, command handling, and exit-code behaviour
- File discovery, exclusions, symlink handling, and path containment
- Secret and source-pattern scanners, rule evaluation, and severity overrides
- Configuration loading, validation, and defaulting
- Staged Git scanning, temporary file handling, and hook installation
- External adapter invocation, argument handling, timeouts, and output limits
- Terminal, JSON, and SARIF report generation, including injection through
  untrusted file names or tool output
- Release workflow, build metadata, and published artifacts
- Supply-chain issues in UXDGuard's own dependency graph

### Out of scope

- False positives, false negatives, or missed vulnerabilities in scanned code
- Vulnerabilities in third-party tools such as Gitleaks or Semgrep; report those
  to the upstream project
- Issues that require an already compromised host, administrator access, a
  modified UXDGuard binary, or a malicious Go toolchain
- Social engineering, physical access, or denial of service against GitHub
- Dependency vulnerabilities without a demonstrated impact on UXDGuard
  (informational reports are still welcome)
- Findings that only affect a configuration the user was explicitly warned
  against, such as enabling an untrusted external tool

If you are unsure whether an issue is in scope, report it privately and let us
triage it.

## Security Model

UXDGuard is designed to reduce risk while scanning untrusted repositories:

- **Local-first:** built-in scanning requires no telemetry, no account, and no
  network service. Scanned code is not uploaded.
- **No code execution:** scanned files are parsed as text and are never executed.
- **Read-only by default:** a normal scan does not modify the target repository.
  `init`, `hook install`, and baseline output write only when explicitly invoked.
- **Root-confined paths:** symlinks are ignored by default. When following is
  enabled, targets must resolve inside the project root and cycles are detected.
- **Private staged content:** staged files are written to a private temporary
  directory with owner-only permissions and removed after the scan.
- **Bounded resources:** external process output, staged blob reads, file sizes,
  and line buffers are capped to avoid unbounded memory use.
- **No shell interpolation:** external tools and Git are invoked with explicit
  argument lists rather than through a shell.
- **Safe output:** terminal output is sanitized, and secret values are not
  printed in reports.
- **Opt-in integrations:** Gitleaks and Semgrep are disabled by default. If
  enabled, they run as separate programs and may have their own network and
  privacy behaviour.

See [Threat Model](docs/threat-model.md) and [Privacy](docs/privacy.md) for more
detail.

## Safe Harbour

We welcome good-faith security research and will not initiate legal action
against researchers who:

- Test only systems they own or have explicit permission to test
- Avoid accessing, modifying, exfiltrating, or deleting data that is not theirs
- Avoid privacy violations, disruption of service, and denial-of-service attacks
- Report findings privately and give us reasonable time to respond
- Comply with applicable law

This statement is not a waiver of any third party's rights and does not authorize
testing of systems or data you do not control.

## Disclosure Policy

We follow coordinated disclosure:

1. You report the issue privately.
2. We confirm the issue, assess severity, and prepare a fix.
3. We agree on a disclosure date with you.
4. We publish a GitHub Security Advisory and release notes after a fixed version
   is available.
5. We credit you unless you ask to remain anonymous.

Please do not publish details, proof-of-concept code, or exploit tooling before a
fix and advisory are ready. If a report is already public, tell us immediately so
we can prioritize accordingly.

## Security Updates

Security fixes are announced through:

- [GitHub Security Advisories](https://github.com/arjunuxd/uxdguard/security/advisories)
- [Releases](https://github.com/arjunuxd/uxdguard/releases)
- The `main` branch and release notes

Published tags and release assets are immutable. A fix that changes shipped code
requires a new version; existing tags are never moved or rewritten.

## No Guarantees

UXDGuard is provided under the MIT License on an "as is" basis, without warranty
of any kind. Pattern-based scanning can produce false positives and false
negatives, and no scanner can guarantee that a codebase is free of
vulnerabilities. Security reports are handled in good faith, but response times
and outcomes are not guaranteed.

# Privacy

- **Local-first:** All core scanning runs locally. No source code is uploaded.
- **No telemetry:** No telemetry, analytics, or tracking.
- **No AI APIs:** No calls to AI APIs.
- **No mandatory network:** No mandatory network requests.
- **Filesystem access:** Only reads files under the project root as discovered, respects exclusions and size limits, does not follow symlinks by default.
- **No code execution:** Scanned files are never executed.
- **Safe reporting:** Reports do not include raw secret values or full source contents. JSON output is serialized without internal stack traces or secrets.
- **No persistent sensitive storage:** No persistent storage of source contents or sensitive scan results.
- **Configuration**: Sensitive values should not be logged.

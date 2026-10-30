# Git Hooks

## Staged Scanning
Scan only staged changes before committing:

```bash
uxdguard check --staged
```

This scans files as they exist in the Git index, respects exclusions and size limits, and uses the same rules as a normal scan.

## Hook Management
Install, check status, or remove the pre-commit hook:

```bash
uxdguard hook install
uxdguard hook status
uxdguard hook uninstall
```

## Behavior
- The hook runs `uxdguard check --staged` before commits.
- If blocking findings are detected, the commit is prevented based on the hook''s exit code.
- Hook installation is safe and only modifies hooks it creates.
- You can bypass hooks with `git commit --no-verify` if needed.

## Requirements
- Git must be installed
- The project must be a Git repository
- Hooks are project-local by default

## Limitations
- Binary files may be skipped
- Complex rename scenarios may have edge cases
- Hook behavior depends on your Git and shell environment
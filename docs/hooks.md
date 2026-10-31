# Git Hooks

## Staged Scanning

Scan only staged changes before committing:

```bash
uxdguard check --staged
```

This reads files as they exist in the Git index, respects exclusions and size limits, and applies the same rules and policy as a normal scan. Unstaged edits are not scanned.

## Hook Management

Install, check, or remove the pre-commit hook:

```bash
uxdguard hook install
uxdguard hook status
uxdguard hook uninstall
```

## Behavior

- The hook runs `uxdguard check --staged` before each commit.
- If the scan exits non-zero under your `policy.fail_on` settings, the commit is blocked.
- Installation honors `core.hooksPath` and linked worktrees by using `git rev-parse --git-path hooks`.
- The executable path is quoted, so installations under paths containing spaces work.
- Existing hooks that are not managed by UXDGuard are never overwritten or removed.
- You can bypass hooks with `git commit --no-verify` when necessary.

## Requirements

- Git must be installed and on `PATH`.
- The project must be a Git repository.
- The hook uses a POSIX `sh` shebang, including on Windows through Git for Windows.

## Limitations

- Binary files are skipped.
- Rename and submodule edge cases depend on Git index reporting.
- Hook behavior depends on your Git and shell environment.

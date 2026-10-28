# Git Staged Scanning and Pre-commit Hooks

## Staged Scanning
Use `uxdguard check --staged` to scan files as they exist in the Git index (not the working tree).

Key behaviors:
- Scans index/blob contents for staged files
- Respects exclusions and file-size limits
- Uses existing Phase 2 rules
- Reports relative paths from repo root
- Returns exit code 1 if blocking findings found, 0 if clean, 2 for config/usage errors, 3 for operational errors
- Works with `--json` for machine-readable output

Correctness:
- If you stage a secret and then remove it from working tree without staging, it is still detected in the index
- If a secret only exists in unstaged changes, it is not reported by `--staged`

## Hook Management
```powershell
uxdguard hook status
uxdguard hook install
uxdguard hook uninstall
```

Safety:
- Never overwrites existing non-UXDGuard hooks
- Never modifies global Git config
- Only removes hooks it created (marked with marker)
- Idempotent install
- Hook invokes `uxdguard check --staged` and exits with its status code

## Limitations
- Requires Git installed and a git repository
- Binary files skipped
- Some edge cases with complex renames may need additional testing
- Windows hook execution uses shell wrapper as created
- Users can bypass hooks with `git commit --no-verify`

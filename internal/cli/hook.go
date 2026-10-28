package root

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/arjunuxd/uxdguard/internal/app"
	"github.com/arjunuxd/uxdguard/internal/gitutil"
)

const hookMarker = "# UXDGuard pre-commit hook"

func cmdHook(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "missing subcommand\n")
		printUsage()
		return app.ExitUsageOrConfig
	}
	if args[0] == "install" {
		return hookInstall(args[1:])
	}
	if args[0] == "uninstall" {
		return hookUninstall(args[1:])
	}
	if args[0] == "status" {
		return hookStatus(args[1:])
	}
	if args[0] == "-h" || args[0] == "--help" {
		fmt.Printf("Usage: uxdguard hook install|uninstall|status [--project-root <path>]\n")
		return app.ExitOK
	}
	fmt.Fprintf(os.Stderr, "unknown hook subcommand\n")
	printUsage()
	return app.ExitUsageOrConfig
}

func getHooksDir(repoRoot string) (string, error) {
	// Check core.hooksPath
	// Use git config
	// Simple approach: use .git/hooks
	gitDir := filepath.Join(repoRoot, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return filepath.Join(gitDir, "hooks"), nil
	}
	// git worktree case, not handling fully
	return filepath.Join(repoRoot, ".git", "hooks"), nil
}

func hookInstall(args []string) int {
	projectRoot := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--project-root" && i+1 < len(args) {
			projectRoot = args[i+1]
			i++
			continue
		}
		if args[i] == "-h" || args[i] == "--help" {
			fmt.Printf("Usage: uxdguard hook install [--project-root <path>]\n")
			return app.ExitOK
		}
		fmt.Fprintf(os.Stderr, "unknown flag: %s\n", args[i])
		return app.ExitUsageOrConfig
	}
	start := projectRoot
	if start == "" {
		start, _ = os.Getwd()
	}
	repoRoot, err := gitutil.RepoRoot(start)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	hooksDir, err := getHooksDir(repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	hookPath := filepath.Join(hooksDir, "pre-commit")
	// Check existing hook
	if _, err := os.Stat(hookPath); err == nil {
		// Read and check if it's ours
		b, _ := os.ReadFile(hookPath)
		if contains(string(b), hookMarker) {
			fmt.Println("Hook already installed")
			return app.ExitOK
		}
		fmt.Fprintf(os.Stderr, "existing pre-commit hook found; refusing to overwrite\n")
		return app.ExitUsageOrConfig
	}
	// Create hook
	exe, _ := os.Executable()
	hookContent := fmt.Sprintf(`#!/bin/sh
%s
%s check --staged
exit $?
`, hookMarker, exe)
	if err := os.WriteFile(hookPath, []byte(hookContent), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	fmt.Printf("Installed pre-commit hook to %s\n", hookPath)
	return app.ExitOK
}

func hookUninstall(args []string) int {
	projectRoot := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--project-root" && i+1 < len(args) {
			projectRoot = args[i+1]
			i++
			continue
		}
		if args[i] == "-h" || args[i] == "--help" {
			fmt.Printf("Usage: uxdguard hook uninstall [--project-root <path>]\n")
			return app.ExitOK
		}
		fmt.Fprintf(os.Stderr, "unknown flag: %s\n", args[i])
		return app.ExitUsageOrConfig
	}
	start := projectRoot
	if start == "" {
		start, _ = os.Getwd()
	}
	repoRoot, err := gitutil.RepoRoot(start)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	hooksDir, err := getHooksDir(repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	hookPath := filepath.Join(hooksDir, "pre-commit")
	b, err := os.ReadFile(hookPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No hook installed")
			return app.ExitOK
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	if !contains(string(b), hookMarker) {
		fmt.Fprintf(os.Stderr, "existing hook is not managed by UXDGuard; refusing to remove\n")
		return app.ExitUsageOrConfig
	}
	if err := os.Remove(hookPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	fmt.Println("Removed pre-commit hook")
	return app.ExitOK
}

func hookStatus(args []string) int {
	projectRoot := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--project-root" && i+1 < len(args) {
			projectRoot = args[i+1]
			i++
			continue
		}
		if args[i] == "-h" || args[i] == "--help" {
			fmt.Printf("Usage: uxdguard hook status [--project-root <path>]\n")
			return app.ExitOK
		}
		fmt.Fprintf(os.Stderr, "unknown flag: %s\n", args[i])
		return app.ExitUsageOrConfig
	}
	start := projectRoot
	if start == "" {
		start, _ = os.Getwd()
	}
	repoRoot, err := gitutil.RepoRoot(start)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	hooksDir, err := getHooksDir(repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	hookPath := filepath.Join(hooksDir, "pre-commit")
	b, err := os.ReadFile(hookPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Hook not installed")
			return app.ExitOK
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	if contains(string(b), hookMarker) {
		fmt.Printf("UXDGuard pre-commit hook installed at %s\n", hookPath)
		return app.ExitOK
	}
	fmt.Printf("Other pre-commit hook installed at %s (not managed by UXDGuard)\n", hookPath)
	return app.ExitOK
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(substr) > 0 && index(s, substr) >= 0))
}

func index(s, sep string) int {
	return len(s) - len(sep) - len(s[:len(s)-len(sep)]) // simple
}

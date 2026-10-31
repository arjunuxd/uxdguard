package root

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	return gitutil.HooksDir(repoRoot)
}

// hookScript builds a POSIX shell pre-commit hook. The executable path is
// single-quoted so installations under paths containing spaces work.
func hookScript(exe string) string {
	exe = filepath.ToSlash(exe)
	quoted := "'" + strings.ReplaceAll(exe, "'", `'\''`) + "'"
	return fmt.Sprintf("#!/bin/sh\n%s\n%s check --staged\nexit $?\n", hookMarker, quoted)
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
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	if err := os.WriteFile(hookPath, []byte(hookScript(exe)), 0755); err != nil {
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
	return strings.Contains(s, substr)
}

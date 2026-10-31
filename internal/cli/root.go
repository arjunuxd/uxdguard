package root

import (
	"fmt"
	"os"

	"github.com/arjunuxd/uxdguard/internal/app"
)

// Execute runs the CLI.
func Execute(args []string) int {
	if len(args) == 0 {
		printUsage()
		return app.ExitUsageOrConfig
	}

	// Global flags
	if args[0] == "--help" || args[0] == "-h" {
		printUsage()
		return app.ExitOK
	}

	if args[0] == "--version" || args[0] == "-v" {
		fmt.Printf("uxdguard version %s\n", CurrentVersion())
		return app.ExitOK
	}

	if args[0] == "version" {
		build := ResolveBuildInfo()
		fmt.Printf(
			"UXDGuard\nSecurity Scanner\nVersion: %s\nCommit: %s\nBuild date: %s\n",
			build.Version,
			build.Commit,
			build.BuildDate,
		)
		return app.ExitOK
	}

	if args[0] == "init" {
		return cmdInit(args[1:])
	}

	if args[0] == "check" {
		return cmdCheck(args[1:])
	}

	if args[0] == "config" {
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "missing subcommand")
			printUsage()
			return app.ExitUsageOrConfig
		}

		if args[1] == "validate" {
			return cmdConfigValidate(args[2:])
		}

		fmt.Fprintln(os.Stderr, "unknown config subcommand")
		printUsage()
		return app.ExitUsageOrConfig
	}

	if args[0] == "rules" {
		return cmdRules(args[1:])
	}

	if args[0] == "hook" {
		return cmdHook(args[1:])
	}

	if args[0] == "uninstall" {
		return cmdUninstall(args[1:])
	}

	fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
	printUsage()
	return app.ExitUsageOrConfig
}

func printUsage() {
	fmt.Printf(`UXDGuard - Local-first security and code-quality scanner

Usage:
  uxdguard --help
  uxdguard --version
  uxdguard version
  uxdguard init
  uxdguard check [--json] [--staged] [--project-root <path>] [--page <n>] [--limit <n>] [--severity <level>] [--rule <id>] [--no-progress]
  uxdguard config validate [--project-root <path>] [--config <path>]
  uxdguard rules
  uxdguard hook install|uninstall|status [--project-root <path>]
  uxdguard uninstall [-y]

Examples:
  uxdguard check                Scan current project
  uxdguard check --staged       Scan staged Git changes
  uxdguard check --json         Output results as JSON
  uxdguard config validate      Validate .uxdguard.yaml

`)
}

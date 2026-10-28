package root

import (
	"fmt"
	"os"

	"github.com/arjunuxd/uxdguard/internal/app"
)

var version = "0.1.0-dev"

// SetVersion allows overriding version at build time
func SetVersion(v string) {
	if v != "" {
		version = v
	}
}

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
		fmt.Printf("uxdguard version %s\n", version)
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
			fmt.Fprintf(os.Stderr, "missing subcommand\n")
			printUsage()
			return app.ExitUsageOrConfig
		}
		if args[1] == "validate" {
			return cmdConfigValidate(args[2:])
		}
		fmt.Fprintf(os.Stderr, "unknown config subcommand\n")
		printUsage()
		return app.ExitUsageOrConfig
	}
	if args[0] == "rules" {
		return cmdRules(args[1:])
	}
	if args[0] == "hook" {
		return cmdHook(args[1:])
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
  uxdguard init
  uxdguard check [--json] [--staged] [--project-root <path>]
  uxdguard config validate [--project-root <path>] [--config <path>]
  uxdguard rules
  uxdguard hook install|uninstall|status [--project-root <path>]

Examples:
  uxdguard check                Scan current project
  uxdguard check --staged       Scan staged Git changes
  uxdguard check --json         Output results as JSON
  uxdguard config validate      Validate .uxdguard.yaml

`)
}

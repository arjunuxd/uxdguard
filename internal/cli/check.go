package root

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/arjunuxd/uxdguard/internal/app"
	"github.com/arjunuxd/uxdguard/internal/config"
	"github.com/arjunuxd/uxdguard/internal/gitutil"
	"github.com/arjunuxd/uxdguard/internal/project"
	"github.com/arjunuxd/uxdguard/internal/report"
	"github.com/arjunuxd/uxdguard/internal/scan"
)

func cmdCheck(args []string) int {
	version := CurrentVersion()
	jsonOut := false
	projectRoot := ""
	staged := false
	page := 1
	limit := 20
	severity := ""
	ruleFilter := ""
	noProgress := false
	for i := 0; i < len(args); i++ {
		if args[i] == "--json" {
			jsonOut = true
			continue
		}
		if args[i] == "--staged" {
			staged = true
			continue
		}
		if args[i] == "--no-progress" {
			noProgress = true
			continue
		}
		if args[i] == "--project-root" && i+1 < len(args) {
			projectRoot = args[i+1]
			i++
			continue
		}
		if args[i] == "--page" && i+1 < len(args) {
			fmt.Sscanf(args[i+1], "%d", &page)
			i++
			continue
		}
		if args[i] == "--limit" && i+1 < len(args) {
			fmt.Sscanf(args[i+1], "%d", &limit)
			i++
			continue
		}
		if args[i] == "--severity" && i+1 < len(args) {
			severity = args[i+1]
			i++
			continue
		}
		if args[i] == "--rule" && i+1 < len(args) {
			ruleFilter = args[i+1]
			i++
			continue
		}
		if args[i] == "-h" || args[i] == "--help" {
			fmt.Printf("Usage: uxdguard check [--json] [--staged] [--project-root <path>] [--page <n>] [--limit <n>] [--severity <level>] [--rule <id>] [--no-progress]\n")
			return app.ExitOK
		}
		fmt.Fprintf(os.Stderr, "unknown flag: %s\n", args[i])
		return app.ExitUsageOrConfig
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if staged {
		// Find repo root
		start := projectRoot
		if start == "" {
			start, _ = os.Getwd()
		}
		repoRoot, err := gitutil.RepoRoot(start)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return app.ExitOperationalFail
		}
		_, cfg, err := project.LoadWithConfig(repoRoot, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return app.ExitUsageOrConfig
		}
		res, err := scan.ScanStaged(ctx, repoRoot, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return app.ExitOperationalFail
		}
		if jsonOut {
			jr := &report.JSONReporter{Out: os.Stdout}
			if err := jr.Report(res, repoRoot, version); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return app.ExitOperationalFail
			}
		} else {
			tr := &report.TerminalReporter{Out: os.Stdout, Err: os.Stderr}
			tr.ReportWithOptions(res, repoRoot, version, report.ReportOptions{
				Page:       page,
				Limit:      limit,
				Severity:   severity,
				Rule:       ruleFilter,
				NoProgress: noProgress,
			})
		}
		exit := scan.EvaluatePolicy(cfg, res)
		return exit
	}
	root, cfg, err := project.LoadWithConfig(projectRoot, "")
	if err != nil {
		var cerr *os.PathError
		if errors.As(err, &cerr) {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return app.ExitOperationalFail
		}
		// config/usage errors
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitUsageOrConfig
	}
	files, err := project.Discover(root, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	sc := scan.NewContext(root, cfg, files).WithContext(ctx)
	res := scan.DefaultRegistry().Run(sc)
	if jsonOut {
		jr := &report.JSONReporter{Out: os.Stdout}
		if err := jr.Report(res, root, version); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return app.ExitOperationalFail
		}
	} else {
		tr := &report.TerminalReporter{Out: os.Stdout, Err: os.Stderr}
		tr.ReportWithOptions(res, root, version, report.ReportOptions{
			Page:       page,
			Limit:      limit,
			Severity:   severity,
			Rule:       ruleFilter,
			NoProgress: noProgress,
		})
	}
	// Evaluate policy
	exit := scan.EvaluatePolicy(cfg, res)
	return exit
}

func cmdInit(args []string) int {
	projectRoot := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--project-root" && i+1 < len(args) {
			projectRoot = args[i+1]
			i++
			continue
		}
		if args[i] == "-h" || args[i] == "--help" {
			fmt.Printf("Usage: uxdguard init [--project-root <path>]\n")
			return app.ExitOK
		}
		fmt.Fprintf(os.Stderr, "unknown flag: %s\n", args[i])
		return app.ExitUsageOrConfig
	}
	root := projectRoot
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return app.ExitOperationalFail
		}
	}
	info, err := os.Stat(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "error: project root is not a directory: %s\n", root)
		return app.ExitUsageOrConfig
	}
	target := filepath.Join(root, ".uxdguard.yaml")
	if _, err := os.Stat(target); err == nil {
		fmt.Fprintf(os.Stderr, "config already exists: %s\n", target)
		return app.ExitUsageOrConfig
	}
	example := `version: 1

scan:
  max_file_size_bytes: 1048576
  follow_symlinks: false
  exclusions:
    - ".git/**"
    - "vendor/**"
    - "node_modules/**"
    - "dist/**"
    - "build/**"

policy:
  fail_on:
    - critical
    - high

# Optional: tune individual rules or scanners.
# rules:
#   disable:
#     - UXD-SEC-203
#   severity:
#     UXD-SEC-203: low
`
	if err := os.WriteFile(target, []byte(example), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	fmt.Printf("Created %s\n", target)
	return app.ExitOK
}

func cmdConfigValidate(args []string) int {
	projectRoot := ""
	configPath := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--project-root" && i+1 < len(args) {
			projectRoot = args[i+1]
			i++
			continue
		}
		if args[i] == "--config" && i+1 < len(args) {
			configPath = args[i+1]
			i++
			continue
		}
		if args[i] == "-h" || args[i] == "--help" {
			fmt.Printf("Usage: uxdguard config validate [--project-root <path>] [--config <path>]\n")
			return app.ExitOK
		}
		fmt.Fprintf(os.Stderr, "unknown flag: %s\n", args[i])
		return app.ExitUsageOrConfig
	}
	root, cfg, err := project.LoadWithConfig(projectRoot, configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitUsageOrConfig
	}
	if err := config.Validate(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitUsageOrConfig
	}
	fmt.Printf("Configuration valid (project: %s)\n", root)
	return app.ExitOK
}

func cmdRules(args []string) int {
	for i := 0; i < len(args); i++ {
		if args[i] == "-h" || args[i] == "--help" {
			fmt.Printf("Usage: uxdguard rules\n")
			return app.ExitOK
		}
		fmt.Fprintf(os.Stderr, "unknown flag: %s\n", args[i])
		return app.ExitUsageOrConfig
	}
	// List registered scanners/rules
	scanners := scan.DefaultRegistry().List()
	if len(scanners) == 0 {
		fmt.Println("No rules registered.")
		return app.ExitOK
	}
	for _, s := range scanners {
		fmt.Printf("- %s: %s (%s)\n", s.ID(), s.Name(), s.Description())
	}
	return app.ExitOK
}

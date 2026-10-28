package app

import (
	"context"

	"github.com/arjunuxd/uxdguard/internal/project"
	"github.com/arjunuxd/uxdguard/internal/scan"
)

// App orchestrates scanning.
type App struct {
	Version string
}

// Run executes the scan pipeline and returns exit code.
func (a *App) Run(ctx context.Context, opts RunOptions) (int, error) {
	// Discovery and loading
	root, cfg, err := project.LoadWithConfig(opts.ProjectRoot, opts.ConfigPath)
	if err != nil {
		return ExitUsageOrConfig, err
	}
	// Discover files
	files, err := project.Discover(root, cfg)
	if err != nil {
		return ExitOperationalFail, err
	}
	// Create scan context
	sc := scan.NewContext(root, cfg, files).WithContext(ctx)
	// Run registry
	res := scan.DefaultRegistry().Run(sc)
	// Evaluate policy
	exit := scan.EvaluatePolicy(cfg, res)
	return exit, nil
}

// RunOptions for the app.
type RunOptions struct {
	ProjectRoot string
	ConfigPath  string
	JSON        bool
}

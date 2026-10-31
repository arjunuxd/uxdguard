package project

import (
	"os"
	"path/filepath"

	"github.com/arjunuxd/uxdguard/internal/config"
)

// LoadWithConfig discovers project root and loads config.
func LoadWithConfig(projectRoot, configPath string) (string, *config.Config, error) {
	var root string
	var err error
	if projectRoot != "" {
		root, err = filepath.Abs(projectRoot)
		if err != nil {
			return "", nil, err
		}
	} else {
		root, err = DiscoverRoot()
		if err != nil {
			return "", nil, err
		}
	}
	info, err := os.Stat(root)
	if err != nil {
		return "", nil, err
	}
	if !info.IsDir() {
		return "", nil, &os.PathError{Op: "load project", Path: root, Err: os.ErrInvalid}
	}
	cfg, _, err := config.Load(root, configPath)
	if err != nil {
		return "", nil, err
	}
	return root, cfg, nil
}

// DiscoverRoot finds project root by looking for common markers.
func DiscoverRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, ".uxdguard.yaml")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// fallback to cwd
	cwd, _ := os.Getwd()
	return cwd, nil
}

package project

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/arjunuxd/uxdguard/internal/config"
)

// Discover finds all eligible files under root.
func Discover(root string, cfg *config.Config) ([]string, error) {
	return discoverInternal(root, cfg, nil)
}

// DiscoverWithProgress finds all eligible files under root with progress callback.
func DiscoverWithProgress(root string, cfg *config.Config, progress func(count int)) ([]string, error) {
	return discoverInternal(root, cfg, progress)
}

func discoverInternal(root string, cfg *config.Config, progress func(count int)) ([]string, error) {
	if cfg == nil {
		cfg = config.Defaults()
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	realRoot := rootAbs
	if resolved, err := filepath.EvalSymlinks(rootAbs); err == nil {
		realRoot = resolved
	}

	var files []string
	visitedDirs := make(map[string]bool)

	addFile := func(path string, size int64) {
		if cfg.Scan.MaxFileSizeBytes > 0 && size > cfg.Scan.MaxFileSizeBytes {
			return
		}
		files = append(files, path)
		if progress != nil {
			progress(len(files))
		}
	}

	var walk func(dir string)
	walk = func(dir string) {
		// Key cycle detection on the resolved directory so symlink loops and
		// duplicate links to the same directory are only traversed once.
		key := dir
		if resolved, err := filepath.EvalSymlinks(dir); err == nil {
			key = resolved
		}
		if visitedDirs[key] {
			return
		}
		visitedDirs[key] = true

		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, entry := range entries {
			path := filepath.Join(dir, entry.Name())
			rel, err := filepath.Rel(rootAbs, path)
			if err != nil {
				continue
			}
			if entry.IsDir() && entry.Name() == ".git" {
				continue
			}
			if Excluded(rel, cfg.Scan.Exclusions) {
				continue
			}

			info, err := os.Lstat(path)
			if err != nil {
				continue
			}
			if info.Mode()&os.ModeSymlink != 0 {
				if !cfg.Scan.FollowSymlinks {
					continue
				}
				target, err := filepath.EvalSymlinks(path)
				if err != nil {
					continue
				}
				if !withinRoot(rootAbs, realRoot, target) {
					continue
				}
				targetInfo, err := os.Stat(path)
				if err != nil {
					continue
				}
				if targetInfo.IsDir() {
					walk(path)
					continue
				}
				if targetInfo.Mode().IsRegular() {
					addFile(path, targetInfo.Size())
				}
				continue
			}
			if info.IsDir() {
				walk(path)
				continue
			}
			if info.Mode().IsRegular() {
				addFile(path, info.Size())
			}
		}
	}

	walk(rootAbs)
	sort.Strings(files)
	return files, nil
}

// withinRoot reports whether target resolves inside the project root. Both the
// lexical root and its resolved form are accepted so symlinked project roots
// work correctly.
func withinRoot(root, realRoot, target string) bool {
	target = filepath.Clean(target)
	for _, base := range []string{root, realRoot} {
		rel, err := filepath.Rel(filepath.Clean(base), target)
		if err != nil {
			continue
		}
		if rel == "." {
			return true
		}
		if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

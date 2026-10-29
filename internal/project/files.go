package project

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/arjunuxd/uxdguard/internal/config"
)

// Discover finds all eligible files under root.
func Discover(root string, cfg *config.Config) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			// Skip unreadable dirs/files gracefully
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			// Check if excluded
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}
			if rel == "." {
				return nil
			}
			// Skip if it's .git
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			if Excluded(rel, cfg.Scan.Exclusions) {
				return filepath.SkipDir
			}
			// Prevent following symlinked directories outside root
			if !cfg.Scan.FollowSymlinks {
				if info, err := os.Lstat(path); err == nil {
					if info.Mode()&os.ModeSymlink != 0 {
						return filepath.SkipDir
					}
				}
			}
			return nil
		}
		// It's a file
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}
		if Excluded(rel, cfg.Scan.Exclusions) {
			return nil
		}
		// Skip symlinks if not following
		if !cfg.Scan.FollowSymlinks {
			if info, err := os.Lstat(path); err == nil {
				if info.Mode()&os.ModeSymlink != 0 {
					return nil
				}
			}
		}
		// Check if it's a regular file
		info, err := os.Lstat(path)
		if err != nil {
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if cfg.Scan.MaxFileSizeBytes > 0 && info.Size() > cfg.Scan.MaxFileSizeBytes {
			return nil // skip oversized
		}
		// Ensure inside root (no traversal)
		cleanRel := filepath.Clean(rel)
		if cleanRel == ".." || strings.HasPrefix(cleanRel, ".."+string(filepath.Separator)) || strings.HasPrefix(filepath.ToSlash(cleanRel), "../") {
			return errors.New("path traversal detected")
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

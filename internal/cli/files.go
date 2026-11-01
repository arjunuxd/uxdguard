package root

import "os"

// createExclusive creates a new file at path without writing through an
// existing symlink. O_EXCL alone is not enough on Windows: CreateFile follows a
// dangling symlink and creates its target, so any existing path (including a
// dangling link) is refused first with a non-following Lstat.
func createExclusive(path string, perm os.FileMode) (*os.File, error) {
	if _, err := os.Lstat(path); err == nil {
		return nil, os.ErrExist
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
}

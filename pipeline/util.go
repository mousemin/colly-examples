package pipeline

import "os"

// isDir returns true if given path is a directory, and returns false when it's
// a file or does not exist.
func isDir(dir string) bool {
	f, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return f.IsDir()
}

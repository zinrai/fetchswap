package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// No directory fsync: after a crash dest is either the old or the new file,
// both consistent, and the next run replaces it again.
func swap(dest string, r io.Reader, mode os.FileMode) error {
	// Next to dest, not in os.TempDir: rename(2) is atomic only within one
	// filesystem.
	f, err := os.CreateTemp(filepath.Dir(dest), "."+filepath.Base(dest)+".*.tmp")
	if err != nil {
		return err
	}
	// Remove unconditionally, not only on failure: after the rename the name
	// no longer exists and this is a no-op.
	defer os.Remove(f.Name())

	if err := writeFile(f, r, mode); err != nil {
		return err
	}
	return os.Rename(f.Name(), dest)
}

func writeFile(f *os.File, r io.Reader, mode os.FileMode) error {
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("copy body: %w", err)
	}
	// Chmod, not CreateTemp's 0600: that would silently tighten dest.
	if err := f.Chmod(mode); err != nil {
		return err
	}
	// No rename before Sync: a crash could leave dest renamed but empty.
	if err := f.Sync(); err != nil {
		return err
	}
	return f.Close()
}

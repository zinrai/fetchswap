package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/iotest"
)

func TestBodyIsWrittenWithMode(t *testing.T) {
	dest := newDest(t)

	if err := swap(dest, strings.NewReader("hello\n"), 0o640); err != nil {
		t.Fatalf("swap: %v", err)
	}

	assertFile(t, dest, "hello\n", 0o640)
}

func TestReadErrorKeepsDest(t *testing.T) {
	dest := newDest(t)
	writeExisting(t, dest, "original", 0o644)

	if err := swap(dest, cutShort("partial"), 0o600); err == nil {
		t.Fatal("swap succeeded, want error")
	}

	assertFile(t, dest, "original", 0o644)
}

func cutShort(s string) io.Reader {
	return io.MultiReader(strings.NewReader(s), iotest.ErrReader(errors.New("connection reset")))
}

func newDest(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "dest")
}

func writeExisting(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
}

func assertFile(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != content || fi.Mode().Perm() != mode {
		t.Errorf("%s = %q mode %o, want %q mode %o", path, got, fi.Mode().Perm(), content, mode)
	}
}

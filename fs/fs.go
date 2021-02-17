package fs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// CopyFS recursively copies the complete content of the given fs.FS object into the directory pointed to by out.
// The function will return immediately as soon as an error occurs which may leave the output directory in an
// incomplete state compared to in.
func CopyFS(in fs.FS, out string) error {
	if err := fs.WalkDir(in, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
		p := filepath.Join(out, path)
		if d.IsDir() {
			if err := os.MkdirAll(p, 0750); err != nil {
				return fmt.Errorf("could not create dir %s: %w", p, err)
			}
		} else {
			out, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0640)
			if err != nil {
				return fmt.Errorf("could not create file %s: %w", p, err)
			}
			in, err := in.Open(path)
			if err != nil {
				return fmt.Errorf("could not open file %s: %w", path, err)
			}
			if _, err := io.Copy(out, in); err != nil {
				return fmt.Errorf("could not copy file contents of %s into %s: %w", path, p, err)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

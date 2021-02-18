package fs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	stripComponents int
}

// Option is used to change the default behaviour of this package's functions.
type Option func(c *config)

// StripComponents strips n leading components from the input FS when copying its content.
func StripComponents(n int) Option {
	return func(c *config) {
		c.stripComponents = n
	}
}

// CopyFSWithOptions is like CopyFS but with the possibility to change its behaviour.
func CopyFSWithOptions(in fs.FS, out string, opts ...Option) error {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}
	if err := fs.WalkDir(in, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
		parts := strings.Split(path, string(filepath.Separator))
		if len(parts) <= c.stripComponents {
			return nil
		}
		strippedPath := filepath.Join(parts[c.stripComponents:]...)
		p := filepath.Join(out, strippedPath)
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

// CopyFS recursively copies the complete content of the given fs.FS object into the directory pointed to by out.
// The function will return immediately as soon as an error occurs which may leave the output directory in an
// incomplete state compared to in.
func CopyFS(in fs.FS, out string) error {
	return CopyFSWithOptions(in, out)
}

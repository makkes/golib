package fs_test

import (
	"embed"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"go.e13.dev/golib/v2/fs"

	"github.com/stretchr/testify/require"
)

//go:embed a
var embFS embed.FS

func TestCopyFS(t *testing.T) {
	target := t.TempDir()
	if err := fs.CopyFS(embFS, target); err != nil {
		t.Fatalf("CopyFS returned an unexpected error: %s", err)
	}
	filepath.Walk("a", func(path string, info os.FileInfo, err error) error {
		require.NoError(t, err, "failed to walk test data files")
		if info.IsDir() {
			return nil
		}
		relPath := path
		unpackedPath := filepath.Join(target, relPath)
		require.FileExists(t, unpackedPath)
		originalFile, err := ioutil.ReadFile(path)
		require.NoError(t, err, "failed to read original file")
		unpackedFile, err := ioutil.ReadFile(unpackedPath)
		require.NoError(t, err, "failed to read unpacked file")
		require.Equal(t, originalFile, unpackedFile, "unpacked file is different from original")
		fi, err := os.Lstat(unpackedPath)
		require.NoError(t, err, "failed to stat unpacked file")
		require.Equal(t, os.FileMode(0640), fi.Mode(), "wrong unpacked file mode on %s", fi.Name())
		return nil
	})
}

func TestCopyFSWithStripComponents(t *testing.T) {
	tests := map[string]struct {
		strip         int
		expectedFiles map[string][]byte
	}{
		"strip 0 components": {
			0,
			map[string][]byte{
				"a/b/c.txt": []byte("This is C\n"),
				"a/b/d.txt": []byte("This is D\n"),
			},
		},
		"strip 1 component": {
			1,
			map[string][]byte{
				"b/c.txt": []byte("This is C\n"),
				"b/d.txt": []byte("This is D\n"),
			},
		},
		"strip 2 components": {
			2,
			map[string][]byte{
				"c.txt": []byte("This is C\n"),
				"d.txt": []byte("This is D\n"),
			},
		},
		"strip 3 components": {
			3,
			nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			target := t.TempDir()
			if err := fs.CopyFSWithOptions(embFS, target, fs.StripComponents(tc.strip)); err != nil {
				t.Fatalf("Unpack returned an unexpected error: %s", err)
			}

			for file, expectedContent := range tc.expectedFiles {
				content, err := os.ReadFile(filepath.Join(target, file))
				require.NoError(t, err, "file %s should exist and be readable", file)
				require.Equal(t, expectedContent, content, "unpacked file is different from original")
			}
		})
	}
}

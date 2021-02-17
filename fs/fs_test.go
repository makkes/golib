package fs_test

import (
	"embed"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"go.e13.dev/golib/fs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed a
var embFS embed.FS

func TestUnpack(t *testing.T) {
	target := t.TempDir()
	if err := fs.CopyFS(embFS, target); err != nil {
		t.Fatalf("Unpack returned an unexpected error: %s", err)
	}
	filepath.Walk("a", func(path string, info os.FileInfo, err error) error {
		require.NoError(t, err, "failed to walk test data files")
		if info.IsDir() {
			return nil
		}
		relPath := path
		unpackedPath := filepath.Join(target, relPath)
		assert.FileExists(t, unpackedPath)
		originalFile, err := ioutil.ReadFile(path)
		assert.NoError(t, err, "failed to read original file")
		unpackedFile, err := ioutil.ReadFile(unpackedPath)
		assert.NoError(t, err, "failed to read unpacked file")
		assert.Equal(t, originalFile, unpackedFile, "unpacked file is different from original")
		fi, err := os.Lstat(unpackedPath)
		assert.NoError(t, err, "failed to stat unpacked file")
		assert.Equal(t, os.FileMode(0640), fi.Mode(), "wrong unpacked file mode on %s", fi.Name())
		return nil
	})
}

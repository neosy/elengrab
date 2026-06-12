package assets

import (
	"io/fs"
	"os"
	"path/filepath"
)

func CopyToDir(dst string, embedded fs.FS) error {
	return fs.WalkDir(embedded, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, p), 0o755)
		}
		data, err := fs.ReadFile(embedded, p)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dst, p), data, 0o644)
	})
}

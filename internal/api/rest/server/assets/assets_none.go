//go:build !embed_assets

package assets

import "io/fs"

var (
	AssetsFS fs.FS = nil
	Embedded       = false
)

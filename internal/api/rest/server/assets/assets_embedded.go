//go:build embed_assets

package assets

import "embed"

//go:embed static/** templates/**
var AssetsFS embed.FS

var Embedded = true

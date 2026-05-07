package sqliterep

import (
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/auth"
	"github.com/neosy/elengrab/internal/repository/sqlite/download"
	"github.com/neosy/elengrab/internal/repository/sqlite/link"
	"github.com/neosy/elengrab/internal/repository/sqlite/media"
)

var (
	MainSchema  persistence.DBSchema = download.NewDBSchema("elengrab")
	AuthSchema  persistence.DBSchema = auth.NewDBSchema("auth")
	MediaSchema persistence.DBSchema = media.NewDBSchema("media")
	LinkSchema  persistence.DBSchema = link.NewDBSchema("link")
)

package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func filePathForMigrate(p string) string {
	if runtime.GOOS == "windows" {
		return "file:" + p
	}

	p = filepath.ToSlash(p)
	return "file://" + p
}

const minSQLiteVersion = "3.53.4" // verified in this project

// requireSQLiteVersion fails early if the bundled engine can't run the migrations.
func requireSQLiteVersion(db *sql.DB) error {
	var v string
	if err := db.QueryRow(`SELECT sqlite_version()`).Scan(&v); err != nil {
		return err
	}
	if compareVersions(v, minSQLiteVersion) < 0 {
		return fmt.Errorf("SQLite %s is too old, need >= %s", v, minSQLiteVersion)
	}
	return nil
}

// latestMigrationVersion returns the highest version among "{version}_{title}.up.sql" files.
func latestMigrationVersion(fsys fs.FS) (uint, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return 0, err
	}

	var latest uint
	for _, e := range entries {
		prefix, _, ok := strings.Cut(e.Name(), "_")
		if !ok {
			continue
		}
		v, err := strconv.ParseUint(prefix, 10, 64)
		if err != nil {
			continue
		}
		latest = max(latest, uint(v))
	}

	return latest, nil
}

// compareVersions compares dotted numeric versions like "3.53.0".
func compareVersions(a, b string) int {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")

	for i := 0; i < max(len(pa), len(pb)); i++ {
		var na, nb int
		if i < len(pa) {
			na, _ = strconv.Atoi(pa[i])
		}
		if i < len(pb) {
			nb, _ = strconv.Atoi(pb[i])
		}
		if na != nb {
			if na < nb {
				return -1
			}
			return 1
		}
	}

	return 0
}

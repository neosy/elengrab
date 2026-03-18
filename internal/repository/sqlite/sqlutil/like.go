package sqlutil

import (
	"strings"

	"github.com/Masterminds/squirrel"
)

func Like(column, value string) squirrel.Sqlizer {
	return squirrel.Expr(column+" LIKE ?", "%"+strings.ToLower(value)+"%")
}

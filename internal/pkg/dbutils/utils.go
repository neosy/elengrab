package dbutils

import (
	"fmt"
	"strings"
)

func UpsertSuffix(columns []string, fields ...string) string {
	conflictFields := strings.Join(fields, ", ")

	strC := ""
	for _, c := range columns {
		strC = fmt.Sprintf("%s%s=%s,", strC, c, "EXCLUDED."+c)
	}
	return fmt.Sprintf(
		"ON CONFLICT (%s) DO UPDATE SET %s",
		conflictFields,
		strings.TrimSuffix(strC, string(strC[len(strC)-1])),
	)
}

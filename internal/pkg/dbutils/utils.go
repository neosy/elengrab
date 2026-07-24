package dbutils

import (
	"fmt"
	"strings"
)

// UpsertSuffix builds an `ON CONFLICT ... DO UPDATE SET ...` SQL clause.
// If no conflict columns are provided, it returns an empty string.
// Each column in updateColumns is assigned from its corresponding EXCLUDED value.
func UpsertSuffix(updateColumns []string, conflictColumns ...string) string {
	if len(conflictColumns) == 0 {
		return ""
	}

	setParts := make([]string, 0, len(updateColumns))
	for _, col := range updateColumns {
		setParts = append(setParts, col+" = EXCLUDED."+col)
	}

	return fmt.Sprintf(
		"ON CONFLICT (%s) DO UPDATE SET %s",
		strings.Join(conflictColumns, ", "),
		strings.Join(setParts, ", "),
	)
}

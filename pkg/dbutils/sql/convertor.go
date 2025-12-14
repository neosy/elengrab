package usql

import (
	"database/sql"
	"time"
)

func String(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

func Int64(v sql.NullInt64) int64 {
	if v.Valid {
		return v.Int64
	}
	return 0
}

func Time(v sql.NullTime) time.Time {
	if v.Valid {
		return v.Time
	}
	return time.Time{}
}

func ToNullString(v string) sql.NullString {
	if v == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: v, Valid: true}
}

func ToNullInt64(v int64) sql.NullInt64 {
	return sql.NullInt64{Int64: v, Valid: true}
}

func ToNullTime(v time.Time) sql.NullTime {
	if v.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: v, Valid: true}
}

func ToNullBool(v bool) sql.NullBool {
	return sql.NullBool{Bool: v, Valid: true}
}

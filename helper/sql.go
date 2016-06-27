package helper

import "database/sql"

// ToNullString invalidates a sql.NullString if empty, validates if notempty
func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

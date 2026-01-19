package handlers

import "database/sql"

func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func toNullInt32FromInt64(v *int64) sql.NullInt32 {
	if v == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(*v), Valid: true}
}

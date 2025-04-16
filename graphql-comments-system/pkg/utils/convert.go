package utils

import (
	"database/sql"
	"strconv"
)

func PtrToString(nullInt sql.NullInt64) *string {
	if !nullInt.Valid {
		return nil
	}
	s := strconv.FormatInt(nullInt.Int64, 10)
	return &s
}

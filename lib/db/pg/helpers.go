package pg

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

// NewQueryBuilder creates a new squirrel query builder with dollar placeholders
func NewQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

// IsDuplicate extracts a unique_constraint_violation from a database error.
func IsDuplicate(err error) (constraint string, ok bool) {
	if err == nil {
		return
	}

	if pgErr, ok1 := err.(*pq.Error); ok1 {
		if pgErr.Severity == "ERROR" && pgErr.Code == "23505" {
			constraint = pgErr.Constraint
			ok = true
			return
		}
	}

	return
}

// IsPLv8Error extracts an error thrown by a PLv8 function.
func IsPLv8Error(err error) (msg string, ok bool) {
	if err == nil {
		return
	}

	if pgErr, ok1 := err.(*pq.Error); ok1 {
		if pgErr.Severity == "ERROR" && pgErr.Code == "XX000" {
			msg = pgErr.Message
			ok = true
			return
		}
	}

	return
}

// Scannable matches the interface for the scannable sql.Row/sql.Rows
type Scannable interface {
	Scan(dest ...interface{}) error
}

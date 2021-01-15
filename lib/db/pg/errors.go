package pg

import (
	"github.com/lib/pq"
)

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

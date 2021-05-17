package pg

import (
	"github.com/jackc/pgconn"
	"github.com/lib/pq"
)

type Error struct {
	Severity   string
	Code       string
	Message    string
	Constraint string
}

// IsDuplicate extracts a unique_constraint_violation from a database error.
func IsDuplicate(err error) (constraint string, ok bool) {
	if err == nil {
		return
	}

	pe := convertError(err)

	if pe.Severity == "ERROR" && pe.Code == "23505" {
		constraint = pe.Constraint
		ok = true
		return
	}

	return
}

// IsPLv8Error extracts an error thrown by a PLv8 function.
func IsPLv8Error(err error) (msg string, ok bool) {
	if err == nil {
		return
	}

	pe := convertError(err)

	if pe.Severity == "ERROR" && pe.Code == "XX000" {
		msg = pe.Message
		ok = true
		return
	}

	return
}

func convertError(err error) (pe Error) {
	switch e := err.(type) {
	case *pq.Error:
		pe = Error{e.Severity, string(e.Code), e.Message, e.Constraint}
	case *pgconn.PgError:
		pe = Error{e.Severity, e.Code, e.Message, e.ConstraintName}
	}
	return
}

package pg

import (
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestIsDuplicate(t *testing.T) {
	tests := []struct {
		Name string

		Error       error
		Constraint  string
		IsDuplicate bool
	}{
		{"NoError", nil, "", false},
		{"NotPgError", errors.New("foo"), "", false},
		{
			"PrimaryKey",
			&pq.Error{
				Severity:   "ERROR",
				Code:       "23505",
				Message:    "duplicate key value violates unique constraint \"events_pkey\"",
				Detail:     "Key (id)=(abc123) already exists.",
				Schema:     "public",
				Table:      "events",
				Constraint: "events_pkey",
				File:       "nbtinsert.c",
				Line:       "434",
				Routine:    "_bt_check_unique",
			},
			"events_pkey", true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			con, ok := IsDuplicate(test.Error)
			if test.IsDuplicate {
				if assert.True(t, ok) {
					assert.Equal(t, test.Constraint, con)
				}
			} else {
				assert.False(t, ok)
			}
		})
	}
}

func TestIsPLv8Error(t *testing.T) {
	tests := []struct {
		Name string

		Error   error
		Message string
		IsPLv8  bool
	}{
		{"NoError", nil, "", false},
		{"NotPgError", errors.New("foo"), "", false},
		{
			"Throw",
			&pq.Error{
				Severity: "ERROR",
				Code:     "XX000",
				Message:  "hello world",
				Where:    "undefined() LINE 1:  throw new Error('hello world'); ",
				File:     "plv8.cc",
				Line:     "1878",
				Routine:  "rethrow",
			},
			"hello world", true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			msg, ok := IsPLv8Error(test.Error)
			if test.IsPLv8 {
				if assert.True(t, ok) {
					assert.Equal(t, test.Message, msg)
				}
			} else {
				assert.False(t, ok)
			}
		})
	}
}

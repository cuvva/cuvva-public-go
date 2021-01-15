package pg

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq" // required for the PGSQL driver to be loaded
)

// Connect opens a database and verifies a connection to the database is alive,
// establishing a connection if necessary.
func Connect(postgresURI string) (*sql.DB, error) {
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// NewQueryBuilder creates a new squirrel query builder with dollar placeholders
func NewQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

// Scannable matches the interface for the scannable sql.Row/sql.Rows
type Scannable interface {
	Scan(dest ...interface{}) error
}

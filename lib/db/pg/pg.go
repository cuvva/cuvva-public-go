package pg

import (
	"database/sql"

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

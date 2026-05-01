package libsql

import (
	"database/sql"
)

type LibSQLAdapter struct {
	db *sql.DB
}
/*
func NewLibSQLAdapter(url string) (*LibSQLAdapter, error) {
	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}
	adapter := &LibSQLAdapter{db: db}
	return adapter, nil
}
*/
// NewLibSQLAdapter now takes a ready-made *sql.DB.
// Because it no longer opens the connection, it can't fail, 
// so we don't even need to return an error anymore!
func NewLibSQLAdapter(db *sql.DB) *LibSQLAdapter {
    return &LibSQLAdapter{db: db}
}

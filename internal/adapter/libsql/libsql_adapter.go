package libsql

import (
	"database/sql"
)

type LibSQLAdapter struct {
	db *sql.DB
}

func NewLibSQLAdapter(url string) (*LibSQLAdapter, error) {
	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}
	adapter := &LibSQLAdapter{db: db}
	return adapter, nil
}



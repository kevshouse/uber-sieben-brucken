package libsql


func (a *LibSQLAdapter) Close() error { return a.db.Close() }
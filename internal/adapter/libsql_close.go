package adapter

import "log"

// Close shuts down LibSQL database connection pool gracefully
func (a *LibSQLAdapter) Close() error {
	log.Println("Closing LibSQL database connection...")
	return a.db.Close()
}

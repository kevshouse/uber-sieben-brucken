package adapter

import "log"

// Close severs the connection to the Neo4j Graph Shore
func (a *Neo4jAdapter) Close() error {
	log.Println("🌿 Closing Neo4j Graph connection...")
	return a.driver.Close(nil)
}
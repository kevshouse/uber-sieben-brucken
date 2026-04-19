# Pothole Forecast

## 1\. Folly: Naming

Changing the folder name later but forgetting to update the go.mod file. The Reality: Go is very sensitive to this. If we name our directory uber-sieben-brucken to something else on the HDD, the code will still work as the module line inside go.mod stays the same. The go.mod is the "identity", the folder is the just the "Location".

## 2\. Folly: The "Global State" Temptation

Writing a main.go that opens a database connection and passes it around as a global variable.

The Reality: This is the quickest way to create "race conditions" and un-testable code.

The Mitigation: We will use Dependency Injection. We will define a "Database Adapter" struct in Go, and pass that struct into our core logic. This keeps the "Bridges" clean and follows ADR 002 regarding Hexagonal Architecture.

---

## 3\. Folly: The "Split Integrity"

Writing code that creates PREVIOUS but forgets to create NEXT.  
**The Reality:** This creates a "One-Way Street" where you can roll back, but you can never roll forward again.

**The Mitigation:** Use the Go Neo4j Driver's Transaction API. The CREATE of both relationships must occur within the same tx.Run() call to ensure the graph never enters a corrupted state.

---


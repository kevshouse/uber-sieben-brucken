# Pothole Forecast

## 1. Naming
   The Folly: Changing the folder name later but forgetting to update the go.mod file. 
   - **The Reality:** Go is very sensitive to this. If we name our directory uber-sieben-brucken to something else on the HDD, the code will still work as the module line inside go.mod stays the same. The go.mod is the "identity", the folder is the just the "Location".

## 2. The "Global State" Temptation

   The Folly: Writing a main.go that opens a database connection and passes it around as a global variable.
   - **The Reality:** This is the quickest way to create "race conditions" and un-testable code.

   - **The Mitigation:** We will use Dependency Injection. We will define a "Database Adapter" struct in Go, and pass that struct into our core logic. This keeps the "Bridges" clean and follows ADR 002 regarding Hexagonal Architecture.

---

## 3. The "Split Integrity"

   The Folly: Writing code that creates PREVIOUS but forgets to create NEXT.  
   - **The Reality:** This creates a "One-Way Street" where you can roll back, but you can never roll forward again.

   - **The Mitigation:** Use the Go Neo4j Driver's Transaction API. The CREATE of both relationships must occur within the same tx.Run() call to ensure the graph never enters a corrupted state.

---

## 4. The "N-Plus-One" Ghost

    The Folly:  Trying to fetch a snippet and all its citations by making separate database calls for every single link.

   - **The Reality:** As the graph grows, your "Oysters" (data) will become sluggish if you don't use Cypher's power to pull the entire "Triangle" in one handshake.

   - **The Mitigation:** We will write a "Deep Fetch" Cypher query that retrieves the Snippet Anchor, its Latest Version, and all active Citation Nodes in a single atomic transaction.

## 5. 🧭 The "Import Ghost"
   
    The Folly: Adding time.Time to the new structs but forgetting that time must be imported at the top.

   - **The Reality:** Go will refuse to compile if import "time" is missing, even if it was there before.

   - **The Mitigation:** Always check the very top of the file after an edit. If you have time.Time in your code, you MUST have import "time" at the top.

## 6. 🚧 The "Sync Gap" 
   

    The Folly: Thinking the two databases will always stay perfectly in sync automatically.
   - **The Reality:** If a Snippet is created in libSQL but the Neo4j connection fails, you end up with a "Ghost Shore"—an identity with no history.
   - **The Mitigation:** 
     In the Service Layer, we will use Go's error handling to ensure that if the second "bridge" fails, we report a failure to the user so they know the handshake was incomplete.

## 7. 🚧 The "Missing Shore"


	The Folly: Running CiteSnippet before the snippets are created in Neo4j.
   - **The Reality:** Unlike MERGE, the ```MATCH``` command requires the nodes to exist. If Snippet A or Snippet B is missing from Neo4j, the query will return nothing, and no citation will be created.

   - **The Mitigation:** We must ensure that whenever a snippet is "born" in libSQL, we also create a "Shadow Anchor" in Neo4j.

## 8. 🚧 The "Unused Import" Error

    The Folly: Keeping an import that you aren't actually using in the code.
   - **The Reality:** Go is one of the few languages that treats an unused import as a compile error. It won't let you build the bridge if there’s a loose stone (unused code) sitting on it.
   - **The Mitigation:** If you ever delete the setupSnippetAnchor helper, you must also delete the neo4j import, or the "walrus" won't even wake up to run the tests.

## 9. 🚧 Need for ```go mod tidy```
    The Folly:Because I added the testify package to make life easier, go will get confused.
	
   - **The Mitigation** One must run ```go mod tidy``` in the terminal before the tests will work:

## 10. 🚧 The "Hidden Character

     The Folly: Having a space or a comment before the package line.

   - **The Reality:** While comments are usually fine, Go is very picky. The package statement should effectively be the first non-comment text the compiler sees.

   - **The Mitigation:** Ensure there are no stray characters or weird encoding artifacts at the very top of the file.

## 11. 🚧 The "Dirty Cache"

     The Folly: Thinking the cache will always catch breaking changes in external dependencies (like your Docker Neo4j instance).

   - **The Reality:** Sometimes the database state changes, but the code doesn't.

   - **The Mitigation:** If you ever want to force a "Fresh Crossing" to be absolutely sure, you can run go test -count=1 -v ./internal/adapter/.... The -count=1 flag tells Go to ignore the cache and run every test manually.

## 12. 🚧 The "Dangling Link"

The Folly: Deleting a snippet without cleaning up its ALNs.

   - **The Reality:** If Snippet B is deleted, the ALN pointing to it becomes a "bridge to nowhere."

   - **The Mitigation:** Our spec should note that ALNs require referential integrity across the graph. When we eventually implement deletion, we must use a DETACH DELETE pattern to ensure the current doesn't leave behind ghostly remains.

## 13. 🚧 The "Interface Mismatch"

The Folly: Adding a function to the Neo4jAdapter but forgetting to add it to the GraphRepository interface in ports.go.

   - **The Reality:** The SnippetService will refuse to compile. It only "sees" what is written in the interface. Even if the adapter has the power to do more, the service can only press the buttons defined in the contract.

   - **The Mitigation:** Always update ports.go first (The Blueprint), then update your Adapters (The Tools).

## 14. 🚧 The "Mock vs. Reality" Gap
The Folly: Thinking that because the Service test passed, the whole app is finished.

- **The Reality:** The Service test proves the logic is right, but it doesn't check if your SQL queries are valid or if Neo4j is online.

- **The Mitigation:** This is why we have the Adapter Tests we ran yesterday. Together, these two layers of testing (Service + Adapter) create an unbreakable "Bridge."

## 15. 🚧 The "JSON Marshalling" Trap
The Folly: Forgetting that Go is strictly typed while JSON is flexible.

- **The Reality:** If a user sends a string where you expect an integer, the bridge will "creak."

- **The Mitigation:** We will define specific Request Structs with JSON tags to ensure the data is validated the moment it hits the gate.

## 16 🚧 The "Request Body" Ghost
The Folly: Trying to read r.Body twice.

- **The Reality:** In Go, the request body is a "Stream." Once you read it (using the JSON decoder), it is consumed. If you try to read it again, it will be empty.

- **The Mitigation:** We decode it once into our createSnippetRequest struct and then immediately move to the logic.

## 17 🚧 The "Missing ID" Creak
The Folly: A user sends a request where source_id or target_id is empty.

- **The Reality:** Our Service will try to tell Neo4j to link "nothing" to "nothing," and the graph will complain.

- **The Mitigation:** In a more advanced version, we would add Validation. For now, if the IDs don't exist in the database, our Service will return an error, and this handler will correctly pass that error back to the user as a 500 Internal Server Error.

## 🚧 The "Import Conflict"
The Folly: Importing net/http and your own internal/adapter/http package at the same time.

- **The Reality:** Go won't allow two packages with the same name.

- **The Mitigation:** In the import block, I suggest aliasing your internal package as api to avoid a collision:
api "github.com/kevshouse/uber-sieben-brucken/internal/adapter/http"
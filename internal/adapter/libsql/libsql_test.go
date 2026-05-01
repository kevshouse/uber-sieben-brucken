package libsql_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/kevshouse/uber-sieben-brucken/internal/adapter/libsql"
	"github.com/kevshouse/uber-sieben-brucken/internal/core"

	_ "github.com/tursodatabase/go-libsql"
)

// setupTestDB gives us a clean in-memory database for each test
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("libsql", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// We must run our migrations here so the schema is set up for our tests
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS snippets (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		owner_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createTableQuery); err != nil {
		t.Fatalf("Failed to create snippets table: %v", err)
	}

	return db
}

func TestLibSQLAdapter_Save(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	// Init our adapter (the Shore)
	adapter := libsql.NewLibSQLAdapter(db)
	
	t.Run("successfully saves a new snippet", func(t *testing.T) {
		// 1. Arrange: Create a test snippet
		snippet := core.Snippet{
			ID:      "snip123",
			Title:   "My First Snippet",
			OwnerID: "user_456",
		}	
		ctx := context.Background()

		// 2. Act: Save it using our adapter
		err := adapter.Save(ctx, &snippet)
		if err != nil {
			t.Fatalf("expected no error on save, got: %v", err)
		}

		// 3. Assert: Query the db directly to ensure the data was actually written
		var savedTitle, savedOwner string
		err = db.QueryRow("SELECT title, owner_id FROM snippets WHERE id = ?", snippet.ID).
			Scan(&savedTitle, &savedOwner)
			
		if err != nil {
			t.Fatalf("failed to query test database: %v", err)
		}
		
		if savedTitle != snippet.Title {
			t.Errorf("expected title %q, got %q", snippet.Title, savedTitle)
		} 
		
		if savedOwner != snippet.OwnerID {
			t.Errorf("expected owner_id %q, got %q", snippet.OwnerID, savedOwner)
		}
	})
	// Save a snippet to the database using the adapter. 
	// Then, modify the same snippet struct in memory 
	// (change the Title and OwnerID, but leave the ID exactly the same).
	t.Run("saves updates to an existing snippet", func(t *testing.T) {
		// Arrange: Create and save an initial snippet
		snippet := core.Snippet{
			ID:      "snip456",
			Title:   "Original Title",
			OwnerID: "user_789",
		}	
		ctx := context.Background()

		if err := adapter.Save(ctx, &snippet); err != nil {
			t.Fatalf("failed to save initial snippet: %v", err)
		}

		//2. Act: Modify the snippet's title and owner, then save again
		snippet.Title = "Updated Title"
		snippet.OwnerID = "user_999"

		if err := adapter.Save(ctx, &snippet); err != nil {
			t.Fatalf("failed to save updated snippet: %v", err)
		}
		
		//3. Assert: Query the db directly to ensure the updates were saved
		var savedTitle, savedOwner string
		err := db.QueryRow("SELECT title, owner_id FROM snippets WHERE id = ?", snippet.ID).
			Scan(&savedTitle, &savedOwner)
		if err != nil {

			t.Fatalf("failed to query test database: %v", err)
		}
		
		if savedTitle != snippet.Title {
			t.Errorf("expected title %q, got %q", snippet.Title, savedTitle)
		}
		
		if savedOwner != snippet.OwnerID {
			t.Errorf("expected owner_id %q, got %q", snippet.OwnerID, savedOwner)
		}
		//4. Bonus Assert: Ensure we didn't accidentally duplicate rows!
		var totalRows int
        err = db.QueryRow("SELECT COUNT(*) FROM snippets").Scan(&totalRows)
        if err != nil {
            t.Fatalf("failed to count total rows: %v", err)
        }
        
        // Since our first test saved 1 snippet, and this test saved 1 snippet, 
        // we should have exactly 2 rows in the entire table right now.
        if totalRows != 2 {
            t.Errorf("expected exactly 2 rows in the table, got %d", totalRows)
        }
	})

}
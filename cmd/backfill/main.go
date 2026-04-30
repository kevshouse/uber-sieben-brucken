package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/kevshouse/uber-sieben-brucken/internal/adapter/libsql"
	"github.com/kevshouse/uber-sieben-brucken/internal/adapter/neo4j"
	"github.com/kevshouse/uber-sieben-brucken/internal/core"
)

// IdentityRepository defines the minimal port required by libSQL for this script.
type IdentityRepository interface {
	GetAll(ctx context.Context) ([]*core.Snippet, error)
}

// GraphRepository defines the minimal port required by Neo4j for this script.
type GraphRepository interface {
	SaveVersion(ctx context.Context, snippet *core.Snippet, version *core.Version) error
}

// Runner orchestrates the backfill migration process for Identity to Graph.
type Runner struct {
	idRepo    IdentityRepository
	graphRepo GraphRepository
}

// NewRunner creates a new migration Runner.
func NewRunner(id IdentityRepository, graph GraphRepository) *Runner {
	return &Runner{
		idRepo:    id,
		graphRepo: graph,
	}
}

// Run executes the 1-way sync logic.
func (r *Runner) Run(ctx context.Context) error {
	// Step 1: Fetch all snippets from Identity.
	snippets, err := r.idRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch snippets from Identity: %w", err)
	}

	// Step 2: Iterate and sync each snippet to Graph database.
	for _, snippet := range snippets {
		// Passing nil for the Version metadata for this backfill step.
		err := r.graphRepo.SaveVersion(ctx, snippet, nil)
		if err != nil {
			return fmt.Errorf("failed to save snippet %s to graph: %w", snippet.ID, err)
		}
	}

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 1. Fix the naming to match terminal command exactly
	libSqlURL := os.Getenv("LIBSQL_URL")
	if libSqlURL == "" {
		libSqlURL = "http://localhost:8081" // Default to Docker shore
	}

	neo4jURI := os.Getenv("NEO4J_URI")
	neo4jUser := os.Getenv("NEO4J_USER")
	neo4jPass := os.Getenv("NEO4J_PASS")

	if neo4jURI == "" {
		log.Fatal("NEO4J_URI environment variable is required.")
	}

	// 2. Initialize libSQL (The Identity Shore)
	log.Printf("🚀 Connecting to libSQL at %s...\n", libSqlURL)
	idRepo, err := libsql.NewLibSQLAdapter(libSqlURL)
	if err != nil {
		log.Fatalf("Failed to initialize libSQL: %v", err)
	}

	// 3. Initialize Neo4j (The Graph Shore)
	log.Printf("🌿 Connecting to Neo4j at %s...\n", neo4jURI)
	graphRepo, err := neo4j.NewNeo4jAdapter(neo4jURI, neo4jUser, neo4jPass)
	if err != nil {
		log.Fatalf("Failed to initialize Neo4j: %v", err)
	}

	runner := NewRunner(idRepo, graphRepo)
	log.Println("🌉 Starting backfill migration...")
	
	if err := runner.Run(ctx); err != nil {
		log.Fatalf("Backfill failed: %v", err)
	}

	log.Println("✨ Backfill completed successfully!")
}
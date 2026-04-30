package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"

	api "github.com/kevshouse/uber-sieben-brucken/internal/adapter/http"
	"github.com/kevshouse/uber-sieben-brucken/internal/adapter/libsql"
	"github.com/kevshouse/uber-sieben-brucken/internal/adapter/neo4j"
	"github.com/kevshouse/uber-sieben-brucken/internal/core"
)

func main() {
	libsqlURL := getEnv("LIBSQL_URL", "file:local.db")
	neo4jURL  := getEnv("NEO4J_URL", "neo4j://localhost:7687")
	neo4jUser := getEnv("NEO4J_USER", "neo4j")
	neo4jPass := getEnv("NEO4J_PASS", "password123")

	// 1. Initialize Adapters
	idRepo, err := libsql.NewLibSQLAdapter(libsqlURL)
	if err != nil {					
		log.Fatalf("Failed to connect to libSQL: %v", err)
	}

	graphRepo, err := neo4j.NewNeo4jAdapter(neo4jURL, neo4jUser, neo4jPass)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}

	// 2. Wiring up the service and HTTP handlers
	service := core.NewSnippetService(idRepo, graphRepo)
	handler := api.NewHandler(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/snippets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CiteSnippet(w, r)
		}
	})

	// mux Search Snippets
	mux.HandleFunc("/search", handler.SearchSnippets )

	srv := &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	// 3. Graceful shutdown Channel
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine to allow for graceful shutdown handling
	go func() {
		log.Println("🌉 Bridge is open on :9090")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for signal
	<-done
	log.Println("Shutting down bridge gracefully...")

	// 4. Cleanup with Timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	// Close database connections using the new Close() contract
	if err := idRepo.Close(); err != nil {
		log.Printf("Error closing LibSQL: %v", err)
	}
	if err := graphRepo.Close(); err != nil {
		log.Printf("Error closing Neo4j: %v", err)
	}

	log.Println("Safe Landing. Server exited.")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}


/*
	http.HandleFunc("/snippets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateSnippet(w, r)
		}
	})

	http.HandleFunc("/citations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CiteSnippet(w, r)
		}
	})

	http.HandleFunc("/search", handler.SearchSnippets)

	log.Println("🌉 Bridge is open on :9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
*/

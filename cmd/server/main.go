package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/ncruces/go-sqlite3/driver"

	"github.com/kevshouse/uber-sieben-brucken/internal/adapter"
	api "github.com/kevshouse/uber-sieben-brucken/internal/adapter/http"
	"github.com/kevshouse/uber-sieben-brucken/internal/core"
)

func main() {
	libsqlURL := getEnv("LIBSQL_URL", "file:local.db")
	neo4jURL  := getEnv("NEO4J_URL", "neo4j://localhost:7687")
	neo4jUser := getEnv("NEO4J_USER", "neo4j")
	neo4jPass := getEnv("NEO4J_PASS", "password123")

	idRepo, err := adapter.NewLibSQLAdapter(libsqlURL)
	if err != nil {
		log.Fatalf("Failed to connect to libSQL: %v", err)
	}

	graphRepo, err := adapter.NewNeo4jAdapter(neo4jURL, neo4jUser, neo4jPass)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}

	service := core.NewSnippetService(idRepo, graphRepo)
	handler := api.NewHandler(service)

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


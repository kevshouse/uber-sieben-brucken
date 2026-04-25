# Project State Document: Day 8 Conclusion

## Architecture Overview
* **Pattern:** Strict Hexagonal (Ports & Adapters).
* **The Brain:** `SnippetService` orchestrates domain logic.
* **The Ports:** `GraphRepository` and `IdentityRepository` enforce strict type safety (`*core.Snippet`).
* **Primary Adapters:** HTTP REST API and a localized `cmd/backfill` utility script.

## The Port Map (Single Source of Truth)
* **`IdentityRepository`** (Relational Store - `internal/adapter/libsql_adapter.go`)
  * `CreateSnippet(ctx context.Context, s *core.Snippet) error`
  * `Search(ctx context.Context, query string) ([]*core.Snippet, error)`
  * `GetAll(ctx context.Context) ([]*core.Snippet, error)` *(Added for Backfill)*
* **`GraphRepository`** (Graph Store - `internal/adapter/neo4j_adapter.go`)
  * `SaveVersion(ctx context.Context, snippet *core.Snippet, version *core.Version) error`
  * `CiteSnippet(ctx context.Context, sourceID, targetID string) error`

## Recent System Alignments
* **Pragmatic TDD:** Successfully built a localized migration script runner (`cmd/backfill`) using a Red-Green-Refactor loop with hand-rolled mocks.
* **Composition Root Wiring:** The `main()` function in the backfill utility successfully handles environment configurations, context timeouts, and database driver instantiation.
* **Interface Satisfaction:** Extended `LibSQLAdapter` with a `GetAll` method to satisfy the script's localized `IdentityRepository` contract, and synchronized the `GraphRepository` signature to handle version metadata.
* **Core Integrity:** The `internal/core` package remains completely untouched by one-off migration logic.

## Immediate Roadmap & Known Technical Debt
1. **Search Implementation:** The HTTP handler logic for `SearchSnippets` needs to be finalized to leverage the `Search` capability we verified in the `LibSQLAdapter`.
2. **Execute the Backfill:** Run the script against the live local databases and verify that Neo4j correctly hydrates legacy nodes without spawning unnecessary "Migration" versions.
3. **Delete Propagation (Garbage Collection):** Establish a clear policy and implementation for removing nodes/edges in the graph when a snippet is deleted from the relational store.
4. **Adapter Teardown (Tech Debt):** Implement graceful `Close()` methods for both the `LibSQLAdapter` and `Neo4jAdapter` to prevent resource leaks in long-running services.
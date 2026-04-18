# Status: Proposed
## Date: April 2026
**Context: We** need a way to track the evolution of code snippets over time. Traditional SQL "history tables" are cumbersome to query for deep lineages, and Git-style storage can be overkill for a lightweight web-API.

## Decision
We will use Neo4j (Cypher) as the primary engine for version history and Go as the orchestrator.

## Rationale
Topology over Records: By representing versions as nodes and edits as relationships (:PREVIOUS_VERSION), we can use native graph traversals to "walk" the history.

 **Cypher Strengths: We** will leverage Cypher’s variable-length path queries (-[:PREVIOUS*]->) to retrieve entire lineages in a single operation, rather than performing multiple joins or recursive CTEs in SQL.

**Self-Healing Logic: By** utilizing an append-only graph model, the "Current" state is simply a pointer. If a corrupted state is entered, the "Bridge" (pointer) can be moved back to a previous valid node without data loss.

## Consequences
**Positive: Highly efficient** history lookups; clear visualization of collaboration "paths"; lower logic complexity in Go.

**Negative: Requires** maintaining a Neo4j instance alongside our primary relational store (libSQL); higher initial setup time for the graph driver.
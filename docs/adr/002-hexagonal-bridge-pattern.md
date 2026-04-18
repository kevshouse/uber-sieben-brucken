# Status: Proposed#
***Context: The UI*** is expected to evolve from a "Minimalist" dashboard to a "Feature-Rich" interactive graph. We need to ensure the backend logic isn't coupled to the initial UI choices.

##Decision##
We will implement Hexagonal Architecture (Ports and Adapters) in Go.

## Rationale
***Isolation: The "Core" (the Königsberg logic)*** will define interfaces (Ports) for persistence and notification.

## Swappability:
***We can build a REST Adapter*** for the minimalist UI today and a WebSocket/GraphRAG Adapter tomorrow without modifying the core Versioning logic.

## Testability: 
**We can mock the Neo4j "Bridge"** to run lightning-fast unit tests on our logic in Go.
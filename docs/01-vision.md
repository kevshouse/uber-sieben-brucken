# Vision
## Project Name: Über sieben Brücken musst du gehn

1. ### The "Königsberg" Problem Statement

 **In** the physical world of 1736, Euler proved that the topology of connections matters more than the distance between them. In the software world of 2026, we often lose this topology. Code reviews and collaboration are treated as static snapshots rather than a path of decisions. This project aims to bridge the gap between "what the code is" and "how the code became," utilizing the inherent strengths of Graph Theory.

2. ### High-Level Objectives

   ***Temporal Integrity: Use*** a Graph-based versioning chain to ensure that every "state" of a code snippet is immutable and reachable.

   ***Engineering Resilience: Demonstrate*** a transition from C-system logic to Go-orchestration, proving that the rigor of 42 London scales to modern cloud architectures.

   ***Operational Clarity: Use** a "Hexagonal" approach to ensure the business logic of "The Path" remains untouched by the changing "Bridges" (UI/Infrastructure) we cross.

3. ### Core Features (Phase 1)

   ***The Immutable Path: A Go-managed backend*** that appends new snippet versions to a Neo4j chain.

   ***The Temporal Walk: A query engine*** that allows a user to "walk" the history of a snippet back to its root.

   ***The Minimalist Gateway: A clean,*** React-based dashboard focused on clarity and "Time-Travel" navigation.

4. ### The "No-Go" Zone (Scope Control)

   ***No multi-cursor live editing*** (Keep it to discrete version commits).

   ***No code execution/sandboxing*** (Focus on the metadata and collaboration).

   ***No microservices*** (Keep it as a clean, single Go binary for Phase 1).
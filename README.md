![Seven Bridges CI](https://github.com/kevshouse/uber-sieben-brucken/actions/workflows/go_tests.yml/badge.svg)
# Über sieben Brücken musst du gehn

> *"The topology of the path matters more than the distance between the shores."*

## 🌉 The Vision
Named after the folk song and the "Seven Bridges of Königsberg" problem solved by Leonhard Euler, this project is a design study in **Temporal Code Integrity**. It explores how code evolution can be managed not just as a series of snapshots, but as an immutable graph of decisions.

Built as a modern evolution of the systems-thinking rigor found at **42 London**, this project demonstrates how to bridge the gap between low-level architectural discipline and high-level cloud-native resilience.

## 🛠 Tech Stack (2026 Edition)
- **Orchestrator:** Go 1.26.2
- **Temporal Engine:** Neo4j 5.15 (Graph Database)
- **Metadata Shore:** libSQL (Primary Relational Store)
- **Infrastructure:** Docker Compose (Virtualized Environment)
- **Architecture:** Hexagonal (Ports and Adapters)

## 🏗 Project Structure
- `internal/core`: The "Königsberg" domain logic—pure, database-agnostic Go.
- `internal/adapter`: The "Bridges"—implementations for Neo4j, libSQL, and HTTP.
- `docs/`: The Turbo Shell—living documentation, schema topologies, and pothole logs.
- `docs/adr/`: Architectural Decision Records—the immutable log of why we built it this way.

## 🚦 Getting Started
1. **Ensure Docker is running.**
2. **Launch the Infrastructure:**
   ```bash
   docker-compose up -d
3. **Verify the Graph:**
   Access the Neo4j Browser at http://localhost:7474 (User: neo4j / Pass: password123).

4. **Verify the Metadata:**
   libSQL is listening on localhost:8081.

## 📝 The Pothole Philosophy
 > *Success is built on falling then climbing, not the other way around!*

- This project maintains an active Pothole Log (docs/pothole_log.md).  
- We document every engineering friction point—from YAML indentation traps to port conflicts—ensuring that the "Bridge" we build is informed by the failures of the past.
---

Developed as an independent engineering study, April 2026.
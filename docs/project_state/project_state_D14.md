# Project State: Uber-Sieben-Brucken
**Date:** 2026-05-03
**Architect:** GoHex Architect & Lead Dev

---

## 🏗️ Architecture Overview
*   **Pattern:** Strict Hexagonal Architecture (Ports & Adapters).
*   **Core:** Domain entities (`Snippet`, `Version`, `Citation`) and Orchestration Services.
*   **Adapters:** 
    *   **Identity Shore:** libSQL (Metadata/Folders).
    *   **Graph Shore:** Neo4j (History/Relationships).
    *   **Driving Adapter:** HTTP/REST Handler.
*   **Testing:** Testify-backed mocks localized in `internal/core/mocks.go` for cross-package accessibility.

## ✅ Recent System Alignments
*   **The Genesis Protocol:** Successfully implemented and tested `CreateGenesis` in `SyncService`. It now generates UUIDs and timestamps in the Core and persists to both Shores simultaneously.
*   **Mocking Infrastructure:** Established `internal/core/mocks.go` as a first-class citizen, providing `MockIdentityRepo` and `MockGraphRepo` that perfectly satisfy the Core Port contracts (including `Close`, `GetAll`, and structured `CiteSnippet`).
*   **TDD Success:** The `services_test` suite is now fully Green, covering both `LiveSync` and `CreateGenesis` flows.

## 🚀 Immediate Roadmap
1.  **Search Implementation:** Wire the `SearchSnippets` method from the HTTP Adapter through the Brain and down to the Identity Shore.
2.  **Citation Logic:** Implement the `Cite` method to link two existing snippets in the Graph Shore.
3.  **Error Compensation:** Tackle the "Ghost Shore" technical debt—adding logic to handle partial failures (e.g., if Identity saves but Graph fails).

## ⚠️ Known Technical Debt
*   **Circular Dependencies:** Mocks are currently inside `package core` to facilitate testing; long-term, consider move to a dedicated `testing` sub-package if the domain grows.
*   **Rollback Mechanism:** `CreateGenesis` currently lacks a transaction-like rollback across the two disparate databases.
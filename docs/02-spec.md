# **02-spec.md: Technical Specification**

## **1\. System Architecture**

We follow a **Hexagonal Architecture** (Ports and Adapters) to ensure the Königsberg logic remains decoupled from specific database implementations.

* **Internal/Core:** Contains the "Path" logic, Snippet entities, and the "Self-Healing" algorithms.  
* **Internal/Adapter:** Contains the specific implementations for:  
  * `neo4j_adapter.go`: Communicating with the Graph.  
  * `libsql_adapter.go`: Communicating with the Metadata store.  
  * `http_adapter.go`: Handling REST requests.

## **2\. The "Handshake" Protocol**

Every operation that modifies the system must follow this sequence:

1. **Metadata Check:** Verify User/Snippet existence in libSQL.  
2. **Graph Mutation:** Append the new `Version` node and link it via `[:PREVIOUS]`.  
3. **Atomic Update:** Shift the `[:HAS_LATEST]` pointer in Neo4j.

## **3\. Performance Constraints**

* **Traversal Depth:** History "walks" should be limited to 50 nodes per request for the initial prototype to prevent browser hanging.  
* **Concurrency:** Go routines must be used to ping both databases simultaneously during the "Heartbeat" check to minimize latency.


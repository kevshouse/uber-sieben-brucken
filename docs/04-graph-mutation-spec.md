This document will serve as the "Bridge Blueprint." By combining the **Cypher** logic with **Mermaid** visualizations, we create a specification that ensures the Go code we write later is just a mechanical implementation of this design.

`docs/04-graph-mutation-spec`

---

# High-performance **Temporal Navigation Engine**.

---

# **The Bidirectional Temporal Path**

## **1\. Structural Overview**

The graph represents code evolution as a **doubly-linked chain** anchored to a static identity. This topography allows for high-speed traversal in both directions:

* **Backward (`PREVIOUS`):** For rollbacks and historical audits.  
* **Forward (`NEXT`):** For "replaying" code evolution and watching the path unfold.

### **The Entities**

* **`Snippet`**: The static anchor (Identity Shore).  
* **`Version`**: The immutable milestone (Temporal Path).

---

## **2\. Visual Topology (Mermaid)**

The following diagram illustrates the bidirectional flow between the anchor and its evolving milestones.

```mermaid  
    graph TD  
    S[Snippet Anchor] --|HAS_LATEST|--> V3((Version 3))  
    V3 --PREVIOUS--> V2((Version 2)) --NEXT--> V3  
    V2 --PREVIOUS--> V1((Version 1 - Genesis))  
    V1 --NEXT--> V2  
    

    style S fill:\#f9f,stroke:\#333,stroke-width:4px  
    style V3 fill:\#00ff00,stroke:\#333,stroke-width:2px  
    style V1 stroke-dasharray: 5 5  
```
---

## **4\. Design Reasoning (Spotnet Legacy)**

* **Traversal Efficiency:** Finding "what happened after Version X" in a standard graph requires a broad search. With the `NEXT` relationship, it is a direct pointer lookup (O(1)).  
* **Chain Integrity:** Every version (except Genesis) must have exactly one `PREVIOUS` and one `NEXT` (except for the Latest).  
* **Self-Healing:** If the system detects a `Version` with a `PREVIOUS` but no `NEXT` pointing back to it, we have identified a "Broken Bridge" that needs manual reconciliation.

---

## **5\. Pothole Forecast: The "Split Integrity"**

**The Folly:** Writing code that creates `PREVIOUS` but forgets to create `NEXT`. **The Reality:** This creates a "One-Way Street" where you can roll back, but you can never roll forward again. **The Mitigation:** Use the Go Neo4j Driver's **Transaction API**. The `CREATE` of both relationships must occur within the same `tx.Run()` call to ensure the graph never enters a corrupted state.

---

### **Next Milestone: `internal/core/ports.go`** *Completed on day 3*

With the spec now reflecting the bidirectional topography, we can define the Go interface that will enforce this "Atomic Handshake."


🏛️ Updating the Topological Map

The ALN Definition:

"The Active Linkage Node (ALN) is a first-class graph entity that represents a versioned relationship between two Snippet Anchors. Unlike a static edge, an ALN maintains its own temporal history, allowing the context of a connection to evolve without losing its lineage."


🔍 Anatomy for Spec Documentation

```mermaid
graph TD
    subgraph libSQL_Shore [Identity Shore: libSQL]
        S1_Meta[Snippet A Metadata]
        S2_Meta[Snippet B Metadata]
    end

    subgraph Neo4j_Current [Relationship Current: Neo4j]
        %% Snippet Anchors
        S1(Snippet A Anchor)
        S2(Snippet B Anchor)

        %% Version History for Snippet A
        V1[Version A.1]
        V2[Version A.2]
        
        %% ALN Logic: Relationship as an Entity
        ALN_Prev[ALN: Initial Link]
        ALN_Latest[ALN: Active Linkage Node]

        %% Connections
        S1 --- V2
        V2 -- PREVIOUS --> V1
        
        %% The Active Linkage
        S1 -- CITES_LATEST --> ALN_Latest
        ALN_Latest -- CITES_TARGET --> S2
        
        %% ALN Lineage
        ALN_Latest -- PREVIOUS --> ALN_Prev
    end

    S1_Meta -.-> S1
    S2_Meta -.-> S2

    style ALN_Latest fill:#f96,stroke:#333,stroke-width:4px
    style ALN_Prev fill:#fff,stroke-dasharray: 5 5

```

## Explanitory Note

1. The Matching phase: The query locates the Source and Target snippet anchors.

2. The Optional Discovery: It identifies if an existing ALN already connects these two shores.

3. The Pivot: It creates a new ALN node (the new relationship context).

4. The Pointer Shift: It deletes the old CITES_LATEST edge and re-anchors the relationship to the new node, while simultaneously linking the new node to the old one via a PREVIOUS edge.
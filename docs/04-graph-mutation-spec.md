This document will serve as the "Bridge Blueprint." By combining the **Cypher** logic with **Mermaid** visualizations, we create a specification that ensures the Go code we write later is just a mechanical implementation of this design.

`docs/04-graph-mutation-spec`

---

# **04-graph-mutation-spec.md: The Temporal Path Design**

## **1\. Structural Overview**

The graph represents code evolution as a directed, acyclic path. Every snippet is an anchor, and every version is a milestone. We avoid the "Mutable Row" trap of SQL by ensuring that no `Version` node is ever modified once created.

### **The Entities**

* **`Snippet`**: The static identity. It holds the `uuid` that links to the **libSQL** metadata.  
* **`Version`**: The immutable content. It holds the code, a hash, and a timestamp.

---

## **2\. Visual Topology (Mermaid)**

This diagram illustrates how a Snippet evolves from a "Genesis" state through multiple iterations.

Code snippet  
```
Mermaid
graph TD  
    S\[Snippet Anchor\] \---|HAS\_LATEST| V3((Version 3))  
    V3 \---|PREVIOUS| V2((Version 2))  
    V2 \---|PREVIOUS| V1((Version 1 \- Genesis))

    style S fill:\#f9f,stroke:\#333,stroke-width:4px  
    style V3 fill:\#00ff00,stroke:\#333,stroke-width:2px  
    style V1 stroke-dasharray: 5 5
```
---

## **3\. The "Handshake" Mutation (Cypher)**

When the Go application receives a new version of code, it performs a three-step atomic transition in Neo4j.

### **Step 1: Identify the Bridge**

We locate the Snippet anchor and the current "Latest" version.

### **Step 2: Create the Milestone**

We create the new Version node.

### **Step 3: Shift the Pointer**

We move the `HAS_LATEST` relationship to the new node and link the new node to the old one.

#### **The Implementation Query:**

Cypher  
// 1\. Find the Anchor  
MATCH (anchor:Snippet {id: $snippetID})

// 2\. Look for the current shore (optional for Genesis versions)  
OPTIONAL MATCH (anchor)-\[oldRel:HAS\_LATEST\]-\>(previousVersion:Version)

// 3\. Construct the new Milestone  
CREATE (newVersion:Version {  
    id: $versionID,  
    content: $content,  
    hash: $hash,  
    timestamp: datetime()  
})

// 4\. Redirect the Bridge  
CREATE (anchor)-\[:HAS\_LATEST\]-\>(newVersion)  
DELETE oldRel

// 5\. Link the History  
WITH previousVersion, newVersion  
WHERE previousVersion IS NOT NULL  
CREATE (newVersion)-\[:PREVIOUS\]-\>(previousVersion)

---

## **4\. Design Reasoning (Spotnet Remnants)**

This design incorporates the "Self-Healing" logic derived from the **Spotnet** project’s logic:

* **Integrity:** If the `CREATE` of the new version fails, the `oldRel` is never deleted, keeping the system on the last known good state.  
* **Traceability:** By using a `[:PREVIOUS]` relationship instead of just an array of IDs, we can use Cypher's **Variable Length Paths** to generate diffs or histories instantly: `MATCH (v:Version {id: $id})-[:PREVIOUS*1..5]->(history) RETURN history`

---

## **5\. Pothole Forecast: The "Floating Milestone"**

**The Folly:** Creating a `Version` node without a `HAS_LATEST` pointer. **The Reality:** A version with no relationship to a Snippet is "orphaned" and invisible to the UI. **The Mitigation:** Our Go adapter will wrap the Cypher query in a transaction. If the relationship creation fails, the entire node creation is rolled back.


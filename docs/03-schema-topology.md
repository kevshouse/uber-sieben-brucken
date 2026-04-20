# **03-schema-topology.md: The Two Shores**

## **1\. Design Reasoning: Hybrid Persistence**

In *Über sieben Brücken musst du gehn*, we utilize a hybrid database strategy to play to the strengths of two different paradigms:

* **Relational (libSQL/SQLite):** Chosen for **identity and ownership**. Relational databases excel at enforcing strict constraints on fixed entities like Users and Snippet metadata.  
* **Graph (Neo4j):** Chosen for **temporal evolution**. The "Seven Bridges" logic requires traversing code history. Graphs excel at following relationships (`:PREVIOUS`) without the performance penalty of recursive SQL joins.

---

## **2\. libSQL Schema (The Identity Shore)**

This database acts as the "Gatekeeper." It stores who a user is and what snippets they have created.

### **Tables**

#### **`users`**

* `id`: UUID (Primary Key)  
* `username`: String (Unique)  
* `email`: String (Unique)  
* `password_hash`: String (Argon2/Bcrypt)  
* `created_at`: Timestamp

#### **`snippets`**

* `id`: UUID (Primary Key)  
* `owner_id`: UUID (Foreign Key \-\> users.id)  
* `title`: String  
* `language`: String  
* `created_at`: Timestamp

---

## **3\. Neo4j Schema (The Temporal Path)**

This database acts as the "History." It tracks the literal path of code changes through time.

### **Nodes (Labels)**

* **`:Snippet`**: A lightweight anchor node that represents the project's existence in the graph.  
* **`:Version`**: An immutable node containing the actual content.  
  * `id`: UUID  
  * `content`: Text (The code)  
  * `timestamp`: DateTime

### **Relationships**

* **`(:Snippet)-[:HAS_LATEST]->(:Version)`**: Points to the current "Self-Healed" state.  
* **`(:Version)-[:PREVIOUS]->(:Version)`**: A directed edge pointing to the chronologically prior state.  
* **`(:User)-[:AUTHORED]->(:Version)`**: Tracks attribution across the graph.

---

## **4\. The "Self-Healing" Logic (Pointer Swapping)**

Traditional systems update a row to "rollback." Our system performs a **Topology Shift**:

1. Identify the "broken" `Version` node currently pointed to by `[:HAS_LATEST]`.  
2. Follow the `[:PREVIOUS]` relationship to find the last known good `Version`.  
3. Delete the `[:HAS_LATEST]` relationship to the broken node.  
4. Create a new `[:HAS_LATEST]` relationship to the good node.

**Result:** The history remains intact (the "broken" node is still in the graph for audit), but the system's "Current" state is restored instantly.

---

## **5\. Pothole Forecast: Data Synchronization**

* **The Folly:** Relying on the Graph to store User passwords or the Relational DB to store Code history.  
* **The Reality:** This creates "Split-Brain" logic.  
* **The Rule:** If it's a **Fact** (User exists), put it in libSQL. If it's a **Relationship** (Version A came before Version B), put it in Neo4j.

---

## Implementation Notes

In our Go code, we will use **UUIDs** as the "Joints." Both libSQL and Neo4j will use the same UUID for a specific Snippet, allowing us to query both shores and join the data in the **Core logic layer**.

##🌉 The Two Shores: Hybrid Persistence Strategy

**To maintain a clean separation between Identity and Evolution,** the system is divided into two distinct data "shores":
1. The Identity Shore (libSQL / Turso)
   - Responsibility: Static Metadata and Ownership.
   - Entity: Snippet Anchor.
   - Attributes: Title, OwnerID, CreatedAt.
   - Why SQL? Fast, ACID-compliant metadata lookups and relational integrity for user ownership.
2. The Relationship Current (Neo4j)
   - Responsibility: Temporal Evolution and Lineage.
   - Entities: Version, Citation.
   - Logic: Bidirectional doubly-linked lists ($O(1)$ history traversal) and the "Relationship as an Entity" pattern for citations.
   - Why Graph? To navigate complex, non-linear paths of code evolution and citation lineage without expensive recursive joins.
---
🔍 Anatomy for Review: "Hybrid Persistence"
   "I chose a Hybrid Persistence model. I use libSQL as the 'Shore' for stable metadata, and Neo4j as the 'Current' for the fluid, versioned relationships. This prevents the Graph from becoming cluttered with static data and keeps the SQL database from struggling with deep relational traversals.
---
"🚧 Pothole Forecast: 
   - The **"Sync Gap"**

     **The Folly:** Thinking the two databases will always stay perfectly in sync automatically.
   - **The Reality:** If a Snippet is created in libSQL but the Neo4j connection fails, you end up with a "Ghost Shore"—an identity with no history.
   - **The Mitigation:** 
     In the Service Layer (which we will build next), we will use Go's error handling to ensure that if the second "bridge" fails, we report a failure to the user so they know the handshake was incomplete.
---
🧭 The Next Crossing: The Bootstrap
   - Now that we've updated the documentation, we are ready to add the **Schema Bootstrap (the CREATE TABLE logic)** to our libsql_adapter.go. This ensures the "Identity Shore" physically exists before we try to land any data on it.-  - We will use <?> placeholders are used to prevent SQL Injection. It ensures that the database treats the input as as (data) rather than "Command" (code). It’s the ultimate safety rail for any SQL bridge.

🏛️ Updating the Topological Map

   The ALN Definition:

   - "The Active Linkage Node (ALN) is a first-class graph entity that represents a versioned relationship between two Snippet Anchors. Unlike a static edge, an ALN maintains its own temporal history, allowing the context of a connection to evolve without losing its lineage."


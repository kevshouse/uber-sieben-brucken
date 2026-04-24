# Pothole Log

## Pothole \#1: The YAML Indentation Trap

* **The Folly:** Assuming YAML is as flexible as C syntax or that Tabs equal Spaces.  
* **The Reality:** Docker Compose will fail silently or with cryptic "additional properties not allowed" errors if indentation is off.  
* **The Mitigation:** Always use 2-space indentation and verify with `docker compose config` before running `up`.

## Pothole 2. Ne04j Adapter fails to update legacy nodes, despite tests passing
- **The plan: **
	1. Refactor the Adapter Query: Change the MATCH on the target to be more resilient.
	2. The "Ghost" Fix: We should use MERGE on the target node within the citation flow. Why? Because if a snippet exists in our SQL "Shore" but was never initialized in the "Graph Current" (a ghost node), the citation should be the moment it is "born" in the graph.
- **First Action point** Update neo4j_adapter.go to make this truly "Green" and fix the legacy data issue.
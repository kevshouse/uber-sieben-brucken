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

## **Day 9 Additions**

| Date | Pothole | Resolution |
| ----- | ----- | ----- |
| 2026-04-27 | **The Zombie Lockup** (Kernel D-State) | libSQL process hung on I/O, preventing container shutdown and eventually locking the Ubuntu OS into Read-Only mode. **Fix:** Hard reboot and fsck to clear filesystem journal. |
| 2026-04-27 | **The Duplicate Bridge** (Redeclared Errors) | Shattering files without deleting logic from the main adapter file caused naming collisions that crashed the Go Language Server (gopls). **Fix:** Empty the main adapter file of methods once moved to specialized files. |
| 2026-04-27 | **The Missing Entry Point Driver** | Adapter requested libsql driver, but main.go hadn't imported it, leading to "unknown driver" errors. **Fix:** Move \_ "github.com/tursodatabase/libsql-client-go/libsql" to the root tool. |
| 2026-04-27 | **Localhost DNS Lag** | Containers were Up, but localhost was resolving to IPv6 ::1 while Docker was on IPv4. **Fix:** Force 127.0.0.1 in connection URIs. |
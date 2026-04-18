# Pothole Log

## Pothole \#1: The YAML Indentation Trap

* **The Folly:** Assuming YAML is as flexible as C syntax or that Tabs equal Spaces.  
* **The Reality:** Docker Compose will fail silently or with cryptic "additional properties not allowed" errors if indentation is off.  
* **The Mitigation:** Always use 2-space indentation and verify with `docker compose config` before running `up`.


## Design Note
1. We are separating internal/core (Domain Logic) from internal /adapter (External tech). This is to prevent us from writing "Logic in Cypher" potholes.
2. **Folly prevention** In Go, the ```internal``` directory is a special keyword. Any code inside it cannot be imported by other projects. This enforces privacy - exactly like static C variables to prevent global pollution.
3. **Pothole Forecast: Naming** 
**The Folly:** Changing the folder name later but forgetting to update the go.mod file.
**The Reality:** Go is **very sensitive** to this. If we name our directory ```uber-sieben-brucken``` to something else on the HDD, the code will still work as the module line inside ```go.mod``` stays the same. The go.mod is the "identity", the folder is the just the "Location".

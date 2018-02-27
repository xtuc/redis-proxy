# Algorithm

- HTTP request: GET /key
  - Lookup(key)
    - hit
      - Return(value)
      - Update item in LRU index
    - Miss
      - fetch the key in redis
      - Return(value)
      - is IsCacheFull()
        - The oldest entry is delete
        - Persist(key, value)
      - is not IsCacheFull()
        - Persist(key, value)

## Definitions

- `Lookup(key)`: Performs a lookup in the cache for a given key
- `Return(value)`: Return the value in the HTTP response (and close it)
- `Persist(key, value)`: Adds a new entry to the cache
- `IsCacheFull()`: the cache is at capacity


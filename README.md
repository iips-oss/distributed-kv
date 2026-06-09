# pacific

This is a persistent key-value store written in Go.
I call is pacific 

It maintains an in-memory map for fast reads and writes, 
and a write-ahead log (WAL) on disk for durability. 
Every mutation is recorded to the WAL before being applied 
to the map. On startup, the WAL is replayed to restore 
the last known state.

## Operations
- PUT key value
- GET key
- DELETE key

## How it works
- PUT writes to the WAL then updates the map
- GET reads directly from the map
- DELETE writes a tombstone to the WAL then removes the key from the map
- On restart, the WAL is replayed from the last snapshot checkpoint
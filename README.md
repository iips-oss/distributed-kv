# Redis Clone

A Redis-like in-memory key-value store built in Go.

## How to Run

Terminal 1 — start the server:

    go run .

Terminal 2 — connect to it:

    nc localhost 8000

## Commands

| Command | Syntax                 | Usage                       |
|---------|------------------------|-----------------------------|
| SET     | SET <key> <value>      | Store a key-value pair      |
| GET     | GET <key>              | Retrieve value by key       |
| DEL     | DEL <key>              | Delete a key-value pair     |
| EXISTS  | EXISTS <key>           | Check if a key exists       |
| KEYS    | KEYS                   | List all keys               |
| FLUSH   | FLUSH                  | Delete all keys             |
| EXPIRE  | EXPIRE <key> <seconds> | Set expiry on a key         |
| TTL     | TTL <key>              | Check remaining expiry time |
| LOG     | LOG                    | Show WAL log (admin only)   |
| HELP    | HELP                   | Show all commands           |
| EXIT    | EXIT                   | Disconnect client           |

## Example Session

    SET name Aishwarya
    OK
    GET name
    Aishwarya
    EXISTS name
    true
    EXPIRE name 10
    OK
    TTL name
    8 seconds remaining
    DEL name
    OK
    
## How it Works

- Data is stored in-memory using a Go map (hashmap)
- sync.RWMutex ensures safe concurrent access
- Each client connection runs in its own goroutine
- WAL (Write Ahead Log) persists every SET, DEL and FLUSH to disk so data survives server restarts
- On startup the WAL file is replayed line by line to rebuild the map
- Expiry uses lazy expiration — when you set EXPIRE on a key, no timer or background job is started. The key stays in memory with a deadline attached. It is only removed when someone tries to GET it after the deadline has passed. This is the same approach Redis uses internally.
- LOG command is password protected — only admins can view the WAL log

## Author

ak-junior3339
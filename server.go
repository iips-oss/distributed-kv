package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func startServer(store *Store) {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	fmt.Println("Server listening on port 8000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn, store)
	}
}

func handleClient(conn net.Conn, store *Store) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}

		var response string

		switch parts[0] {
		case "SET":
			if len(parts) < 3 {
				response = "ERR usage: SET <key> <value>"
			} else {
				store.Set(parts[1], parts[2])
				response = "OK"
			}
		case "GET":
			if len(parts) < 2 {
				response = "ERR usage: GET <key>"
			} else {
				val, ok := store.Get(parts[1])
				if ok {
					response = val
				} else {
					response = "key not found or key expired "
				}
			}
		case "DEL":
			if len(parts) < 2 {
				response = "ERR usage: DEL <key>"
			} else {
				store.Delete(parts[1])
				response = "OK"
			}
		case "EXISTS":
			if len(parts) < 2 {
				response = "ERR usage: EXISTS <key>"
			} else {
				if store.Exists(parts[1]) {
					response = "true"
				} else {
					response = "false"
				}
			}
		case "KEYS":
			keys := store.Keys()
			if len(keys) == 0 {
				response = "(empty)"
			} else {
				for i, k := range keys {
					response += fmt.Sprintf("%d) %s\n", i+1, k)
				}
			}
		case "FLUSH":
			store.Flush()
			response = "OK"
		case "EXPIRE":
			if len(parts) < 3 {
				response = "ERR usage: EXPIRE <key> <seconds>"
			} else {
				seconds, err := strconv.Atoi(parts[2])
				if err != nil {
					response = "ERR seconds must be a number"
				} else {
					store.Expire(parts[1], seconds)
					response = "OK"
				}
			}
		case "TTL":
			if len(parts) < 2 {
				response = "ERR usage: TTL <key>"
			} else {
				remaining := store.TTL(parts[1])
				if remaining == -1 {
					response = "no expiry set or key expired"
				} else {
					response = fmt.Sprintf("%d seconds remaining", remaining)
				}
			}
		case "LOG":
			// ask for password
			fmt.Fprintln(conn, "Enter password:")

			// read the password from client
			if !scanner.Scan() {
				response = "ERR could not read password"
				continue
			}

			password := scanner.Text()

			if password != "admin123" {
				response = "ERR wrong password access denied"
			} else {
				log, err := readWAL()
				if err != nil {
					response = "ERR could not read WAL file"
				} else {
					response = log
				}
			}
		case "HELP":
			response = "SET    SET <key> <value>    Store a key-value pair\n" +
				"GET    GET <key>            Retrieve value by key\n" +
				"DEL    DEL <key>            Delete a key-value pair\n" +
				"EXISTS EXISTS <key>         Check if a key exists\n" +
				"KEYS   KEYS                 List all keys\n" +
				"FLUSH  FLUSH                Delete all keys\n" +
				"EXPIRE EXPIRE <key> <sec>   Set expiry on a key\n" +
				"TTL    TTL <key>            Check remaining expiry time\n" +
				"EXIT   EXIT                 Quit the program"
		case "EXIT":
			fmt.Fprintln(conn, "Bye!")
			return
		default:
			response = "ERR unknown command"
		}

		fmt.Fprintln(conn, response)
	}
}

package main

import (
	"fmt"
)

func main() {
	fmt.Println("Redis Clone Starting...")
	store := NewStore()

	// replay WAL to restore previous state
	if err := replayWAL(store); err != nil {
		fmt.Println("Error replaying WAL:", err)
	} else {
		fmt.Println("WAL replayed successfully!")
	}

	startServer(store)
}

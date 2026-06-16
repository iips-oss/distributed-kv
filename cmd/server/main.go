package main

import (
	"fmt"
	"github.com/ShantanuPaliwal2419/distributed-kv/internal/store"
	"log"
)

func main() {
	wal, err := store.NewWAL("wal.log")
	if err != nil {
		log.Fatal(err)
	}

	s := store.NewStore(wal)

	//  recover old data first
	err = s.Replay()
	if err != nil {
		log.Fatal(err)
	}

	// test operations
	_ = s.Put("name", "shantanu")
	_ = s.Put("city", "mumbai")

	val, ok := s.Get("name")
	fmt.Println(val, ok)

	_ = s.Delete("city")
}

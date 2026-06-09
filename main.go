// COPILOT GENERATE FILE
package main
import (
    "fmt"
    "log"
    "github.com/Rudraksha-007/pacific/store"
)

func main() {
    // Initialize store with WAL file
    s, err := store.NewStore("test.wal")
    if err != nil {
        log.Fatalf("Failed to create store: %v", err)
    }

    fmt.Println("=== Store Test Suite ===\n")

    // Test 1: Put and Get operations
    fmt.Println("Test 1: Put and Get")
    if err := s.Put("name", "rudra"); err != nil {
        log.Fatalf("Put failed: %v", err)
    }
    if err := s.Put("city", "indore"); err != nil {
        log.Fatalf("Put failed: %v", err)
    }
    if err := s.Put("age", "25"); err != nil {
        log.Fatalf("Put failed: %v", err)
    }

    val, ok := s.Get("name")
    if ok {
        fmt.Printf("  ✓ GET name: %s\n", val)
    } else {
        log.Fatal("  ✗ Failed to get name")
    }

    val, ok = s.Get("city")
    if ok {
        fmt.Printf("  ✓ GET city: %s\n", val)
    } else {
        log.Fatal("  ✗ Failed to get city")
    }

    // Test 2: Delete operation
    fmt.Println("\nTest 2: Delete")
    if err := s.Delete("city"); err != nil {
        log.Fatalf("Delete failed: %v", err)
    }

    _, ok = s.Get("city")
    if !ok {
        fmt.Println("  ✓ city was deleted correctly")
    } else {
        log.Fatal("  ✗ Failed to delete city")
    }

    // Test 3: Verify other keys still exist
    fmt.Println("\nTest 3: Verify remaining keys")
    val, ok = s.Get("name")
    if ok && val == "rudra" {
        fmt.Printf("  ✓ name still exists: %s\n", val)
    } else {
        log.Fatal("  ✗ name key is missing or corrupted")
    }

    val, ok = s.Get("age")
    if ok && val == "25" {
        fmt.Printf("  ✓ age still exists: %s\n", val)
    } else {
        log.Fatal("  ✗ age key is missing or corrupted")
    }

    // Test 4: Get non-existent key
    fmt.Println("\nTest 4: Non-existent key")
    _, ok = s.Get("nonexistent")
    if !ok {
        fmt.Println("  ✓ Non-existent key correctly returns false")
    } else {
        log.Fatal("  ✗ Should not find non-existent key")
    }

    fmt.Println("\n=== All tests passed! ===")
}
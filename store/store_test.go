package store

import (
	"testing"
)

func TestSet(t *testing.T) {
	store := NewStore()
	store.Set("key1", "value1")

	value, ok := store.Get("key1")
	if !ok {
		t.Errorf("Expected: key1 to exist")
	}
	if value != "value1" {
		t.Errorf("Expected: value1, got: %s", value)
	}
}

func TestGet(t *testing.T) {
	store := NewStore()

	value, ok := store.Get("nonexistent")
	if ok {
		t.Errorf("Expected: nonexistent key to not exist")
	}
	if value != "" {
		t.Errorf("Expected: empty string, got: %s", value)
	}

	store.Set("key2", "value2") 
	value, ok = store.Get("key2")
	if !ok {
		t.Errorf("Expected: key2 to exist")
	}
	if value != "value2" {
		t.Errorf("Expected: value2, got: %s", value)
	}
}	

func TestDelete(t *testing.T) {
	store := NewStore()

	store.Set("key3", "value3")
	deleted := store.Delete("key3")
	if !deleted {
		t.Errorf("Expected: key3 to be deleted")
	}
	// We are not wrapping the below code in an if deleted block because 
	// we want to test the behavior of Get after Delete regardless of whether
	//  Delete returns true or false. This ensures that even if Delete fails to delete the key,
	//  we can verify that Get does not return the deleted value.
	value, ok := store.Get("key3")
	if ok {
		t.Errorf("Expected: key3 to be deleted")
	}
	if value != "" {
		t.Errorf("Expected: empty string, got: %s", value)
	}


	deleted = store.Delete("nonexistent")
	if deleted {
		t.Errorf("Expected: nonexistent key to not be deleted")
	}
}
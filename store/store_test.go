package store

import (
	"testing"
)

func TestSet(t *testing.T) {
	store := NewStore()
	revision := store.Revision()

	rev := store.Set("key1", "value1")
	value, ok := store.Get("key1")
	if (rev != revision+1) {
		t.Errorf("Expected: revision to have incremented")
	}
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
	revision := store.Revision()
	deleted, rev := store.Delete("key3")
	if (rev != revision+1) {
		t.Errorf("Expected: Revision to have incremented")
	}
	if !deleted {
		t.Errorf("Expected: key3 to be deleted")
	}
	// We are not wrapping the below code in an if deleted block because 
	// we want to test the behavior of Get after Delete regardless of whether, Delete returns true or false.
	value, ok := store.Get("key3")
	if ok {
		t.Errorf("Expected: key3 to be deleted")
	}
	if value != "" {
		t.Errorf("Expected: empty string, got: %s", value)
	}


	// _ is a blank identifier, used to ignore the return values that aren't needed.
	deleted, _ = store.Delete("nonexistent")
	if deleted {
		t.Errorf("Expected: nonexistent key to not be deleted")
	}
}
	

package kvstore // Use the same package name as your hashtable.go

import (
	"fmt"     // Needed for collision test keys
	"reflect" // Needed for comparing slices/maps robustly
	"testing"
)

// Helper function to check Get results accurately
func assertGetValue(t *testing.T, ht *HashTable, key string, expectedValue interface{}, expectedFound bool) {
	// t.Helper() marks this function as a test helper.
	// When t.Errorf is called inside a helper, the line number reported
	// is from the calling function, making debugging easier.
	t.Helper()

	value, found := ht.Get(key)

	if found != expectedFound {
		t.Errorf("Get(%q): found = %v; want %v", key, found, expectedFound)
		// No point comparing value if found status is wrong
		return
	}

	// Only compare values if we expected to find the key
	if expectedFound {
		// Use reflect.DeepEqual for robust comparison, especially for slices, maps, etc.
		if !reflect.DeepEqual(value, expectedValue) {
			t.Errorf("Get(%q): value = %v (%T); want %v (%T)", key, value, value, expectedValue, expectedValue)
		}
	}
}

// Test Basic Insert, Get, Update, Delete operations
func TestHashTable_BasicOperations(t *testing.T) {
	ht := NewHashTable() // Assumes initialCapacity is set within NewHashTable

	// Use t.Run to structure tests into sub-tests
	t.Run("Initial State", func(t *testing.T) {
		if ht.size != 0 {
			t.Errorf("Initial size = %d; want 0", ht.size)
		}
		if ht.capacity != initialCapacity { // Check initial capacity matches constant
			t.Errorf("Initial capacity = %d; want %d", ht.capacity, initialCapacity)
		}
		// Try getting a key from an empty table
		assertGetValue(t, ht, "anyKey", nil, false)
	})

	t.Run("Insert and Get Single Item", func(t *testing.T) {
		ht.Insert("name", "Alice")
		if ht.size != 1 {
			t.Errorf("Size after first insert = %d; want 1", ht.size)
		}
		assertGetValue(t, ht, "name", "Alice", true)
		// Check non-existent key again
		assertGetValue(t, ht, "age", nil, false)
	})

	t.Run("Insert Multiple Items", func(t *testing.T) {
		ht.Insert("age", 30)
		ht.Insert("city", "New York")
		byteSlice := []byte{1, 2, 3}
		ht.Insert("data", byteSlice)

		if ht.size != 4 {
			t.Errorf("Size after multiple inserts = %d; want 4", ht.size)
		}
		assertGetValue(t, ht, "name", "Alice", true)
		assertGetValue(t, ht, "age", 30, true)
		assertGetValue(t, ht, "city", "New York", true)
		assertGetValue(t, ht, "data", byteSlice, true) // Uses reflect.DeepEqual in helper
	})

	t.Run("Update Existing Item", func(t *testing.T) {
		initialSize := ht.size
		ht.Insert("age", 31) // Update age

		if ht.size != initialSize { // Size should NOT change on update
			t.Errorf("Size after update = %d; want %d", ht.size, initialSize)
		}
		assertGetValue(t, ht, "age", 31, true) // Check updated value
		// Ensure other values weren't affected
		assertGetValue(t, ht, "name", "Alice", true)
		assertGetValue(t, ht, "city", "New York", true)
	})

	t.Run("Delete Existing Item", func(t *testing.T) {
		ht.Delete("city")
		if ht.size != 3 {
			t.Errorf("Size after delete = %d; want 3", ht.size)
		}
		assertGetValue(t, ht, "city", nil, false) // Verify deletion
		// Ensure other values remain
		assertGetValue(t, ht, "name", "Alice", true)
		assertGetValue(t, ht, "age", 31, true)
	})

	t.Run("Delete Non-Existent Item", func(t *testing.T) {
		initialSize := ht.size
		ht.Delete("country") // Key doesn't exist

		if ht.size != initialSize { // Size should NOT change
			t.Errorf("Size after deleting non-existent key = %d; want %d", ht.size, initialSize)
		}
		assertGetValue(t, ht, "country", nil, false)
	})

	t.Run("Delete All Remaining Items", func(t *testing.T) {
		ht.Delete("name")
		ht.Delete("age")
		ht.Delete("data")

		if ht.size != 0 {
			t.Errorf("Size after deleting all = %d; want 0", ht.size)
		}
		assertGetValue(t, ht, "name", nil, false)
		assertGetValue(t, ht, "age", nil, false)
		assertGetValue(t, ht, "data", nil, false)
	})
}

// Test Collision Handling (by inserting many items)
// This test relies on the hash function distributing keys into buckets,
// inevitably causing collisions within the initial capacity.
func TestHashTable_Collisions(t *testing.T) {
	ht := NewHashTable()
	numItems := initialCapacity * 3 // Insert enough items to likely cause many collisions

	keys := make([]string, numItems)
	values := make([]int, numItems)

	// Insert items
	for i := 0; i < numItems; i++ {
		keys[i] = fmt.Sprintf("collision_key_%d", i)
		values[i] = i * 10 // Use distinct values
		ht.Insert(keys[i], values[i])
	}

	// Check size after mass insert
	if ht.size != numItems {
		// Use Fatalf because subsequent checks depend on correct insertion
		t.Fatalf("After inserting %d items, size = %d; want %d", numItems, ht.size, numItems)
	}

	// Verify all inserted items can be retrieved correctly
	t.Run("Retrieve All Amidst Collisions", func(t *testing.T) {
		for i := 0; i < numItems; i++ {
			assertGetValue(t, ht, keys[i], values[i], true)
		}
	})

	// Delete roughly half the items and verify others remain
	t.Run("Delete Some Amidst Collisions", func(t *testing.T) {
		deletedCount := 0
		for i := 0; i < numItems; i += 2 { // Delete items with even indices
			ht.Delete(keys[i])
			deletedCount++
		}

		expectedSize := numItems - deletedCount
		if ht.size != expectedSize {
			t.Errorf("After deleting %d items, size = %d; want %d", deletedCount, ht.size, expectedSize)
		}

		// Check remaining items and ensure deleted items are gone
		for i := 0; i < numItems; i++ {
			if i%2 == 0 { // These should have been deleted
				assertGetValue(t, ht, keys[i], nil, false)
			} else { // These should still exist
				assertGetValue(t, ht, keys[i], values[i], true)
			}
		}
	})
}

// Test storing and retrieving nil as a value
func TestHashTable_NilValue(t *testing.T) {
	ht := NewHashTable()

	t.Run("Insert Nil Value", func(t *testing.T) {
		ht.Insert("myNilKey", nil)
		if ht.size != 1 {
			t.Errorf("Size after inserting nil = %d; want 1", ht.size)
		}
		assertGetValue(t, ht, "myNilKey", nil, true) // Should be found, value is nil
	})

	t.Run("Update To Nil Value", func(t *testing.T) {
		ht.Insert("anotherKey", "someValue")
		ht.Insert("anotherKey", nil) // Update to nil
		if ht.size != 2 {            // Should have "myNilKey" and "anotherKey"
			t.Errorf("Size after updating to nil = %d; want 2", ht.size)
		}
		assertGetValue(t, ht, "anotherKey", nil, true)
	})

	t.Run("Delete Key With Nil Value", func(t *testing.T) {
		ht.Delete("myNilKey")
		if ht.size != 1 { // Only "anotherKey" should remain
			t.Errorf("Size after deleting key with nil value = %d; want 1", ht.size)
		}
		assertGetValue(t, ht, "myNilKey", nil, false)  // Should not be found
		assertGetValue(t, ht, "anotherKey", nil, true) // Other key should still exist
	})
}

// Note: If the resize functionality were enabled, additional tests would be needed:
// 1. Test that resizing triggers at the correct load factor.
// 2. Test that all elements are correctly rehashed and retrievable after resizing.
// 3. Test insert/get/delete operations immediately after a resize occurs.

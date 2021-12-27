package main

import (
	"time"

	"github.com/Squwid/cache/lru"
)

type DatabaseObject struct {
	ID          string
	Name        string
	Category    string
	CreatedTime time.Time
}

// Struct MUST have a Key() method to avoid collisions in the hashmap. The hashmap lets
// access of the objects in the linked list O(1) time instead of O(N).
func (dbo DatabaseObject) Key() string { return dbo.ID }

func main() {
	// Make an LRU cache of type `DatabaseObject` with a size of 3.
	// The cache is thread safe, so its a great use for a web server
	var cache = lru.NewCache[DatabaseObject](3)

	// Get some object from the database and cache it
	obj := DatabaseObject{ID: "abc123", Name: "Database Object", Category: "Color", CreatedTime: time.Now().UTC()}
	cache.Add(obj) // Add the item to the cache

	// Get the item from the cache
	cacheObj := cache.Get("abc123")
	if cacheObj == nil {
		// Cache object was not found, get from database and store in cache
	}
}

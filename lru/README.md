# Least Recently Used Cache (LRU)

NOTE: This is a generic cache and is only compatible with Go 1.18 and newer

## About

The LRU cache is the most widespread caching method. The cache is backed by a doubly linked list for each item, with a hashmap for quick access by using a key. When storing items, items with the same key will be overwritten and moved to the front of the cache. When the cache is full and an item is inserted, the least recently used item will be evicted, and the new item will be placed at the start of the linked list.

The weakness of this implementation of the LRU is that each element must have a unique identifier. The TTL also falls on the implementor for this cache, as there is no built in item timeout.

| Action | Worst Case |
| ------ | ---------- |
| Space | `O(N)` |
| Add Item | `O(1)` |
| Access Item | `O(1)` |
| Remove Item | `O(1)` |

## Installation

`go get -u github.com/Squwid/cache/lru`

## How to Use

The cache is thread-safe, meaning there is no need for any mutexes, since the package handles it.

```go
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

```
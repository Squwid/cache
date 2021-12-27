# Go Generic Caches

This is a Go 1.18 package with multiple generic caches. 

Note: This will not work on any version older than Go 1.18

## Caches

The goal of this project is to offer a variety of dependency free, thread-safe, generic caches for different needs. Below is the list of current caches that are available.

### Least Recently Used (LRU)

The LRU cache is the most widespread caching method. The cache is backed by a doubly linked list for each item, with a hashmap for quick access by using a key. When storing items, items with the same key will be overwritten and moved to the front of the cache. When the cache is full and an item is inserted, the least recently used item will be evicted, and the new item will be placed at the start of the linked list.

The weakness of this implementation of the LRU is that each element must have a unique identifier. The TTL also falls on the implementor for this cache, as there is no built in item timeout.

| Action | Worst Case |
| ------ | ---------- |
| Space | `O(N)` |
| Add Item | `O(1)` |
| Access Item | `O(1)` |
| Remove Item | `O(1)` |



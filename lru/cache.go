package lru

import "sync"

type CacheConstraint interface{ Key() string }

type Cache[T CacheConstraint] struct {
	m map[string]*Node[T]

	maxSize int
	size    int

	head *Node[T]
	tail *Node[T]

	lock *sync.Mutex
}

type Node[T CacheConstraint] struct {
	Prev *Node[T]
	Next *Node[T]

	Item T
}

func NewCache[T CacheConstraint](size int) *Cache[T] {
	return &Cache[T]{
		m:       make(map[string]*Node[T]),
		maxSize: size,
		head:    nil,
		tail:    nil,
		lock:    &sync.Mutex{},
	}
}

// Size returns the current size of the cache
func (c *Cache[T]) Size() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.cacheSize()
}

func (c *Cache[T]) cacheSize() int { return c.size }

func (c *Cache[T]) Add(item T) {
	key := item.Key()

	c.lock.Lock()
	defer c.lock.Unlock()

	n, ok := c.m[key]
	if ok {
		n.Item = item
		c.moveNodeToFront(n)
	} else {
		if c.cacheSize() == c.maxSize {
			c.evictTail()
		}
		n = &Node[T]{Item: item}
		c.m[key] = n
		c.insertToFront(n)
	}
}

// Remove removes an item from the cache. It returns the object if it was removed,
// otherwise nil will be returned
func (c *Cache[T]) Remove(key string) *T {
	c.lock.Lock()
	defer c.lock.Unlock()

	n, ok := c.m[key]
	if !ok {
		return nil
	}
	c.removeNode(n)
	c.removeFromMap(key)

	return &n.Item
}

func (c *Cache[T]) Get(key string) *T {
	c.lock.Lock()
	defer c.lock.Unlock()

	if node, ok := c.m[key]; ok {
		c.moveNodeToFront(node) // Move most recent GET to front
		return &node.Item
	}
	return nil
}

// PeekHead will return the item at the head without moving its position
// in the cache
func (c *Cache[T]) PeekHead() *T {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.head == nil {
		return nil
	}
	return &c.head.Item
}

// PeekTail will return the item at the tail without moving its position
// in the cache.
func (c *Cache[T]) PeekTail() *T {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.tail == nil {
		return nil
	}
	return &c.tail.Item
}

func (c *Cache[T]) evictTail() {
	if c.tail == nil {
		return
	}

	n := c.tail
	c.removeNode(n)
	c.removeFromMap(n.Item.Key())
}

// removeNode removes a node from the cache. removeNode does NOT remove
// the item from the map
func (c *Cache[T]) removeNode(node *Node[T]) {
	if node == nil {
		return
	}

	if c.cacheSize() == 1 {
		c.head = nil
		c.tail = nil
	} else if c.cacheSize() == 2 {
		if node == c.head {
			c.head = c.tail
			c.head.Prev = nil
		} else {
			c.tail = c.head
			c.tail.Next = nil
		}
	} else {
		if node == c.head {
			c.head = node.Next
			c.head.Prev = nil
		} else if node == c.tail {
			c.tail = node.Prev
			c.tail.Next = nil
		} else {
			node.Prev.Next = node.Next
			node.Next.Prev = node.Prev
		}
	}

	node.Prev = nil
	node.Next = nil
	c.size--
}

func (c *Cache[T]) removeFromMap(id string) { delete(c.m, id) }

// insertToFront will insert a node to the front of the cache. NOTE: this does not
// actually add the node to the cache via map. That will be done upstream. This function
// expects the node to already to be added to c.m
func (c *Cache[T]) insertToFront(node *Node[T]) {
	if c.cacheSize() == 0 {
		c.head = node
		c.tail = node
	} else {
		node.Next = c.head
		c.head.Prev = node
		c.head = node
	}

	c.size++
}

// moveNodeToFront moves a node that already exists in the cache to the front
func (c *Cache[T]) moveNodeToFront(node *Node[T]) {
	if node == nil || node.Prev == nil {
		return
	}

	c.removeNode(node)
	c.insertToFront(node)
}

// ForEach will run the given function for each item in the cache in order of the cache
func (c *Cache[T]) ForEach(f func(item T, index int)) {
	c.lock.Lock()
	defer c.lock.Unlock()

	node := c.head
	var i = 0
	for node != nil {
		f(node.Item, i)
		node = node.Next
		i++
	}
}

package lru

type CacheConstraint interface {
	ID() string
}

type Cache[T CacheConstraint] struct {
	m map[string]*Node[T]

	MaxSize int
	size    int

	Head *Node[T]
	Tail *Node[T]
}

type Node[T CacheConstraint] struct {
	Prev *Node[T]
	Next *Node[T]

	Item T
}

func NewCache[T CacheConstraint](size int) *Cache[T] {
	return &Cache[T]{
		m:       make(map[string]*Node[T]),
		MaxSize: size,
		Head:    nil,
		Tail:    nil,
	}
}

// Size returns the current size of the cache
func (c *Cache[T]) Size() int { return c.size }

func (c *Cache[T]) Add(item T) {
	id := item.ID()
	// If key already exists in cache, overwrite value and move to front
	n, ok := c.m[id]
	if ok {
		n.Item = item
		c.moveNodeToFront(n)
	} else {
		if c.Size() == c.MaxSize {
			c.evictTail()
		}
		n = &Node[T]{Item: item}
		c.m[id] = n
		c.insertToFront(n)
	}
}

// Remove removes an item from the cache. It returns the object if it was removed,
// otherwise nil will be returned
func (c *Cache[T]) Remove(id string) *T {
	n, ok := c.m[id]
	if !ok {
		return nil
	}
	c.removeNode(n)
	c.removeFromMap(id)

	return &n.Item
}

func (c *Cache[T]) Get(key string) *T {
	if node, ok := c.m[key]; ok {
		c.moveNodeToFront(node) // Move most recent GET to front
		return &node.Item
	}
	return nil
}

// PeekHead will return the item at the head without moving its position
// in the cache
func (c *Cache[T]) PeekHead() *T {
	if c.Head == nil {
		return nil
	}
	return &c.Head.Item
}

// PeekTail will return the item at the tail without moving its position
// in the cache.
func (c *Cache[T]) PeekTail() *T {
	if c.Tail == nil {
		return nil
	}
	return &c.Tail.Item
}

func (c *Cache[T]) evictTail() {
	if c.Tail == nil {
		return
	}

	n := c.Tail
	c.removeNode(n)
	c.removeFromMap(n.Item.ID())
}

// removeNode removes a node from the cache. removeNode does NOT remove
// the item from the map
func (c *Cache[T]) removeNode(node *Node[T]) {
	if node == nil {
		return
	}

	if c.Size() == 1 {
		c.Head = nil
		c.Tail = nil
	} else if c.Size() == 2 {
		if node == c.Head {
			c.Head = c.Tail
			c.Head.Prev = nil
		} else {
			c.Tail = c.Head
			c.Tail.Next = nil
		}
	} else {
		if node == c.Head {
			c.Head = node.Next
			c.Head.Prev = nil
		} else if node == c.Tail {
			c.Tail = node.Prev
			c.Tail.Next = nil
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
	if c.Size() == 0 {
		c.Head = node
		c.Tail = node
	} else {
		node.Next = c.Head
		c.Head.Prev = node
		c.Head = node
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
	node := c.Head
	var i = 0
	for node != nil {
		f(node.Item, i)
		node = node.Next
		i++
	}
}

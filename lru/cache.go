package lru

type CacheConstraint interface {
	ID() string
}

type Cache[T CacheConstraint] struct {
	m map[string]*Node[T]

	MaxSize int

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

func (c *Cache[T]) Add(item T) {
	// If key already exists in cache, overwrite value and move to front
	n, ok := c.m[item.ID()]
	if ok {
		n.Item = item
		c.moveNodeToFront(n)
	} else {
		if len(c.m) == c.MaxSize {
			c.evictTail()
		}
		n = &Node[T]{Item: item}
		c.m[item.ID()] = n
		c.insertToFront(n)
	}
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
	newtail := c.Tail.Prev

	newtail.Prev.Next = nil
	c.Tail = newtail
}

// insertToFront will insert a node to the front of the cache. NOTE: this does not
// actually add the node to the cache via map. That will be done upstream. This function
// expects the node to already to be added to c.m
func (c *Cache[T]) insertToFront(node *Node[T]) {
	if len(c.m) <= 1 {
		c.Head = node
		c.Tail = node
	} else {
		node.Next = c.Head
		c.Head.Prev = node
		c.Head = node
	}

}

// moveNodeToFront moves a node that already exists in the cache to the front
func (c *Cache[T]) moveNodeToFront(node *Node[T]) {
	// If previous does not exist then node must be in the front, or if the cache has only one element
	if node.Prev == nil || len(c.m) <= 1 {
		return
	}

	// Modify previous and next nodes
	prev := node.Prev
	next := node.Next
	prev.Next = next
	if c.Tail == node {
		c.Tail = prev
	} else {
		next.Prev = prev
	}

	node.Next = c.Head
	c.Head.Prev = node
	c.Head = node
}

// ForEach will run the given function for each item in the cache in order of the cache
func (c *Cache[T]) ForEach(f func(item T)) {
	node := c.Head
	for node != nil {
		f(node.Item)
		node = node.Next
	}
}

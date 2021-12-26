package lru

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestLRU struct {
	Value string
}

func (l TestLRU) ID() string {
	return l.Value
}

func TestCache(t *testing.T) {
	cache := NewCache[TestLRU](4)

	assert.Nil(t, cache.PeekHead())
	assert.Nil(t, cache.PeekTail())

	item1 := TestLRU{Value: "1"}
	item2 := TestLRU{Value: "2"}

	cache.Add(item1)
	cache.Add(item2)

	assert.Equal(t, cache.PeekHead().Value, item2.Value)
	assert.Equal(t, cache.PeekTail().Value, item1.Value)
	cache.ForEach(func(item TestLRU) {
		fmt.Printf("%v ", item.Value)
	})
	fmt.Println("")

	assert.Equal(t, cache.Get(item1.Value).Value, item1.Value) // Moves item1 to front
	assert.Equal(t, cache.PeekHead().Value, item1.Value)
	assert.Equal(t, cache.PeekTail().Value, item2.Value)
}

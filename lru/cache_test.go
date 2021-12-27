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

func TestSmallCache(t *testing.T) {
	cache := NewCache[TestLRU](4)

	assert.Nil(t, cache.PeekHead())
	assert.Nil(t, cache.PeekTail())

	item1 := TestLRU{Value: "1"}
	item2 := TestLRU{Value: "2"}

	cache.ForEach(func(item TestLRU, index int) { fmt.Printf("%v ", item.Value) })
	fmt.Println("")
	assert.Equal(t, cache.Size(), 0)
	assert.Equal(t, len(cache.m), 0)
	cache.Add(item1)
	assert.Equal(t, cache.Size(), 1)
	assert.Equal(t, len(cache.m), 1)
	cache.Add(item2)
	assert.Equal(t, cache.Size(), 2)
	assert.Equal(t, len(cache.m), 2)

	assert.Equal(t, cache.PeekHead().Value, item2.Value)
	assert.Equal(t, cache.PeekTail().Value, item1.Value)
	cache.ForEach(func(item TestLRU, index int) { fmt.Printf("%v ", item.Value) })
	fmt.Println("")

	cache.ForEach(func(item TestLRU, index int) { fmt.Printf("%v ", item.Value) })
	fmt.Println("")
	assert.Equal(t, cache.Get(item1.Value).Value, item1.Value) // Moves item1 to front
	assert.Equal(t, cache.PeekHead().Value, item1.Value)
	assert.Equal(t, cache.PeekTail().Value, item2.Value)
	assert.Equal(t, len(cache.m), 2)

	cache.ForEach(func(item TestLRU, index int) { fmt.Printf("%v ", item.Value) })
	fmt.Println("")
	item3 := TestLRU{Value: "3"}
	cache.Add(item3)
	assert.Equal(t, cache.Size(), 3)
	assert.Equal(t, len(cache.m), 3)
	assert.Equal(t, cache.PeekHead().Value, item3.Value)
	assert.Nil(t, cache.Get("abc"))

	cache.ForEach(func(item TestLRU, index int) { fmt.Printf("%v ", item.Value) })
	fmt.Println("")
	assert.Equal(t, cache.Get(item2.Value).Value, item2.Value)
	assert.Equal(t, cache.PeekHead().Value, item2.Value)
	assert.Equal(t, cache.Size(), 3)
	assert.Equal(t, len(cache.m), 3)
	assert.Equal(t, cache.PeekTail().Value, item1.Value)

	cache.ForEach(func(item TestLRU, index int) { fmt.Printf("%v ", item.Value) })
	fmt.Println("")
	assert.Nil(t, cache.Remove("abc"))
	assert.Equal(t, cache.Remove(item1.Value).Value, item1.Value)
	assert.Equal(t, cache.PeekHead().Value, item2.Value)
	assert.Equal(t, cache.PeekTail().Value, item3.Value)
	assert.Equal(t, cache.Size(), 2)
	assert.Equal(t, len(cache.m), 2)

	cache.ForEach(func(item TestLRU, index int) { fmt.Printf("%v ", item.Value) })
	fmt.Println("")
	assert.Equal(t, cache.Remove(item2.Value).Value, item2.Value)
	assert.Equal(t, cache.PeekHead().Value, item3.Value)
	assert.Equal(t, cache.PeekTail().Value, item3.Value)
	assert.Equal(t, cache.Size(), 1)
	assert.Equal(t, len(cache.m), 1)

	cache.ForEach(func(item TestLRU, index int) { fmt.Printf("%v ", item.Value) })
	fmt.Println("")
	assert.Equal(t, cache.Remove(item3.Value).Value, item3.Value)
	assert.Nil(t, cache.PeekHead())
	assert.Nil(t, cache.PeekTail())
	assert.Equal(t, cache.Size(), 0)
	assert.Equal(t, len(cache.m), 0)
}

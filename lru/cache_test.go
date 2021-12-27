package lru

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestLRU struct {
	Value string
}

func (l TestLRU) ID() string { return l.Value }

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

func TestLargerCache(t *testing.T) {
	const MAX = 10

	cache := NewCache[TestLRU](MAX)
	for i := 0; i < 100; i++ {
		lru := TestLRU{Value: fmt.Sprintf("%v", i)}
		cache.Add(lru)

		assert.LessOrEqual(t, cache.Size(), MAX)
		assert.LessOrEqual(t, len(cache.m), MAX)
		assert.Equal(t, cache.Size(), len(cache.m))
	}

	cache.ForEach(func(item TestLRU, index int) { fmt.Printf("%v ", item.Value) })
	fmt.Println("")
	assert.Equal(t, "99", cache.PeekHead().Value)
	assert.Equal(t, "90", cache.PeekTail().Value)
	assert.Equal(t, MAX, cache.Size())
	assert.Equal(t, MAX, len(cache.m))

	assert.Equal(t, "95", cache.Get("95").Value)
	assert.Equal(t, "90", cache.Get("90").Value)
	cache.ForEach(func(item TestLRU, index int) { fmt.Printf("%v ", item.Value) })
	fmt.Println("")
	assert.Equal(t, "90", cache.PeekHead().Value)
	assert.Equal(t, "91", cache.PeekTail().Value)
	assert.Equal(t, MAX, cache.Size())
	assert.Equal(t, MAX, len(cache.m))
}

func TestRaceConditions(t *testing.T) {
	var cache = NewCache[TestLRU](100)
	for i := 0; i < 5000; i++ {
		go func(c *Cache[TestLRU]) {
			for j := 0; j < 9999; j++ {
				time.Sleep(10 * time.Nanosecond)
				key := fmt.Sprintf("%v", time.Now().UTC().Unix())
				lru := TestLRU{Value: key}
				cache.Get(lru.Value)
				cache.Add(lru)
				cache.PeekHead()
				cache.PeekTail()
				cache.Size()

				time.Sleep(10 * time.Nanosecond)
				cache.Remove(lru.Value)
			}
		}(cache)
	}
}

package cache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xtuc/redis-proxy/cache"
)

const (
	capacity = 100
)

func TestLRUInsertItem(t *testing.T) {
	c := cache.NewLRUIndex(capacity)

	c.Update("a")

	assert.Equal(t, "a", c.String())
}

func TestLRUUpdateItem(t *testing.T) {
	c := cache.NewLRUIndex(capacity)

	c.Update("a")
	c.Update("b")

	assert.Equal(t, "b,a", c.String())
}

func TestLRUUpdateItemMultipleTimes(t *testing.T) {
	c := cache.NewLRUIndex(capacity)

	c.Update("a")
	c.Update("b")
	c.Update("a")

	assert.Equal(t, "a,b", c.String())
}

func TestLRUGetOldestEmpty(t *testing.T) {
	c := cache.NewLRUIndex(capacity)

	assert.Nil(t, c.GetOldest())
}

func TestLRUGetOldest(t *testing.T) {
	c := cache.NewLRUIndex(capacity)

	c.Update("a")
	c.Update("b")

	assert.Equal(t, c.GetOldest(), newString("a"))
}

func TestLRUUpdateAndGetOldest(t *testing.T) {
	c := cache.NewLRUIndex(capacity)

	c.Update("a")
	c.Update("b")
	c.Update("a")

	assert.Equal(t, newString("b"), c.GetOldest())
}

func TestLRURemove(t *testing.T) {
	c := cache.NewLRUIndex(capacity)

	c.Update("a")
	c.Remove("a")

	assert.Equal(t, "", c.String())
}

func TestIsAtCapacity(t *testing.T) {
	c := cache.NewLRUIndex(1)

	assert.False(t, c.IsAtCapacity())

	c.Update("a") // Add one element

	assert.True(t, c.IsAtCapacity())
}

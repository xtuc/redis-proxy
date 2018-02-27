package cache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xtuc/redis-proxy/cache"
)

func newString(str string) *string {
	return &str
}

func TestCacheNewCache(t *testing.T) {
	c := cache.NewCache()
	assert.Equal(t, c.Size(), 0)
}

func TestCacheSetAndGet(t *testing.T) {
	c := cache.NewCache()

	c.Set("foo", "bar")

	assert.Equal(t, newString("bar"), c.Get("foo"))
	assert.Equal(t, 1, c.Size())
}

func TestCacheGetUnexisting(t *testing.T) {
	c := cache.NewCache()
	assert.Nil(t, c.Get("foo"))
}

func TestCacheSetMultipleTimes(t *testing.T) {
	c := cache.NewCache()

	c.Set("foo", "a")
	c.Set("foo", "b")

	assert.Equal(t, newString("b"), c.Get("foo"))
}

func TestCacheRemove(t *testing.T) {
	c := cache.NewCache()
	c.Set("foo", "b")

	c.Remove("foo")

	assert.Nil(t, c.Get("foo"))
}

package cache

type CacheBackend map[string]string

type InMemoryCache struct {
	backend CacheBackend
}

func NewCache() *InMemoryCache {
	return &InMemoryCache{
		backend: make(CacheBackend),
	}
}

func (cache *InMemoryCache) Size() int {
	return len(cache.backend)
}

func (cache *InMemoryCache) Get(k string) *string {
	if value, ok := cache.backend[k]; ok {
		return &value
	}

	return nil
}

func (cache *InMemoryCache) Set(k string, v string) {
	cache.backend[k] = v
}

func (cache *InMemoryCache) Remove(k string) {
	delete(cache.backend, k)
}

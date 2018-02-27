package http

import (
	"errors"
	"io"
	"log"
	httppkg "net/http"
	"sync"

	"github.com/xtuc/redis-proxy/cache"
	"github.com/xtuc/redis-proxy/types"
)

const (
	CACHE_HEADER_KEY = "X-Cache"
)

type HTTPService struct {
	cache     *cache.InMemoryCache
	lruIndex  *cache.LRUIndex
	backend   types.BackendInterface
	ttlWorker cache.TTLWorker

	rwMutex sync.RWMutex
}

func NewHTTPService(
	cache *cache.InMemoryCache,
	lruIndex *cache.LRUIndex,
	backend types.BackendInterface,
	ttlWorker cache.TTLWorker,
	lock sync.RWMutex,
) *HTTPService {

	return &HTTPService{
		cache:     cache,
		lruIndex:  lruIndex,
		backend:   backend,
		ttlWorker: ttlWorker,

		rwMutex: lock,
	}
}

/*
  Cache helper methods
*/

// Add an item to the cache and LRUIndex
func (c HTTPService) SetInCache(key, value string) error {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	/*
	  Cache is at capacity.
	  Remove the oldest to be able to insert a new one
	*/
	if c.lruIndex.IsAtCapacity() {
		oldestItem := c.lruIndex.GetOldest()

		if oldestItem == nil {
			return errors.New("Can not remove oldest item from cache")
		}

		c.lruIndex.Remove(*oldestItem)
		c.cache.Remove(*oldestItem)
	}

	c.ttlWorker.Update(key)
	c.lruIndex.Update(key)
	c.cache.Set(key, value)

	return nil
}

// Get an item from the cache or from the backend (and caches it)
func (c HTTPService) GetFromCacheOrBackend(key string) (value string, wasCached bool, err error) {
	wasCached = false

	// Lock only for reading for now
	c.rwMutex.RLock()

	cachedValue := c.cache.Get(key)
	c.rwMutex.RUnlock()

	// Value is cached, return it and update the LRUIndex
	if cachedValue != nil {
		c.lruIndex.Update(key)
		c.ttlWorker.Update(key)

		wasCached = true
		return *cachedValue, wasCached, nil
	}

	// Value is not cache, fetch it from backend
	value, backendErr := c.backend.GetValue(key)

	if backendErr != nil {
		return "", wasCached, backendErr
	}

	// Set in cache for the next time
	cacheErr := c.SetInCache(key, value)

	if cacheErr != nil {
		return "", wasCached, cacheErr
	}

	// Return the value from the backend
	return value, wasCached, nil
}

/*
  Static helper functions
*/

func SendKeyNotFoundError(key string, w httppkg.ResponseWriter) {
	w.Header().Add(CACHE_HEADER_KEY, "miss")
	w.WriteHeader(httppkg.StatusNotFound)
	io.WriteString(w, "")

	log.Printf("404 - unexisting key %s", key)
}

func SendKeyOk(key, value string, w httppkg.ResponseWriter) {

	// Check if it's not already set
	if w.Header().Get(CACHE_HEADER_KEY) == "" {
		w.Header().Add(CACHE_HEADER_KEY, "hit")
	}

	w.WriteHeader(httppkg.StatusOK)
	io.WriteString(w, value)

	log.Printf("200 - existing key %s", key)
}

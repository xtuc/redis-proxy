package main

import (
	"flag"
	"log"
	"sync"

	"github.com/gorilla/mux"

	"github.com/xtuc/redis-proxy/cache"
	"github.com/xtuc/redis-proxy/http"
	"github.com/xtuc/redis-proxy/redis"
)

var (
	redisHost = flag.String("redis-host", "localhost", "Redis server host")

	httpHost = flag.String("http-host", "localhost", "HTTP server host")
	httpPort = flag.String("http-port", "8080", "HTTP server port")

	cacheCapacity = flag.Int("cache-capacity", 10^5, "Cache capacity")
	cacheTTL      = flag.Int("cache-ttl", 60*60, "Cache items time to life")
)

func main() {
	flag.Parse()

	cacheLock := sync.RWMutex{}

	cacheInstance := cache.NewCache()
	lruIndex := cache.NewLRUIndex(*cacheCapacity)
	redisClient, redisErr := redis.NewRedisClient(*redisHost)
	ttlWorker := cache.MakeTTLWorker(*cacheTTL)

	if redisErr != nil {
		panic(redisErr)
	}

	log.Printf("Connected to Redis at %s\n", *redisHost)

	// Start the TTL worker
	go ttlWorker.Run()

	// Handle cache invalidations
	go func() {
		for {
			item := <-ttlWorker.Events()

			cacheLock.Lock()

			log.Printf("Got invalidation for key %s\n", item.Name)

			lruIndex.Remove(item.Name)
			cacheInstance.Remove(item.Name)

			cacheLock.Unlock()
		}
	}()

	service := http.NewHTTPService(
		cacheInstance,
		lruIndex,
		redisClient,
		ttlWorker,

		cacheLock,
	)

	r := mux.NewRouter()

	r.
		HandleFunc("/{key}", service.PutHandler).
		Methods("PUT")

	r.
		HandleFunc("/{key}", service.LookupHandler).
		Methods("GET")

	server := http.NewHTTPServer(*httpHost, *httpPort, r)

	waitChan, err := server.Start()

	if err != nil {
		panic(err)
	}

	log.Printf("HTTP server listening at %s:%s\n", *httpHost, *httpPort)

	<-waitChan
}

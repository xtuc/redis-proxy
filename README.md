# redis-proxy

> transparent Redis proxy service.

The document seems to priviledge caching rather than sharding (or any other logic), the proxy focuses on access performance, thus it's eventually consistent.

## Usage

```sh
./cmd/redis-proxy/redis-proxy
```

### Testing stack

Alternatively you can use the local testing stack by running the following command:

```sh
docker-compose up
```

### Options

| name | descr |
|---|---|
| `cache-ttl` | Cache entry time to live (in second, defaults to 1h) |
| `cache-capacity` | Number of entries in the cache (defaults to 100k) |
| `redis-host` | Host of the Redis server (defaults to `localhost` ) |
| `http-host` | Host of the HTTP server (defaults to `localhost` ) |
| `http-port` | Port of the HTTP server (defaults to `8080` ) |
| `max-connections` | Number of connections allowed at the same time. Less than two means that requests will be processed sequentially (defaults to `1`) |
| `max-connection-queue` | Number of connections in the queue (defaults to 1k) |

### Get a key

```sh
curl -D - http://localhost:8080/KEY
```

### Put a key

```sh
curl -XPUT -D - http://localhost:8080/KEY -d 'VALUE'
```

## Cache

The cache is the most critical component of the proxy, we need to guarantee good read performance.

I first though about a Binary tree that could hold a wide range of keys and provide efficient lookups.

It seems that Golang's map implementation offers an average time complexity of O(1) for lookups.

Having a string type for the keys is not ideal and decreases the cache access performance, but conforms to Redis semantics.

No persistence can be configured, the cache is only in-memory.

### Key expiration

Each entry in the cache has a TTL (configurable, see the options) and will be deleted after that delay.

Note that the collection is done in background to avoid degrading the performance of the requests.

### Number of entry

The number of entry is limited (configurable, see the options).

Once the limit is reached the least recently used items will be deleted.

## Components

### Caching

Contains caching logic.

These structures are not thread safe!

- `cache/cache`: contains the cache
- `cache/lru`: contains the LRU index
- `cache/ttl`: contains the Time To Live logic

### HTTP

Contains the HTTP interface.

- `http/server`: HTTP server
- `http/service`: instance that processing HTTP requests

### Redis

Contains the communication with the backing Redis server.

- `redis/redis.go`: Redis client

## Redis dashboard

Go to [localhost:8081](http://localhost:8081).

## Caveats

- Only URL safe characters are permitted as key.
- The go-redis package https://github.com/go-redis/redis/issues/641
- I'm not using go 1.9, the code could be clearer with type aliases.
- With very bad timing you can:
  - See a dead cache entry (no stale data)
- Logging and error handling are intentionally simple
- Unexisting keys in the proxy cache will be fetched from Redis every time
- Value in Redis cannot be an empty string
- I'm still using < Go1.9 which is annoying when you want to use type aliases
- The `http/put-handler` has a DOS vulnerability, it's using a ReadCloser to read the request body, an attacker could take advantage of not terminating the content and filling the http server's capacity.
- I didn't set up a cleanup mechanism on shutdown

## Personal thoughts

- I priviledged doing my own implementation rather than using existing libs
- I splited the cache into seperated pieces, it allows more granular locking and it's extensible but could be more error prone.

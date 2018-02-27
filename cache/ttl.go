package cache

import (
	"log"
	"time"
)

type TTLWorker struct {
	ticker  *time.Ticker
	backend map[string]time.Time

	events chan KeyExpired
	ttl    time.Duration
}

const (
	LOOP_INTERVAL = 1 * time.Second
)

type KeyExpired struct{ Name string }

// Return a TTLWorker (the ttl is expressed second)
func MakeTTLWorker(ttl int) TTLWorker {
	return TTLWorker{
		backend: make(map[string]time.Time),
		events:  make(chan KeyExpired),
		ticker:  time.NewTicker(LOOP_INTERVAL),
		ttl:     time.Duration(ttl) * time.Second,
	}
}

func (worker TTLWorker) Update(key string) {
	worker.backend[key] = time.Now()
}

func (worker TTLWorker) Run() {
	for {
		<-worker.ticker.C
		worker.tick()
	}
}

func (worker TTLWorker) tick() {
	now := time.Now()
	backendCopy := worker.backend

	for key, time := range backendCopy {

		// Check if the item has expired
		if time.Add(worker.ttl).Before(now) {
			select {
			case worker.events <- KeyExpired{key}:
				delete(worker.backend, key)
			default:
				log.Println("Dropped cache invalidation request")
			}
		}
	}
}

func (worker TTLWorker) Events() chan KeyExpired {
	return worker.events
}

package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExpired(t *testing.T) {
	worker := MakeTTLWorker(0)

	waitChan := make(chan bool)

	go func() {
		select {
		case notif := <-worker.Events():
			assert.Equal(t, notif.Name, "a")
			waitChan <- true
		}
	}()

	worker.Update("a")
	worker.tick()

	<-waitChan
}

func TestIsNotExpired(t *testing.T) {
	worker := MakeTTLWorker(100)

	go func() {
		select {
		case <-worker.Events():
			t.Fatal("Item shoudln't have expired")
		}
	}()

	worker.Update("a")
	worker.tick()
	worker.tick()
	worker.tick()
}

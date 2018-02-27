package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/xtuc/redis-proxy/cache"
	"github.com/xtuc/redis-proxy/http"
)

/*
  Backend interface mock
*/

type BackendMock struct {
	mock.Mock
}

func (m BackendMock) GetValue(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (BackendMock) PutValue(key, value string) error { return nil }

/*
  Actual tests
*/

func TestSendKeyNotFoundError(t *testing.T) {
	w := httptest.NewRecorder()

	http.SendKeyNotFoundError("key", w)

	assert.Equal(t, w.Header().Get(http.CACHE_HEADER_KEY), "miss")
}

func TestSendOk(t *testing.T) {
	w := httptest.NewRecorder()

	http.SendKeyOk("key", "value", w)

	assert.Equal(t, w.Header().Get(http.CACHE_HEADER_KEY), "hit")
}

func TestSendOkAlreadySet(t *testing.T) {
	w := httptest.NewRecorder()

	w.Header().Add(http.CACHE_HEADER_KEY, "foobar")
	http.SendKeyOk("key", "value", w)

	assert.Equal(t, w.Header().Get(http.CACHE_HEADER_KEY), "foobar")
}

func TestSendInCache(t *testing.T) {
	fakeBackend := new(BackendMock)

	cacheInstance := cache.NewCache()
	capacity := 1

	service := http.NewHTTPService(
		cacheInstance,
		cache.NewLRUIndex(capacity),
		fakeBackend,
	)

	assert.Nil(t, service.SetInCache("foo", "bar"))
	assert.Nil(t, service.SetInCache("foo1", "bar1"))
	assert.Nil(t, service.SetInCache("foo2", "bar2"))

	assert.Equal(t, capacity, cacheInstance.Size())
}

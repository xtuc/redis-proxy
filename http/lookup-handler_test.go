package http_test

import (
	httppkg "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/xtuc/redis-proxy/cache"
	"github.com/xtuc/redis-proxy/http"
)

var (
	lookupReq, _ = httppkg.NewRequest("GET", "/{key}", nil)
)

func TestLookupHandlerUncachedKey(t *testing.T) {
	fakeBackend := new(BackendMock)
	fakeBackend.On("GetValue", "foo").Return("bar", nil)

	service := http.NewHTTPService(
		cache.NewCache(),
		cache.NewLRUIndex(10),
		fakeBackend,
	)

	lookupReq = mux.SetURLVars(lookupReq, map[string]string{
		"key": "foo",
	})

	rr := httptest.NewRecorder()
	handler := httppkg.HandlerFunc(service.LookupHandler)
	handler.ServeHTTP(rr, lookupReq)

	assert.Equal(t, httppkg.StatusOK, rr.Code)
	assert.Equal(t, "bar", rr.Body.String())
	fakeBackend.AssertExpectations(t)
}

func TestLookupHandlerUnexistingKey(t *testing.T) {
	fakeBackend := new(BackendMock)
	fakeBackend.On("GetValue", "foo").Return("", nil)

	service := http.NewHTTPService(
		cache.NewCache(),
		cache.NewLRUIndex(10),
		fakeBackend,
	)

	lookupReq = mux.SetURLVars(lookupReq, map[string]string{
		"key": "foo",
	})

	rr := httptest.NewRecorder()
	handler := httppkg.HandlerFunc(service.LookupHandler)
	handler.ServeHTTP(rr, lookupReq)

	assert.Equal(t, httppkg.StatusNotFound, rr.Code)
	fakeBackend.AssertExpectations(t)
}

func TestLookupHandlerCacheKey(t *testing.T) {
	fakeBackend := new(BackendMock)
	cacheInstance := cache.NewCache()

	cacheInstance.Set("foo", "bar")

	service := http.NewHTTPService(
		cacheInstance,
		cache.NewLRUIndex(10),
		fakeBackend,
	)

	lookupReq = mux.SetURLVars(lookupReq, map[string]string{
		"key": "foo",
	})

	rr := httptest.NewRecorder()
	handler := httppkg.HandlerFunc(service.LookupHandler)
	handler.ServeHTTP(rr, lookupReq)

	assert.Equal(t, httppkg.StatusOK, rr.Code)
	assert.Equal(t, "bar", rr.Body.String())
	fakeBackend.AssertExpectations(t)
}

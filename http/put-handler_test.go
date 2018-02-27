package http_test

import (
	"io/ioutil"
	httppkg "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/xtuc/redis-proxy/cache"
	"github.com/xtuc/redis-proxy/http"
)

var (
	putReq, _ = httppkg.NewRequest("PUT", "/{key}", nil)
)

func TestPutHandlerNewKey(t *testing.T) {
	value := "a"

	fakeBackend := new(BackendMock)

	cacheInstance := cache.NewCache()

	service := http.NewHTTPService(
		cacheInstance,
		cache.NewLRUIndex(10),
		fakeBackend,
	)

	putReq.Body = ioutil.NopCloser(strings.NewReader(value))

	putReq = mux.SetURLVars(putReq, map[string]string{
		"key": "foo",
	})

	rr := httptest.NewRecorder()
	handler := httppkg.HandlerFunc(service.PutHandler)
	handler.ServeHTTP(rr, putReq)

	assert.Equal(t, httppkg.StatusOK, rr.Code)
	assert.Equal(t, value, rr.Body.String())
	fakeBackend.AssertExpectations(t)

	assert.Equal(t, value, *cacheInstance.Get("foo"))
}

func TestPutHandlerUpdateKey(t *testing.T) {
	fakeBackend := new(BackendMock)
	newValue := "b"

	cacheInstance := cache.NewCache()
	cacheInstance.Set("foo", "oldvalue")

	service := http.NewHTTPService(
		cacheInstance,
		cache.NewLRUIndex(10),
		fakeBackend,
	)

	putReq.Body = ioutil.NopCloser(strings.NewReader(newValue))

	putReq = mux.SetURLVars(putReq, map[string]string{
		"key": "foo",
	})

	rr := httptest.NewRecorder()
	handler := httppkg.HandlerFunc(service.PutHandler)
	handler.ServeHTTP(rr, putReq)

	assert.Equal(t, httppkg.StatusOK, rr.Code)
	assert.Equal(t, newValue, rr.Body.String())
	fakeBackend.AssertExpectations(t)

	assert.Equal(t, newValue, *cacheInstance.Get("foo"))
}

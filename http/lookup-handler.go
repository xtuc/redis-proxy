package http

import (
	httppkg "net/http"

	"github.com/gorilla/mux"
)

func (c HTTPService) LookupHandler(w httppkg.ResponseWriter, r *httppkg.Request) {
	vars := mux.Vars(r)

	key, hasKey := vars["key"]

	if !hasKey {
		SendKeyNotFoundError("", w)
		return
	}

	value, wasCached, err := c.GetFromCacheOrBackend(key)

	// FIXME(sven): do something better here
	if err != nil {
		panic(err)
	}

	// Key was not found
	if value == "" {
		SendKeyNotFoundError(key, w)
		return
	}

	// Set the miss header is we missed the cache
	if !wasCached {
		w.Header().Add(CACHE_HEADER_KEY, "miss")
	}

	SendKeyOk(key, value, w)
}

// Foo

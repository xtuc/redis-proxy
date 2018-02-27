package http

import (
	"errors"
	"io/ioutil"
	"log"
	httppkg "net/http"

	"github.com/gorilla/mux"
)

var (
	EMPTY_BODY_ERR = errors.New("Body can not be empty")
)

func (c HTTPService) PutHandler(w httppkg.ResponseWriter, r *httppkg.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if r.Body == nil {
		panic(EMPTY_BODY_ERR)
	}

	bytesValue, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	value := string(bytesValue)

	cachedValue := c.cache.Get(key)

	if cachedValue != nil {
		log.Printf("200 - updated value for key %s\n", key)
	} else {
		log.Printf("200 - added value for key %s\n", key)
	}

	SendKeyOk(key, value, w)
	cacheErr := c.SetInCache(key, value)

	if cacheErr != nil {
		panic(cacheErr)
	}
}

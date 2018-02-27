GO = go

.PHONY: build

install-deps:
	go get -v ./...

build: install-deps
	cd cmd/redis-proxy && \
	  $(GO) build

test:
	go test -v ./...

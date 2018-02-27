FROM golang:1.8-alpine3.5 as builder

RUN apk update
RUN apk add git make

ENV APP_HOME $GOPATH/src/github.com/xtuc/redis-proxy

COPY ./ $APP_HOME

WORKDIR $APP_HOME

RUN make install-deps build

RUN cp -f $APP_HOME/cmd/redis-proxy/redis-proxy /usr/bin

FROM alpine:3.5

COPY --from=builder /usr/bin/redis-proxy /usr/bin/redis-proxy

CMD ["/usr/bin/redis-proxy"]

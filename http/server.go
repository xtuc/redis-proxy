package http

import "net/http"

type HTTPServer struct {
	waitChan chan interface{}
	server   *http.Server
}

func NewHTTPServer(host, port string, handler http.Handler) *HTTPServer {
	waitChan := make(chan interface{})

	server := &http.Server{
		Addr:    host + ":" + port,
		Handler: handler,
	}

	return &HTTPServer{
		server:   server,
		waitChan: waitChan,
	}
}

func (instance *HTTPServer) Start() (chan interface{}, error) {
	err := instance.server.ListenAndServe()

	if err != nil {
		return instance.waitChan, err
	}

	return instance.waitChan, nil
}

func (instance *HTTPServer) Stop() error {
	instance.waitChan <- nil
	close(instance.waitChan)

	return instance.server.Close()
}

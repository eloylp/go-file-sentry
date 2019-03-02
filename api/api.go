package api

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
)

type apiServer struct {
	address string
	server  *http.Server
	errors  chan error
	wg      *sync.WaitGroup
}

func (a *apiServer) Errors() chan error {
	return a.errors
}

func NewApiServer(address string, wg *sync.WaitGroup) *apiServer {
	return &apiServer{
		address: address,
		errors:  make(chan error),
		wg:      wg,
	}
}

func (a *apiServer) StartServer() {
	a.server = &http.Server{
		Handler: a.router(),
	}
	unixListener, err := net.Listen("unix", a.address)
	if err != nil {
		log.Fatal(err)

	}
	a.wg.Add(1)
	err = a.server.Serve(unixListener)
	a.wg.Done()
	if err != nil {
		a.errors <- err
	}
}

func (a *apiServer) Shutdown() {
	err := a.server.Shutdown(context.Background())
	if err != nil {
		a.errors <- err
	}
}

func (a *apiServer) router() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/status", a.statusHandler())
	return router
}

func (a *apiServer) writeResponse(w http.ResponseWriter, b []byte) {
	_, err := w.Write(b)
	if err != nil {
		a.errors <- err
	}
}

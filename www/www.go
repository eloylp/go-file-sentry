package www

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
)

type www struct {
	address string
	server  *http.Server
	errors  chan error
	wg      *sync.WaitGroup
}

func (s *www) Errors() chan error {
	return s.errors
}

func NewApiServer(address string, wg *sync.WaitGroup) *www {
	return &www{
		address: address,
		errors:  make(chan error),
		wg:      wg,
	}
}

func (s *www) StartServer() {
	s.server = &http.Server{
		Handler: s.router(),
	}
	unixListener, err := net.Listen("unix", s.address)
	if err != nil {
		log.Fatal(err)

	}
	s.wg.Add(1)
	err = s.server.Serve(unixListener)
	s.wg.Done()
	if err != nil {
		s.errors <- err
	}
}

func (s *www) Shutdown() {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		s.errors <- err
	}
}

func (s *www) router() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/status", s.statusHandler())
	return router
}

func (s *www) writeResponse(w http.ResponseWriter, b []byte) {
	_, err := w.Write(b)
	if err != nil {
		s.errors <- err
	}
}

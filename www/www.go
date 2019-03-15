package www

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
)

type www struct {
	socket *url.URL
	server *http.Server
	errors chan error
	wg     *sync.WaitGroup
}

func (s *www) Errors() chan error {
	return s.errors
}

func NewApiServer(address *url.URL, wg *sync.WaitGroup) *www {
	return &www{
		socket: address,
		errors: make(chan error),
		wg:     wg,
	}
}

func (s *www) StartServer() {
	s.server = &http.Server{
		Handler: s.router(),
	}
	a := s.address()
	l, err := net.Listen(s.socket.Scheme, a)
	if err != nil {
		log.Fatal(err)

	}
	s.wg.Add(1)
	err = s.server.Serve(l)
	s.wg.Done()
	if err != nil {
		s.errors <- err
	}
}

func (s *www) address() string {
	var a string
	if s.socket.Scheme == "unix" {
		a = s.socket.Path
	} else {
		a = s.socket.Host
	}
	return a
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

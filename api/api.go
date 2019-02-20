package api

import (
	"context"
	"net"
	"net/http"
)

type ApiServer struct {
	Address  string
	Shutdown chan struct{}
	server   *http.Server
	Ctx      context.Context
	Error    chan error
}

func (a *ApiServer) StartServer() {

	a.server = &http.Server{
		Handler: a.router(),
	}
	unixListener, err := net.Listen("unix", a.Address)
	if err != nil {
		a.Error <- err
		return
	}
	go a.handleShutdown()
	err = a.server.Serve(unixListener)
	if err != nil {
		a.Error <- err
		return
	}
}

func (a *ApiServer) router() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/status", a.statusHandler())
	return router
}

func (a *ApiServer) handleShutdown() {
	<-a.Shutdown
	err := a.server.Shutdown(a.Ctx)
	if err != nil {
		a.Error <- err
	}
}

func (a *ApiServer) writeResponse(w http.ResponseWriter, b []byte) {
	_, err := w.Write(b)
	if err != nil {
		a.Error <- err
	}
}

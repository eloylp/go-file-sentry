package main

import (
	"github.com/eloylp/go-file-sentry/config"
	"github.com/eloylp/go-file-sentry/factory"
	"github.com/eloylp/go-file-sentry/program"
	"github.com/eloylp/go-file-sentry/term"
	"github.com/eloylp/go-file-sentry/www"
	"sync"
)

func main() {

	mainShutdown := make(chan struct{})
	term.Listen(mainShutdown)
	cfg := config.NewConfigFromParams()
	wg := new(sync.WaitGroup)
	apiServer := www.NewApiServer(cfg.Socket(), wg)
	watchers := factory.Watchers(cfg, wg)
	p := program.Program{
		Api:      apiServer,
		Watchers: watchers,
		Config:   cfg,
		Wg:       wg,
	}
	p.Start()
	<-mainShutdown
	p.Shutdown()
}

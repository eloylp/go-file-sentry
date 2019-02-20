package main

import (
	"github.com/eloylp/go-file-sentry/config"
	"github.com/eloylp/go-file-sentry/program"
	"github.com/eloylp/go-file-sentry/sentry"
	"github.com/eloylp/go-file-sentry/term"
	"sync"
)

func main() {

	mainShutdown := make(chan struct{})
	term.Listen(mainShutdown)
	cfg := config.NewConfigFromParams()
	wg := new(sync.WaitGroup)
	watchers := sentry.Watchers(cfg, wg)
	p := program.Program{
		Watchers: watchers,
		Config:   cfg,
		Wg:       wg,
	}
	p.Start()
	<-mainShutdown
	p.Shutdown()
}

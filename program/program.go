package program

import (
	"github.com/eloylp/go-file-sentry/config"
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/version"
	"log"
	"sync"
)

type ApiServer interface {
	StartServer()
	Shutdown()
	Errors() chan error
}

type Watcher interface {
	WFile(handler func(f *file.File))
	Shutdown()
	Err() chan error
	Info() chan string
}

type Program struct {
	Api      ApiServer
	Watchers []Watcher
	Config   *config.Config
	Wg       *sync.WaitGroup
}

func (p *Program) Shutdown() {
	p.Api.Shutdown()
	for _, w := range p.Watchers {
		w.Shutdown()
	}
	p.Wg.Wait()
}

func (p *Program) Start() {
	p.startApi()
	p.startWatching()
	p.startLogging()
}

func (p *Program) startApi() {
	go p.Api.StartServer()
}

func (p *Program) startWatching() {
	for _, w := range p.Watchers {
		go w.WFile(func(f *file.File) {
			log.Printf("Saving new version of file %s", f.Path())
			version.NewVersion(p.Config.StoragePath(), f)
		})
	}
}

func (p *Program) startLogging() {
	go func(server ApiServer) {
		for err := range server.Errors() {
			log.Print(err)
		}
	}(p.Api)
	for _, w := range p.Watchers {
		go func(w Watcher) {
			for err := range w.Err() {
				log.Print(err)
			}
		}(w)
		go func(w Watcher) {
			for info := range w.Info() {
				log.Print(info)
			}
		}(w)
	}
}

package program

import (
	"github.com/eloylp/go-file-sentry/config"
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/version"
	"github.com/eloylp/go-file-sentry/watcher"
	"log"
	"sync"
)

type Program struct {
	Watchers []*watcher.Watcher
	Config   *config.Config
	Wg       *sync.WaitGroup
}

func (p *Program) Shutdown() {
	for _, w := range p.Watchers {
		w.Shutdown <- struct{}{}
	}
	p.Wg.Wait()
}

func (p *Program) Start() {
	p.startWatching()
	p.startLogging()
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
	for _, w := range p.Watchers {
		go func(w *watcher.Watcher) {
			for err := range w.Err {
				log.Print(err)
			}
		}(w)
		go func(w *watcher.Watcher) {
			for info := range w.Infos {
				log.Print(info)
			}
		}(w)
	}
}

package factory

import (
	"github.com/eloylp/go-file-sentry/config"
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/program"
	"github.com/eloylp/go-file-sentry/watcher"
	"sync"
)

func Watchers(cfg *config.Config, wg *sync.WaitGroup) []program.Watcher {
	var watchers []program.Watcher
	for _, wFile := range cfg.WFiles() {
		wFile := file.NewFile(wFile)
		watchers = append(watchers, watcher.NewWatcher(wFile, wg))
	}
	return watchers
}

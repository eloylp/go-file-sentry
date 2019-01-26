package sentry

import (
	"github.com/eloylp/go-file-sentry/config"
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/version"
	"github.com/eloylp/go-file-sentry/watcher"
)

func Start(cfg *config.Config) {
	for _, watchedFile := range cfg.WFiles() {
		wFile := file.NewFile(watchedFile)
		go watcher.WFile(wFile, func(file *file.File) {
			version.NewVersion(cfg.StoragePath(), file)
		})
	}
}

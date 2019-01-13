package sentry

import (
	"github.com/eloylp/go-file-sentry/config"
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/version"
	"github.com/eloylp/go-file-sentry/watcher"
)

func StartSentry(cfg *config.Config) {
	for _, watchedFile := range cfg.WatchedFiles() {
		wFile := file.NewFile(watchedFile)
		go watcher.WatchFile(wFile, func(file *file.File) {
			version.MakeNewVersion(cfg.StoragePath(), file)
		})
	}
}

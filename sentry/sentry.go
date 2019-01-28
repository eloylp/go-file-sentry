package sentry

import (
	"github.com/eloylp/go-file-sentry/config"
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/version"
	"github.com/eloylp/go-file-sentry/watcher"
	"log"
)

func Start(cfg *config.Config) {
	for _, wFile := range cfg.WFiles() {
		wFile := file.NewFile(wFile)
		go watcher.WFile(wFile, func(file *file.File) {
			log.Printf("Saving new version of file %s", wFile.Path())
			version.NewVersion(cfg.StoragePath(), file)
		})
	}
}
package watcher

import (
	"github.com/eloylp/go-file-sentry/file"
	"github.com/rjeczalik/notify"
	"log"
)

func WFile(wFile *file.File, handler func(file *file.File)) {
	events := make(chan notify.EventInfo, 1)
	err := notify.Watch(wFile.Path(), events, notify.InCloseWrite)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Starting watching file %s", wFile.Path())
	for event := range events {
		if event.Event() == notify.InCloseWrite {
			log.Printf("Changes watched in %s, handling a new version ...", wFile.Path())
			wFile.LoadMetadata()
			handler(wFile)
		}
	}
}

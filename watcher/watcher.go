package watcher

import (
	"github.com/eloylp/go-file-sentry/file"
	"github.com/rjeczalik/notify"
	"log"
)

var wEvents = []notify.Event{
	notify.InCloseWrite,
	notify.InModify,
	notify.Write,
	notify.Remove,
}

func WFile(wFile *file.File, handler func(file *file.File)) {
	events := make(chan notify.EventInfo)
	err := notify.Watch(wFile.Path(), events, notify.All)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Starting watching file %s", wFile.Path())
	for event := range events {
		if isWEvent(event) {
			log.Printf("Changes watched in %s, handling a new version ...", wFile.Path())
			wFile.LoadMetadata()
			handler(wFile)
		} else {
			log.Printf("Ignoring event %s for file %s", event.Event(), wFile.Path())
		}
	}
}

func isWEvent(event notify.EventInfo) bool {
	for _, wEvent := range wEvents {
		if event.Event() == wEvent {
			return true
		}
	}
	return false
}

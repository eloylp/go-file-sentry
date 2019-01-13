package watcher

import (
	"github.com/eloylp/go-file-sentry/file"
	"github.com/rjeczalik/notify"
	"log"
)

func WatchFile(file *file.File, handler func(file *file.File)) {
	events := make(chan notify.EventInfo, 1)
	err := notify.Watch(file.GetPath(), events, notify.InCloseWrite)
	if err != nil {
		log.Fatal(err)
	}
	for event := range events {
		if event.Event() == notify.InCloseWrite {
			file.LoadMetadata()
			handler(file)
		}
	}
}

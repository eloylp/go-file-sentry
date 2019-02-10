package watcher

import (
	"github.com/eloylp/go-file-sentry/file"
	"github.com/fsnotify/fsnotify"
	"log"
)

var wEvents = []fsnotify.Op{
	fsnotify.Write,
	fsnotify.Create,
	fsnotify.Remove,
}

func WFile(wFile *file.File, handler func(file *file.File)) {
	watcher, err := fsnotify.NewWatcher()
	defer watcher.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Add(wFile.Path())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Starting watching file %s", wFile.Path())
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if changedEvent(event) {
				log.Printf("Changes watched in %s, handling a new version ...", wFile.Path())
				wFile.LoadMetadata()
				handler(wFile)
			} else if event.Op&fsnotify.Rename == fsnotify.Rename {
				log.Printf("Adding new inode watcher for file %s due to inode change", wFile.Path())
				_ = watcher.Add(wFile.Path())
			} else {
				log.Printf("Ignoring event %s", event.String())
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Error in watcher for file %s: %v", wFile.Path(), err)
		}
	}
}

func changedEvent(event fsnotify.Event) bool {
	for _, op := range wEvents {
		if event.Op&op == op {
			return true
		}
	}
	return false
}

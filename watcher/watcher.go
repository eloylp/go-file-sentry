package watcher

import (
	"fmt"
	"github.com/eloylp/go-file-sentry/file"
	"github.com/fsnotify/fsnotify"
	"sync"
)

var wEvents = []fsnotify.Op{
	fsnotify.Write,
	fsnotify.Create,
	fsnotify.Remove,
}

type Watcher struct {
	Shutdown chan struct{}
	Err      chan error
	Infos    chan string
	File     *file.File
	Wg       *sync.WaitGroup
}

func NewWatcher(file *file.File, wg *sync.WaitGroup) *Watcher {
	w := &Watcher{File: file}
	w.Shutdown = make(chan struct{})
	w.Infos = make(chan string)
	w.Err = make(chan error)
	w.Wg = wg
	return w
}

func (w *Watcher) WFile(handler func(file *file.File)) {
	w.Wg.Add(1)
	defer close(w.Infos)
	defer close(w.Err)
	watcher, err := fsnotify.NewWatcher()
	defer watcher.Close()
	if err != nil {
		w.Err <- err
		return
	}
	err = watcher.Add(w.File.Path())
	if err != nil {
		w.Err <- err
		return
	}
	w.Infos <- fmt.Sprintf("Starting watching file %s", w.File.Path())
	//// TODO EVALUATE WATCHER CLOSE INSTEAD OF THIS.
wLoop:
	for {
		select {
		case <-w.Shutdown:
			w.Infos <- fmt.Sprintf("Shutting down watcher  of file %s", w.File.Path())
			break wLoop
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if changedEvent(event) {
				w.Infos <- fmt.Sprintf("Changes watched in %s, handling a new version ...", w.File.Path())
				w.File.LoadMetadata()
				handler(w.File)
			} else if event.Op&fsnotify.Rename == fsnotify.Rename {
				w.Infos <- fmt.Sprintf("Adding new inode watcher for file %s due to inode change", w.File.Path())
				err = watcher.Add(w.File.Path())
				if err != nil {
					w.Err <- err
				}
			} else {
				w.Infos <- fmt.Sprintf("Ignoring event %s", event.String())
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			w.Err <- err
		}
	}
	w.Wg.Done()
}

func changedEvent(event fsnotify.Event) bool {
	for _, op := range wEvents {
		if event.Op&op == op {
			return true
		}
	}
	return false
}

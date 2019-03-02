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
	shutdown chan struct{}
	err      chan error
	info     chan string
	file     *file.File
	wg       *sync.WaitGroup
}

func (w *Watcher) Info() chan string {
	return w.info
}

func (w *Watcher) Err() chan error {
	return w.err
}

func (w *Watcher) Shutdown() {
	w.shutdown <- struct{}{}
}

func NewWatcher(file *file.File, wg *sync.WaitGroup) *Watcher {
	w := &Watcher{file: file}
	w.shutdown = make(chan struct{})
	w.info = make(chan string)
	w.err = make(chan error)
	w.wg = wg
	return w
}

func (w *Watcher) WFile(handler func(f *file.File)) {
	w.wg.Add(1)
	defer close(w.info)
	defer close(w.err)
	watcher, err := fsnotify.NewWatcher()
	defer watcher.Close()
	if err != nil {
		w.err <- err
		return
	}
	err = watcher.Add(w.file.Path())
	if err != nil {
		w.err <- err
		return
	}
	w.info <- fmt.Sprintf("Starting watching file %s", w.file.Path())
	//// TODO EVALUATE WATCHER CLOSE INSTEAD OF THIS.
wLoop:
	for {
		select {
		case <-w.shutdown:
			w.info <- fmt.Sprintf("Shutting down watcher  of file %s", w.file.Path())
			break wLoop
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if changedEvent(event) {
				w.info <- fmt.Sprintf("Changes watched in %s, handling a new version ...", w.file.Path())
				w.file.LoadMetadata()
				handler(w.file)
			} else if event.Op&fsnotify.Rename == fsnotify.Rename {
				w.info <- fmt.Sprintf("Adding new inode watcher for file %s due to inode change", w.file.Path())
				err = watcher.Add(w.file.Path())
				if err != nil {
					w.err <- err
				}
			} else {
				w.info <- fmt.Sprintf("Ignoring event %s", event.String())
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			w.err <- err
		}
	}
	w.wg.Done()
}

func changedEvent(event fsnotify.Event) bool {
	for _, op := range wEvents {
		if event.Op&op == op {
			return true
		}
	}
	return false
}

package main

import (
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	Jitome      *Jitome
	Task        *Task
	WatchConfig *WatchConfig
	w           *fsnotify.Watcher
}

func (watcher *Watcher) Wait() {
	for {
		select {
		case event := <-watcher.w.Events:
			path := event.Name
			if event.Op&fsnotify.Chmod != 0 {
				// chmod is ignored
				continue
			}

			if event.Op&fsnotify.Create != 0 && isDir(path) {
				for _, watchConfig := range watcher.Task.Watch {
					err := watch(path, watchConfig.IgnoreDir, watcher.w)
					if err != nil {
						panic(err)
					}
				}
			}

			watcher.Jitome.Event <- &Event{
				Watcher: watcher,
				Ev:      event,
			}

		case err := <-watcher.w.Errors:
			panic(err)
		}
	}
}

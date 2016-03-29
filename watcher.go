package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

type Watcher struct {
	Jitome      *Jitome
	Target      *Target
	WatchConfig *WatchConfig
	w           *fsnotify.Watcher
	index       int
}

func NewWatcher(jitome *Jitome, target *Target, watchConfig *WatchConfig, w *fsnotify.Watcher, index int) (*Watcher, error) {
	watcher := &Watcher{
		Jitome:        jitome,
		Target:          target,
		WatchConfig:   watchConfig,
		w:             w,

		index: index,
	}

	return watcher, nil
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
				for _, watchConfig := range watcher.Target.Watch {
					err := watch(path, watchConfig.IgnorePatterns, watcher.w)
					if err != nil {
						panic(err)
					}
				}
			}

			path = normalizePath(path)

			if debug {
				log.Printf("detected changing '%s'", path)
			}

			// check pattern.
			if !watcher.Match(path) {
				if debug {
					log.Printf("target '%s' watcher %d detected changing '%s' but it was unmatched to pattern config", watcher.Target.Name, watcher.index, path)
				}
				continue
			}

			watcher.Target.events <- &Event{
				Watcher: watcher,
				Ev:      event,
			}

		case err := <-watcher.w.Errors:
			panic(err)
		}
	}
}

func (watcher *Watcher) Match(path string) bool {
	if len(watcher.WatchConfig.Patterns) == 0 {
		return true
	}

	for _, ptn := range watcher.WatchConfig.IgnorePatterns {
		if debug {
			log.Printf("cheking ignore pattern '%s'", ptn.String())
		}

		if ptn.MatchString(path) {
			if debug {
				log.Printf("matched ignore pattern '%s'", ptn.String())
			}
			return false
		}
	}

	for _, ptn := range watcher.WatchConfig.Patterns {
		if debug {
			log.Printf("cheking pattern '%s'", ptn.String())
		}

		if ptn.MatchString(path) {
			if debug {
				log.Printf("matched pattern '%s'", ptn.String())
			}

			return true
		}
	}

	return false
}

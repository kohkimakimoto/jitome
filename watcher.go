package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"reflect"
	"regexp"
)

type Watcher struct {
	Jitome              *Jitome
	Task                *Task
	WatchConfig         *WatchConfig
	w                   *fsnotify.Watcher
	watchPatternRegexps []*regexp.Regexp

	index int
}

func NewWatcher(jitome *Jitome, task *Task, watchConfig *WatchConfig, w *fsnotify.Watcher, index int) (*Watcher, error) {
	watcher := &Watcher{
		Jitome:              jitome,
		Task:                task,
		WatchConfig:         watchConfig,
		w:                   w,
		watchPatternRegexps: []*regexp.Regexp{},
		index:               index,
	}

	pattern := watchConfig.Pattern
	if pattern != nil {
		if patternStr, ok := pattern.(string); ok {
			reg, err := regexp.Compile(patternStr)
			if err != nil {
				return nil, fmt.Errorf("invalid pattern '%s': %v", patternStr, err)
			}
			watcher.watchPatternRegexps = append(watcher.watchPatternRegexps, reg)
		} else if e, ok := pattern.([]interface{}); ok {
			for _, patternStr := range e {
				reg, err := regexp.Compile(patternStr.(string))
				if err != nil {
					return nil, fmt.Errorf("invalid pattern '%s': %v", patternStr, err)
				}
				watcher.watchPatternRegexps = append(watcher.watchPatternRegexps, reg)
			}
		} else {
			v := reflect.ValueOf(pattern)
			return nil, fmt.Errorf("invalid format pattern: %v", v.Type())
		}
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
				for _, watchConfig := range watcher.Task.Watch {
					err := watch(path, watchConfig.IgnoreDir, watcher.w)
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
					log.Printf("'%s' watcher %d detected changing '%s' but it was unmatched to pattern config", watcher.Task.Name, watcher.index, path)
				}
				continue
			}

			watcher.Task.events <- &Event{
				Watcher: watcher,
				Ev:      event,
			}

		case err := <-watcher.w.Errors:
			panic(err)
		}
	}
}

func (watcher *Watcher) Match(path string) bool {
	ret := false
	if len(watcher.watchPatternRegexps) == 0 {
		return true
	}

	for _, reg := range watcher.watchPatternRegexps {
		if debug {
			log.Printf("cheking pattern '%s'", reg.String())
		}
		if reg != nil && reg.MatchString(path) {
			ret = true
			break
		}
	}
	return ret
}

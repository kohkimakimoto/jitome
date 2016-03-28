package main

import (
	"github.com/fsnotify/fsnotify"
	"regexp"
	"reflect"
	"fmt"
)

type Watcher struct {
	Jitome      *Jitome
	Task        *Task
	WatchConfig *WatchConfig
	w           *fsnotify.Watcher
	watchPatternRegexps   []*regexp.Regexp
}

func NewWatcher(jitome *Jitome, task *Task, watchConfig *WatchConfig, w *fsnotify.Watcher) (*Watcher, error) {
	watcher := &Watcher{
		Jitome:      jitome,
		Task:        task,
		WatchConfig: watchConfig,
		w:           w,
		watchPatternRegexps:[]*regexp.Regexp{},
	}

	pattern := watchConfig.Pattern
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

			// check pattern.

			watcher.Task.events <- &Event{
				Watcher: watcher,
				Ev:      event,
			}

		case err := <-watcher.w.Errors:
			panic(err)
		}
	}
}

//func (task *Task) Match(path string) bool {
//	ret := false
//	for i, reg := range task.watchRegexps {
//		if reg != nil && reg.MatchString(path) {
//
//			ret = true
//			if debug {
//				printDebugLog("Matched a watch string: '" + task.watchStrings[i] + "' (" + path + ")")
//			}
//
//			break
//		} else {
//			if debug {
//				printDebugLog("Unmatched a watch string: '" + task.watchStrings[i] + "' (" + path + ")")
//			}
//		}
//	}
//
//	if ret {
//		for i, reg := range task.excludeRegexps {
//			if reg != nil && reg.MatchString(path) {
//				if debug {
//					printDebugLog("Matched a exclude string'" + task.excludeStrings[i] + "' (" + path + ")")
//				}
//				ret = false
//			} else {
//				if debug {
//					printDebugLog("Unmatched a exclude string'" + task.excludeStrings[i] + "' (" + path + ")")
//				}
//			}
//		}
//	}
//
//	return ret
//}

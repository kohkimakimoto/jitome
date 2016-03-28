package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

type Jitome struct {
	Config   *Config
	Watchers []*Watcher

	// Event is queue that receives file change event.
	Events chan *Event
}

type Event struct {
	Watcher *Watcher
	Ev      fsnotify.Event
}

func NewJitome(config *Config) *Jitome {
	w := &Jitome{
		Config:   config,
		Watchers: []*Watcher{},
		Events:   make(chan *Event),
	}

	return w
}

func (jitome *Jitome) Start() error {
	// register watchers
	for name, task := range jitome.Config.Tasks {
		task.Name = name
		task.events = make(chan *Event, 30)
		task.jitome = jitome

		if debug {
			log.Printf("setup '%s'", name)
		}


		// register watched directories
		for i, watchConfig := range task.Watch {
			w, err := fsnotify.NewWatcher()
			if err != nil {
				return err
			}

			err = watch(watchConfig.Base, watchConfig.IgnoreDir, w)
			if err != nil {
				return err
			}

			watcher, err := NewWatcher(jitome, task, watchConfig, w, i)
			if err != nil {
				return err
			}

			jitome.Watchers = append(jitome.Watchers, watcher)
			go watcher.Wait()
		}
		go task.Wait()
	}
	defer jitome.Close()

	log.Print(FgGB("starting jitome..."))

	for {
		event := <-jitome.Events
		runTask(event)
	}

	return nil
}

func runTask(event *Event) {
	log.Printf("'%s' detected changing '%s' [%s].", FgCB(event.Watcher.Task.Name), FgYB(event.Ev.Name), FgGB(eventOpStr(&event.Ev)))

	path := event.Ev.Name
	task := event.Watcher.Task
	code := task.Script

	code = os.Expand(code, func(s string) string {
		switch s {
		case "JITOME_FILE":
			return path
		}
		return os.Getenv(s)
	})

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", code)
	} else {
		cmd = exec.Command("sh", "-c", code)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Printf("[warning] %v", err)
	}
}

func watch(base string, ignoreDir interface{}, watcher *fsnotify.Watcher) error {
	if base == "" {
		base = "."
	}

	if debug {
		log.Printf("walks watched directories '%s'", base)
	}

	ignores := []string{}
	if ignoreDir != nil {
		if e, ok := ignoreDir.(string); ok {
			ignores = append(ignores, e)
		} else if e, ok := ignoreDir.([]interface{}); ok {
			for _, i := range e {
				ignores = append(ignores, i.(string))
			}
		} else {
			v := reflect.ValueOf(ignoreDir)
			return fmt.Errorf("invalid format ignore_dir: %v", v.Type())
		}
		if debug {
			log.Printf("ignore_dir: %v", ignores)
		}
	}

	// register watched directories.
	err := filepath.Walk(base, func(path string, fi os.FileInfo, err error) error {
		if err != nil || !fi.IsDir() {
			return nil
		}

		//path = normalizePath(path)

		for _, ig := range ignores {
			if strings.HasPrefix(ig, "/") {
				ig = strings.TrimPrefix(ig, "/")
				if strings.HasPrefix(path, ig) {
					return nil
				}
			} else {
				if strings.Contains(path, ig) {
					return nil
				}
			}
		}

		err = watcher.Add(path)
		if err != nil {
			return err
		}

		if debug {
			log.Printf("added watched dir '%s'", path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (jitome *Jitome) Close() {
	for _, watcher := range jitome.Watchers {
		watcher.w.Close()
	}
}

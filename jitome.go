package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type Jitome struct {
	Config   *Config
	Watchers []*Watcher

	// Event is queue that receives file change event.
	Events chan *Event

	InitTask *Target
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
	log.Print("starting jitome...")

	// check init task
	for name, target := range jitome.Config.Targets {
		if name == "init" {
			target.Name = name
			target.jitome = jitome
			jitome.InitTask = target
			if debug {
				log.Print("registered 'init' task")
			}
			break
		}
	}

	if jitome.InitTask != nil {
		if jitome.InitTask.Script != "" {
			log.Print("running '" + FgCB("init") + "' target.")
			err := runScript(jitome.InitTask.Script)
			if err != nil {
				return err
			}
			log.Print("finished '" + FgCB("init") + "' target.")
		}
	}

	// register watchers
	for name, target := range jitome.Config.Targets {
		if name == "init" {
			continue
		}

		target.Name = name
		target.events = make(chan *Event, 30)
		target.jitome = jitome

		log.Print("evaluate target '" + FgCB(name) + "'.")

		// register watched directories
		for i, watchConfig := range target.Watch {
			watchConfig.InitPatterns()
			w, err := fsnotify.NewWatcher()
			if err != nil {
				return err
			}

			err = watch(watchConfig.Base, watchConfig.IgnorePatterns, w)
			if err != nil {
				return err
			}

			watcher, err := NewWatcher(jitome, target, watchConfig, w, i)
			if err != nil {
				return err
			}

			jitome.Watchers = append(jitome.Watchers, watcher)
			go watcher.Wait()
		}
		go target.Wait()
	}
	defer jitome.Close()

	log.Print("watching files...")

	for {
		event := <-jitome.Events
		runTask(event)
	}

	return nil
}

func runTask(event *Event) {
	log.Printf("'%s' target detected '%s' changing [%s]. running script.", FgCB(event.Watcher.Target.Name), FgYB(event.Ev.Name), FgGB(eventOpStr(&event.Ev)))

	path := event.Ev.Name
	task := event.Watcher.Target
	script := task.Script

	script = os.Expand(script, func(s string) string {
		switch s {
		case "JITOME_FILE":
			return path
		}
		return os.Getenv(s)
	})

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", script)
	} else {
		cmd = exec.Command("sh", "-c", script)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Printf("[warning] %v", err)
	}

	log.Printf("'%s' target finished script.", FgCB(event.Watcher.Target.Name))
}

func runScript(script string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", script)
	} else {
		cmd = exec.Command("sh", "-c", script)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func watch(base string, ignorePatterns []*regexp.Regexp, watcher *fsnotify.Watcher) error {
	if base == "" {
		base = "."
	}


	if debug {
		log.Printf("walks watched directories '%s'", base)
	}

	// register watched directories.
	err := filepath.Walk(base, func(path string, fi os.FileInfo, err error) error {
		if err != nil || !fi.IsDir() {
			// watche only directries
			return nil
		}

		path = normalizePath(path)

		for _, pattern := range ignorePatterns {
			if !strings.HasPrefix(path, "/") {
				// add "/" to the path
				if ok := pattern.MatchString("/" + path); ok {
					return nil
				}
			} else {
				if ok := pattern.MatchString(path); ok {
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

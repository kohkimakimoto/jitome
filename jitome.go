package main

import (
	"fmt"
	"github.com/deckarep/gosx-notifier"
	"github.com/fsnotify/fsnotify"
	"github.com/kohkimakimoto/jitome/bindata"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type Jitome struct {
	config   *Config
	watchers []*Watcher
	events   chan *Event
	cmd      *exec.Cmd
}

type Event struct {
	Watcher *Watcher
	Ev      fsnotify.Event
}

func NewJitome(config *Config) *Jitome {
	w := &Jitome{
		config:   config,
		watchers: []*Watcher{},
		events:   make(chan *Event),
	}

	return w
}

func (jitome *Jitome) Start() error {
	// check init task
	// register watchers
	for name, target := range jitome.config.Targets {
		if name == "init" {
			continue
		}

		target.Name = name
		target.events = make(chan *Event, 30)
		target.jitome = jitome

		log.Print("evaluating target '" + FgCB(name) + "'.")

		// register watched directories
		for i, watchConfig := range target.Watch {
			watchConfig.InitPatterns()
			w, err := fsnotify.NewWatcher()
			if err != nil {
				return err
			}

			err = watch(watchConfig.Base, watchConfig.ignoreRegs, w)
			if err != nil {
				return err
			}

			watcher, err := NewWatcher(jitome, target, watchConfig, w, i)
			if err != nil {
				return err
			}

			jitome.watchers = append(jitome.watchers, watcher)
			go watcher.Wait()
		}
		go target.Wait()
	}
	defer jitome.Close()

	go func() {
		log.Print("watching files...")
		for {
			event := <-jitome.events
			runTarget(event)
		}
	}()

	go func() {
		err := jitome.restartCommand()
		if err != nil {
			log.Print(FgRB(fmt.Sprintf("[warning] %v", err)))
		}
	}()

	// wait infinitely.
	select {}
}

func (jitome *Jitome) restartCommand() error {
	if len(jitome.config.commandArgs) == 0 {
		return nil
	}

	if debug {
		log.Printf("restart command '%s'", jitome.config.Command)
	}

	if err := jitome.terminate(); err != nil {
		return err
	}

	return jitome.spawn()
}

//
// spawn and terminate refers to https://github.com/mattn/goemon
// MIT License: Yasuhiro Matsumoto
//

func (jitome *Jitome) spawn() error {
	log.Printf("starting command '%s'...", FgYB(jitome.config.Command))
	jitome.cmd = exec.Command(jitome.config.commandArgs[0], jitome.config.commandArgs[1:]...)
	jitome.cmd.Stdout = os.Stdout
	jitome.cmd.Stderr = os.Stderr
	err := jitome.cmd.Start()
	if err != nil {
		return err
	}

	if jitome.cmd.Process != nil {
		log.Printf("started command (pid: %d)", jitome.cmd.Process.Pid)
	}

	return jitome.cmd.Wait()
}

func (jitome *Jitome) terminate() error {
	if jitome.cmd != nil && jitome.cmd.Process != nil {
		pid := jitome.cmd.Process.Pid
		log.Printf("terminating command (pid: %d)...", pid)

		if err := jitome.cmd.Process.Signal(os.Interrupt); err != nil {
			log.Print(FgRB(fmt.Sprintf("[warning] %v", err)))
			return nil
		} else {
			cd := 5
			for cd > 0 {
				if jitome.cmd.ProcessState != nil && jitome.cmd.ProcessState.Exited() {
					break
				}
				time.Sleep(time.Second)
				cd--
			}
		}

		if jitome.cmd.ProcessState != nil && jitome.cmd.ProcessState.Exited() {
			jitome.cmd.Process.Kill()
		}

		log.Printf("terminated command. pid: %d", pid)
	}

	return nil
}

func runTarget(event *Event) {
	log.Printf("'%s' target detected '%s' changing by event '%s'.", FgCB(event.Watcher.Target.Name), FgYB(event.Ev.Name), FgYB(eventOpStr(&event.Ev)))

	target := event.Watcher.Target

	if runtime.GOOS == "darwin" && target.Notification {
		// desktop notification is supported only darwin.

		// tmp image file.
		tmpFile, err := ioutil.TempFile("", "jitome.icon.")
		if err != nil {
			log.Print(err)
		}
		defer func() {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
		}()

		appIcon := tmpFile.Name()

		if debug {
			log.Printf("generated icon file: %s", appIcon)
		}

		err = ioutil.WriteFile(appIcon, bindata.MustAsset("logo.png"), 0644)
		if err != nil {
			log.Print(err)
		}

		notification := gosxnotifier.NewNotification(fmt.Sprintf("'%s' target detected '%s' changing.", event.Watcher.Target.Name, event.Ev.Name))
		notification.Title = "Jitome"
		notification.Sound = gosxnotifier.Default
		notification.AppIcon = appIcon

		err = notification.Push()
		if err != nil {
			log.Print(err)
		}
	}

	path := event.Ev.Name
	script := target.Script

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

	log.Printf("'%s' target running script...", FgCB(event.Watcher.Target.Name))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Print(FgRB(fmt.Sprintf("[warning] %v", err)))
	}

	log.Printf("'%s' target finished script.", FgCB(event.Watcher.Target.Name))

	if target.Restart {
		log.Print("restarting...")
		err = target.jitome.restartCommand()
		if err != nil {
			log.Print(FgRB(fmt.Sprintf("[warning] %v", err)))
		}
	}
}

func runCommand(command string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
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
	for _, watcher := range jitome.watchers {
		watcher.w.Close()
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/deckarep/gosx-notifier"
	"github.com/fsnotify/fsnotify"
	"github.com/kohkimakimoto/jitome/bindata"
)

type Jitome struct {
	config   *Config
	watchers []*Watcher
	events   chan *Event
	cmd      *exec.Cmd
	cmds      []*exec.Cmd
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
		cmds: []*exec.Cmd{},
	}

	return w
}

func (jitome *Jitome) Start() error {
	// start commands
	err := jitome.startCommands()
	if err != nil {
		log.Print(FgRB(fmt.Sprintf("[warning] %v", err)))
	}

	// register watchers
	for name, task := range jitome.config.Tasks {
		task.Name = name
		task.events = make(chan *Event, 30)
		task.jitome = jitome

		log.Print("activating task '" + FgCB(name) + "'.")

		// register watched directories
		for i, watchConfig := range task.Watch {
			watchConfig.InitPatterns()
			w, err := fsnotify.NewWatcher()
			if err != nil {
				return err
			}

			err = watch(watchConfig.Base, watchConfig.ignoreRegs, w)
			if err != nil {
				return err
			}

			watcher, err := NewWatcher(jitome, task, watchConfig, w, i)
			if err != nil {
				return err
			}

			jitome.watchers = append(jitome.watchers, watcher)
			go watcher.Wait()
		}
		go task.Wait()
	}
	defer jitome.Close()

	go func() {
		log.Print("watching files...")
		for {
			event := <-jitome.events
			runTask(event)
		}
	}()

	// wait infinitely.
	select {}
}

func (jitome *Jitome) startCommands() error {
	config := jitome.config

	for _, c := range config.Commands {
		go func (c string) {
			log.Printf("running command '%s'...", FgYB(c))
			err := jitome.spawn(c)
			if err != nil {
				panic(err)
			}
		}(c)
	}

	return nil
}

func (jitome *Jitome) spawn(command string) (error) {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "bash"
		flag = "-c"
	}
	cmd := exec.Command(shell, flag, command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	if cmd.Process != nil {
		log.Printf("process: %s (pid: %d).", command, cmd.Process.Pid)
	}
	jitome.cmds = append(jitome.cmds, cmd)

	return cmd.Wait()
}

func (jitome *Jitome) terminate(cmd *exec.Cmd) error {
	if cmd != nil && cmd.Process != nil {
		pid := cmd.Process.Pid
		log.Printf("terminating command (pid: %d)...", pid)

		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			log.Print(FgRB(fmt.Sprintf("[warning] %v", err)))
			return nil
		} else {
			cd := 5
			for cd > 0 {
				if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
					break
				}
				time.Sleep(time.Second)
				cd--
			}
		}

		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			cmd.Process.Kill()
		}

		log.Printf("terminated command (pid: %d).", pid)
	}

	return nil
}

func runTask(event *Event) {
	log.Printf("'%s' task detected '%s' changing by event '%s'.", FgCB(event.Watcher.Task.Name), FgYB(event.Ev.Name), FgYB(eventOpStr(&event.Ev)))

	task := event.Watcher.Task

	if runtime.GOOS == "darwin" && task.Notification {
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

		notification := gosxnotifier.NewNotification(fmt.Sprintf("'%s' task detected '%s' changing.", event.Watcher.Task.Name, event.Ev.Name))
		notification.Title = "Jitome"
		notification.Sound = gosxnotifier.Default
		notification.AppIcon = appIcon

		err = notification.Push()
		if err != nil {
			log.Print(err)
		}
	}

	path := event.Ev.Name
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

	log.Printf("'%s' task running script...", FgCB(event.Watcher.Task.Name))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Print(FgRB(fmt.Sprintf("[warning] %v", err)))
	}

	log.Printf("'%s' task finished script.", FgCB(event.Watcher.Task.Name))
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

	for _, cmd := range jitome.cmds {
		jitome.terminate(cmd)
	}
}

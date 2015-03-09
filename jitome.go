package main

import (
	"github.com/codegangsta/cli"
	"github.com/go-fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"time"
)

var runCommand = cli.Command{
	Name:        "run",
	Usage:       "Runs jitome.",
	Description: "",
	Action:      doRun,
}

var initCommand = cli.Command{
	Name:        "init",
	Usage:       "Generatea an initial jitome configuration file '.jitome'.",
	Description: "",
	Action:      doInit,
}

var debug bool = false
var config *AppConfig

func main() {
	app := cli.NewApp()
	app.Name = "jitome"
	app.Usage = "Jitome runs a command when files change."
	app.Version = "0.2.0"
	app.Author = "Kohki Makimoto"
	app.Email = "kohki.makimoto@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
            Value: ".jitome",
			Usage: "Configuration file (default .jitome)",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Runs on debug mode",
		},
	}
	app.Commands = []cli.Command{
		runCommand,
		initCommand,
	}
	app.Action = doRun
	app.Run(os.Args)
}

func doInit(c *cli.Context) {
    path := c.String("config")
    if path == "" {
        path = ".jitome"
    }
    config = WriteAppConfig(path)

	printLogWithoutTimestamp("<info:bold>Generated file:</info:bold> <comment>" + path + "</comment>")
}

func doRun(c *cli.Context) {
	debug = c.Bool("debug")
	if debug {
		printDebugLog("You runs jitome on the debug mode")
	}

	path := c.String("config")
	if path == "" {
		for _, p := range []string{".jitome", ".jitome.toml", ".jigome/config"} {
			if isFile(p) {
				path = p
				break
			}
		}
		if path == "" {
			log.Fatal("Configuration File is not found: .jitome")
		}
	}

	config = NewAppConfig(path)

	printLog("<info:bold>Booted using</info:bold> <comment>" + path + "</comment>")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	events := make(chan fsnotify.Event)
	bufferedEvents := make(chan fsnotify.Event, 30)

	go runEventsRegister(events, watcher)
	go runBufferedEventsRegister(bufferedEvents, events)

	for {
		event := <-bufferedEvents
		startTaskWithPath(&event, event.Name)
	}
}

func runBufferedEventsRegister(bufferedEvents chan fsnotify.Event, events chan fsnotify.Event) {

	// event loop
	for {
		event := <-events

		buffer := []fsnotify.Event{event}
		bufferedFilesMap := map[string]int{event.Name: 1}

		timer := time.NewTimer(300 * time.Millisecond)

	outer:
		for {
			select {
			case nextEvent := <-events:
				if event.Name != nextEvent.Name {
					if _, exists := bufferedFilesMap[nextEvent.Name]; !exists {
						buffer = append(buffer, nextEvent)
						bufferedFilesMap[nextEvent.Name] = 1
					}
				}
			case <-timer.C:
				for _, be := range buffer {
					if debug {
						printDebugLog("Got a event (" + eventOpStr(&be) + ") about: " + be.Name)
					}
					bufferedEvents <- be
				}
				break outer
			}
		}
	}
}

func startTaskWithPath(event *fsnotify.Event, path string) {
	for name, task := range config.Tasks {
		if task.Match(path) {
			printLog("<info:bold>Detected changing:</info:bold> " + path)
			printLog("<info:bold>Starting:</info:bold> <comment>" + name + "</comment>")
			task.RunCommand(path)
			printLog("<info:bold>Finished:</info:bold> <comment>" + name + "</comment>")
		}
	}
}

func runEventsRegister(events chan<- fsnotify.Event, watcher *fsnotify.Watcher) {
	// Retrieves initial watching files
	watch(".", watcher)

	if debug {
		printDebugLog("Booted event loop for watching.")
	}

	// event loop
	for {
		select {
		case event := <-watcher.Events:
			path := event.Name
			if event.Op&fsnotify.Chmod != 0 {
				continue
			}

			if event.Op&fsnotify.Create != 0 && isDir(path) {
				watch(path, watcher)
			}

			events <- event

		case err := <-watcher.Errors:
			log.Fatal(err)
		}
	}
}

func watch(root string, watcher *fsnotify.Watcher) {
	if debug {
		printDebugLog("Walks watched directories: " + root)
	}

	err := filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err != nil || !fi.IsDir() {
			return nil
		}

		err = watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}

		if debug {
			printDebugLog("Added Watched dir: " + path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func eventOpStr(event *fsnotify.Event) string {
	opStr := "unknown"
	if event.Op&fsnotify.Create != 0 {
		opStr = "Create"
	} else if event.Op&fsnotify.Write != 0 {
		opStr = "Write"
	} else if event.Op&fsnotify.Remove != 0 {
		opStr = "Remove"
	} else if event.Op&fsnotify.Rename != 0 {
		opStr = "Rename"
	} else if event.Op&fsnotify.Chmod != 0 {
		opStr = "Chmod"
	}

	return opStr
}

func init() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[global options] {{end}}command{{if .Flags}} [command options]{{end}} [arguments...]

VERSION:
   {{.Version}}{{if or .Author .Email}}

AUTHOR:{{if .Author}}
   {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
   {{.Email}}{{end}}{{end}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`
}

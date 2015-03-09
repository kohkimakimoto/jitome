package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
)

type AppConfig struct {
	Tasks map[string]Task
}

type Task struct {
	Watch   interface{}
	Exclude interface{}
	Command interface{}
}

func WriteAppConfig(path string) *AppConfig {
	if isExist(path) {
		log.Fatal("'" + path + "' is already exists.")
	}

	content := []byte("[build]\n" +
		"watch=[\"*.go\"]\n" +
		"command=[\"go build\"]\n",
	)

	err := ioutil.WriteFile(path, content, os.ModePerm)
	if err != nil {
		log.Fatal("Unable create file: '" + path + "'")
	}

	return NewAppConfig(path)
}

func NewAppConfig(path string) *AppConfig {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	config := &AppConfig{}

	_, err = toml.Decode(string(content), &config.Tasks)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func (task *Task) Watches() []string {
	watches := make([]string, 0)
	t := reflect.TypeOf(task.Watch)
	if t.Kind() == reflect.String {
		watches = append(watches, task.Watch.(string))
	} else {
		for _, v := range task.Watch.([]interface{}) {
			watches = append(watches, v.(string))
		}
	}

	return watches
}

func (task *Task) Commands() []string {
	commands := make([]string, 0)
	t := reflect.TypeOf(task.Command)
	if t.Kind() == reflect.String {
		commands = append(commands, task.Command.(string))
	} else {
		for _, v := range task.Command.([]interface{}) {
			commands = append(commands, v.(string))
		}
	}

	return commands
}

func (task *Task) Match(path string) bool {
	ret := false
	for _, pattern := range task.Watches() {
		match, _ := filepath.Match(pattern, path)
		if match {
			ret = true

			if debug {
				printDebugLog("Matched '" + pattern + "' (" + path + ")")
			}

			break
		} else {
			if debug {
				printDebugLog("Unmatched '" + pattern + "' (" + path + ")")
			}
		}
	}
	return ret
}

func (task *Task) RunCommand(path string) {
	env := []string{"FILE=" + path}

	for _, cmdline := range task.Commands() {
		printLog("<info:bold>Command: </info:bold><magenta>" + cmdline + "</magenta>")
		cmd := exec.Command("sh", "-c", cmdline)
		cmd.Env = env
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}
}

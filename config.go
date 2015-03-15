package main

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
)

type AppConfig struct {
	Tasks map[string]*Task
}

type Task struct {
	Watch   interface{} `yaml:"watch"`
	Exclude interface{} `yaml:"execlude"`
	Command interface{} `yaml:"command"`
	regWatch      *regexp.Regexp
}

func WriteAppConfig(path string) *AppConfig {
	if isExist(path) {
		log.Fatal("'" + path + "' is already exists.")
	}

	var content []byte
	if IsYaml(path) {
		content = []byte("build:\n" +
			"    watch: \"*.go\"\n" +
			"    command: \"go build\"\n")
	} else {
		content = []byte("[build]\n" +
			"watch=[\"*.go\"]\n" +
			"command=[\"go build\"]\n")
	}

	err := ioutil.WriteFile(path, content, os.ModePerm)
	if err != nil {
		log.Fatal("Unable create file: '" + path + "'")
	}

	return NewAppConfig(path)
}

func NewAppConfig(path string) *AppConfig {
	if debug {
		printDebugLog("Loading config" + path)
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	config := &AppConfig{}

	if IsYaml(path) {
		err = yaml.Unmarshal(content, &config.Tasks)
	} else {
		_, err = toml.Decode(string(content), &config.Tasks)
	}

	if err != nil {
		log.Fatal(err)
	}

	for name, task := range config.Tasks {
		_ = name
		if task.Watch == nil || task.Watch == "" {
			if debug {
				printDebugLog("watch is not defined")
			}
			continue
		}


	}

	if err != nil {
		log.Fatal(err)
	}

	return config
}

func IsYaml(path string) bool {
	reg := regexp.MustCompile("\\.yml$")
	return reg.MatchString(path)
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
	for _, cmdline := range task.Commands() {
        env := append(os.Environ(), "FILE=" + path)
		printLog("<info:bold>Command: </info:bold><magenta>" + cmdline + "</magenta>")
		cmd := exec.Command("sh", "-c", cmdline)
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", cmdline)
		}
		cmd.Env = env
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}
}

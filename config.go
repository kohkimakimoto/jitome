package main

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
)

type AppConfig struct {
	Tasks map[string]*Task
}

type Task struct {
	Watch          interface{} `yaml:"watch"`
	Exclude        interface{} `yaml:"execlude"`
	Command        interface{} `yaml:"command"`
	watchStrings   []string
	watchRegexps   []*regexp.Regexp
	excludeStrings []string
	excludeRegexps []*regexp.Regexp
	commandStrings []string
}

func WriteAppConfig(path string) *AppConfig {
	if isExist(path) {
		log.Fatal("'" + path + "' is already exists.")
	}

	var content []byte
	if IsYaml(path) {
		content = []byte("build:\n" +
			"    watch: '.+\\.go$'" +
			"    command: 'go build'\n")
	} else {
		content = []byte("[build]\n" +
			"watch=['.+\\.go$']\n" +
			"command=['go build']\n")
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

	for _, task := range config.Tasks {
		// load watch
		if task.Watch != nil && task.Watch != "" {
			watches := make([]string, 0)
			t := reflect.TypeOf(task.Watch)
			if t.Kind() == reflect.String {
				watches = append(watches, task.Watch.(string))
			} else {
				for _, v := range task.Watch.([]interface{}) {
					watches = append(watches, v.(string))
				}
			}
			task.watchStrings = watches

			for _, pattern := range task.watchStrings {
				reg, err := regexp.Compile(pattern)
				if err != nil {
					log.Fatal(err)
					continue
				}
				task.watchRegexps = append(task.watchRegexps, reg)
			}

		}

		// load command
		if task.Command != nil && task.Command != "" {
			commands := make([]string, 0)
			t := reflect.TypeOf(task.Command)
			if t.Kind() == reflect.String {
				commands = append(commands, task.Command.(string))
			} else {
				for _, v := range task.Command.([]interface{}) {
					commands = append(commands, v.(string))
				}
			}
			task.commandStrings = commands
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

func (task *Task) Match(path string) bool {
	ret := false
	for i, reg := range task.watchRegexps {
		if reg != nil && reg.MatchString(path) {
			ret = true

			if debug {
				printDebugLog("Matched '" + task.watchStrings[i] + "' (" + path + ")")
			}

			break
		} else {
			if debug {
				printDebugLog("Unmatched '" + task.watchStrings[i] + "' (" + path + ")")
			}
		}
	}
	return ret
}

func (task *Task) RunCommandWithPath(path string) {
	for _, cmdline := range task.commandStrings {
		//env := append(os.Environ(), "FILE="+path)
		cmdline = os.Expand(cmdline, func(s string) string {
			switch s {
			case "JITOME_FILE":
				return path
			}
			return os.Getenv(s)
		})

		printLog("<info:bold>Command: </info:bold><magenta>" + cmdline + "</magenta>")
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", cmdline)
		} else {
			cmd = exec.Command("sh", "-c", cmdline)
		}
		//		cmd.Env = env
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}
}

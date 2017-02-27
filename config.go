package main

type Config struct {
	Commands    []string         `yaml:"commands"`
	Tasks       map[string]*Task `yaml:"tasks"`
}

func NewConfig() *Config {
	return &Config{
		Tasks:       map[string]*Task{},
	}
}

var initialConfig = `# Jitome is a simple file watcher. - https://github.com/kohkimakimoto/jitome
# commands.
# commands:
#   - echo foo

# tasks.
tasks:
  build:
    notification: true
    watch:
      - base: ""
        ignore: ["\.git$"]
        pattern: "\.go$"
    script: |
      go build .

`

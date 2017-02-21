package main

type Config struct {
	Command     string           `yaml:"command"`
	Commands    []string         `yaml:"commands"`
	Tasks       map[string]*Task `yaml:"tasks"`
	commandArgs []string
}

func NewConfig() *Config {
	return &Config{
		Tasks:       map[string]*Task{},
		commandArgs: []string{},
	}
}

var initialConfig = `# Jitome is a simple file watcher. - https://github.com/kohkimakimoto/jitome

# tasks.
tasks:
  build:
    notification: true
    watch:
      - base: ""
        ignore: [".git"]
        pattern: "*.go"
    script: |
      go build .

`

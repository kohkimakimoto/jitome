package main

type Config struct {
	Command     string             `yaml:"command"`
	Targets     map[string]*Target `yaml:"targets"`
	commandArgs []string
}

func NewConfig() *Config {
	return &Config{
		Targets:     map[string]*Target{},
		commandArgs: []string{},
	}
}

var initialConfig = `# Jitome is a simple file watcher. - https://github.com/kohkimakimoto/jitome
# command (optional).
# command: "your/server/start/command"

# targets.
targets:
  build:
    notification: true
    restart: true
    watch:
      - base: ""
        ignore: [".git"]
        pattern: "*.go"
    script: |
      go build .

`

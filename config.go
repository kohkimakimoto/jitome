package main

type Config struct {
	Targets map[string]*Target
}

var initialConfig = `# Jitome is a simple file watcher. - https://github.com/kohkimakimoto/jitome

# init is a special target that runs when jitome starts.
#init:
#  script: |
#    echo "current directory is $(pwd)"

# targets.
build:
  notification: true
  watch:
    - base: ""
      ignore: [".git"]
      pattern: "*.go"
  script: |
    go build .
`

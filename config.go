package main

type Config struct {
	Tasks map[string]*Task
}

var initialConfig = `build:
  watch:
    - base: ""
      ignore_dir: [".git"]
      pattern: '.+\.go$'
  script: |
    cat $JITOME_FILE
    go test .
`

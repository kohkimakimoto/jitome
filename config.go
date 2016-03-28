package main

type Config struct {
	Tasks map[string]*Task
}

var initialConfig = `build:
  watch:
    - base: "."
      pattern: "*.go"
  script: |
    go test ./...
    go build
`
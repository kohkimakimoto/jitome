package main

type Config struct {
	Targets map[string]*Target
}

var initialConfig = `build:
  watch:
    - base: ""
      ignore: [".git"]
      pattern: "*.go"
  script: |
    go build .
`

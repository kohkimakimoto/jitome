package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(StderrWriter, FgRB("[error] %v\n"), err)
			os.Exit(1)
		}
	}()

	os.Exit(realMain())
}

var debug bool
var configFile string

func realMain() int {
	log.SetPrefix("")
	log.SetOutput(StdoutWriter)
	log.SetFlags(0)

	var initFlag bool

	flag.StringVar(&configFile, "c", "jitome.yml", "")
	flag.StringVar(&configFile, "config", "jitome.yml", "")
	flag.BoolVar(&initFlag, "i", false, "")
	flag.BoolVar(&initFlag, "init", false, "")
	flag.BoolVar(&debug, "d", false, "")
	flag.BoolVar(&debug, "debug", false, "")

	flag.Usage = printUsage
	flag.Parse()

	if initFlag {
		err := generateConfig()
		if err != nil {
			panic(err)
		}

		log.Printf("Created '%s'", configFile)
		return 0
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		panic(fmt.Errorf("'%s' is not found. If you want to get help, run 'jitome -h'.", configFile))
	}

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	config := &Config{}
	err = yaml.Unmarshal(b, &config.Tasks)
	if err != nil {
		panic(err)
	}

	j := NewJitome(config)
	err = j.Start()
	if err != nil {
		panic(err)
	}

	return 0
}

func generateConfig() error {
	if _, err := os.Stat(configFile); err == nil {
		return fmt.Errorf("%s is already existed.", configFile)
	}

	err := ioutil.WriteFile(configFile, []byte(`build:
  watch:
    - base: "."
      pattern: "*.go"
  code: |
    go test ./...
    go build
`), 0644)

	if err != nil {
		return err
	}

	return nil
}

func printUsage() {
	fmt.Println(`Usage: jitome [<options>]

Jitome is a watcher for file changing.

Options:
  -c|-config    Specify a config file.
  -i|-init      Create initial config file.
  -d|-debug     Run on debug mode.
`)
}

package main

import (
	"flag"
	"fmt"
	"github.com/kballard/go-shellquote"
	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, FgRB("jitome error: %v\n", err))
			os.Exit(1)
		}
	}()

	os.Exit(realMain())
}

var debug bool
var configFile string

func realMain() int {
	log.SetPrefix(FgGB("[jitome] "))
	log.SetOutput(os.Stdout)
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

		log.Printf("created '%s'", configFile)
		return 0
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		panic(fmt.Errorf("'%s' is not found. If you want to get help, run 'jitome -h'.", configFile))
	}

	log.Print("starting jitome...")
	log.Printf("loading config '%s'", FgYB(configFile))

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	config := NewConfig()
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}

	if config.Command != "" {
		args, err := shellwords.Parse(config.Command)
		if err != nil {
			panic(err)
		}
		config.commandArgs = args
	}

	if nargs := flag.NArg(); nargs > 0 {
		config.Command = shellquote.Join(flag.Args()...)
		config.commandArgs = flag.Args()
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

	err := ioutil.WriteFile(configFile, []byte(initialConfig), 0644)
	if err != nil {
		return err
	}

	return nil
}

func printUsage() {
	fmt.Println(`Usage: jitome [<options>] [<command>]

  Jitome is a simple file watcher.

Options:
  -c, -config    Load configuration from a file. default 'jitome.yml'
  -i, -init      Create initial config file.
  -d, -debug     Run on debug mode.
  -h, -help      Show help.
`)
}

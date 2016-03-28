package main

type WatchConfig struct {
	Base string `yaml:"base"`
	// string or []string
	IgnoreDir interface{} `yaml:"ignore_dir"`
	// string or []string
	Pattern interface{} `yaml:"pattern"`
}

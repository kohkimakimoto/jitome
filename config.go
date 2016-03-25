package main

type Config struct {
	Tasks map[string]*Task
}

type Task struct {
	Name string `yaml:"-"`
	Watch []*WatchConfig `yaml:"watch"`
	Code  string         `yaml:"code"`
}

type WatchConfig struct {
	Base      string      `yaml:"base"`
	IgnoreDir interface{} `yaml:"ignore_dir"`
	Pattern   interface{} `yaml:"pattern"`
}

package main

import (
	"testing"
)

func TestIsYaml(t *testing.T) {
	ret := IsYaml("tests/jitome.test.toml")
	if ret {
		t.Error("Wrong check. This is not a yaml file")
	}
	ret = IsYaml("tests/jitome.test.yaml.aaa")
	if ret {
		t.Error("Wrong check. This is not a yaml file")
	}
	ret = IsYaml("tests/jitome.test.yml")
	if !ret {
		t.Error("Wrong check. This is a yaml file")
	}
	ret = IsYaml("jitome.test.yml")
	if !ret {
		t.Error("Wrong check. This is a yaml file")
	}
	ret = IsYaml(".jitome.yml")
	if !ret {
		t.Error("Wrong check. This is a yaml file")
	}
}

func TestNewAppConfig(t *testing.T) {

	// toml
	config := NewAppConfig("tests/jitome.test.toml")
	if config == nil {
		t.Error("Can not load a file.")
	}
	n := len(config.Tasks)
	if n != 2 {
		t.Error("Parse error. unmatch number of tasks.")
	}

	// yaml
	config = NewAppConfig("tests/jitome.test.yml")
	if config == nil {
		t.Error("Can not load a file.")
	}
	n = len(config.Tasks)
	if n != 2 {
		t.Error("Parse error. unmatch number of tasks.")
	}
}

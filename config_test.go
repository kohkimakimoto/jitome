package main

import (
	"testing"
)

func TestNewAppConfig(t *testing.T) {

	config := NewAppConfig("tests/jitome.test.yml")
	if config == nil {
		t.Error("aaa")
	}

}

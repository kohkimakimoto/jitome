package main

import (
	"fmt"
	"reflect"
	"regexp"
)

type WatchConfig struct {
	Base string `yaml:"base"`
	// string or []string
	Ignore         interface{}      `yaml:"ignore"`
	IgnorePatterns []*regexp.Regexp `yaml:"-"`
	// string or []string
	Pattern  interface{}      `yaml:"pattern"`
	Patterns []*regexp.Regexp `yaml:"-"`
}

func (wc *WatchConfig) InitPatterns() {
	wc.IgnorePatterns = []*regexp.Regexp{}

	if wc.Ignore != nil {
		if e, ok := wc.Ignore.(string); ok {
			p := compilePattern(e)
			wc.IgnorePatterns = append(wc.IgnorePatterns, p)
		} else if e, ok := wc.Ignore.([]interface{}); ok {
			for _, i := range e {
				p := compilePattern(i.(string))
				wc.IgnorePatterns = append(wc.IgnorePatterns, p)
			}
		} else {
			v := reflect.ValueOf(wc.Ignore)
			panic(fmt.Errorf("invalid format ignore: %v", v.Type()))
		}
	}

	wc.Patterns = []*regexp.Regexp{}
	if wc.Pattern != nil {
		if patternStr, ok := wc.Pattern.(string); ok {
			p := compilePattern(patternStr)
			wc.Patterns = append(wc.Patterns, p)
		} else if e, ok := wc.Pattern.([]interface{}); ok {
			for _, patternStr := range e {
				p := compilePattern(patternStr.(string))
				wc.Patterns = append(wc.Patterns, p)
			}
		} else {
			v := reflect.ValueOf(wc.Pattern)
			panic(fmt.Errorf("invalid format pattern: %v", v.Type()))
		}
	}
}

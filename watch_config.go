package main

import (
	"fmt"
	"reflect"
	"regexp"
)

type WatchConfig struct {
	Base string `yaml:"base"`
	// string or []string
	Ignore     interface{}      `yaml:"ignore"`
	ignoreRegs []*regexp.Regexp `yaml:"-"`
	// string or []string
	Pattern     interface{}      `yaml:"pattern"`
	PatternRegs []*regexp.Regexp `yaml:"-"`
}

func (wc *WatchConfig) InitPatterns() {
	wc.ignoreRegs = []*regexp.Regexp{}

	if wc.Ignore != nil {
		if e, ok := wc.Ignore.(string); ok {
			p := compilePattern(e)
			wc.ignoreRegs = append(wc.ignoreRegs, p)
		} else if e, ok := wc.Ignore.([]interface{}); ok {
			for _, i := range e {
				p := compilePattern(i.(string))
				wc.ignoreRegs = append(wc.ignoreRegs, p)
			}
		} else {
			v := reflect.ValueOf(wc.Ignore)
			panic(fmt.Errorf("invalid format ignore: %v", v.Type()))
		}
	}

	wc.PatternRegs = []*regexp.Regexp{}
	if wc.Pattern != nil {
		if patternStr, ok := wc.Pattern.(string); ok {
			p := compilePattern(patternStr)
			wc.PatternRegs = append(wc.PatternRegs, p)
		} else if e, ok := wc.Pattern.([]interface{}); ok {
			for _, patternStr := range e {
				p := compilePattern(patternStr.(string))
				wc.PatternRegs = append(wc.PatternRegs, p)
			}
		} else {
			v := reflect.ValueOf(wc.Pattern)
			panic(fmt.Errorf("invalid format pattern: %v", v.Type()))
		}
	}
}

package main

import (
	"fmt"
	"testing"
)



func TestCompliePattern(t *testing.T) {
	p := compilePattern(".git")
	fmt.Printf("regex: %s (from '.git')\n", p.String())
	if !p.MatchString(".git") {
		t.Error("should be matched")
	}
	if !p.MatchString("/.git") {
		t.Error("should be matched")
	}
	if !p.MatchString("/.git/aaaa") {
		t.Error("should be matched")
	}

	p = compilePattern("*.go")
	fmt.Printf("regex: %s (from '*.go')\n", p.String())
	if !p.MatchString("aaa.go") {
		t.Error("should be matched")
	}
	if !p.MatchString("aa/aaa.go") {
		t.Error("should be matched")
	}
	if p.MatchString("aaa.goaa") {
		t.Error("should not be matched")
	}

	p = compilePattern("hoge.go")
	fmt.Printf("regex: %s (from 'hoge.go')\n", p.String())
	if !p.MatchString("hoge.go") {
		t.Error("should be matched")
	}
	if !p.MatchString("aaa/hoge.go") {
		t.Error("should be matched")
	}
	if p.MatchString("aaahoge.go") {
		t.Error("should not be matched")
	}

	p = compilePattern("/hoge.go")
	fmt.Printf("regex: %s (from '/hoge.go')\n", p.String())
	if !p.MatchString("/hoge.go") {
		t.Error("should be matched")
	}
	if !p.MatchString("hoge.go") {
		t.Error("should not be matched")
	}

}

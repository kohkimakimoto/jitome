package main

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"regexp"
)

var regexpForNormalizing = regexp.MustCompile("^\\./")

func normalizePath(path string) string {
	path = filepath.ToSlash(path)
	// remove "./"
	return regexpForNormalizing.ReplaceAllString(path, "")
}

func isDir(path string) (ret bool) {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func eventOpStr(event *fsnotify.Event) string {
	opStr := "unknown"
	if event.Op&fsnotify.Create != 0 {
		opStr = "create"
	} else if event.Op&fsnotify.Write != 0 {
		opStr = "write"
	} else if event.Op&fsnotify.Remove != 0 {
		opStr = "remove"
	} else if event.Op&fsnotify.Rename != 0 {
		opStr = "rename"
	} else if event.Op&fsnotify.Chmod != 0 {
		opStr = "chmod"
	}

	return opStr
}

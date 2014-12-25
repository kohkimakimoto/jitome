package main

import (
	"os"
)

func isExist(path string) (ret bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ret = false
	} else {
		ret = true
	}

	return ret
}

//
// IsNotExist returns true if the path doesn't exists.
//
func isNotExist(path string) (ret bool) {
	return !isExist(path)
}

//
// IsDir returns true if the path is a directory.
//
func isDir(path string) (ret bool) {
	if isNotExist(path) {
		return false
	}

	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

//
// IsDir returns true if the path is a file.
//
func isFile(path string) (ret bool) {
	if isNotExist(path) {
		return false
	}

	return !isDir(path)
}

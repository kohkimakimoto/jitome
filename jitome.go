package main

import (
	"os"
)

func main() {
	defer func() {
	}()

	os.Exit(realMain())
}

func realMain() int {

	return 0
}

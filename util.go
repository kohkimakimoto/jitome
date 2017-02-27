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

func compilePattern(line string) *regexp.Regexp {
	if line == "" {
		return nil
	}

	return regexp.MustCompile(line)
}

//func compilePattern(pattern string) (*regexp.Regexp, error) {
//	if pattern[0] == '%' {
//		return regexp.Compile(pattern[1:])
//	}
//
//	var buf bytes.Buffer
//
//	for n, pat := range strings.Split(pattern, "|") {
//		if n == 0 {
//			buf.WriteString("^")
//		} else {
//			buf.WriteString("$|")
//		}
//		if fs, err := filepath.Abs(pat); err == nil {
//			pat = filepath.ToSlash(fs)
//		}
//		rs := []rune(pat)
//		for i := 0; i < len(rs); i++ {
//			if rs[i] == '/' {
//				if runtime.GOOS == "windows" {
//					buf.WriteString(`[/\\]`)
//				} else {
//					buf.WriteRune(rs[i])
//				}
//			} else if rs[i] == '*' {
//				if i < len(rs)-1 && rs[i+1] == '*' {
//					i++
//					if i < len(rs)-1 && rs[i+1] == '/' {
//						i++
//						buf.WriteString(`.*`)
//					} else {
//						return nil, fmt.Errorf("invalid wildcard: %s", pattern)
//					}
//				} else {
//					buf.WriteString(`[^/]+`)
//				}
//			} else if rs[i] == '?' {
//				buf.WriteString(`\S`)
//			} else {
//				buf.WriteString(fmt.Sprintf(`[\x%x]`, rs[i]))
//			}
//		}
//		buf.WriteString("$")
//	}
//
//	return regexp.Compile(buf.String())
//}

package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"regexp"
	"time"
)

func printColored(s interface{}) {
	as, b := s.(string)
	if !b {
		fmt.Print(s)
		return
	}

	reg := regexp.MustCompile("(<[^>]+>)?([^<]*)(</[^>]+>)?")
	matches := reg.FindAllStringSubmatch(as, -1)

	for _, v := range matches {
		if v[1] == "<info>" {
			ct.ChangeColor(ct.Green, false, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else if v[1] == "<info:bold>" {
			ct.ChangeColor(ct.Green, true, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else if v[1] == "<green>" {
			ct.ChangeColor(ct.Green, false, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else if v[1] == "<green:bold>" {
			ct.ChangeColor(ct.Green, true, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else if v[1] == "<comment>" {
			ct.ChangeColor(ct.Yellow, false, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else if v[1] == "<comment:bold>" {
			ct.ChangeColor(ct.Yellow, true, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else if v[1] == "<error>" {
			ct.ChangeColor(ct.Red, false, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else if v[1] == "<error:bold>" {
			ct.ChangeColor(ct.Red, true, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else if v[1] == "<blue>" {
			ct.ChangeColor(ct.Blue, false, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else if v[1] == "<blue:bold>" {
			ct.ChangeColor(ct.Blue, true, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else if v[1] == "<cyan>" {
			ct.ChangeColor(ct.Cyan, false, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else if v[1] == "<cyan:bold>" {
			ct.ChangeColor(ct.Cyan, true, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else if v[1] == "<magenta>" {
			ct.ChangeColor(ct.Magenta, false, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else if v[1] == "<magenta:bold>" {
			ct.ChangeColor(ct.Magenta, true, ct.None, false)
			fmt.Print(v[2])
			ct.ResetColor()
		} else {
			fmt.Print(v[0])
		}
	}
}

func printColoredln(s interface{}) {
	printColored(s)
	fmt.Println("")
}

func printLog(s interface{}) {
	time := time.Now()
	printColoredln("<green:bold>[</green:bold><cyan:bold>" + time.Format("2006-01-02T15:04:05Z07:00") + "</cyan:bold><green:bold>]</green:bold> " + s.(string))
}

func printDebugLog(s interface{}) {
	time := time.Now()
	printColoredln("[" + time.Format("2006-01-02T15:04:05Z07:00") + "] " + s.(string))
}

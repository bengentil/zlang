// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import (
	"fmt"
	"github.com/axw/gollvm/llvm"
	"os"
	"strings"
)

const (
	ResetCode  = "\033[0m" // reset color
	RedCode    = "\033[31m"
	GreenCode  = "\033[32m"
	YellowCode = "\033[33m"
)

var DEBUG = false

func EnableDebug() {
	DEBUG = true
}

func DisableDebug() {
	DEBUG = false
}

func Error(f string, v ...interface{}) {
	format := RedCode + f + ResetCode
	fmt.Fprintf(os.Stderr, format, v...)
}

func Errorln(f string, v ...interface{}) {
	Error(f+"\n", v...)
}

func Debug(f string, v ...interface{}) {
	if DEBUG {
		fmt.Fprintf(os.Stderr, f, v...)
	}
}

func DebugDump(v *llvm.Value) {
	if DEBUG {
		EnableColor()
		if v != nil {
			v.Dump()
		} else {
			Debug("<nil>\n")
		}
		DisableColor()
	}
}

func DebugDumpMod(m *llvm.Module) {
	if DEBUG {
		EnableColor()
		if m != nil {
			m.Dump()
		} else {
			Debug("<nil>\n")
		}
		DisableColor()
	}
}

func EnableColor() {
	fmt.Fprintf(os.Stderr, YellowCode)
}

func DisableColor() {
	fmt.Fprintf(os.Stderr, ResetCode)
}

// replace <nil> by null to have valid JSON
// fix also array prepresentation by separating object with , instead of space
func JNil(s string) string {
	return strings.Replace(strings.Replace(s, "<nil>", "null", -1), "} {\"__type", "},{\"__type", -1)
}

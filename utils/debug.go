package utils

import (
	"fmt"
	"os"
)

var Verbose bool

func DebugPrintln(message string) {
	if Verbose {
		StdErrorPrintln(message)
	}
}

func StdErrorPrintln(message string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, message+"\n", a...)
}

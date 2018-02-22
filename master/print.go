// +build DEBUG

package main

import (
	"fmt"
	"os"
)

//func output() {
//fmt.Fprintf(os.Stdout, "test")
//}

func outputln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stdout, a...)
}

func outputf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stdout, format, a...)
}

// +build DEBUG

package shared

import (
	"fmt"
	"os"
)

func Outputln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stdout, a...)
}

func Outputf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stdout, format, a...)
}

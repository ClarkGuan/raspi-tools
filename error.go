package rtools

import (
	"fmt"
	"os"
)

func Fail(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

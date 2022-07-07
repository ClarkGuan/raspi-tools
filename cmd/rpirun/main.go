package main

import (
	"fmt"
	"os"

	rtools "github.com/ClarkGuan/raspi-tools"
)

func main() {
	if len(os.Args) < 2 {
		rtools.Fail("not enough arguments")
	}

	target := os.Args[1]
	args := os.Args[2:]
	if err := rtools.RunCmd(target, args...); err != nil {
		_, _ = fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
}

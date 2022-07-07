package main

import (
	"os"
	"os/exec"

	rtools "github.com/ClarkGuan/raspi-tools"
)

func main() {
	if len(os.Args) < 2 {
		rtools.Fail("not enough arguments")
	}

	cmdName := os.Args[1]
	args := os.Args[2:]
	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = collectEnvs()
	if err := cmd.Run(); err != nil {
		os.Exit(cmd.ProcessState.ExitCode())
	}
}

func collectEnvs() []string {
	allEnvs := append(os.Environ(),
		"GOOS=linux",
		"GOARCH=arm64",
		"CC="+rtools.CC(),
		"CXX="+rtools.CXX(),
		"CGO_ENABLED=1")

	return allEnvs
}

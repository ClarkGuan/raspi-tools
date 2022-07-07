package main

import (
	"os"
	"os/exec"

	rtools "github.com/ClarkGuan/raspi-tools"
)

const (
	targetName    = "aarch64-unknown-linux-musl"
	targetEnvName = "AARCH64_UNKNOWN_LINUX_MUSL"
)

func main() {
	if len(os.Args) < 3 {
		rtools.Fail("not enough arguments")
	}
	args := os.Args[2:]
	cmd := exec.Command("cargo", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = collectEnvs()
	if err := cmd.Run(); err != nil {
		os.Exit(cmd.ProcessState.ExitCode())
	}
}

func collectEnvs() []string {
	allEnvs := append(os.Environ(),
		"CARGO_TARGET_"+targetEnvName+"_LINKER="+rtools.CC(),
		"CARGO_TARGET_"+targetEnvName+"_RUNNER=rpirun",
		"CARGO_BUILD_TARGET="+targetName)
	return allEnvs
}

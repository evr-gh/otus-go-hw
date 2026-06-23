package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cntx := context.Background()
	//nolint:gosec
	ec := exec.CommandContext(cntx, cmd[0], cmd[1:]...)

	ec.Stdout = os.Stdout
	ec.Stderr = os.Stderr
	ec.Stdin = os.Stdin

	ec.Env = []string{}

	for _, val := range os.Environ() {
		parts := strings.Split(val, "=")
		if _, ok := env[parts[0]]; !ok {
			ec.Env = append(ec.Env, val)
		}
	}

	for name, ev := range env {
		switch {
		case (name == "") || (strings.Contains(name, "=")):
			fmt.Printf("Error executing command: invalid env variable name (%s)\n", name)
			return -1
		case (!ev.NeedRemove):
			ec.Env = append(ec.Env, name+"="+ev.Value)
		}
	}

	err := ec.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		fmt.Println("Error executing command:", err)
		return -1
	}
	return ec.ProcessState.ExitCode()
}

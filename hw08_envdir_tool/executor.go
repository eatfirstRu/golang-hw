package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for k, v := range env {
		if v.NeedRemove {
			err := os.Unsetenv(k)
			if err != nil {
				return 100
			}
		} else {
			os.Setenv(k, v.Value)
		}
	}

	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin

	err := command.Run()
	if err != nil {
		// fmt.Printf("command.Run() has error: %T, %v", err, err).
		return 101
	}

	return 0
}

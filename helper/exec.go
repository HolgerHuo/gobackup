package helper

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	spaceRegexp = regexp.MustCompile("[\\s]+")
)

// Exec cli commands
func Exec(command string, args ...string) (output string, err error) {
	commands := spaceRegexp.Split(command, -1)
	command = commands[0]
	commandArgs := []string{}
	if len(commands) > 1 {
		commandArgs = commands[1:]
	}
	if len(args) > 0 {
		commandArgs = append(commandArgs, args...)
	}

	fullCommand, err := exec.LookPath(command)
	if err != nil {
		return "", fmt.Errorf("%s cannot be found", command)
	}

	cmd := exec.Command(fullCommand, commandArgs...)
	cmd.Env = os.Environ()

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	// Log command at debug level
	slog.Debug("Executing command",
		"component", "exec",
		"command", fullCommand,
		"args", strings.Join(commandArgs, " "))

	out, err := cmd.Output()
	if err != nil {
		slog.Debug("Command execution failed",
			"component", "exec",
			"command", fullCommand,
			"args", strings.Join(commandArgs, " "),
			"error", stdErr.String())
		err = errors.New(stdErr.String())
		return
	}

	output = strings.Trim(string(out), "\n")
	return
}

// ExecWithCustomEnv executes a command with additional environment variables
func ExecWithCustomEnv(command string, envVars []string, args ...string) (output string, err error) {
	commands := spaceRegexp.Split(command, -1)
	command = commands[0]
	commandArgs := []string{}
	if len(commands) > 1 {
		commandArgs = commands[1:]
	}
	if len(args) > 0 {
		commandArgs = append(commandArgs, args...)
	}

	fullCommand, err := exec.LookPath(command)
	if err != nil {
		return "", fmt.Errorf("%s cannot be found", command)
	}

	cmd := exec.Command(fullCommand, commandArgs...)
	cmd.Env = append(os.Environ(), envVars...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	// Log command at debug level
	slog.Debug("Executing command with custom env", 
		"component", "exec",
		"command", fullCommand,
		"args", strings.Join(commandArgs, " "))

	out, err := cmd.Output()
	if err != nil {
		slog.Debug("Command execution failed", 
			"component", "exec",
			"command", fullCommand,
			"args", strings.Join(commandArgs, " "),
			"error", stdErr.String())
		err = errors.New(stdErr.String())
		return
	}

	output = strings.Trim(string(out), "\n")
	return
}

func ExecWithStdio(command string, stdout bool, args ...string) (output string, err error) {
	commands := spaceRegexp.Split(command, -1)
	command = commands[0]
	commandArgs := []string{}
	if len(commands) > 1 {
		commandArgs = commands[1:]
	}
	if len(args) > 0 {
		commandArgs = append(commandArgs, args...)
	}

	fullCommand, err := exec.LookPath(command)
	if err != nil {
		return "", fmt.Errorf("%s cannot be found", command)
	}

	cmd := exec.Command(fullCommand, commandArgs...)
	cmd.Env = os.Environ()

	var stdErr bytes.Buffer
	var stdOut bytes.Buffer
	cmd.Stderr = &stdErr

	if stdout {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = &stdOut
	}

	err = cmd.Run()
	if err != nil {
		slog.Debug("Command execution failed",
			"component", "exec",
			"command", fullCommand,
			"args", strings.Join(commandArgs, " "),
			"error", stdErr.String())
		err = errors.New(stdErr.String())
	}
	output = strings.Trim(stdOut.String(), "\n")

	return
}

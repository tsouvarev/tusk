package task

import (
	"os"
	"os/exec"

	"github.com/rliebz/tusk/config/marshal"
	"github.com/rliebz/tusk/ui"
)

const (
	shellEnvVar  = "SHELL"
	defaultShell = "sh"
)

// Command is a command passed to the shell.
type Command struct {
	Do    string `yaml:"do"`
	Print string `yaml:"print"`
}

// UnmarshalYAML allows strings to be interpreted as Do actions.
func (c *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var do string
	doCandidate := marshal.UnmarshalCandidate{
		Unmarshal: func() error { return unmarshal(&do) },
		Assign: func() {
			*c = Command{
				Do:    do,
				Print: do,
			}
		},
	}

	type commandType Command // Use new type to avoid recursion
	var commandItem commandType
	commandCandidate := marshal.UnmarshalCandidate{
		Unmarshal: func() error { return unmarshal(&commandItem) },
		Assign: func() {
			*c = Command(commandItem)
			if c.Print == "" {
				c.Print = c.Do
			}
		},
	}

	return marshal.UnmarshalOneOf(doCandidate, commandCandidate)
}

// CommandList is a list of commands with custom yaml unamrshaling.
type CommandList []Command

// UnmarshalYAML allows single items to be used as lists.
func (cl *CommandList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var commandSlice []Command
	sliceCandidate := marshal.UnmarshalCandidate{
		Unmarshal: func() error { return unmarshal(&commandSlice) },
		Assign:    func() { *cl = commandSlice },
	}

	var commandItem Command
	itemCandidate := marshal.UnmarshalCandidate{
		Unmarshal: func() error { return unmarshal(&commandItem) },
		Assign:    func() { *cl = CommandList{commandItem} },
	}

	return marshal.UnmarshalOneOf(sliceCandidate, itemCandidate)
}

// execCommand executes a shell command.
func execCommand(command string) error {
	shell := getShell()
	cmd := exec.Command(shell, "-c", command) // nolint: gosec
	cmd.Stdin = os.Stdin
	if ui.Verbosity > ui.VerbosityLevelSilent {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}

// getShell returns the value of the `SHELL` environment variable, or `sh`.
func getShell() string {
	if shell := os.Getenv(shellEnvVar); shell != "" {
		return shell
	}

	return defaultShell
}
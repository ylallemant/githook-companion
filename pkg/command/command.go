package command

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func New(binary string) *command {
	instance := new(command)

	instance.binary = binary

	return instance
}

var (
	_ Command = &command{}
)

type command struct {
	binary    string
	arguments []string
}

func (i *command) Execute() (string, error) {
	cmd := exec.Command(i.binary, i.arguments...)
	stdout, err := cmd.Output()

	if err != nil {
		return "", errors.Wrapf(err, "command execution failed")
	}

	return strings.Trim(string(stdout), "\n"), nil
}

func (i *command) AddArg(argument string) {
	i.arguments = append(i.arguments, argument)
}

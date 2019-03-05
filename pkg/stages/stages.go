package stages

import (
	"bytes"
	"os/exec"

	"github.com/puppetlabs/nebula/pkg/errors"
)

// Stages

type Driver interface {
	Name() string
	Execute(string) errors.Error
}

type TriggerType int

const (
	StageTrigger TriggerType = 0
	UserTrigger  TriggerType = 1
)

type Trigger interface {
	Name() string
	Type() TriggerType
}

type Stage interface {
	SetDriver() errors.Error
	Execute() errors.Error
	Next() Stage
	SetTrigger(Trigger)
	Trigger() Trigger
}

type ShellDriver struct {
	name string
}

func NewShellDriver() *ShellDriver {
	return &ShellDriver{}
}

func (d *ShellDriver) Execute(cmdString string) errors.Error {
	cmd := exec.Command(cmdString)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return errors.NewCommandUnknownCommandExecutionError(err.Error(), cmdString)
	}

	return nil
}

func (d *ShellDriver) Name() string {
	return "ShellDriver"
}

type Command interface {
	SetDriver(d Driver)
	SetCommand(string)
	Execute() errors.Error
}

type ShellCommand struct { // implements Command interface
	cmd    string
	driver Driver
}

func NewShellCommand(d Driver, cmd string) *ShellCommand {
	shellCommand := &ShellCommand{}

	shellCommand.SetDriver(d)
	shellCommand.SetCommand(cmd)

	return shellCommand
}

func (sc *ShellCommand) SetDriver(d Driver) {
	sc.driver = d
}

func (sc *ShellCommand) SetCommand(cmd string) {
	sc.cmd = cmd
}

func (sc *ShellCommand) Execute() errors.Error {
	if sc.cmd == "" {
		return errors.NewCommandNoCommandToExecuteError()
	}

	// execute it

	err := sc.driver.Execute(sc.cmd)

	if err != nil {
		return err
	}

	return nil
}

type CommandStage struct { // implements Stage interface
	commands []Command
	driver   Driver
}

func NewCommandStage(d Driver, cmds []Command) *CommandStage {
	cs := &CommandStage{}

	cs.SetDriver(d)

	return cs
}

func (cs *CommandStage) SetCommands(cmds []string) {
	for _, cmd := range cmds {
		command := &ShellCommand{} // default to ShellCommand

		command.SetCommand(cmd)
		command.SetDriver(cs.driver)

		cs.commands = append(cs.commands, command)
	}
}

func (cs *CommandStage) SetDriver(d Driver) {
	cs.driver = d
}

func (cs *CommandStage) Execute() errors.Error {
	for _, cmd := range cs.commands {
		err := cmd.Execute()

		if err != nil {
			return err
		}
	}

	return nil
}

package execution

import (
	"bytes"
	"html/template"
	"os/exec"

	logging "github.com/puppetlabs/insights-logging"
	"github.com/puppetlabs/nebula/pkg/errors"
)

func interfaceSlice(strSlice []string) []interface{} {
	interfaceSlice := make([]interface{}, 0)

	for _, s := range strSlice {
		interfaceSlice = append(interfaceSlice, s)
	}

	return interfaceSlice
}

func ExecuteCommand(rawCommand string, variables map[string]string, l logging.Logger) errors.Error {
	var deferr error
	var err errors.Error

	tmpl, deferr := template.New("shell-command").Parse(rawCommand)

	if err != nil {
		return errors.NewExecutionInvalidShellTemplateError(rawCommand).WithCause(deferr)
	}

	var templateOutput bytes.Buffer

	deferr = tmpl.Execute(&templateOutput, variables)

	if err != nil {
		return errors.NewExecutionShellTemplateExecutionError(rawCommand).WithCause(deferr)
	}

	processedRawCommand := templateOutput.String()

	l.Warn("$ " + processedRawCommand)

	commands, err := parseCommandArguments(processedRawCommand)

	if err != nil {
		return err
	}

	command := commands[0]
	args := commands[1:]

	cmd := exec.Command(command, args...)

	var outBuff bytes.Buffer
	var errBuff bytes.Buffer

	cmd.Stdout = &outBuff
	cmd.Stderr = &errBuff
	deferr = cmd.Run()

	stdOutStr := outBuff.String()
	stdErrStr := errBuff.String()

	if stdErrStr != "" {
		l.Error(stdErrStr)
	}

	if stdOutStr != "" {
		l.Info(stdOutStr)
	}

	if deferr != nil {
		return errors.NewExecutionShellCommandNonZeroExitError(processedRawCommand, stdOutStr, stdErrStr).WithCause(deferr)
	}

	return err
}

func parseCommandArguments(command string) ([]string, errors.Error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, errors.NewExecutionInvalidShellCommandError(command)
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}

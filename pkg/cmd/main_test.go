package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func parseArgs(str string) []string {
	args := strings.Split(str, " ")

	var filter []string

	for i, arg := range args {
		if i == 0 && arg == CommandName {
			continue
		}

		filter = append(filter, arg)
	}

	return filter
}

func ExecuteCommand(args string) (string, string) {
	var stdout, stderr bytes.Buffer

	cmd := getCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(parseArgs(args))
	cmd.Execute()

	return stdout.String(), stderr.String()
}

func TestCommands(t *testing.T) {
	t.Run("`relay` should present help text", func(t *testing.T) {
		stdout, _ := ExecuteCommand("relay")

		assert.True(t, strings.HasPrefix(stdout, "Relay connects your tools"))
	})
}

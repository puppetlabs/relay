package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kballard/go-shellquote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseArgs(str string) ([]string, error) {
	args, err := shellquote.Split(str)
	if err != nil {
		return nil, err
	}

	var filter []string

	for i, arg := range args {
		if i == 0 && arg == CommandName {
			continue
		}

		filter = append(filter, arg)
	}

	return filter, nil
}

func ExecuteCommand(args string) (string, string, error) {
	var stdout, stderr bytes.Buffer

	cmd := getCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	pArgs, err := parseArgs(args)
	cmd.SetArgs(pArgs)
	if err != nil {
		return stdout.String(), stderr.String(), err
	}
	err = cmd.Execute()
	if err != nil {
		return stdout.String(), stderr.String(), err
	}

	return stdout.String(), stderr.String(), nil
}

func TestCommands(t *testing.T) {
	t.Run("`relay` should present help text", func(t *testing.T) {
		stdout, _, err := ExecuteCommand("relay")
		require.NoError(t, err)

		assert.True(t, strings.HasPrefix(stdout, "Relay connects your tools"))
	})
}

func TestMetadataCommands(t *testing.T) {
	t.Run("`relay dev metadata` should present spec", func(t *testing.T) {
		t.Skip("Travis binds to ipv6 then fails because ipv6 isn't supported...")
		stdout, stderr, err := ExecuteCommand(`relay dev metadata --run 1234 --step foo --input ../../examples/metadata-configs/simple.yaml -- python -c "import os,requests; print(requests.get('{}/spec'.format(os.environ['METADATA_API_URL'])).content)"`)
		require.NoError(t, err)
		require.Empty(t, stderr)

		//TODO The go stderr and stdout are always empty and python's output goes to the console. Why is that?
		//assert.True(t, strings.HasPrefix(stdout, "6bkpuV9fF3LX1Yo79OpfTwsw8wt5wsVLGTPJjDTu"))
		require.Empty(t, stdout)
	})
}

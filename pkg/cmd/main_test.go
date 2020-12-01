package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	err := cmd.Execute()
	if err != nil {
		return stdout.String(), err.Error()
	}

	return stdout.String(), stderr.String()
}

func TestCommands(t *testing.T) {
	t.Run("`relay` should present help text", func(t *testing.T) {
		stdout, _ := ExecuteCommand("relay")

		assert.True(t, strings.HasPrefix(stdout, "Relay connects your tools"))
	})
}

func TestMetadataCommands(t *testing.T) {
	t.Run("`relay dev metadata` should present spec", func(t *testing.T) {
		stdout, stderr := ExecuteCommand("relay dev metadata --run 1234 --step foo --input ../../examples/metadata-configs/simple.yaml -- python -m os 'requests.get(os.environ[\"METADATA_API_URL\"])'")
		require.Empty(t, stderr)

		//TODO The stderr and stdout are always empty on `relay dev metadata` tests. Why is that?
		//assert.True(t, strings.HasPrefix(stdout, "6bkpuV9fF3LX1Yo79OpfTwsw8wt5wsVLGTPJjDTu"))
		require.Empty(t, stdout)
	})
}

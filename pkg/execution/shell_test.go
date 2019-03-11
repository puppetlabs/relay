package execution

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCommandArguments(t *testing.T) {

	type Assertion struct {
		Input       string
		Output      []string
		ShouldError bool
	}

	assertions := []Assertion{
		Assertion{
			Input: "ls -ltra",
			Output: []string{
				"ls",
				"-ltra",
			},
		},
		Assertion{
			Input: "docker build . -t pcr-internal:5050/my-image:my-tag",
			Output: []string{
				"docker",
				"build",
				".",
				"-t",
				"pcr-internal:5050/my-image:my-tag",
			},
		},
		Assertion{
			Input: "go test ./...",
			Output: []string{
				"go",
				"test",
				"./...",
			},
		},
		Assertion{
			Input: "docker push pcr-internal:5050/my-image:my-tag",
			Output: []string{
				"docker",
				"push",
				"pcr-internal:5050/my-image:my-tag",
			},
		},
	}

	for _, assertion := range assertions {
		output, err := parseCommandArguments(assertion.Input)

		if assertion.ShouldError {
			require.Error(t, err)
		}

		require.Equal(t, assertion.Output, output)
	}
}

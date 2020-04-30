package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/puppetlabs/relay/pkg/config"
	"github.com/puppetlabs/relay/pkg/errors"
)

// Confirm prompts users for confirmation, interfacing with global config through --yes flag
func Confirm(prompt string, cfg *config.Config) (bool, errors.Error) {
	if cfg.Yes {
		return true, nil
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(fmt.Sprintf("%v [y/N] ", prompt))
	prompt, err := reader.ReadString('\n')

	if err != nil {
		return false, errors.NewGeneralUnknownError().WithCause(err)
	}

	return strings.Contains(strings.TrimSpace(prompt), "y"), nil
}

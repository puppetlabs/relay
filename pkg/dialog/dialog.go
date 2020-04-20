// package diaglog encapsulates standard user messaging for all standard cli behavior.
// It is distinct from https://github.com/puppetlabs/horsehead/blob/master/logging/logger.go
// which should be used for structured debug logging, intended for developers or for semi-technical
// users running in debug-mode. This package is for polished messages that are leveled
// but unstructured. All messages are hidden in json output mode, under the assumption
// that users will want to pipe json output to a file or another process
package dialog

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/puppetlabs/relay/pkg/config"
)

type Dialog struct {
	config *config.Config
}

func NewDialog(config *config.Config) *Dialog {
	return &Dialog{
		config,
	}
}

// Info does not print a prefix
func (d *Dialog) Info(message string) {
	if d.config.Out == config.OutputTypeText {
		fmt.Println(message)
	}
}

func (d *Dialog) Warn(message string) {
	if d.config.Out == config.OutputTypeText {
		fmt.Println(color.YellowString("Warning:"), message)
	}
}

func (d *Dialog) Error(message string) {
	if d.config.Out == config.OutputTypeText {
		fmt.Println(color.RedString("Error:"), message)
	}
}

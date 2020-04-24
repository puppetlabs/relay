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

type Dialog interface {
	Info(message string)
	Warn(message string)
	Error(message string)
}

type TextDialog struct{}

// Info does not print a prefix
func (d *TextDialog) Info(message string) {
	fmt.Println(message)
}

func (d *TextDialog) Warn(message string) {
	fmt.Println(color.YellowString("Warning:"), message)
}

func (d *TextDialog) Error(message string) {
	fmt.Println(color.RedString("Error:"), message)
}

type JSONDialog struct{}

// right now json dialog methods do nothing
func (d *JSONDialog) Info(message string) {
}

func (d *JSONDialog) Warn(message string) {
}

func (d *JSONDialog) Error(message string) {
}

func NewDialog(cfg *config.Config) Dialog {
	if cfg.Out == config.OutputTypeJSON {
		return &JSONDialog{}
	}

	return &TextDialog{}
}

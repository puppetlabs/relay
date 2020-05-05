// package diaglog encapsulates standard user messaging for all standard cli behavior.
// It is distinct from https://github.com/puppetlabs/horsehead/blob/master/logging/logger.go
// which should be used for structured debug logging, intended for developers or for semi-technical
// users running in debug-mode. This package is for polished messages that are leveled
// but unstructured. All messages are hidden in json output mode, under the assumption
// that users will want to pipe json output to a file or another process
package dialog

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/puppetlabs/relay/pkg/config"
)

type Dialog interface {
	WithWriter(io.Writer) Dialog

	Info(string)
	Infof(string, ...interface{})

	Error(string)
	Errorf(string, ...interface{})

	// Table returns a table for formatting for output.
	Table() Table
}

type TextDialog struct {
	w io.Writer
}

func (d *TextDialog) WithWriter(w io.Writer) Dialog {
	return &TextDialog{w}
}

func withNewLine(str string) string {
	if len(str) == 0 {
		return ""
	}

	if str[len(str)-1] != '\n' {
		return str + "\n"
	}

	return str
}

// Info does not print a prefix
func (d *TextDialog) Info(message string) {
	fmt.Fprintf(d.w, withNewLine(message))
}

func (d *TextDialog) Infof(message string, args ...interface{}) {
	fmt.Fprintf(d.w, withNewLine(message), args...)
}

func (d *TextDialog) Error(msg string) {
	fmt.Fprintf(d.w, "%s%s", color.RedString("Error:"), msg)
}

func (d *TextDialog) Errorf(msg string, args ...interface{}) {
	str := fmt.Sprintf(msg, args...)
	fmt.Fprintf(d.w, "%s%s", color.RedString("Error:"), str)
}

func (d *TextDialog) Table() Table {
	return &textTable{}
}

type JSONDialog struct {
	w io.Writer
}

func (d *JSONDialog) WithWriter(w io.Writer) Dialog {
	return &JSONDialog{w}
}

func (d *JSONDialog) Info(message string) {
	// noop
}

func (d *JSONDialog) Infof(message string, args ...interface{}) {
	// noop
}

func (d *JSONDialog) Error(message string) {
	// noop
}

func (d *JSONDialog) Errorf(message string, args ...interface{}) {
	// noop
}

func (d *JSONDialog) Table() Table {
	return &jsonTable{}
}

func FromConfig(cfg *config.Config) Dialog {
	switch cfg.Out {
	case config.OutputTypeJSON:
		return &JSONDialog{os.Stdout}
	default:
		return &TextDialog{os.Stdout}
	}
}

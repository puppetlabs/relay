package logger

import (
	"github.com/inconshreveable/log15"
	"github.com/puppetlabs/horsehead/v2/logging"
)

const (
	defaultAt = "nebula"
)

type Options struct {
	At    []string
	Debug bool
}

func New(opts Options) logging.Logger {
	lvl := log15.LvlInfo

	if opts.Debug {
		lvl = log15.LvlDebug
	}

	if len(opts.At) == 0 {
		opts.At = []string{defaultAt}
	}

	logging.SetLevel(lvl)

	logger := logging.Builder().At(opts.At...).Build()

	return logger
}

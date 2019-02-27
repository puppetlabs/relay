package logger

import (
	"github.com/inconshreveable/log15"
	log "github.com/puppetlabs/insights-logging"
)

const (
	defaultAt = "nebula"
)

type Options struct {
	At    []string
	Debug bool
}

func New(opts Options) log.Logger {
	lvl := log15.LvlInfo

	if opts.Debug {
		lvl = log15.LvlDebug
	}

	if len(opts.At) == 0 {
		opts.At = []string{defaultAt}
	}

	logger := log.Builder().At(opts.At...).Build()
	log.SetLevel(lvl)

	return logger
}

package config

import (
	"os"

	logging "github.com/puppetlabs/insights-logging"
	"github.com/puppetlabs/nebula/pkg/logger"
	"github.com/spf13/viper"
)

const (
	defaultConfigName       = "config"
	defaultSystemConfigPath = "/etc/puppet/nebula/"
	defaultUserConfigPath   = "$HOME/.config/nebula/"
)

type CLIRuntime interface {
	Config() *Config
	IO() *IO
	Logger() logging.Logger
}

func NewCLIRuntime() (CLIRuntime, error) {
	return NewStandardRuntime()
}

type StandardRuntime struct {
	config *Config
	io     *IO
	logger logging.Logger
}

func (sr *StandardRuntime) Config() *Config {
	return sr.config
}

func (sr *StandardRuntime) IO() *IO {
	return sr.io
}

func (sr *StandardRuntime) Logger() logging.Logger {
	return sr.logger
}

func NewStandardRuntime() (*StandardRuntime, error) {
	loader := viper.New()

	loader.SetConfigName(defaultConfigName)
	loader.AddConfigPath(defaultSystemConfigPath)
	loader.AddConfigPath(defaultUserConfigPath)

	var cfg Config

	if err := loader.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	r := StandardRuntime{
		config: &cfg,
		io:     &IO{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		logger: logger.New(logger.Options{Debug: cfg.Debug}),
	}

	return &r, nil
}

package config

import (
	"os"

	logging "github.com/puppetlabs/insights-logging"
	"github.com/puppetlabs/nebula/pkg/logger"
	"github.com/puppetlabs/nebula/pkg/workflow/loader"
	"github.com/spf13/viper"
)

const (
	defaultConfigName       = "config"
	defaultConfigType       = "yaml"
	defaultSystemConfigPath = "/etc/puppet/nebula/"
	defaultUserConfigPath   = "$HOME/.config/nebula/"
)

type CLIRuntime interface {
	Config() *Config
	IO() *IO
	Logger() logging.Logger
	WorkflowLoader() loader.Loader
	SetConfig(*Config)
	SetIO(*IO)
	SetLogger(logging.Logger)
	SetWorkflowLoader(loader.Loader)
}

func NewCLIRuntime() (CLIRuntime, error) {
	return NewStandardRuntime()
}

type StandardRuntime struct {
	config         *Config
	io             *IO
	logger         logging.Logger
	workflowLoader loader.Loader
}

func (sr *StandardRuntime) Config() *Config {
	return sr.config
}

func (sr *StandardRuntime) SetConfig(cfg *Config) {
	sr.config = cfg
}

func (sr *StandardRuntime) IO() *IO {
	return sr.io
}

func (sr *StandardRuntime) SetIO(streams *IO) {
	sr.io = streams
}

func (sr *StandardRuntime) Logger() logging.Logger {
	return sr.logger
}

func (sr *StandardRuntime) SetLogger(l logging.Logger) {
	sr.logger = l
}

func (sr *StandardRuntime) WorkflowLoader() loader.Loader {
	return sr.workflowLoader
}

func (sr *StandardRuntime) SetWorkflowLoader(l loader.Loader) {
	sr.workflowLoader = l
}

func NewStandardRuntime() (*StandardRuntime, error) {
	v := viper.New()

	v.SetConfigName(defaultConfigName)
	v.SetConfigType(defaultConfigType)
	v.AddConfigPath(defaultSystemConfigPath)
	v.AddConfigPath(defaultUserConfigPath)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	r := StandardRuntime{
		config:         &cfg,
		io:             &IO{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		logger:         logger.New(logger.Options{Debug: cfg.Debug}),
		workflowLoader: loader.ImpliedWorkflowFileLoader{},
	}

	return &r, nil
}

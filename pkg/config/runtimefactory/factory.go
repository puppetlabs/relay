package runtimefactory

import (
	"os"
	"path/filepath"

	"github.com/puppetlabs/horsehead/logging"
	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/puppetlabs/nebula/pkg/io"
	"github.com/puppetlabs/nebula/pkg/logger"
	"github.com/spf13/viper"
)

const (
	defaultConfigName           = "config"
	defaultConfigType           = "yaml"
	defaultSystemConfigPath     = "/etc/puppet/nebula/"
	defaultDockerHostSocketPath = "/var/run/docker.sock"
)

type RuntimeFactory interface {
	Config() *config.Config
	IO() *io.IO
	Logger() logging.Logger
	SetConfig(*config.Config)
	SetIO(*io.IO)
	SetLogger(logging.Logger)
}

func NewRuntimeFactory() (RuntimeFactory, error) {
	return NewStandardRuntime()
}

type StandardRuntime struct {
	config *config.Config
	io     *io.IO
	logger logging.Logger
}

func (sr *StandardRuntime) Config() *config.Config {
	return sr.config
}

func (sr *StandardRuntime) SetConfig(cfg *config.Config) {
	sr.config = cfg
}

func (sr *StandardRuntime) IO() *io.IO {
	return sr.io
}

func (sr *StandardRuntime) SetIO(streams *io.IO) {
	sr.io = streams
}

func (sr *StandardRuntime) Logger() logging.Logger {
	return sr.logger
}

func (sr *StandardRuntime) SetLogger(l logging.Logger) {
	sr.logger = l
}

func NewStandardRuntime() (*StandardRuntime, error) {
	v := viper.New()

	v.SetConfigName(defaultConfigName)
	v.SetConfigType(defaultConfigType)
	v.AddConfigPath(defaultSystemConfigPath)
	v.AddConfigPath(userConfigPath())

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg config.Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.CachePath = userCachePath()

	if err := os.MkdirAll(cfg.CachePath, 0750); err != nil {
		return nil, err
	}

	cfg.TokenPath = filepath.Join(cfg.CachePath, "auth-token")

	r := StandardRuntime{
		config: &cfg,
		io:     &io.IO{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		logger: logger.New(logger.Options{Debug: cfg.Debug}),
	}

	return &r, nil
}

func userConfigPath() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "nebula")
	}

	return filepath.Join(os.Getenv("HOME"), ".config", "nebula")
}

func userCachePath() string {
	if os.Getenv("XDG_CACHE_HOME") != "" {
		return filepath.Join(os.Getenv("XDG_CACHE_HOME"), "nebula")
	}

	return filepath.Join(os.Getenv("HOME"), ".cache", "nebula")
}

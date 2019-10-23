package runtimefactory

import (
	"os"
	"path/filepath"

	"github.com/puppetlabs/horsehead/v2/logging"
	"github.com/puppetlabs/nebula-cli/pkg/config"
	"github.com/puppetlabs/nebula-cli/pkg/io"
	"github.com/puppetlabs/nebula-cli/pkg/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConfigName           = "config"
	defaultConfigType           = "yaml"
	defaultSystemConfigPath     = "/etc/puppet/nebula/"
	defaultDockerHostSocketPath = "/var/run/docker.sock"
)

type RuntimeFactory interface {
	Config() (*config.Config, error)
	IO() *io.IO
	Logger() (logging.Logger, error)
}

func NewRuntimeFactory(flags *pflag.FlagSet) RuntimeFactory {
	return NewStandardRuntime(flags)
}

type StandardRuntime struct {
	flags  *pflag.FlagSet
	config *config.Config
	io     *io.IO
	logger logging.Logger
}

func (sr *StandardRuntime) Config() (*config.Config, error) {
	if sr.config == nil {
		cp, err := sr.flags.GetString("config")
		if err != nil {
			return nil, err
		}

		v := viper.New()

		v.SetConfigName(defaultConfigName)
		v.SetConfigType(defaultConfigType)

		if cp != "" {
			// SetConfigFile will check of path is not empty. If it is set, then it
			// will force viper to attempt loading the configuration from that file only.
			// If the file doesn't exist, then we want to bail and inform the user that something
			// went wrong as an explicit file path for configuration seems important.
			v.SetConfigFile(cp)
		} else {
			v.AddConfigPath(defaultSystemConfigPath)
			v.AddConfigPath(userConfigDir())
		}

		if err := v.ReadInConfig(); err != nil {
			if cp != "" {
				return nil, err
			}

			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, err
			}
		}

		var cfg config.Config

		if err := v.Unmarshal(&cfg); err != nil {
			return nil, err
		}

		if cfg.CacheDir == "" {
			cfg.CacheDir = userCacheDir()
		}

		if err := os.MkdirAll(cfg.CacheDir, 0750); err != nil {
			return nil, err
		}

		if cfg.TokenPath == "" {
			cfg.TokenPath = filepath.Join(cfg.CacheDir, "auth-token")
		}

		sr.config = &cfg
	}

	return sr.config, nil
}

func (sr *StandardRuntime) IO() *io.IO {
	return sr.io
}

func (sr *StandardRuntime) Logger() (logging.Logger, error) {
	if sr.logger == nil {
		cfg, err := sr.Config()
		if err != nil {
			return nil, err
		}

		sr.logger = logger.New(logger.Options{Debug: cfg.Debug})
	}

	return sr.logger, nil
}

func NewStandardRuntime(flags *pflag.FlagSet) *StandardRuntime {
	r := StandardRuntime{
		flags: flags,
		io:    &io.IO{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	}

	return &r
}

func userConfigDir() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "nebula")
	}

	return filepath.Join(os.Getenv("HOME"), ".config", "nebula")
}

func userCacheDir() string {
	if os.Getenv("XDG_CACHE_HOME") != "" {
		return filepath.Join(os.Getenv("XDG_CACHE_HOME"), "nebula")
	}

	return filepath.Join(os.Getenv("HOME"), ".cache", "nebula")
}

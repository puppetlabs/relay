package runtimefactory

import (
	"os"
	"path/filepath"

	logging "github.com/puppetlabs/insights-logging"
	"github.com/puppetlabs/nebula/pkg/config"
	"github.com/puppetlabs/nebula/pkg/execution/executor"
	"github.com/puppetlabs/nebula/pkg/io"
	"github.com/puppetlabs/nebula/pkg/loader"
	"github.com/puppetlabs/nebula/pkg/logger"
	"github.com/puppetlabs/nebula/pkg/state"
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
	PlanLoader() loader.Loader
	StateManager() state.Manager
	ActionExecutor() executor.ActionExecutor
	SetConfig(*config.Config)
	SetIO(*io.IO)
	SetLogger(logging.Logger)
	SetPlanLoader(loader.Loader)
	SetStateManager(state.Manager)
	SetActionExecutor(executor.ActionExecutor)
}

func NewRuntimeFactory() (RuntimeFactory, error) {
	return NewStandardRuntime()
}

type StandardRuntime struct {
	config       *config.Config
	io           *io.IO
	logger       logging.Logger
	planLoader   loader.Loader
	stateManager state.Manager
	executor     executor.ActionExecutor
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

func (sr *StandardRuntime) PlanLoader() loader.Loader {
	return sr.planLoader
}

func (sr *StandardRuntime) SetPlanLoader(l loader.Loader) {
	sr.planLoader = l
}

func (sr *StandardRuntime) StateManager() state.Manager {
	return sr.stateManager
}

func (sr *StandardRuntime) SetStateManager(sm state.Manager) {
	sr.stateManager = sm
}

func (sr *StandardRuntime) ActionExecutor() executor.ActionExecutor {
	return sr.executor
}

func (sr *StandardRuntime) SetActionExecutor(e executor.ActionExecutor) {
	sr.executor = e
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

	if cfg.DockerExecutor.HostSocketPath == "" {
		cfg.DockerExecutor.HostSocketPath = defaultDockerHostSocketPath
	}

	cfg.CachePath = userCachePath()

	if err := os.MkdirAll(cfg.CachePath, 0750); err != nil {
		return nil, err
	}

	cfg.TokenPath = filepath.Join(cfg.CachePath, "auth-token")

	sm, err := state.NewFilesystemStateManager("")
	if err != nil {
		return nil, err
	}

	exec, eerr := executor.NewExecutor(executor.RegistryCredentials{
		Registry: cfg.DockerExecutor.Registry,
		User:     cfg.DockerExecutor.RegistryUser,
		Pass:     cfg.DockerExecutor.RegistryPass,
	})
	if err != nil {
		return nil, eerr
	}

	r := StandardRuntime{
		config:       &cfg,
		io:           &io.IO{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		logger:       logger.New(logger.Options{Debug: cfg.Debug}),
		planLoader:   loader.ImpliedPlanFileLoader{},
		stateManager: sm,
		executor:     exec,
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

// Package config defines global configuration values
package config

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// OutputType may be either text or json
type OutputType string

const (
	OutputTypeText OutputType = "text"
	OutputTypeJSON OutputType = "json"
)

func (ot OutputType) String() string {
	return string(ot)
}

type AuthTokenType string

const (
	AuthTokenTypeAPI     AuthTokenType = "api"
	AuthTokenTypeSession AuthTokenType = "session"
)

func (att AuthTokenType) String() string {
	return string(att)
}

func AuthTokenTypes() []AuthTokenType {
	return []AuthTokenType{AuthTokenTypeAPI, AuthTokenTypeSession}
}

func AuthTokenTypesAsString() []string {
	authTokenTypes := make([]string, len(AuthTokenTypes()))
	for index, tokenType := range AuthTokenTypes() {
		authTokenTypes[index] = tokenType.String()
	}

	return authTokenTypes
}

const (
	RelayEnvironment = "relay"

	defaultConfigName     = "config"
	defaultConfigType     = "yaml"
	defaultCurrentContext = "relaysh"
)

var defaultContexts = map[string]*ContextConfig{
	"relaysh": {
		Domains: &APIContext{
			Name:      "relaysh",
			APIDomain: &url.URL{Scheme: "https", Host: "api.relay.sh"},
			UIDomain:  &url.URL{Scheme: "https", Host: "app.relay.sh"},
			WebDomain: &url.URL{Scheme: "https", Host: "relay.sh"},
		},
	},
	"dev": {
		Domains: &APIContext{
			Name:      "dev",
			APIDomain: &url.URL{Scheme: "http", Host: "relay-api.local:8080"},
			UIDomain:  &url.URL{Scheme: "http", Host: "relay-ui.local:8080"},
			WebDomain: &url.URL{Scheme: "http", Host: "relay-ui.local:8080"},
		},
	},
}

type APIContext struct {
	Name      string
	APIDomain *url.URL
	UIDomain  *url.URL
	WebDomain *url.URL
}

func (ac *APIContext) Merge(target *APIContext) *APIContext {
	if target == nil {
		return ac
	}

	return &APIContext{
		APIDomain: coalesceURL(ac.APIDomain, target.APIDomain),
		UIDomain:  coalesceURL(ac.UIDomain, target.UIDomain),
		WebDomain: coalesceURL(ac.WebDomain, target.WebDomain),
	}
}

type InstallerConfig struct {
	InstallerImage                            string
	LogServiceImage                           string
	MetadataAPIImage                          string
	OperatorImage                             string
	OperatorVaultInitImage                    string
	OperatorWebhookCertificateControllerImage string
	VaultServerImage                          string
	VaultSidecarImage                         string
}

type LogServiceConfig struct {
	CredentialsKey        string
	CredentialsSecretName string
	Project               string
	Dataset               string
	Table                 string
}

type AuthConfig struct {
	Tokens map[AuthTokenType]string
}

type ContextConfig struct {
	Auth    *AuthConfig
	Domains *APIContext
}

type Config struct {
	Debug          bool
	Yes            bool
	Out            OutputType
	CacheDir       string
	TokenPath      string
	CurrentContext string

	ContextConfig map[string]*ContextConfig

	InstallerConfig  *InstallerConfig
	LogServiceConfig *LogServiceConfig
}

// GetDefaultConfig returns a config set used for error formatting when the user's config set cannot be read
func GetDefaultConfig() *Config {
	return &Config{
		Debug:          true,
		Yes:            false,
		Out:            OutputTypeText,
		CacheDir:       userCacheDir(),
		CurrentContext: defaultCurrentContext,

		ContextConfig: defaultContexts,
	}
}

func NewAPIContext(v *viper.Viper) (*APIContext, error) {
	apiDomain, err := url.Parse(v.GetString("apiDomain"))
	if err != nil {
		return nil, err
	}

	uiDomain, err := url.Parse(v.GetString("uiDomain"))
	if err != nil {
		return nil, err
	}

	webDomain, err := url.Parse(v.GetString("webDomain"))
	if err != nil {
		return nil, err
	}

	return &APIContext{
		APIDomain: apiDomain,
		UIDomain:  uiDomain,
		WebDomain: webDomain,
	}, nil
}

func NewInstallerConfig(v *viper.Viper) *InstallerConfig {
	return &InstallerConfig{
		InstallerImage:         v.GetString("installerImage"),
		LogServiceImage:        v.GetString("logServiceImage"),
		MetadataAPIImage:       v.GetString("metadataAPIImage"),
		OperatorImage:          v.GetString("operatorImage"),
		OperatorVaultInitImage: v.GetString("operatorVaultInitImage"),
		OperatorWebhookCertificateControllerImage: v.GetString("operatorWebhookCertificateControllerImage"),
		VaultServerImage:  v.GetString("vaultServerImage"),
		VaultSidecarImage: v.GetString("vaultSidecarImage"),
	}
}

func NewLogServiceConfig(v *viper.Viper) *LogServiceConfig {
	return &LogServiceConfig{
		CredentialsKey:        v.GetString("credentialsKey"),
		CredentialsSecretName: v.GetString("credentialsSecretName"),
		Project:               v.GetString("project"),
		Dataset:               v.GetString("dataset"),
		Table:                 v.GetString("table"),
	}
}

// FromFlags uses viper to read global configuration from persistent flags,
// environment variables, and / or yaml config read from $HOME/.config/relay
func FromFlags(flags *pflag.FlagSet) (*Config, error) {
	v := viper.New()

	v.SetEnvPrefix(RelayEnvironment)
	v.AutomaticEnv()

	v.SetDefault("debug", false)
	v.BindPFlag("debug", flags.Lookup("debug"))

	v.SetDefault("yes", false)
	v.BindPFlag("yes", flags.Lookup("yes"))

	v.SetDefault("out", string(OutputTypeText))
	v.BindPFlag("out", flags.Lookup("out"))

	v.SetDefault("cache_dir", userCacheDir())
	v.SetDefault("data_dir", userDataDir())

	v.SetDefault("context", defaultCurrentContext)
	v.BindPFlag("context", flags.Lookup("context"))

	if err := readInConfigFile(v, flags); err != nil {
		return nil, err
	}

	context := v.GetString("context")

	output, err := readOutput(v)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Debug:    v.GetBool("debug"),
		Yes:      v.GetBool("yes"),
		Out:      output,
		CacheDir: v.GetString("cache_dir"),

		CurrentContext: context,
		ContextConfig:  defaultContexts,
	}

	if config.ContextConfig[context] == nil {
		config.ContextConfig[context] = &ContextConfig{}
	}

	// FIXME This will likely change to read in the entire context section
	// to enable switching context on demand without necessarily reloading
	// the configuration
	if context != "" {
		installerConfigSection := v.Sub(fmt.Sprintf("config.%s.installer", context))
		if installerConfigSection != nil {
			config.InstallerConfig = NewInstallerConfig(installerConfigSection)
		}

		logServiceConfigSection := v.Sub(fmt.Sprintf("config.%s.logService", context))
		if logServiceConfigSection != nil {
			config.LogServiceConfig = NewLogServiceConfig(logServiceConfigSection)
		}

		contextSection := v.Sub(fmt.Sprintf("contexts.%s", context))
		if contextSection != nil {
			if config.ContextConfig[context].Domains == nil {
				config.ContextConfig[context].Domains = &APIContext{}
			}

			domainConfig, err := NewAPIContext(contextSection)
			if err != nil {
				return nil, err
			}

			config.ContextConfig[context].Domains =
				config.ContextConfig[context].Domains.Merge(domainConfig)

			authSection := contextSection.Sub("auth")
			if authSection != nil {
				config.ContextConfig[context].Auth = &AuthConfig{
					Tokens: make(map[AuthTokenType]string),
				}

				for _, tokenType := range AuthTokenTypes() {
					token := authSection.GetString(fmt.Sprintf("tokens.%s", tokenType))
					config.ContextConfig[context].Auth.Tokens[tokenType] = token
				}
			}
		}
	}

	return config, nil
}

func WriteConfig(cfg *Config, flags *pflag.FlagSet) error {
	v := viper.New()

	v.SetEnvPrefix(RelayEnvironment)
	v.AutomaticEnv()

	readInConfigFile(v, flags)

	if cfg.CurrentContext != "" {
		v.Set("context", cfg.CurrentContext)
	}

	if cfg.ContextConfig != nil {
		for context, config := range cfg.ContextConfig {
			if config.Auth != nil && config.Auth.Tokens != nil {
				tokens := config.Auth.Tokens
				for name, value := range tokens {
					v.Set(fmt.Sprintf("contexts.%s.auth.tokens.%s", context, name), value)
				}
			}
		}
	}

	if err := v.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func WriteGlobalConfig(cfg *Config, flags *pflag.FlagSet) error {
	v := viper.New()

	v.SetEnvPrefix(RelayEnvironment)
	v.AutomaticEnv()

	readInConfigFile(v, flags)

	v.Set("debug", cfg.Debug)
	v.Set("out", cfg.Out)
	v.Set("yes", cfg.Yes)

	if err := v.WriteConfig(); err != nil {
		return err
	}

	return nil
}

// readInConfigFile reads config file location from viper flags, then
// reads in config from specified location or the default
func readInConfigFile(v *viper.Viper, flags *pflag.FlagSet) error {
	cp, err := flags.GetString("config")
	if err != nil {
		return errors.NewConfigInvalidConfigFlag().WithCause(err)
	}

	v.SetConfigName(defaultConfigName)
	v.SetConfigType(defaultConfigType)

	if cp != "" {
		v.SetConfigFile(cp)
	} else {
		v.AddConfigPath(userConfigDir())
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			p := cp
			if p == "" {
				p = path.Join(userConfigDir(),
					strings.Join([]string{defaultConfigName, defaultConfigType}, "."))
			}

			if err := os.MkdirAll(path.Dir(p), 0750); err != nil {
				return err
			}

			if err := v.WriteConfigAs(p); err != nil {
				return err
			}
		} else {
			// Config file was found but another error was produced
			return errors.NewConfigInvalidConfigFile(cp).WithCause(err)
		}
	}

	return nil
}

// userConfigDir gets default user config dir
func userConfigDir() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return filepath.Join(os.Getenv("XDG_CONFIG_HOME"), RelayEnvironment)
	}

	return filepath.Join(os.Getenv("HOME"), ".config", RelayEnvironment)
}

// userCacheDir gets default user cache dir, used as directory for storing tokens
func userCacheDir() string {
	if os.Getenv("XDG_CACHE_HOME") != "" {
		return filepath.Join(os.Getenv("XDG_CACHE_HOME"), RelayEnvironment)
	}

	return filepath.Join(os.Getenv("HOME"), ".cache", RelayEnvironment)
}

// userDataDir gets default user data dir. The data dir is used to store long term
// data generated by the cli.
func userDataDir() string {
	if os.Getenv("XDG_DATA_HOME") != "" {
		return filepath.Join(os.Getenv("XDG_DATA_HOME"), RelayEnvironment)
	}

	return filepath.Join(os.Getenv("HOME"), ".local", "share", RelayEnvironment)
}

// readOutput reads and validates output config value
func readOutput(v *viper.Viper) (OutputType, error) {
	output := OutputType(v.GetString("out"))

	if output != OutputTypeText && output != OutputTypeJSON {

		return "", errors.NewConfigInvalidOutputFlag(v.GetString("out"))
	}

	return output, nil
}

// readAPIDomain reads and validates api domain config value
func readAPIDomain(v *viper.Viper) (*url.URL, error) {
	urlString := v.GetString("api_domain")
	url, err := url.Parse(urlString)

	if err != nil {
		return nil, errors.NewConfigInvalidAPIDomain(urlString)
	}

	return url, nil
}

// readUIDomain reads and validates ui domain config value
func readUIDomain(v *viper.Viper) (*url.URL, error) {
	urlString := v.GetString("ui_domain")
	url, err := url.Parse(urlString)

	if err != nil {
		return nil, errors.NewConfigInvalidUIDomain(urlString)
	}

	return url, nil
}

// readWebDomain reads and validates web domain config value
func readWebDomain(v *viper.Viper) (*url.URL, error) {
	urlString := v.GetString("web_domain")
	url, err := url.Parse(urlString)

	if err != nil {
		return nil, errors.NewConfigInvalidWebDomain(urlString)
	}

	return url, nil
}

func coalesceURL(src *url.URL, dst *url.URL) *url.URL {
	if dst != nil && dst.String() != "" {
		return dst
	}

	return src
}

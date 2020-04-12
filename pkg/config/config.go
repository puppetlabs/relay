package config

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type OutputType string

const (
	OutputTypeText OutputType = "text"
	OutputTypeJSON OutputType = "json"
)

const (
	defaultAPIDomain  = "https://api.nebula.puppet.com"
	defaultUIDomain   = "https://nebula.puppet.com"
	defaultWebDomain  = "https://relay.sh"
	defaultConfigName = "config"
	defaultConfigType = "yaml"
)

type Config struct {
	Debug     bool
	Verbose   bool
	Out       OutputType
	APIDomain *url.URL
	UIDomain  *url.URL
	WebDomain *url.URL
	CacheDir  string
	TokenPath string
}

// GetConfig uses viper to read global configuration from persistent flags,
// environment variables, and / or yaml config read from $HOME/.config/relay
func GetConfig(flags *pflag.FlagSet) (*Config, error) {
	v := viper.New()

	v.SetEnvPrefix("relay")
	v.AutomaticEnv()

	v.SetDefault("debug", false)
	v.BindPFlag("debug", flags.Lookup("debug"))

	v.SetDefault("verbose", false)
	v.BindPFlag("verbose", flags.Lookup("verbose"))

	v.SetDefault("out", string(OutputTypeText))
	v.BindPFlag("out", flags.Lookup("out"))

	v.SetDefault("api_domain", defaultAPIDomain)
	v.SetDefault("ui_domain", defaultUIDomain)
	v.SetDefault("web_domain", defaultUIDomain)
	v.SetDefault("cache_dir", userCacheDir())
	v.SetDefault("token_path", filepath.Join(userCacheDir(), "auth-token"))

	if err := readInConfigFile(v, flags); err != nil {
		return nil, err
	}

	output, oerr := readOutput(v)

	if oerr != nil {
		return nil, oerr
	}

	apiDomain, aderr := readAPIDomain(v)

	if aderr != nil {
		return nil, aderr
	}

	uiDomain, uderr := readUIDomain(v)

	if uderr != nil {
		return nil, uderr
	}

	webDomain, wderr := readWebDomain(v)

	if wderr != nil {
		return nil, wderr
	}

	config := &Config{
		Debug:     v.GetBool("debug"),
		Verbose:   v.GetBool("verbose"),
		Out:       output,
		APIDomain: apiDomain,
		UIDomain:  uiDomain,
		WebDomain: webDomain,
		CacheDir:  v.GetString("cache_dir"),
		TokenPath: v.GetString("token_path"),
	}

	return config, nil
}

func readInConfigFile(v *viper.Viper, flags *pflag.FlagSet) error {
	cp, err := flags.GetString("config")
	if err != nil {
		return errors.NewConfigInvalidConfigFlag().WithCause(err)
	}

	v.SetConfigName(defaultConfigName)
	v.SetConfigType(defaultConfigType)

	if cp != "" {
		// SetConfigFile will check of path is not empty. If it is set, then it
		// will force viper to attempt loading the configuration from that file only.
		// If the file doesn't exist, then we want to bail and inform the user that something
		// went wrong as an explicit file path for configuration seems important.
		v.SetConfigFile(cp)
	} else {
		v.AddConfigPath(userConfigDir())
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; This is fine if they didn't specify a custom path
			// but we want to alert them if the path they specified doesn't exist

			if cp == "" {
				return nil
			} else {
				return errors.NewConfigFileNotFound(cp).WithCause(err)
			}
		} else {
			// Config file was found but another error was produced
			return errors.NewConfigInvalidConfigFile(cp).WithCause(err)
		}
	}

	return nil
}

func userConfigDir() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "relay")
	}

	return filepath.Join(os.Getenv("HOME"), ".config", "relay")
}

func userCacheDir() string {
	if os.Getenv("XDG_CACHE_HOME") != "" {
		return filepath.Join(os.Getenv("XDG_CACHE_HOME"), "nebula")
	}

	return filepath.Join(os.Getenv("HOME"), ".cache", "nebula")
}

func readOutput(v *viper.Viper) (OutputType, error) {
	output := OutputType(v.GetString("out"))

	if output != OutputTypeText && output != OutputTypeJSON {
		return "", errors.NewConfigInvalidOutputFlag(v.GetString("out"))
	}

	return output, nil
}

func readAPIDomain(v *viper.Viper) (*url.URL, error) {
	urlString := v.GetString("api_domain")
	url, err := url.Parse(urlString)

	if err != nil {
		return nil, errors.NewConfigInvalidAPIDomain(urlString)
	}

	return url, nil
}

func readUIDomain(v *viper.Viper) (*url.URL, error) {
	urlString := v.GetString("ui_domain")
	url, err := url.Parse(urlString)

	if err != nil {
		return nil, errors.NewConfigInvalidUIDomain(urlString)
	}

	return url, nil
}

func readWebDomain(v *viper.Viper) (*url.URL, error) {
	urlString := v.GetString("web_domain")
	url, err := url.Parse(urlString)

	if err != nil {
		return nil, errors.NewConfigInvalidWebDomain(urlString)
	}

	return url, nil
}

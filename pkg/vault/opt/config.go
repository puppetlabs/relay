// Package opt houses structs for server configuration.
package opt

import (
	"net/url"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

const (
	DefaultVaultURL = "http://localhost:8200"

	VaultAddrConfigOption = "vault_addr"

	VaultAddrEnvironmentVariable = "VAULT_ADDR"
)

var ProviderSet = wire.NewSet(
	NewConfig,
)

type Config struct {
	VaultAddr *url.URL
}

func NewConfig() (*Config, error) {
	viper.AutomaticEnv()

	viper.BindEnv(VaultAddrConfigOption, VaultAddrEnvironmentVariable)
	viper.SetDefault(VaultAddrConfigOption, DefaultVaultURL)

	vaultURL, err := url.Parse(viper.GetString(VaultAddrConfigOption))
	if err != nil {
		return nil, err
	}

	conf := &Config{
		VaultAddr: vaultURL,
	}

	return conf, nil
}

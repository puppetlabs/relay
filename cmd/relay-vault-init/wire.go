// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/puppetlabs/relay/pkg/vault/op/vault"
	"github.com/puppetlabs/relay/pkg/vault/opt"
)

func vaultConfig(cfg *opt.Config) vault.Config {
	return vault.Config{
		Addr: cfg.VaultAddr.String(),
	}
}

func InitializeServices(ctx context.Context, cfg *opt.Config) (services, error) {
	wire.Build(
		vaultConfig,
		vault.ProviderSet,
		wire.Struct(new(services), "*"),
	)

	return services{}, nil
}

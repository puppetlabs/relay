package vault

import (
	"context"

	"github.com/google/wire"
	vaultapi "github.com/hashicorp/vault/api"
)

var ProviderSet = wire.NewSet(
	NewClient,
)

type Config struct {
	Addr string
}

func NewClient(ctx context.Context, cfg Config) (*vaultapi.Client, error) {
	vc, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		return nil, err
	}

	if err := vc.SetAddress(cfg.Addr); err != nil {
		return nil, err
	}

	return vc, nil
}

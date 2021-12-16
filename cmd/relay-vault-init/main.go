package main

import (
	"context"
	"log"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/puppetlabs/relay/pkg/dev"
	"github.com/puppetlabs/relay/pkg/vault/opt"
)

type services struct {
	vault *vaultapi.Client
}

func main() {
	ctx := context.Background()
	conf, err := opt.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	svcs, err := InitializeServices(ctx, conf)
	if err != nil {
		log.Fatal(err)
	}

	vi, err := dev.NewVaultInitializer(svcs.vault)
	if err != nil {
		log.Fatal(err)
	}

	err = vi.Initialize(ctx)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

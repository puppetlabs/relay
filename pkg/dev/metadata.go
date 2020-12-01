package dev

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/inconshreveable/log15"
	"github.com/puppetlabs/errawr-go/v2/pkg/errawr"
	"github.com/puppetlabs/horsehead/v2/logging"
	"github.com/puppetlabs/relay-core/pkg/metadataapi/opt"
	"github.com/puppetlabs/relay-core/pkg/metadataapi/sample"
	"github.com/puppetlabs/relay-core/pkg/metadataapi/server"
	"github.com/puppetlabs/relay-core/pkg/metadataapi/server/middleware"
	"github.com/puppetlabs/relay-core/pkg/util/lifecycleutil"
)

type MetadataAPIManager struct {
	cfg Config
}

type MetadataMockOptions struct {
	RunID string
	StepName string
	Input string
}

func (m *MetadataAPIManager) InitializeMetadataApi(ctx context.Context, mockOptions MetadataMockOptions) (string, error) {
	dynamicAddr := make(chan string)
	s, token, err := m.initializeMetadataServer(ctx, dynamicAddr, mockOptions)
	if err != nil {
		return "", err
	}

	var listenOpts []lifecycleutil.ListenWaitHTTPOption
	go func() {
		// This will end by the closer when the ctx is marked done
		err = lifecycleutil.ListenWaitHTTP(ctx, s, listenOpts...)
		if err != nil {
			println(fmt.Errorf("couldn't start metadata service: %v", err.Error()))
			os.Exit(1)
		}
	}()

	return fmt.Sprintf("http://:%s@%s", token, <-dynamicAddr), nil
}

func (m *MetadataAPIManager) initializeMetadataServer(ctx context.Context, addr chan string, mockOptions MetadataMockOptions) (*http.Server, string, error) {
	var auth middleware.Authenticator
	var tm *sample.TokenMap
	cfg := opt.NewConfig()
	cfg.ListenPort = 0
	cfg.SampleConfigFiles = []string{mockOptions.Input}
	log := m.cfg.Dialog

	if cfg.Debug {
		logging.SetLevel(log15.LvlDebug)
	}

	if sc, err := cfg.SampleConfig(); err != nil {
		return nil, "", err
	} else if sc != nil {
		var key []byte

		if ek := cfg.SampleHS256SigningKey; ek != "" {
			var err error

			key, err = base64.StdEncoding.DecodeString(ek)
			if err != nil {
				return nil, "", fmt.Errorf("could not decode signing key: %+v", err)
			}
		}

		tg, err := sample.NewHS256TokenGenerator(key)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create token generator: %+v", err)
		}

		tm = tg.GenerateAll(ctx, sc)

		auth = sample.NewAuthenticator(sc, tg.Key())
	}
	var serverOpts []server.Option
	serverOpts = append(serverOpts, server.WithErrorSensitivity(errawr.ErrorSensitivityAll))

	s := &http.Server{
		Handler: server.NewHandler(auth, serverOpts...),
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.ListenPort),
		BaseContext: func(l net.Listener) context.Context {
			addr <- l.Addr().String()
			return context.Background()
		},
	}
	token, found := tm.ForStep(mockOptions.RunID, mockOptions.StepName)
	if !found {
		return nil, "", fmt.Errorf("failed to find run ID %s with a step named %s in %s", mockOptions.RunID, mockOptions.StepName, mockOptions.Input)
	}

	log.Infof("startup finished, listening for metadata connections")
	log.Infof("----------------------------------------------------")
	return s, token, nil
}

func NewMetadataAPIManager(cfg Config) *MetadataAPIManager {
	return &MetadataAPIManager{
		cfg: cfg,
	}
}

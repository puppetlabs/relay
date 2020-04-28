package outputs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/puppetlabs/horsehead/v2/encoding/transfer"
)

var (
	ErrOutputsClientKeyEmpty      = errors.New("key is required but was empty")
	ErrOutputsClientValueEmpty    = errors.New("value is required but was empty")
	ErrOutputsClientTaskNameEmpty = errors.New("taskName is required but was empty")
	ErrOutputsClientEnvVarMissing = errors.New(MetadataAPIURLEnvName + " was expected but was empty")
	ErrOutputsClientNotFound      = errors.New("output was not found")
)

// OutputsClient is a client for storing task outputs in
// the nebula outputs storage.
type OutputsClient interface {
	SetOutput(ctx context.Context, key string, value interface{}) error
	GetOutput(ctx context.Context, taskName, key string) (interface{}, error)
}

// DefaultOutputsClient uses the default net/http.Client to
// store task output values.
type DefaultOutputsClient struct {
	apiURL *url.URL
}

func (c DefaultOutputsClient) SetOutput(ctx context.Context, key string, value interface{}) error {
	if key == "" {
		return ErrOutputsClientKeyEmpty
	}

	if value == "" {
		return ErrOutputsClientValueEmpty
	}

	loc := *c.apiURL
	loc.Path = path.Join(loc.Path, key)

	encoded, err := json.Marshal(transfer.JSONInterface{Data: value})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", loc.String(), bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}

func (c DefaultOutputsClient) GetOutput(ctx context.Context, taskName, key string) (interface{}, error) {
	if key == "" {
		return "", ErrOutputsClientKeyEmpty
	}

	if taskName == "" {
		return "", ErrOutputsClientTaskNameEmpty
	}

	loc := *c.apiURL
	loc.Path = path.Join(loc.Path, taskName, key)

	req, err := http.NewRequest("GET", loc.String(), nil)
	if err != nil {
		return "", err
	}

	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return "", ErrOutputsClientNotFound
		}

		return "", fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var output Output

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return "", err
	}

	return output.Value.Data, nil
}

func NewDefaultOutputsClient(location *url.URL) OutputsClient {
	return &DefaultOutputsClient{apiURL: location}
}

func NewDefaultOutputsClientFromNebulaEnv() (OutputsClient, error) {
	locStr := os.Getenv(MetadataAPIURLEnvName)

	if locStr == "" {
		return nil, ErrOutputsClientEnvVarMissing
	}

	loc, err := url.Parse(locStr)
	if err != nil {
		return nil, err
	}

	loc.Path = "/outputs"

	return NewDefaultOutputsClient(loc), nil
}

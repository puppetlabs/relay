package secrets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

var (
	ErrClientKeyEmpty = errors.New("key is required but was empty")
	ErrClientNotFound = errors.New("secret was not found")
)

type Client interface {
	GetSecret(ctx context.Context, key string) (string, error)
}

type DefaultClient struct {
	apiURL *url.URL
}

func (dc *DefaultClient) GetSecret(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", ErrClientKeyEmpty
	}

	loc := *dc.apiURL
	loc.Path = path.Join(loc.Path, key)

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

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrClientNotFound
	} else if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var secret Secret
	if err := json.NewDecoder(resp.Body).Decode(&secret); err != nil {
		return "", err
	}

	v, err := secret.Value.Decode()
	if err != nil {
		return "", err
	}

	return string(v), nil
}

func NewDefaultClient(location *url.URL) Client {
	return &DefaultClient{
		apiURL: location,
	}
}

package resolve

import (
	"context"
	"net/url"

	"github.com/puppetlabs/nebula-sdk/pkg/outputs"
	"github.com/puppetlabs/nebula-sdk/pkg/secrets"
)

type MetadataAPISecretTypeResolver struct {
	client secrets.Client
}

var _ SecretTypeResolver = &MetadataAPISecretTypeResolver{}

func (sr *MetadataAPISecretTypeResolver) ResolveSecret(ctx context.Context, name string) (string, error) {
	v, err := sr.client.GetSecret(ctx, name)
	if err == secrets.ErrClientNotFound {
		return "", &SecretNotFoundError{Name: name}
	} else if err != nil {
		return "", err
	}

	return v, nil
}

func NewMetadataAPISecretTypeResolver(api *url.URL) *MetadataAPISecretTypeResolver {
	return &MetadataAPISecretTypeResolver{
		client: secrets.NewDefaultClient(api.ResolveReference(&url.URL{Path: "/secrets"})),
	}
}

type MetadataAPIOutputTypeResolver struct {
	client outputs.OutputsClient
}

var _ OutputTypeResolver = &MetadataAPIOutputTypeResolver{}

func (sr *MetadataAPIOutputTypeResolver) ResolveOutput(ctx context.Context, from, name string) (interface{}, error) {
	v, err := sr.client.GetOutput(ctx, from, name)
	if err == outputs.ErrOutputsClientNotFound {
		return "", &OutputNotFoundError{From: from, Name: name}
	} else if err != nil {
		return "", err
	}

	return v, nil
}

func NewMetadataAPIOutputTypeResolver(api *url.URL) *MetadataAPIOutputTypeResolver {
	return &MetadataAPIOutputTypeResolver{
		client: outputs.NewDefaultOutputsClient(api.ResolveReference(&url.URL{Path: "/outputs"})),
	}
}

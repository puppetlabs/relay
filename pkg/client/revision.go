package client

import (
	"fmt"
	"net/http"

	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/model"
)

func (c *Client) CreateRevision(workflowName string, YAML string) (*model.RevisionEntity, errors.Error) {
	response := &model.RevisionEntity{}

	var headers = map[string]string{
		"Content-Type": fmt.Sprintf("application/vnd.puppet.relay.%s+yaml", APIVersion),
	}

	if err := c.Request(
		WithMethod(http.MethodPost),
		WithPath(fmt.Sprintf("/api/workflows/%s/revisions", workflowName)),
		WithBodyEncodingType(BodyEncodingTypeYAML),
		WithHeaders(headers),
		WithBody(YAML),
		WithResponseInto(response),
	); err != nil {
		return nil, err
	}

	return response, nil
}

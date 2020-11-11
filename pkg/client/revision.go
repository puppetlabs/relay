package client

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/model"
)

func (c *Client) Validate(YAML string) (*model.RevisionEntity, errors.Error) {
	response := &model.RevisionEntity{}

	var headers = map[string]string{
		"Content-Type": fmt.Sprintf("application/vnd.puppet.relay.%s+yaml", APIVersion),
	}

	if err := c.Request(
		WithMethod(http.MethodPost),
		WithPath(fmt.Sprintf("/api/revisions/validate")),
		WithBodyEncodingType(BodyEncodingTypeYAML),
		WithHeaders(headers),
		WithBody(YAML),
		WithResponseInto(response),
	); err != nil {
		return nil, err
	}

	return response, nil
}

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

func (c *Client) GetRevision(workflowName, revisionID string) (*model.RevisionEntity, errors.Error) {
	response := &model.RevisionEntity{}

	if err := c.Request(
		WithPath(path.Join("/api/workflows", url.PathEscape(workflowName), "revisions", url.PathEscape(revisionID))),
		WithResponseInto(response),
	); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetLatestRevision(workflowName string) (*model.RevisionEntity, errors.Error) {
	wf, err := c.GetWorkflow(workflowName)
	if err != nil {
		return nil, err
	}

	if wf.Workflow.LatestRevision == nil || wf.Workflow.LatestRevision.ID == "" {
		return nil, errors.NewClientResponseNotFound()
	}

	return c.GetRevision(workflowName, wf.Workflow.LatestRevision.ID)
}

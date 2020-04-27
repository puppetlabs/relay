package client

import (
	"fmt"
	"net/http"

	"github.com/puppetlabs/relay/pkg/errors"
	"github.com/puppetlabs/relay/pkg/model"
)

type CreateWorkflowParameters struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *Client) CreateWorkflow(name string) (*model.WorkflowEntity, errors.Error) {
	params := &CreateWorkflowParameters{
		Name:        name,
		Description: "",
	}

	response := &model.WorkflowEntity{}

	if err := c.Request(
		WithMethod(http.MethodPost),
		WithPath("/api/workflows"),
		WithBody(params),
		WithResponseInto(response),
	); err != nil {
		return nil, err
	}

	return response, nil
}

type DeleteWorkflowResponse struct {
	Success    bool   `json:"success"`
	ResourceId string `json:"resource_id"`
}

func (c *Client) DeleteWorkflow(name string) (*DeleteWorkflowResponse, errors.Error) {
	response := &DeleteWorkflowResponse{}

	if err := c.Request(
		WithMethod(http.MethodDelete),
		WithPath(fmt.Sprintf("/api/workflows/%v", name)),
		WithResponseInto(response),
	); err != nil {
		return nil, err
	}

	return response, nil
}

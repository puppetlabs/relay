package client

import (
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

	if err := c.post("/api/workflows", nil, params, response); err != nil {
		return nil, err
	}

	return response, nil
}

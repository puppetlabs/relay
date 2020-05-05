package client

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

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

type ListWorkflowsResponse struct {
	Workflows []model.Workflow `json:"workflows"`
}

func (c *Client) ListWorkflows() (*ListWorkflowsResponse, errors.Error) {
	resp := &ListWorkflowsResponse{}

	if err := c.Request(WithPath("/api/workflows/"), WithResponseInto(&resp)); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetWorkflow(name string) (*model.WorkflowEntity, errors.Error) {
	response := &model.WorkflowEntity{}

	if err := c.Request(
		WithPath(fmt.Sprintf("/api/workflows/%v", name)),
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

type RunWorkflowRequest struct {
	Parameters map[string]string `json:"parameters"`
}

type RunWorkflowWorkflowResponse struct {
	Name string `json:"name"`
}

type RunWorkflowParameterValueResponse struct {
	Value string `json:"value"`
}

type RunWorkflowRevisionResponse struct {
	Id string `json:"id"`
}

type RunWorkflowStateResponse struct {
	Status    string     `json:"status"`
	StartedAt *time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`

	// TODO: Add steps here, in case we really care about that.
}

type RunWorkflowRunResponse struct {
	CreatedAt  time.Time                                    `json:"created_at"`
	RunNumber  int                                          `json:"run_number"`
	Revision   RunWorkflowRevisionResponse                  `json:"revision"`
	State      RunWorkflowStateResponse                     `json:"state"`
	Parameters map[string]RunWorkflowParameterValueResponse `json:"parameters"`
	Workflow   RunWorkflowWorkflowResponse                  `json:"workflow"`
}

type RunWorkflowResponse struct {
	Run RunWorkflowRunResponse `json:"run"`
}

func (c *Client) RunWorkflow(name string, params map[string]string) (*RunWorkflowResponse, errors.Error) {
	req := &RunWorkflowRequest{
		Parameters: params,
	}

	resp := &RunWorkflowResponse{}

	if err := c.Request(
		WithMethod(http.MethodPost),
		WithPath(fmt.Sprintf("/api/workflows/%v/runs", name)),
		WithBody(req),
		WithResponseInto(resp),
	); err != nil {
		return nil, err
	}

	return resp, nil
}

// DownloadWorkflow gets the latest configuration (as a YAML string) for a
// given workflow name. This is very purppose-built and likely rather frail. We
// should probably not be doing this this way.
func (c *Client) DownloadWorkflow(name string) (string, errors.Error) {
	workflow, err := c.GetWorkflow(name)

	if err != nil {
		return "", err
	}

	// TODO: Do we really want this to blow up or...
	revId := workflow.Workflow.LatestRevision.Id
	rev := &model.RevisionEntity{}

	if err := c.Request(
		WithPath(fmt.Sprintf("/api/workflows/%s/revisions/%s", name, revId)),
		WithResponseInto(rev),
	); err != nil {
		return "", err
	}

	dec, berr := base64.URLEncoding.DecodeString(rev.Revision.Raw)

	if berr != nil {
		return "", errors.NewClientUnkownError().WithCause(berr)
	}

	return string(dec), nil
}

// Code generated by go-swagger; DO NOT EDIT.

package workflow_runs_v1

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetWorkflowRunStepLogsParams creates a new GetWorkflowRunStepLogsParams object
// with the default values initialized.
func NewGetWorkflowRunStepLogsParams() *GetWorkflowRunStepLogsParams {
	var ()
	return &GetWorkflowRunStepLogsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetWorkflowRunStepLogsParamsWithTimeout creates a new GetWorkflowRunStepLogsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetWorkflowRunStepLogsParamsWithTimeout(timeout time.Duration) *GetWorkflowRunStepLogsParams {
	var ()
	return &GetWorkflowRunStepLogsParams{

		timeout: timeout,
	}
}

// NewGetWorkflowRunStepLogsParamsWithContext creates a new GetWorkflowRunStepLogsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetWorkflowRunStepLogsParamsWithContext(ctx context.Context) *GetWorkflowRunStepLogsParams {
	var ()
	return &GetWorkflowRunStepLogsParams{

		Context: ctx,
	}
}

// NewGetWorkflowRunStepLogsParamsWithHTTPClient creates a new GetWorkflowRunStepLogsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetWorkflowRunStepLogsParamsWithHTTPClient(client *http.Client) *GetWorkflowRunStepLogsParams {
	var ()
	return &GetWorkflowRunStepLogsParams{
		HTTPClient: client,
	}
}

/*GetWorkflowRunStepLogsParams contains all the parameters to send to the API endpoint
for the get workflow run step logs operation typically these are written to a http.Request
*/
type GetWorkflowRunStepLogsParams struct {

	/*Accept
	  The version of the API, in this case should be "application/nebula-api.v1+json"

	*/
	Accept string
	/*RunNumber
	  Inrecrmented run number of the associated workflow

	*/
	RunNumber int64
	/*StepName
	  Unique workflow step name

	*/
	StepName string
	/*WorkflowName
	  Workflow name. Must be unique within a user account

	*/
	WorkflowName string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) WithTimeout(timeout time.Duration) *GetWorkflowRunStepLogsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) WithContext(ctx context.Context) *GetWorkflowRunStepLogsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) WithHTTPClient(client *http.Client) *GetWorkflowRunStepLogsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAccept adds the accept to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) WithAccept(accept string) *GetWorkflowRunStepLogsParams {
	o.SetAccept(accept)
	return o
}

// SetAccept adds the accept to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) SetAccept(accept string) {
	o.Accept = accept
}

// WithRunNumber adds the runNumber to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) WithRunNumber(runNumber int64) *GetWorkflowRunStepLogsParams {
	o.SetRunNumber(runNumber)
	return o
}

// SetRunNumber adds the runNumber to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) SetRunNumber(runNumber int64) {
	o.RunNumber = runNumber
}

// WithStepName adds the stepName to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) WithStepName(stepName string) *GetWorkflowRunStepLogsParams {
	o.SetStepName(stepName)
	return o
}

// SetStepName adds the stepName to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) SetStepName(stepName string) {
	o.StepName = stepName
}

// WithWorkflowName adds the workflowName to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) WithWorkflowName(workflowName string) *GetWorkflowRunStepLogsParams {
	o.SetWorkflowName(workflowName)
	return o
}

// SetWorkflowName adds the workflowName to the get workflow run step logs params
func (o *GetWorkflowRunStepLogsParams) SetWorkflowName(workflowName string) {
	o.WorkflowName = workflowName
}

// WriteToRequest writes these params to a swagger request
func (o *GetWorkflowRunStepLogsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// header param Accept
	if err := r.SetHeaderParam("Accept", o.Accept); err != nil {
		return err
	}

	// path param run_number
	if err := r.SetPathParam("run_number", swag.FormatInt64(o.RunNumber)); err != nil {
		return err
	}

	// path param step_name
	if err := r.SetPathParam("step_name", o.StepName); err != nil {
		return err
	}

	// path param workflow_name
	if err := r.SetPathParam("workflow_name", o.WorkflowName); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
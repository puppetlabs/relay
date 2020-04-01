// Code generated by go-swagger; DO NOT EDIT.

package workflow_secrets

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewDeleteWorkflowSecretParams creates a new DeleteWorkflowSecretParams object
// with the default values initialized.
func NewDeleteWorkflowSecretParams() *DeleteWorkflowSecretParams {
	var ()
	return &DeleteWorkflowSecretParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteWorkflowSecretParamsWithTimeout creates a new DeleteWorkflowSecretParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDeleteWorkflowSecretParamsWithTimeout(timeout time.Duration) *DeleteWorkflowSecretParams {
	var ()
	return &DeleteWorkflowSecretParams{

		timeout: timeout,
	}
}

// NewDeleteWorkflowSecretParamsWithContext creates a new DeleteWorkflowSecretParams object
// with the default values initialized, and the ability to set a context for a request
func NewDeleteWorkflowSecretParamsWithContext(ctx context.Context) *DeleteWorkflowSecretParams {
	var ()
	return &DeleteWorkflowSecretParams{

		Context: ctx,
	}
}

// NewDeleteWorkflowSecretParamsWithHTTPClient creates a new DeleteWorkflowSecretParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDeleteWorkflowSecretParamsWithHTTPClient(client *http.Client) *DeleteWorkflowSecretParams {
	var ()
	return &DeleteWorkflowSecretParams{
		HTTPClient: client,
	}
}

/*DeleteWorkflowSecretParams contains all the parameters to send to the API endpoint
for the delete workflow secret operation typically these are written to a http.Request
*/
type DeleteWorkflowSecretParams struct {

	/*WorkflowName
	  Workflow name

	*/
	WorkflowName string
	/*WorkflowSecretName
	  The name of a workflow secret

	*/
	WorkflowSecretName string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the delete workflow secret params
func (o *DeleteWorkflowSecretParams) WithTimeout(timeout time.Duration) *DeleteWorkflowSecretParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete workflow secret params
func (o *DeleteWorkflowSecretParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete workflow secret params
func (o *DeleteWorkflowSecretParams) WithContext(ctx context.Context) *DeleteWorkflowSecretParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete workflow secret params
func (o *DeleteWorkflowSecretParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete workflow secret params
func (o *DeleteWorkflowSecretParams) WithHTTPClient(client *http.Client) *DeleteWorkflowSecretParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete workflow secret params
func (o *DeleteWorkflowSecretParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithWorkflowName adds the workflowName to the delete workflow secret params
func (o *DeleteWorkflowSecretParams) WithWorkflowName(workflowName string) *DeleteWorkflowSecretParams {
	o.SetWorkflowName(workflowName)
	return o
}

// SetWorkflowName adds the workflowName to the delete workflow secret params
func (o *DeleteWorkflowSecretParams) SetWorkflowName(workflowName string) {
	o.WorkflowName = workflowName
}

// WithWorkflowSecretName adds the workflowSecretName to the delete workflow secret params
func (o *DeleteWorkflowSecretParams) WithWorkflowSecretName(workflowSecretName string) *DeleteWorkflowSecretParams {
	o.SetWorkflowSecretName(workflowSecretName)
	return o
}

// SetWorkflowSecretName adds the workflowSecretName to the delete workflow secret params
func (o *DeleteWorkflowSecretParams) SetWorkflowSecretName(workflowSecretName string) {
	o.WorkflowSecretName = workflowSecretName
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteWorkflowSecretParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param workflowName
	if err := r.SetPathParam("workflowName", o.WorkflowName); err != nil {
		return err
	}

	// path param workflowSecretName
	if err := r.SetPathParam("workflowSecretName", o.WorkflowSecretName); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

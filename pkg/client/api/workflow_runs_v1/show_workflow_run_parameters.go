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

	strfmt "github.com/go-openapi/strfmt"
)

// NewShowWorkflowRunParams creates a new ShowWorkflowRunParams object
// with the default values initialized.
func NewShowWorkflowRunParams() *ShowWorkflowRunParams {
	var ()
	return &ShowWorkflowRunParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewShowWorkflowRunParamsWithTimeout creates a new ShowWorkflowRunParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewShowWorkflowRunParamsWithTimeout(timeout time.Duration) *ShowWorkflowRunParams {
	var ()
	return &ShowWorkflowRunParams{

		timeout: timeout,
	}
}

// NewShowWorkflowRunParamsWithContext creates a new ShowWorkflowRunParams object
// with the default values initialized, and the ability to set a context for a request
func NewShowWorkflowRunParamsWithContext(ctx context.Context) *ShowWorkflowRunParams {
	var ()
	return &ShowWorkflowRunParams{

		Context: ctx,
	}
}

// NewShowWorkflowRunParamsWithHTTPClient creates a new ShowWorkflowRunParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewShowWorkflowRunParamsWithHTTPClient(client *http.Client) *ShowWorkflowRunParams {
	var ()
	return &ShowWorkflowRunParams{
		HTTPClient: client,
	}
}

/*ShowWorkflowRunParams contains all the parameters to send to the API endpoint
for the show workflow run operation typically these are written to a http.Request
*/
type ShowWorkflowRunParams struct {

	/*Accept
	  The version of the API, in this case should be "application/nebula-api.v1+json"

	*/
	Accept string
	/*Authorization
	  The JWT bearer token

	*/
	Authorization string
	/*Rid
	  ID of the workflow run we want to know about

	*/
	Rid string
	/*Wid
	  ID of the workflow whose runs we want to view

	*/
	Wid string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the show workflow run params
func (o *ShowWorkflowRunParams) WithTimeout(timeout time.Duration) *ShowWorkflowRunParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the show workflow run params
func (o *ShowWorkflowRunParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the show workflow run params
func (o *ShowWorkflowRunParams) WithContext(ctx context.Context) *ShowWorkflowRunParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the show workflow run params
func (o *ShowWorkflowRunParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the show workflow run params
func (o *ShowWorkflowRunParams) WithHTTPClient(client *http.Client) *ShowWorkflowRunParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the show workflow run params
func (o *ShowWorkflowRunParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAccept adds the accept to the show workflow run params
func (o *ShowWorkflowRunParams) WithAccept(accept string) *ShowWorkflowRunParams {
	o.SetAccept(accept)
	return o
}

// SetAccept adds the accept to the show workflow run params
func (o *ShowWorkflowRunParams) SetAccept(accept string) {
	o.Accept = accept
}

// WithAuthorization adds the authorization to the show workflow run params
func (o *ShowWorkflowRunParams) WithAuthorization(authorization string) *ShowWorkflowRunParams {
	o.SetAuthorization(authorization)
	return o
}

// SetAuthorization adds the authorization to the show workflow run params
func (o *ShowWorkflowRunParams) SetAuthorization(authorization string) {
	o.Authorization = authorization
}

// WithRid adds the rid to the show workflow run params
func (o *ShowWorkflowRunParams) WithRid(rid string) *ShowWorkflowRunParams {
	o.SetRid(rid)
	return o
}

// SetRid adds the rid to the show workflow run params
func (o *ShowWorkflowRunParams) SetRid(rid string) {
	o.Rid = rid
}

// WithWid adds the wid to the show workflow run params
func (o *ShowWorkflowRunParams) WithWid(wid string) *ShowWorkflowRunParams {
	o.SetWid(wid)
	return o
}

// SetWid adds the wid to the show workflow run params
func (o *ShowWorkflowRunParams) SetWid(wid string) {
	o.Wid = wid
}

// WriteToRequest writes these params to a swagger request
func (o *ShowWorkflowRunParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// header param Accept
	if err := r.SetHeaderParam("Accept", o.Accept); err != nil {
		return err
	}

	// header param Authorization
	if err := r.SetHeaderParam("Authorization", o.Authorization); err != nil {
		return err
	}

	// path param rid
	if err := r.SetPathParam("rid", o.Rid); err != nil {
		return err
	}

	// path param wid
	if err := r.SetPathParam("wid", o.Wid); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

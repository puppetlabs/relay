// Code generated by go-swagger; DO NOT EDIT.

package access_control

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

// NewGetInvitesParams creates a new GetInvitesParams object
// with the default values initialized.
func NewGetInvitesParams() *GetInvitesParams {
	var ()
	return &GetInvitesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetInvitesParamsWithTimeout creates a new GetInvitesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetInvitesParamsWithTimeout(timeout time.Duration) *GetInvitesParams {
	var ()
	return &GetInvitesParams{

		timeout: timeout,
	}
}

// NewGetInvitesParamsWithContext creates a new GetInvitesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetInvitesParamsWithContext(ctx context.Context) *GetInvitesParams {
	var ()
	return &GetInvitesParams{

		Context: ctx,
	}
}

// NewGetInvitesParamsWithHTTPClient creates a new GetInvitesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetInvitesParamsWithHTTPClient(client *http.Client) *GetInvitesParams {
	var ()
	return &GetInvitesParams{
		HTTPClient: client,
	}
}

/*GetInvitesParams contains all the parameters to send to the API endpoint
for the get invites operation typically these are written to a http.Request
*/
type GetInvitesParams struct {

	/*Status
	  A filter to describe the invite status

	*/
	Status *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get invites params
func (o *GetInvitesParams) WithTimeout(timeout time.Duration) *GetInvitesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get invites params
func (o *GetInvitesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get invites params
func (o *GetInvitesParams) WithContext(ctx context.Context) *GetInvitesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get invites params
func (o *GetInvitesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get invites params
func (o *GetInvitesParams) WithHTTPClient(client *http.Client) *GetInvitesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get invites params
func (o *GetInvitesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithStatus adds the status to the get invites params
func (o *GetInvitesParams) WithStatus(status *string) *GetInvitesParams {
	o.SetStatus(status)
	return o
}

// SetStatus adds the status to the get invites params
func (o *GetInvitesParams) SetStatus(status *string) {
	o.Status = status
}

// WriteToRequest writes these params to a swagger request
func (o *GetInvitesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Status != nil {

		// query param status
		var qrStatus string
		if o.Status != nil {
			qrStatus = *o.Status
		}
		qStatus := qrStatus
		if qStatus != "" {
			if err := r.SetQueryParam("status", qStatus); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

package integration_providers

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

// NewAuthorizeProviderParams creates a new AuthorizeProviderParams object
// with the default values initialized.
func NewAuthorizeProviderParams() *AuthorizeProviderParams {
	var ()
	return &AuthorizeProviderParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewAuthorizeProviderParamsWithTimeout creates a new AuthorizeProviderParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewAuthorizeProviderParamsWithTimeout(timeout time.Duration) *AuthorizeProviderParams {
	var ()
	return &AuthorizeProviderParams{

		timeout: timeout,
	}
}

// NewAuthorizeProviderParamsWithContext creates a new AuthorizeProviderParams object
// with the default values initialized, and the ability to set a context for a request
func NewAuthorizeProviderParamsWithContext(ctx context.Context) *AuthorizeProviderParams {
	var ()
	return &AuthorizeProviderParams{

		Context: ctx,
	}
}

// NewAuthorizeProviderParamsWithHTTPClient creates a new AuthorizeProviderParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewAuthorizeProviderParamsWithHTTPClient(client *http.Client) *AuthorizeProviderParams {
	var ()
	return &AuthorizeProviderParams{
		HTTPClient: client,
	}
}

/*AuthorizeProviderParams contains all the parameters to send to the API endpoint
for the authorize provider operation typically these are written to a http.Request
*/
type AuthorizeProviderParams struct {

	/*ProviderID
	  The provider ID to reference

	*/
	ProviderID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the authorize provider params
func (o *AuthorizeProviderParams) WithTimeout(timeout time.Duration) *AuthorizeProviderParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the authorize provider params
func (o *AuthorizeProviderParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the authorize provider params
func (o *AuthorizeProviderParams) WithContext(ctx context.Context) *AuthorizeProviderParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the authorize provider params
func (o *AuthorizeProviderParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the authorize provider params
func (o *AuthorizeProviderParams) WithHTTPClient(client *http.Client) *AuthorizeProviderParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the authorize provider params
func (o *AuthorizeProviderParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithProviderID adds the providerID to the authorize provider params
func (o *AuthorizeProviderParams) WithProviderID(providerID string) *AuthorizeProviderParams {
	o.SetProviderID(providerID)
	return o
}

// SetProviderID adds the providerId to the authorize provider params
func (o *AuthorizeProviderParams) SetProviderID(providerID string) {
	o.ProviderID = providerID
}

// WriteToRequest writes these params to a swagger request
func (o *AuthorizeProviderParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param providerId
	if err := r.SetPathParam("providerId", o.ProviderID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

package integration_providers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/puppetlabs/nebula-cli/pkg/client/api/models"
)

// AuthorizeProviderReader is a Reader for the AuthorizeProvider structure.
type AuthorizeProviderReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AuthorizeProviderReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewAuthorizeProviderOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewAuthorizeProviderDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewAuthorizeProviderOK creates a AuthorizeProviderOK with default headers values
func NewAuthorizeProviderOK() *AuthorizeProviderOK {
	return &AuthorizeProviderOK{}
}

/*AuthorizeProviderOK handles this case with default header values.

Authorization information for the provider
*/
type AuthorizeProviderOK struct {
	Payload *AuthorizeProviderOKBody
}

func (o *AuthorizeProviderOK) Error() string {
	return fmt.Sprintf("[POST /api/providers/{providerId}/authorize][%d] authorizeProviderOK  %+v", 200, o.Payload)
}

func (o *AuthorizeProviderOK) GetPayload() *AuthorizeProviderOKBody {
	return o.Payload
}

func (o *AuthorizeProviderOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(AuthorizeProviderOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAuthorizeProviderDefault creates a AuthorizeProviderDefault with default headers values
func NewAuthorizeProviderDefault(code int) *AuthorizeProviderDefault {
	return &AuthorizeProviderDefault{
		_statusCode: code,
	}
}

/*AuthorizeProviderDefault handles this case with default header values.

An error occurred
*/
type AuthorizeProviderDefault struct {
	_statusCode int

	Payload *AuthorizeProviderDefaultBody
}

// Code gets the status code for the authorize provider default response
func (o *AuthorizeProviderDefault) Code() int {
	return o._statusCode
}

func (o *AuthorizeProviderDefault) Error() string {
	return fmt.Sprintf("[POST /api/providers/{providerId}/authorize][%d] authorizeProvider default  %+v", o._statusCode, o.Payload)
}

func (o *AuthorizeProviderDefault) GetPayload() *AuthorizeProviderDefaultBody {
	return o.Payload
}

func (o *AuthorizeProviderDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(AuthorizeProviderDefaultBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*AuthorizeProviderDefaultBody Error response
swagger:model AuthorizeProviderDefaultBody
*/
type AuthorizeProviderDefaultBody struct {

	// error
	Error *models.Error `json:"error,omitempty"`
}

// Validate validates this authorize provider default body
func (o *AuthorizeProviderDefaultBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateError(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *AuthorizeProviderDefaultBody) validateError(formats strfmt.Registry) error {

	if swag.IsZero(o.Error) { // not required
		return nil
	}

	if o.Error != nil {
		if err := o.Error.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("authorizeProvider default" + "." + "error")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *AuthorizeProviderDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *AuthorizeProviderDefaultBody) UnmarshalBinary(b []byte) error {
	var res AuthorizeProviderDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*AuthorizeProviderOKBody authorize provider o k body
swagger:model AuthorizeProviderOKBody
*/
type AuthorizeProviderOKBody struct {
	models.Entity

	// authorization
	Authorization models.ProviderAuth `json:"authorization,omitempty"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (o *AuthorizeProviderOKBody) UnmarshalJSON(raw []byte) error {
	// AuthorizeProviderOKBodyAO0
	var authorizeProviderOKBodyAO0 models.Entity
	if err := swag.ReadJSON(raw, &authorizeProviderOKBodyAO0); err != nil {
		return err
	}
	o.Entity = authorizeProviderOKBodyAO0

	// AuthorizeProviderOKBodyAO1
	var dataAuthorizeProviderOKBodyAO1 struct {
		Authorization models.ProviderAuth `json:"authorization,omitempty"`
	}
	if err := swag.ReadJSON(raw, &dataAuthorizeProviderOKBodyAO1); err != nil {
		return err
	}

	o.Authorization = dataAuthorizeProviderOKBodyAO1.Authorization

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (o AuthorizeProviderOKBody) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	authorizeProviderOKBodyAO0, err := swag.WriteJSON(o.Entity)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, authorizeProviderOKBodyAO0)

	var dataAuthorizeProviderOKBodyAO1 struct {
		Authorization models.ProviderAuth `json:"authorization,omitempty"`
	}

	dataAuthorizeProviderOKBodyAO1.Authorization = o.Authorization

	jsonDataAuthorizeProviderOKBodyAO1, errAuthorizeProviderOKBodyAO1 := swag.WriteJSON(dataAuthorizeProviderOKBodyAO1)
	if errAuthorizeProviderOKBodyAO1 != nil {
		return nil, errAuthorizeProviderOKBodyAO1
	}
	_parts = append(_parts, jsonDataAuthorizeProviderOKBodyAO1)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this authorize provider o k body
func (o *AuthorizeProviderOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with models.Entity
	if err := o.Entity.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (o *AuthorizeProviderOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *AuthorizeProviderOKBody) UnmarshalBinary(b []byte) error {
	var res AuthorizeProviderOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

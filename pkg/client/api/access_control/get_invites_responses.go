// Code generated by go-swagger; DO NOT EDIT.

package access_control

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/puppetlabs/nebula-cli/pkg/client/api/models"
)

// GetInvitesReader is a Reader for the GetInvites structure.
type GetInvitesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetInvitesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetInvitesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetInvitesDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetInvitesOK creates a GetInvitesOK with default headers values
func NewGetInvitesOK() *GetInvitesOK {
	return &GetInvitesOK{}
}

/*GetInvitesOK handles this case with default header values.

The list of invites
*/
type GetInvitesOK struct {
	Payload *GetInvitesOKBody
}

func (o *GetInvitesOK) Error() string {
	return fmt.Sprintf("[GET /api/invites][%d] getInvitesOK  %+v", 200, o.Payload)
}

func (o *GetInvitesOK) GetPayload() *GetInvitesOKBody {
	return o.Payload
}

func (o *GetInvitesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetInvitesOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetInvitesDefault creates a GetInvitesDefault with default headers values
func NewGetInvitesDefault(code int) *GetInvitesDefault {
	return &GetInvitesDefault{
		_statusCode: code,
	}
}

/*GetInvitesDefault handles this case with default header values.

An error occurred
*/
type GetInvitesDefault struct {
	_statusCode int

	Payload *GetInvitesDefaultBody
}

// Code gets the status code for the get invites default response
func (o *GetInvitesDefault) Code() int {
	return o._statusCode
}

func (o *GetInvitesDefault) Error() string {
	return fmt.Sprintf("[GET /api/invites][%d] getInvites default  %+v", o._statusCode, o.Payload)
}

func (o *GetInvitesDefault) GetPayload() *GetInvitesDefaultBody {
	return o.Payload
}

func (o *GetInvitesDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetInvitesDefaultBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*GetInvitesDefaultBody Error response
swagger:model GetInvitesDefaultBody
*/
type GetInvitesDefaultBody struct {

	// error
	Error *models.Error `json:"error,omitempty"`
}

// Validate validates this get invites default body
func (o *GetInvitesDefaultBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateError(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetInvitesDefaultBody) validateError(formats strfmt.Registry) error {

	if swag.IsZero(o.Error) { // not required
		return nil
	}

	if o.Error != nil {
		if err := o.Error.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getInvites default" + "." + "error")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetInvitesDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetInvitesDefaultBody) UnmarshalBinary(b []byte) error {
	var res GetInvitesDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*GetInvitesOKBody Response type for all pending and/or accepted invites for this account
swagger:model GetInvitesOKBody
*/
type GetInvitesOKBody struct {

	// A list of invites
	Invites []*models.Invite `json:"invites"`
}

// Validate validates this get invites o k body
func (o *GetInvitesOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateInvites(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetInvitesOKBody) validateInvites(formats strfmt.Registry) error {

	if swag.IsZero(o.Invites) { // not required
		return nil
	}

	for i := 0; i < len(o.Invites); i++ {
		if swag.IsZero(o.Invites[i]) { // not required
			continue
		}

		if o.Invites[i] != nil {
			if err := o.Invites[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getInvitesOK" + "." + "invites" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetInvitesOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetInvitesOKBody) UnmarshalBinary(b []byte) error {
	var res GetInvitesOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

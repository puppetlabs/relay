// Code generated by go-swagger; DO NOT EDIT.

package workflow_secrets

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/puppetlabs/nebula-cli/pkg/client/api/models"
)

// CreateWorkflowSecretReader is a Reader for the CreateWorkflowSecret structure.
type CreateWorkflowSecretReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateWorkflowSecretReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewCreateWorkflowSecretCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewCreateWorkflowSecretDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCreateWorkflowSecretCreated creates a CreateWorkflowSecretCreated with default headers values
func NewCreateWorkflowSecretCreated() *CreateWorkflowSecretCreated {
	return &CreateWorkflowSecretCreated{}
}

/*CreateWorkflowSecretCreated handles this case with default header values.

Summary of newly created secret
*/
type CreateWorkflowSecretCreated struct {
	Payload *CreateWorkflowSecretCreatedBody
}

func (o *CreateWorkflowSecretCreated) Error() string {
	return fmt.Sprintf("[POST /api/workflows/{workflowName}/secrets][%d] createWorkflowSecretCreated  %+v", 201, o.Payload)
}

func (o *CreateWorkflowSecretCreated) GetPayload() *CreateWorkflowSecretCreatedBody {
	return o.Payload
}

func (o *CreateWorkflowSecretCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(CreateWorkflowSecretCreatedBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateWorkflowSecretDefault creates a CreateWorkflowSecretDefault with default headers values
func NewCreateWorkflowSecretDefault(code int) *CreateWorkflowSecretDefault {
	return &CreateWorkflowSecretDefault{
		_statusCode: code,
	}
}

/*CreateWorkflowSecretDefault handles this case with default header values.

An error occurred
*/
type CreateWorkflowSecretDefault struct {
	_statusCode int

	Payload *CreateWorkflowSecretDefaultBody
}

// Code gets the status code for the create workflow secret default response
func (o *CreateWorkflowSecretDefault) Code() int {
	return o._statusCode
}

func (o *CreateWorkflowSecretDefault) Error() string {
	return fmt.Sprintf("[POST /api/workflows/{workflowName}/secrets][%d] createWorkflowSecret default  %+v", o._statusCode, o.Payload)
}

func (o *CreateWorkflowSecretDefault) GetPayload() *CreateWorkflowSecretDefaultBody {
	return o.Payload
}

func (o *CreateWorkflowSecretDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(CreateWorkflowSecretDefaultBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*CreateWorkflowSecretBody The request type for creating a workflow secret
swagger:model CreateWorkflowSecretBody
*/
type CreateWorkflowSecretBody struct {

	// name
	// Required: true
	Name *string `json:"name"`

	// value
	// Required: true
	Value models.BinaryString `json:"value"`
}

// Validate validates this create workflow secret body
func (o *CreateWorkflowSecretBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateValue(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *CreateWorkflowSecretBody) validateName(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"name", "body", o.Name); err != nil {
		return err
	}

	return nil
}

func (o *CreateWorkflowSecretBody) validateValue(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"value", "body", o.Value); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *CreateWorkflowSecretBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateWorkflowSecretBody) UnmarshalBinary(b []byte) error {
	var res CreateWorkflowSecretBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*CreateWorkflowSecretCreatedBody create workflow secret created body
swagger:model CreateWorkflowSecretCreatedBody
*/
type CreateWorkflowSecretCreatedBody struct {
	models.Entity

	// secret
	Secret *models.WorkflowSecretSummary `json:"secret,omitempty"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (o *CreateWorkflowSecretCreatedBody) UnmarshalJSON(raw []byte) error {
	// CreateWorkflowSecretCreatedBodyAO0
	var createWorkflowSecretCreatedBodyAO0 models.Entity
	if err := swag.ReadJSON(raw, &createWorkflowSecretCreatedBodyAO0); err != nil {
		return err
	}
	o.Entity = createWorkflowSecretCreatedBodyAO0

	// CreateWorkflowSecretCreatedBodyAO1
	var dataCreateWorkflowSecretCreatedBodyAO1 struct {
		Secret *models.WorkflowSecretSummary `json:"secret,omitempty"`
	}
	if err := swag.ReadJSON(raw, &dataCreateWorkflowSecretCreatedBodyAO1); err != nil {
		return err
	}

	o.Secret = dataCreateWorkflowSecretCreatedBodyAO1.Secret

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (o CreateWorkflowSecretCreatedBody) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	createWorkflowSecretCreatedBodyAO0, err := swag.WriteJSON(o.Entity)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, createWorkflowSecretCreatedBodyAO0)
	var dataCreateWorkflowSecretCreatedBodyAO1 struct {
		Secret *models.WorkflowSecretSummary `json:"secret,omitempty"`
	}

	dataCreateWorkflowSecretCreatedBodyAO1.Secret = o.Secret

	jsonDataCreateWorkflowSecretCreatedBodyAO1, errCreateWorkflowSecretCreatedBodyAO1 := swag.WriteJSON(dataCreateWorkflowSecretCreatedBodyAO1)
	if errCreateWorkflowSecretCreatedBodyAO1 != nil {
		return nil, errCreateWorkflowSecretCreatedBodyAO1
	}
	_parts = append(_parts, jsonDataCreateWorkflowSecretCreatedBodyAO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this create workflow secret created body
func (o *CreateWorkflowSecretCreatedBody) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with models.Entity
	if err := o.Entity.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateSecret(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *CreateWorkflowSecretCreatedBody) validateSecret(formats strfmt.Registry) error {

	if swag.IsZero(o.Secret) { // not required
		return nil
	}

	if o.Secret != nil {
		if err := o.Secret.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("createWorkflowSecretCreated" + "." + "secret")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *CreateWorkflowSecretCreatedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateWorkflowSecretCreatedBody) UnmarshalBinary(b []byte) error {
	var res CreateWorkflowSecretCreatedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*CreateWorkflowSecretDefaultBody Error response
swagger:model CreateWorkflowSecretDefaultBody
*/
type CreateWorkflowSecretDefaultBody struct {

	// error
	Error *models.Error `json:"error,omitempty"`
}

// Validate validates this create workflow secret default body
func (o *CreateWorkflowSecretDefaultBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateError(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *CreateWorkflowSecretDefaultBody) validateError(formats strfmt.Registry) error {

	if swag.IsZero(o.Error) { // not required
		return nil
	}

	if o.Error != nil {
		if err := o.Error.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("createWorkflowSecret default" + "." + "error")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *CreateWorkflowSecretDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *CreateWorkflowSecretDefaultBody) UnmarshalBinary(b []byte) error {
	var res CreateWorkflowSecretDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
